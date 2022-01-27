package remote_log_go

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Logger struct {
	appName     string
	storageDays int32
	serviceName string
	transport   []interface{}
}

/**
 * @description: 创建Logger类
 * @param {string} appName
 * @param {int32} storageDays
 * @param {[]interface{}} transport
 * @return {*}
 */
func NewLogger(appName string, storageDays int32, transport ...interface{}) *Logger {
	hostname, _ := os.Hostname()

	if appName == "" {
		panic(errors.New("appname cannot be empty"))
	}

	RemoteLogApiUrl = os.Getenv("REMOTE_LOG_API_URL")
	if RemoteLogApiUrl == "" {
		panic(errors.New("invalid env REMOTE_LOG_API_URL"))
	}

	goPath := os.Getenv("GO_APP_DATA")
	if goPath != "" {
		ErrorLogPath = fmt.Sprintf("%v/%v/remote_logs", goPath, appName)
		os.MkdirAll(ErrorLogPath, os.ModePerm)
	} else {
		panic(errors.New("invalid env GO_APP_DATA"))
	}

	return &Logger{
		appName:     appName,
		storageDays: storageDays,
		transport:   transport,
		serviceName: hostname,
	}
}

func (l *Logger) Init() error {
	params := url.Values{}
	params.Add("app", l.appName)
	params.Add("storageDays", fmt.Sprintf("%v", l.storageDays))
	params.Add("pwd", "b3981ef7-694b-11ec-a673-00163e1357b3")
	baseUrl, _ := url.Parse(RemoteLogApiUrl)
	baseUrl.Path = "api/appStorageDays"
	baseUrl.RawQuery = params.Encode()

	res, err := http.Get(baseUrl.String())
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if response.Code != 200 {
		return errors.New(response.Message)
	}

	return nil
}

func (l *Logger) Debug(message string) {
	l.log(Debug, message)
}

func (l *Logger) Info(message string) {
	l.log(Info, message)
}

func (l *Logger) Warn(message string) {
	l.log(Warn, message)
}

func (l *Logger) Error(message error) {
	l.log(Error, message.Error())
}

func (l *Logger) Access(message string) {
	l.log(Access, message)
}

func (l *Logger) log(level Level, message string) {
	if message == "" {
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	logInfo := &LogInfo{
		ServiceName: l.serviceName,
		AppName:     l.appName,
		Level:       string(level),
		LogTime:     now,
		Message:     message,
	}

	isLog := false //  // 是否日志已记录，未记录的日志console打印

	for _, item := range l.transport {
		switch v := item.(type) {
		case *HttpTransport:
			if v.ShouldLog(level) {
				isLog = true
				v.Log(logInfo)
			}
		case *ConsoleTransport:
			if v.ShouldLog(level) {
				isLog = true
				v.Log(logInfo)
			}
		}
	}

	if !isLog {
		consoleTransport := NewConsoleTransport()
		consoleTransport.Log(logInfo)
	}
}
