package remote_log_go

import (
	"sync"
	"testing"
	"time"
)

func TestMany(t *testing.T) {

	httpTransport := NewHttpTransport(Info, Warn, Error, Access)
	consoleTransport := NewConsoleTransport(Debug)
	log := NewLogger("go_app", 40, httpTransport, consoleTransport)

	err := log.Init()

	if err != nil {
		t.Errorf("初始化错误:%v", err.Error())
	}

	var wait sync.WaitGroup
	wait.Add(100)
	workM(log, &wait)
	wait.Wait()

	t.Log("执行完成")

	time.Sleep(time.Second * 20)
}

func workM(log *Logger, w *sync.WaitGroup) {
	for i := 0; i < 100; i++ {
		go func() {
			defer w.Done()
			for i := 0; i < 100; i++ {
				log.Info("go应用消息测试")
			}
		}()
	}
}

func TestOne(t *testing.T) {
	httpTransport := NewHttpTransport(Info, Warn, Error, Access)
	consoleTransport := NewConsoleTransport(Debug)
	log := NewLogger("go_app", 40, httpTransport, consoleTransport)

	// err := log.Init()

	// if err != nil {
	// 	t.Errorf("初始化错误:%v", err.Error())
	// }

	log.Info("go消息测试2")

	log.Info("go消息测试3")

	time.Sleep(time.Second * 5)
}
