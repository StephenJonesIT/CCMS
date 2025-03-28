/*
 * @File: config.config.go
 * @Description: Defines config information of the service
 * @Author: Tran Thanh Sang (tranthanhsang.it.la@gmail.com)
 */
package config

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sync"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Configuration chứa tất cả cấu hình ứng dụng
type Configuration struct {
	Port                string `json:"port"`
	Environment         string `json:"environment"` // "development" hoặc "production"
	EnableGinConsoleLog bool   `json:"enableConsoleLog"`
	EnableGinFileLog    bool   `json:"enableGinFileLog"`
	LogFileName         string `json:"logFileName"`
	LogMaxSize         int    `json:"logMaxSize"`    // MB
	LogMaxBackups      int    `json:"logMaxBackups"` // Số lượng file log tối đa
	LogMaxAge          int    `json:"logMaxAge"`     // Số ngày lưu trữ
	LogLevel           string `json:"logLevel"`      // "debug", "info", "warn", "error"
}

var (
	Config *Configuration
	once   sync.Once
)

// LoadConfig tải cấu hình từ file (chỉ thực hiện 1 lần)
func LoadConfig(configPath string) error {
	var initErr error
	once.Do(func() {
		// Đảm bảo đường dẫn tuyệt đối
		absPath, err := filepath.Abs(configPath)
		if err != nil {
			initErr = err
			return
		}

		file, err := os.Open(absPath)
		if err != nil {
			initErr = err
			return
		}
		defer file.Close()

		Config = new(Configuration)
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&Config); err != nil {
			initErr = err
			return
		}

		// Thiết lập logging
		if err := setupLogging(); err != nil {
			initErr = err
			return
		}
	})

	return initErr
}

// setupLogging cấu hình hệ thống logging
func setupLogging() error {
	// Thiết lập log level
	level, err := log.ParseLevel(Config.LogLevel)
	if err != nil {
		level = log.InfoLevel // Mặc định là Info nếu có lỗi
	}
	log.SetLevel(level)

	// Thiết lập formatter
	if Config.Environment == "production" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
		})
	}

	// Cấu hình file logging nếu được bật
	if Config.EnableGinFileLog {
		logFile := &lumberjack.Logger{
			Filename:   Config.LogFileName,
			MaxSize:    Config.LogMaxSize,
			MaxBackups: Config.LogMaxBackups,
			MaxAge:     Config.LogMaxAge,
			Compress:   true, // Nén các file log cũ
		}

		// Ghi log ra cả file và console nếu được bật
		if Config.EnableGinConsoleLog {
			mw := io.MultiWriter(os.Stdout, logFile)
			log.SetOutput(mw)
		} else {
			log.SetOutput(logFile)
		}
	} else if Config.EnableGinConsoleLog {
		// Chỉ ghi log ra console
		log.SetOutput(os.Stdout)
	}

	log.Info("Logging configuration initialized")
	return nil
}

// GetConfig trả về bản sao của cấu hình để tránh thay đổi ngẫu nhiên
func GetConfig() Configuration {
	return *Config
}
