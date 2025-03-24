package model

import (
	"fmt"
	"io"
	"os"

	"github.com/gopxl/beep/v2"
	"github.com/pkg/errors"
)

type Speech struct {
	Sfx
	Meaning string       `json:"meaning"`
	Skip    []int        `json:"skip"`
	Format  *beep.Format `json:"-"`
	dbuf    *[][2]float64
	pos     int
	size    int
	err     error
}

func (s *Speech) Prepare(conf *Config) error {
	var err error
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "sam")
		}
	}()
	file, err := os.Open(conf.Path + "/" + s.Path)
	if err != nil {
		s.err = err
		return err
	}

	// read sam
	var buf []byte
	wholeBuf, err := io.ReadAll(file)
	if err != nil {
		s.err = err
		return err
	}
	if len(s.Skip) > 0 {
		buf = wholeBuf[s.Skip[0] : len(wholeBuf)-int(s.Skip[1])]
	} else {
		buf = wholeBuf
	}

	// fmt.Println("read", file)

	numChannels := s.NumChannels
	if numChannels == 0 {
		numChannels = 1
	}

	s.Format = &beep.Format{
		SampleRate:  beep.SampleRate(s.SampleRate),
		NumChannels: numChannels,
		Precision:   1,
	}

	decoded := make([][2]float64, len(buf))
	for i, by := range buf {
		val := (float64(by)/(1<<8-1)*2 - 1)
		decoded[i] = [2]float64{val, val}
	}

	s.dbuf = &decoded
	s.size = len(decoded)

	// fmt.Println("s.Format", s.Format)
	// fmt.Println("s.size", s.size)

	return nil
}

func (s Speech) Resample(to beep.SampleRate) *beep.Resampler {
	// fmt.Println("resampling", s)
	return beep.Resample(3, s.Format.SampleRate, to, &s)
}

func (s *Speech) Stream(samples [][2]float64) (n int, ok bool) {
	if s.err != nil {
		return 0, false
	}
	if s.pos >= s.size {
		s.pos = 0
		// return 0, false
	}

	pointer := s.pos
	var val [2]float64
	ok = true
	for i := range samples {
		if pointer == s.size {
			// s.pos = 0
			pointer = 0
			// ok = false
			break
		}
		val = (*s.dbuf)[pointer]
		samples[i] = val
		pointer++
		n++
	}
	s.pos = pointer

	// fmt.Println(s.Meaning, s.pos)

	return n, ok
}

func (s *Speech) Len() int {
	return s.size
}

func (s *Speech) Position() int {
	return s.pos
}

func (s *Speech) Seek(p int) error {
	if p > s.size {
		return fmt.Errorf("invalid position %d", p)
	}
	s.pos = p

	return nil
}

func (s *Speech) Err() error {
	return s.err
}
