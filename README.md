# remote-log-go

remote-log sdk go 版本。将日志内容按照统一格式通过 http 发送到日志采集层，支持缓存和压缩上传。

## 安装

```bash
go get -u github.com/doubility/remote-log-go
```

## 快速开始

可拷贝 examples 中的例子

**重点：Logger 申明为全局变量，初始化一次！！！**

logger/logger.go (根据情况替换 [应用名称])

```go
package logger

import (
	"log"

	remote_log_go "github.com/doubility/remote-log-go"
)

var Logger *remote_log_go.Logger

func init() {
    httpTransport := remote_log_go.NewHttpTransport(remote_log_go.Info, remote_log_go.Warn, remote_log_go.Error, remote_log_go.Access)
	consoleTransport := remote_log_go.NewConsoleTransport(remote_log_go.Debug)
	Logger = remote_log_go.NewLogger("应用名称", 120, httpTransport, consoleTransport)

	err := Logger.Init()
	if err != nil {
		log.Fatal("logger初始化失败" + err.Error())
	}
}
```

main.go (go-test 为 go.mod module)

```go
package main

import (
    "go-test/logger"
)

func main() {
	logger.Logger.Info("记录info日志") // http上传日志
	logger.Logger.Warn("记录warn日志") // http上传日志

	logger.Logger.Debug("debug日志") // console打印日志
}
```

## 详细说明

```code
// 日志类型 可在查询时筛选
-remote_log_go.Debug
-remote_log_go.Info
-remote_log_go.Warn
-remote_log_go.Error
-remote_log_go.Access
```

```go
import (
    remote_log_go "github.com/doubility/remote-log-go"
)

// 申明日志存储方式，一种日志类型可选择多种存储方式

// 日志上传到服务器
// (info、warn、error、access日志使用http上传到服务器)
httpTransport := remote_log_go.NewHttpTransport(remote_log_go.Info, remote_log_go.Warn, remote_log_go.Error, remote_log_go.Access)

// 日志输出到控制台
// (debug日志使用console打印)
consoleTransport := remote_log_go.NewConsoleTransport(remote_log_go.Debug)

// 实例化
// appName string 应用的名称（查询日志时可使用）
// storageDays number 日志存储天数 (最小30天，最大360天)
// transport transport ...interface{} 日志处理方式 接受HttpTransport和ConsoleTransport
Logger := remote_log_go.NewLogger(appName, storageDays, transport)

// 初始化
err := Logger.init()

// 记录各种类型的日志
Logger.Debug(string);
Logger.Info(string);
Logger.Warn(string);
Logger.Error(error);
Logger.Access(string);
```

## 注意

1、需要环境变量`NODE_APP_DATA`，上传失败的日志将保存在此目录下。

2、需要环境变量`REMOTE_LOG_API_URL`，上传日志的地址。

3、http 上传日志使用缓存上传，默认满足以下条件时上传 1 次：间隔 1 秒、缓存日志条数>=100、缓存日志总长度>=50000。如有特殊需求，可联系架构部。

4、缓存日志长度大于 1000 时压缩上传。
