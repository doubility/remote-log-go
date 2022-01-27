package remote_log_go

import (
	"fmt"
	"os"
)

type ConsoleTransport struct {
	allowLevel []Level // 允许使用此方式的日志级别
}

/**
 * @description: 创建ConsoleTransport类
 * @param {...Level} allowLevel
 * @return {*}
 */
func NewConsoleTransport(allowLevel ...Level) *ConsoleTransport {
	return &ConsoleTransport{
		allowLevel: allowLevel,
	}
}

/**
 * @description: 是否允许此方式记录日志
 * @param {Level} level
 * @return {*}
 */
func (c *ConsoleTransport) shouldLog(level Level) bool {
	for _, v := range c.allowLevel {
		if v == level {
			return true
		}
	}

	return false
}

/**
 * @description: 记录日志
 * @param {*LogInfo} log
 * @return {*}
 */
func (c *ConsoleTransport) log(log *LogInfo) {
	logStr := formatConsole(log)
	if log.Level == string(Error) {
		fmt.Fprintln(os.Stderr, logStr)
	} else {
		fmt.Println(logStr)
	}
}

/**
 * @description: 日志格式化
 * @param {*logger.LogInfo} log
 * @return {*}
 */
func formatConsole(log *LogInfo) string {
	return fmt.Sprintf("%v %v %v %v %v", log.LogTime, log.Level, log.ServiceName, log.AppName, log.Message)
}
