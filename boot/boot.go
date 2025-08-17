package boot

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"template/global"
	"template/router"
	"time"

	"go.uber.org/zap"
)

func Startup() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	if err := runHooks(hookPreStart, map[string]string{
		"PID":  fmt.Sprintf("%d", os.Getpid()),
		"PORT": fmt.Sprintf("%d", global.Config.Env.Port),
	}, HookModeLenient); err != nil {
		global.Logger.Error("run hooks failed", zap.Error(err))
	}

	app := router.Routers()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", global.Config.Env.Port),
		Handler: app,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
			os.Exit(1)
		}
	}()

	if err := runHooks(hookStarted, map[string]string{
		"PID":  fmt.Sprintf("%d", os.Getpid()),
		"PORT": fmt.Sprintf("%d", global.Config.Env.Port),
	}, HookModeLenient); err != nil {
		global.Logger.Error("run hooks failed", zap.Error(err))
	}
	sig := <-stop
	if err := runHooks(hookExit, map[string]string{
		"signal": fmt.Sprintf("%s", sig),
	}, HookModeLenient); err != nil {
		global.Logger.Error("run hooks failed", zap.Error(err))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("HTTP server shutdown error: %v\n", err)
	} else {
		fmt.Println("HTTP server shutdown successfully")
	}

	if global.DB != nil {
		if conn, err := global.DB.DB(); err == nil {
			_ = conn.Close()
			fmt.Println("Database connection closed")
		}
	}

	os.Exit(0)
}
