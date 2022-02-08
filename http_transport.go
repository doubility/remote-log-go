package remote_log_go

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
	"unsafe"
)

type HttpTransport struct {
	allowLevel      []Level          // 允许使用此方式的日志级别
	maxBufferLength int64            // 最大缓存字符串长度
	maxBufferSize   int64            // 最大日志条数
	bufferLog       []unsafe.Pointer // 缓存日志
	bufferLength    int64            // 缓存日志长度
	bufferChan      chan string
	t               *time.Ticker
}

/**
 * @description: 创建HttpTransport类
 * @param {...string} allowLevel
 * @return {*}
 */
func NewHttpTransport(allowLevel ...Level) *HttpTransport {
	h := &HttpTransport{
		allowLevel:      allowLevel,
		maxBufferLength: 50000,
		maxBufferSize:   100,
		bufferLog:       make([]unsafe.Pointer, 0, 200),
		bufferLength:    0,
		bufferChan:      make(chan string, 10000),
		t:               time.NewTicker(time.Millisecond * 1000),
	}
	// 定时执行任务、接受chan中的日志
	go h.createInterval()

	return h
}

func (h *HttpTransport) createInterval() {
	for {
		select {
		case <-h.t.C:
			h.flush()
		case logStr := <-h.bufferChan:
			h.bufferLength += int64(len(logStr))
			h.bufferLog = append(h.bufferLog, unsafe.Pointer(&logStr))
			if len(h.bufferLog) >= int(h.maxBufferSize) || h.bufferLength >= h.maxBufferLength {
				h.flush()
			}
		}
	}
}

/**
 * @description: 设置自动上传间隔
 * @param {int64} ms
 * @return {*}
 */
func (h *HttpTransport) SetFlushInterval(ms int64) {
	h.t.Reset(time.Millisecond * time.Duration(ms))
}

/**
 * @description: 设置最大缓存字符串长度
 * @param {int64} length
 * @return {*}
 */
func (h *HttpTransport) SetMaxBufferLength(length int64) {
	h.maxBufferLength = length
}

/**
 * @description: 设置最大缓存条数
 * @param {int64} size
 * @return {*}
 */
func (h *HttpTransport) SetMaxBufferSize(size int64) {
	h.maxBufferSize = size
	h.bufferLog = make([]unsafe.Pointer, 0, size*2)
}

/**
 * @description: 是否允许此方式记录日志
 * @param {Level} level
 * @return {*}
 */
func (h *HttpTransport) shouldLog(level Level) bool {
	for _, v := range h.allowLevel {
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
func (h *HttpTransport) log(log *LogInfo) {
	logStr := formatHttp(log)
	h.bufferChan <- logStr
}

/**
 * @description: 处理日志
 * @param {*}
 * @return {*}
 */
func (h *HttpTransport) flush() {
	if len(h.bufferLog) > 0 {
		arrStrBufferLog := []string{}
		for _, v := range h.bufferLog {
			arrStrBufferLog = append(arrStrBufferLog, *(*string)(v))
		}
		// 长度大于1000时压缩上传
		// 压缩失败时，原字符串上传
		if h.bufferLength > 1000 {
			bytesData, err := json.Marshal(arrStrBufferLog)
			if err != nil {
				sendLog(1, arrStrBufferLog, "", 0)
			}
			strLog, err := doZlibCompress(bytesData)
			if err != nil {
				sendLog(1, arrStrBufferLog, "", 0)
			} else {
				sendLog(2, []string{}, strLog, 0)
			}

		} else {
			sendLog(1, arrStrBufferLog, "", 0)
		}

		h.bufferLog = h.bufferLog[:0]
		h.bufferLength = 0
	}
}

/**
 * @description: 请求接口 上传日志 失败重试3次
 * @param {int32} _type
 * @param {[]string} data1
 * @param {string} data2
 * @return {*}
 */
func sendLog(_type int32, data1 []string, data2 string, tryNum int32) {
	defer func() {
		if err := recover(); err != nil {
			if tryNum >= 3 {
				data := make(map[string]interface{})
				data["type"] = _type
				data["data1"] = data1
				data["data2"] = data2
				bytesData, _ := json.Marshal(data)

				httpErrorLog(string(bytesData))

				fmt.Println(err)
			} else {
				go func() {
					time.Sleep(time.Second)
					tryNum++
					sendLog(_type, data1, data2, tryNum)
				}()
			}

		}
	}()

	data := make(map[string]interface{})
	data["type"] = _type
	data["data1"] = data1
	data["data2"] = data2
	bytesData, _ := json.Marshal(data)

	res, err := http.Post(RemoteLogApiUrl+"/api/collectLog?pwd=b3981ef7-694b-11ec-a673-00163e1357b3", "application/json", bytes.NewBuffer(bytesData))
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}

	if response.Code != 200 {
		panic(errors.New(response.Message))
	}
}

// 记录上传失败的日志
func httpErrorLog(log string) {
	file, _ := os.OpenFile(fmt.Sprintf("%v/error_log_%v.log", ErrorLogPath, time.Now().Format("2006-01-02")), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(log + "\n")
	write.Flush()
}

/**
 * @description: 日志格式化
 * @param {*LogInfo} log
 * @return {*}
 */
func formatHttp(log *LogInfo) string {
	var s strings.Builder
	s.WriteString(log.LogTime)
	s.WriteString("|**|")
	s.WriteString(log.Level)
	s.WriteString("|**|")
	s.WriteString(log.ServiceName)
	s.WriteString("|**|")
	s.WriteString(log.AppName)
	s.WriteString("|**|")
	s.WriteString(log.Message)
	return s.String()
}

/**
 * @description: 压缩字符串
 * @param {[]byte} src
 * @return {*}
 */
func doZlibCompress(src []byte) (string, error) {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	n, err := w.Write(src)
	if err != nil || n == 0 {
		return "", err
	}
	err = w.Close()
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(in.Bytes()), nil
}
