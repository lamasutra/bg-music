package player

import "github.com/gopxl/beep/v2"

type resampler struct {
	streamer  beep.Streamer
	resampler *beep.Resampler
}

func resampleStream(quality int, old, new beep.SampleRate, s beep.Streamer) *resampler {
	return &resampler{
		streamer:  s,
		resampler: beep.Resample(quality, old, new, s),
	}
}

func (r *resampler) Stream(samples [][2]float64) (n int, ok bool) {
	return r.resampler.Stream(samples)
}

func (r *resampler) Err() error {
	return r.resampler.Err()
}

func (r *resampler) Ratio() float64 {
	return r.resampler.Ratio()
}

func (r *resampler) Len() int {
	val, ok := r.streamer.(beep.StreamSeeker)
	if !ok {
		return 0
	}

	return val.Len()
}

func (r *resampler) Position() int {
	val, ok := r.streamer.(beep.StreamSeeker)
	if !ok {
		return 0
	}

	return val.Position()
}

func (r *resampler) Seek(p int) error {
	val, ok := r.streamer.(beep.StreamSeeker)
	if !ok {
		return nil
	}

	return val.Seek(p)
}

func (r *resampler) Close() error {
	val, ok := r.streamer.(beep.StreamCloser)
	if !ok {
		return nil
	}

	return val.Close()
}
