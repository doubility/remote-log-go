package remote_log_go

import "testing"

func BenchmarkLogger(b *testing.B) {
	httpTransport := NewHttpTransport(Info, Warn, Error, Access)
	consoleTransport := NewConsoleTransport(Debug)
	log := NewLogger("go_app", 40, httpTransport, consoleTransport)

	log.Init()

	for i := 0; i < b.N; i++ {
		log.Info("go消息测试")
	}
}
