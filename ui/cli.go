package ui

import (
	"fmt"
	"strings"
	"time"
)

type cliState struct {
}

func NewCli() *cliState {
	return &cliState{}
}

func (s *cliState) Debug(args ...any) {
	var newArgs []any
	hasR := false
	if str, ok := args[0].(string); ok && strings.Contains(str, "\r") {
		args[0] = str[1:]
		hasR = true
	}
	if hasR {
		newArgs = []any{"\r", time.Now().Format("15:04:05.000")}
	} else {
		newArgs = []any{time.Now().Format("15:04:05.000")}
	}
	newArgs = append(newArgs, args...)
	fmt.Println(newArgs...)
}

func (s *cliState) Error(args ...any) {
	newArgs := []any{"ERR:"}
	newArgs = append(newArgs, args...)
	s.Debug(newArgs...)
}
