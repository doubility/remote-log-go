package main

import (
	"time"

	"github.com/doubility/remote-log-go/examples/logger"
)

func main() {
	logger.Logger.Info("记录info日志") // http上传日志
	logger.Logger.Warn("记录warn日志") // http上传日志

	logger.Logger.Debug("debug日志") // console打印日志

	time.Sleep(time.Second * 3)
}

// go run examples/main.go
