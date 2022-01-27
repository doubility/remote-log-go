package remote_log_go

type Level string

const (
	Debug  Level = "debug"
	Info   Level = "info"
	Warn   Level = "warn"
	Error  Level = "error"
	Access Level = "access"
)
