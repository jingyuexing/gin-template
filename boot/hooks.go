package boot

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"template/global"
	"time"

	"go.uber.org/zap"
)

const (
	defaultHookDir     = "hooks"
	defaultHookTimeout = 30 * time.Second // 可用 HOOK_TIMEOUT 覆盖
	envHookTimeoutKey  = "HOOK_TIMEOUT"   // 秒
	envHookDirKey      = "HOOK_DIR"
	hookPreStart       = "prestart"
	hookStarted        = "started"
	hookExit           = "exit"
)

type HookMode int

const (
	HookModeLenient HookMode = iota // 宽松模式：错误仅记录，不阻断
	HookModeStrict                  // 严格模式：出错中断
)

// 公共解释器映射
var commonInterpreters = map[string][]string{
	".py": {"python"},
	".js": {"node"},
	".ts": {"ts-node"},
	".rb": {"ruby"},
}

// 平台解释器映射
var platformInterpreters = map[string]map[string][]string{
	"windows": {
		".bat": {"cmd.exe", "/c"},
		".cmd": {"cmd.exe", "/c"},
		".ps1": {"powershell.exe", "-ExecutionPolicy", "Bypass", "-File"},
	},
	"linux": {
		".sh":   {"sh"},
		".bash": {"sh"},
	},
	"darwin": {
		".sh":   {"sh"},
		".bash": {"sh"},
	},
}

func hookTimeout() time.Duration {
	if s := os.Getenv(envHookTimeoutKey); s != "" {
		if v, err := time.ParseDuration(s); err == nil && v > 0 {
			return v
		}
		if v, err := time.ParseDuration(s + "s"); err == nil && v > 0 {
			return v
		}
	}
	return defaultHookTimeout
}

func hookDir() string {
	if d := os.Getenv(envHookDirKey); d != "" {
		return d
	}
	return defaultHookDir
}

func runHooks(phase string, extraEnv map[string]string, mode HookMode) error {
	path := filepath.Join(hookDir(), phase)

	info, err := os.Stat(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("stat hook path %s: %w", path, err)
		}
		if err := os.MkdirAll(hookDir(), 0755); err != nil {
			return fmt.Errorf("failed to create hooks dir: %w", err)
		}
	}

	if info.Mode().IsRegular() {
		// 如果是文件，直接执行
		return runHookScript(path, phase, extraEnv, mode)
	}

	if info.IsDir() {
		// 如果是目录，执行目录下所有脚本
		entries, err := os.ReadDir(path)
		if err != nil {
			return fmt.Errorf("read dir %s: %w", path, err)
		}
		// 排序保证顺序执行
		sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
		var wg sync.WaitGroup
		errCh := make(chan error, len(entries))
		// 并发执行 防止阻塞
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			script := filepath.Join(path, e.Name())
			wg.Add(1)
			go func(script string) {
				defer wg.Done()
				if err := runHookScript(script, phase, extraEnv, mode); err != nil {
					errCh <- err
				}
			}(script)
		}

		wg.Wait()
		close(errCh)

		// 如果严格模式，任何一个出错都返回
		if mode == HookModeStrict {
			if err, ok := <-errCh; ok {
				return err
			}
		}
	}

	return nil
}

// 构建命令
func buildHookCommand(ctx context.Context, path, ext string, info os.FileInfo) (*exec.Cmd, error) {
	ext = strings.ToLower(ext)

	// 公共解释器
	if args, ok := commonInterpreters[ext]; ok {
		if _, err := exec.LookPath(args[0]); err != nil {
			return nil, fmt.Errorf("interpreter not found: %s", args[0])
		}
		return exec.CommandContext(ctx, args[0], append(args[1:], path)...), nil
	}

	// 平台特定解释器
	if args, ok := platformInterpreters[runtime.GOOS][ext]; ok {
		if _, err := exec.LookPath(args[0]); err != nil {
			return nil, fmt.Errorf("interpreter not found: %s", args[0])
		}
		return exec.CommandContext(ctx, args[0], append(args[1:], path)...), nil
	}

	// Windows 特殊处理
	if runtime.GOOS == "windows" {
		if info.Mode()&0111 != 0 {
			return exec.CommandContext(ctx, path), nil
		}
		return exec.CommandContext(ctx, "cmd.exe", "/c", path), nil
	}

	// Linux / macOS：可执行文件
	if info.Mode()&0111 != 0 {
		return exec.CommandContext(ctx, path), nil
	}

	return nil, fmt.Errorf("unsupported hook file: %s", path)
}

func runHookScript(path string, phase string, extraEnv map[string]string, mode HookMode) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		return nil
	}

	ext := strings.ToLower(filepath.Ext(path))
	global.Logger.Info("hook", zap.String("phase", phase), zap.String("path", path))

	ctx, cancel := context.WithTimeout(context.Background(), hookTimeout())
	defer cancel()

	cmd, err := buildHookCommand(ctx, path, ext, info)
	if err != nil {
		global.Logger.Error("build hook command failed", zap.Error(err))
		return err
	}
	// 注入环境变量
	env := os.Environ()
	for k, v := range extraEnv {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	env = append(env, fmt.Sprintf("HOOK_PHASE=%s", phase))
	cmd.Env = env

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		exitCode := -1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
		global.Logger.Error("hook",
			zap.String("phase", phase),
			zap.String("path", path),
			zap.Int("exit_code", exitCode),
			zap.Error(err),
		)
		if mode == HookModeStrict {
			return fmt.Errorf("hook failed: %s exit_code=%d", path, exitCode)
		}
	}
	return nil
}
