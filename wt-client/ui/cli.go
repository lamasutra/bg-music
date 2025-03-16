package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/lamasutra/bg-music/wt-client/types"
)

type cliState struct {
}

func NewCli() *cliState {
	return &cliState{}
}

func (s *cliState) Debug(args ...any) {
	length := len(args) + 1
	buf := make([]string, length)
	buf[0] = time.Now().Format("15:04:05.000")
	for i, val := range args {
		buf[i+1] = fmt.Sprint(val)
	}
	fmt.Println(strings.Join(buf, " "))
}

func (s *cliState) Error(args ...any) {
	newArgs := []any{"ERR:"}
	newArgs = append(newArgs, args...)
	s.Debug(newArgs...)
}

func (s *cliState) Input(in types.WtInput) {
	s.Debug(in)
}
