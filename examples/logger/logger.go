package logger

import (
	"log"

	remote_log_go "github.com/doubility/remote-log-go"
)

var Logger *remote_log_go.Logger

func init() {
	httpTransport := remote_log_go.NewHttpTransport(remote_log_go.Info, remote_log_go.Warn, remote_log_go.Error, remote_log_go.Access)
	consoleTransport := remote_log_go.NewConsoleTransport(remote_log_go.Debug)
	Logger = remote_log_go.NewLogger("go_app", 120, httpTransport, consoleTransport)

	err := Logger.Init()
	if err != nil {
		log.Fatal("logger初始化失败" + err.Error())
	}
}
