package audio

import (
	"fmt"

	"github.com/gopxl/beep/v2"
	"github.com/lamasutra/bg-music/pkg/logger"
)

type sequencer struct {
	locked bool
	dummy  *silentStreamer
	code   string
	size   int
	closed bool
	buffer []beep.Streamer
}

type silentStreamer struct{}

func (d *silentStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	i := 0
	for i < len(samples) {
		samples[i] = [2]float64{0, 0}
		i++
	}

	return len(samples), true
}

func (d *silentStreamer) Err() error {
	return nil
}

func NewBeepSequencer(maxSize int, code string) sequencer {
	return sequencer{
		code:   code,
		size:   0,
		buffer: make([]beep.Streamer, maxSize),
	}
}

func (s *sequencer) shift() {
	i := 0
	bufferSize := len(s.buffer)
	newBuffer := make([]beep.Streamer, bufferSize)
	// ui.Debug(fmt.Sprintf("seq(%s):relocating %d of seq: %p", s.code, 1, s))
	for i < bufferSize {
		if i+1 < bufferSize {
			newBuffer[i] = s.buffer[i+1]
		}
		i++
	}
	s.size--
	s.buffer = newBuffer
	// ui.Debug(fmt.Sprintf(" = new size %d", s.size))
	if s.buffer[0] == nil {
		s.buffer[0] = s.dummy
	}
}

func (s *sequencer) Stream(samples [][2]float64) (n int, ok bool) {
	if s.closed {
		logger.Debug(fmt.Sprintf("seq(%s):closed", s.code))
		return 0, false
	}

	// stream silence
	if s.size == 0 || s.locked {
		// ui.Debug(fmt.Sprintf("seq(%s):stream:empty", s.code))
		return s.dummy.Stream(samples)
	}

	// if s.code == "narrator" {
	// ui.Debug(fmt.Sprintf("seq(%s):stream:buffer", s.code))
	// }
	// ui.Debug(fmt.Sprintf("seq(%s):stream:buffer", s.code))
	for len(samples) > 0 {
		if s.buffer[0] == nil {
			// ui.Debug(fmt.Sprintf("seq(%s):buffer[0]:nil n(%d) samples", s.code, n))
			s.size = 0
			return s.dummy.Stream(samples)
			// return 32, true
			// return 0, true
			// continue
		}
		sn, sok := s.buffer[0].Stream(samples)
		samples = samples[sn:]
		// ui.Debug(fmt.Sprintf("seq(%s):stream:buffer[0]:sn %d, ok %t", s.code, sn, sok))

		n, ok = n+sn, ok || sok
		if !sok {
			logger.Debug(fmt.Sprintf("seq(%s):switching stream", s.code))
			if !s.locked {
				s.shift()
			}
			ok = true
		}
		// if !ok {
		// 	ui.Debug("! ok")
		// }
	}

	// if s.code == "narrator" {
	// ui.Debug("ok %d %d", n, ok)
	// }

	return n, ok
	// return len(samples), ok
}

func (s *sequencer) Err() error {
	return nil
}

func (s *sequencer) ReplaceCurrent(stream beep.Streamer) {
	logger.Debug(fmt.Sprintf("seq(%s):replacing stream %p, buffer size: %d", s.code, stream, s.size))
	closer, ok := s.buffer[0].(beep.StreamCloser)
	if ok {
		logger.Debug(fmt.Sprintf("seq(%s):closing stream %p, buffer size: %d", s.code, stream, s.size))
		closer.Close()
	}
	s.buffer[0] = stream
	i := 1
	for i < len(s.buffer) {
		s.buffer[i] = nil
		i++
	}
	s.size = 1
}

func (s *sequencer) Append(stream beep.Streamer) {
	s.buffer[s.size] = stream
	s.size++
	logger.Debug(fmt.Sprintf("seq(%s):appending stream %p, buffer size: %d", s.code, stream, s.size))
}

func (s *sequencer) GetCurrentStreamer() beep.Streamer {
	return s.GetStreamer(0)
}

func (s *sequencer) GetStreamer(index int) beep.Streamer {
	return s.buffer[index]
}

func (s *sequencer) Close() {
	s.closed = true
	i := 0
	for i < len(s.buffer) {
		closer, ok := s.buffer[i].(beep.StreamCloser)
		if ok {
			closer.Close()
		}
		i++
	}
}

func (s *sequencer) Lock() {
	s.locked = true
}

func (s *sequencer) Unlock() {
	s.locked = false
}
