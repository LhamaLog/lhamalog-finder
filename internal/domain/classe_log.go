package domain

type ClasseLog string

const (
	ERROR   ClasseLog = "ERROR"
	WARNING ClasseLog = "WARNING"
	INFO    ClasseLog = "INFO"
)

func (c ClasseLog) Valida() bool { return c == ERROR || c == WARNING || c == INFO }
func (c ClasseLog) gravidade() int {
	switch c {
	case ERROR:
		return 3
	case WARNING:
		return 2
	case INFO:
		return 1
	}
	return 0
}
func MaisGrave(a, b ClasseLog) ClasseLog {
	if b.gravidade() > a.gravidade() {
		return b
	}
	return a
}
