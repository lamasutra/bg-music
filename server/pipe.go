package server

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func handlePipeFile(ch chan string, sleepTime time.Duration) {
	var buffer []byte
	pipeFile, err := os.OpenFile("event.pipe", os.O_CREATE|os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		fmt.Println(err)
		return
	}
	pipeFileReader := bufio.NewReader(pipeFile)

	for {
		buffer, err = pipeFileReader.ReadBytes('\n')
		if err != nil {
			// fmt.Println("pipe error", err)
		} else {
			ch <- strings.TrimRight(string(buffer), "\n")
		}
		time.Sleep(sleepTime)
	}
}
