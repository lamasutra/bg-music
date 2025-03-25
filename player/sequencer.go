package player

import (
	"fmt"

	"github.com/gopxl/beep/v2"
	"github.com/lamasutra/bg-music/ui"
)

type sequencer struct {
	code   string
	size   int
	closed bool
	buffer []beep.Streamer
}

func NewBeepSequencer(maxSize int, code string) sequencer {
	return sequencer{
		code:   code,
		size:   0,
		buffer: make([]beep.Streamer, maxSize),
	}
}

func (s *sequencer) relocate(offset int) {
	// i := 0
	ui.Debug(fmt.Sprintf("seq(%s):relocating %d of seq: %p", s.code, offset, s))
	remaining := s.buffer[offset:]
	remainingSize := len(remaining)
	s.buffer = remaining
	// for i < len(s.buffer) {
	// 	if i < remainingSize {
	// 		s.buffer[i] = remaining[i]
	// 	} else {
	// 		s.buffer[i] = nil
	// 	}
	// 	i++
	// }
	s.size = remainingSize
	ui.Debug(fmt.Sprintf(" = new size %d", s.size))
}

func (s *sequencer) Stream(samples [][2]float64) (n int, ok bool) {
	if s.closed {
		ui.Debug(fmt.Sprintf("seq(%s):closed", s.code))
		return 0, false
	}
	// stream silence
	if s.size == 0 {
		ui.Debug(fmt.Sprintf("seq(%s):stream:empty", s.code))
		// return 0, true

		// for i < len(samples) {
		// 	samples[i] = [2]float64{0, 0}
		// 	i++
		// }
		// return len(samples), true
		// i := 0
		// for i < 32 {
		// 	samples[i] = [2]float64{0, 0}
		// 	i++
		// }
		// return 32, true
		return len(samples), true
	}

	ui.Debug(fmt.Sprintf("seq(%s):stream:buffer", s.code))
	for len(samples) > 0 {
		if s.buffer[0] == nil {
			ui.Debug(fmt.Sprintf("seq:buffer[0]:nil"))
			continue
		}
		sn, sok := s.buffer[0].Stream(samples)
		samples = samples[sn:]
		ui.Debug(fmt.Sprintf("seq(%s):stream:buffer[0]:sn %d, ok %t", s.code, sn, sok))

		n, ok = n+sn, ok || sok
		if !sok {
			ui.Debug("switching stream")
			s.relocate(1)
		}
	}

	// ui.Debug("seq:stream", samples)

	return n, ok
}

func (s *sequencer) Err() error {
	return nil
}

func (s *sequencer) Append(stream beep.Streamer) {
	s.buffer[s.size] = stream
	s.size++
	ui.Debug(fmt.Sprintf("appending stream %p to sequencer %p, buffer size: %d", stream, s, s.size))
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
