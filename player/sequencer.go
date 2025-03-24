package player

import (
	"fmt"

	"github.com/gopxl/beep/v2"
	"github.com/lamasutra/bg-music/ui"
)

type sequencer struct {
	size   int
	buffer []*beep.Streamer
}

func NewBeepSequencer(maxSize int) *sequencer {
	return &sequencer{
		size:   0,
		buffer: make([]*beep.Streamer, maxSize),
	}
}

func (s *sequencer) relocate(offset int) {
	i := 0
	ui.Debug(fmt.Sprintf("relocating %d of seq: %p", offset, s))
	remaining := s.buffer[offset:]
	remainingSize := len(remaining)
	for i < len(s.buffer) {
		if i < remainingSize {
			s.buffer[i] = remaining[i]
		} else {
			s.buffer[i] = nil
		}
		i++
	}
	s.size = remainingSize
	ui.Debug(fmt.Sprintf(" = new size %d", s.size))
}

func (s *sequencer) Stream(samples [][2]float64) (n int, ok bool) {
	if s.size == 0 {
		ui.Debug("seq:stream:empty")
		return 0, true
		i := 0
		for i < len(samples) {
			samples[i] = [2]float64{0, 0}
			i++
		}
		return len(samples), true
	}

	ui.Debug("seq:stream:buffer", s.buffer)
	i := 0
	for i < len(s.buffer) && len(samples) > 0 {
		if s.buffer[i] == nil {
			// ui.Debug(fmt.Sprintf("seq:buffer[%d]:nil", i))
			continue
		}
		sn, sok := (*s.buffer[i]).Stream(samples)
		samples = samples[sn:]
		n, ok = n+sn, ok || sok
		if !sok {
			i++
		}
	}

	// ui.Debug("seq:stream", samples)

	if i > 0 {
		s.relocate(i)
	}

	return n, ok
}

func (s *sequencer) Err() error {
	return nil
}

func (s *sequencer) Append(stream *beep.Streamer) {
	s.buffer[s.size] = stream
	s.size++
	ui.Debug(fmt.Sprintf("appending stream %p to sequencer %p, buffer size: %d", stream, s, s.size))
}
