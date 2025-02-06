package config

import (
	"encoding/json"
	"fmt"
	"os"
	"template/common/dotenv"

	"github.com/jingyuexing/go-utils"
)

type BaseEnv struct {
	ConfigurePath string `json:"config" env:"config"`
	Mode string `json:"mode" env:"mode"`
	GinMode string `json:"gin_mode" env:"gin_mode"`
	Port int `json:"port" env:"port"`
	EnvMode string `json:"env_mode" env:"env_mode"`
	LoggerLanguage string `json:"logger_lang" env:"logger_lang"`
}

// TokenConfig represents the configuration for token expiration times
type TokenConfig struct {
	Type                   string `json:"type"`                     // token类型
	AccessTokenExpiration  string `json:"access_token_expiration"`  // 访问令牌的过期时间
	RefreshTokenExpiration string `json:"refresh_token_expiration"` // 刷新令牌的过期时间
}

type Database struct {
	Type     string `json:"type" env:"type"`
	Username string `json:"username" env:"username"`
	Password string `json:"password" env:"password"`
	Host     string `json:"host" env:"host"`
	Port     int    `json:"port" env:"port"`
	DBName   string `json:"dbname" env:"dbname"`
	Config   string `json:"config" env:"config"`
}

type SystemUserConfig struct {
	Sign string `json:"sign"`
}

type SystemAdminConfig struct {
	Sign string `json:"sign"`
}

type SystemVersion struct {
	System string `json:"system"`
}

type System struct {
	Language      string            `json:"lang"`
	User          SystemUserConfig  `json:"user"`
	Admin         SystemAdminConfig `json:"admin"`
	Version       SystemVersion     `json:"version"`
	Locale        string            `json:"locale" env:"locale"`
	HideVersion   bool              `json:"hide_version" env:"hide_version"`
	Token         TokenConfig       `json:"token"`
	Static        string            `json:"static" env:"static_path"`
	SwaggerEnable bool              `json:"swagger_enable" env:"swagger_enable"`
	LoggerLever   string            `json:"logger_level" env:"logger_level"`
}

type Configure struct {
	Env      BaseEnv  `json:"env"`
	Database Database `json:"database"`
	Redis    Database `json:"redis"`
	System   System   `json:"system"`
}

func LoadBaseConfig() *BaseEnv {
	fmt.Println("loading .env file")
	data, err := os.ReadFile(".env")
	if err != nil {
		fmt.Printf("Error reading .env file: %v\n", err)
		return nil
	}
	env := dotenv.New(string(data), "_.")
	baseEnv := &BaseEnv{}
	if err := env.Bind(baseEnv); err != nil {
		fmt.Printf("Error binding base env: %v\n", err)
		return nil
	}
	
	// 如果指定了mode，加载对应的环境配置
	if baseEnv.Mode != "" {
		modeEnvData, err := os.ReadFile(".env." + baseEnv.Mode)
		if (err == nil) {
			env.Load(string(modeEnvData))
			env.Bind(baseEnv)
		}
	}
	return baseEnv
}

func LoadingConfigure() *Configure {
	// 首先加载基础环境配置
	baseEnv := LoadBaseConfig()
	if baseEnv == nil {
		fmt.Println("Failed to load base configuration")
		return nil
	}

	// 使用ConfigurePath加载JSON配置
	config := utils.LoadConfig[Configure](baseEnv.ConfigurePath)
	
	// 重新加载环境变量以确保所有配置都被正确覆盖
	env := LoadEnv()
	if env != nil {
		env.Bind(&config.Env)
		env.Bind(&config.Database)
		env.Bind(&config.System)
	}

	// 保存合并后的配置
	WriteConfigure(config, baseEnv.ConfigurePath)
	return config
}

func WriteConfigure(config *Configure, configPath string) {
	// 将配置结构体编码为 JSON 格式
	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Printf("Error encoding config to JSON: %v", err)
	}

	// 打开或创建文件
	file, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Error opening or creating file: %v", err)
	}
	defer file.Close()

	// 将 JSON 数据写入文件
	if _, err := file.Write(jsonData); err != nil {
		fmt.Printf("Error writing JSON data to file: %v", err)
	}

}

func LoadEnv() *dotenv.DotENV {
	data, err := os.ReadFile(".env")
	if err != nil {
		return nil
	}
	env := dotenv.New(string(data), "_.")
	
	baseEnv := &BaseEnv{}
	if err := env.Bind(baseEnv); err != nil {
		return nil
	}

	if baseEnv.Mode != "" {
		modeData, err := os.ReadFile(".env." + baseEnv.Mode)
		if err == nil {
			env.Load(string(modeData))
		}
	}
	return env
}
