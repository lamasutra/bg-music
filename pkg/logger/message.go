package logger

import (
	"fmt"
	"strings"
	"time"
)

type MessageRenderer interface {
	Render() string
}

type message struct {
	level uint8
	time  time.Time
	msg   []any
}

func (m message) Render() string {
	length := len(m.msg) + 1
	buf := make([]string, length)
	buf[0] = m.time.Format("15:04:05.000")
	for i, val := range m.msg {
		buf[i+1] = fmt.Sprint(val)
	}
	return strings.Join(buf, " ")
}
