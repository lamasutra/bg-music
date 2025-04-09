package audio

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/flac"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/vorbis"
	"github.com/lamasutra/bg-music/internal/ui"
	"github.com/lamasutra/bg-music/model"
)

type sequencers struct {
	music    sequencer
	sfx      sequencer
	narrator sequencer
}

type Beep struct {
	initialized   bool
	muted         bool
	paused        bool
	currentMusic  *model.Music
	volumePercent uint8
	volumed       *effects.Volume
	musicStreamer beep.StreamSeekCloser
	sfxStreamer   beep.StreamSeekCloser
	mixer         *beep.Mixer
	sequencers    sequencers
	format        beep.Format
	sampleRate    beep.SampleRate
	crossfadeNum  int
	musicEnded    chan (bool)
	stopWatchEnd  chan (bool)
	playlist      *[]model.Music
	currentConfig *model.Config
}

var speechCache = make(map[string]*effects.Volume, 0)

func CreateBeepPlayer() *Beep {
	return &Beep{
		format: beep.Format{
			SampleRate:  44100,
			NumChannels: 2,
			Precision:   2,
		},
	}
}

func (p *Beep) Init() {
	p.musicEnded = make(chan bool)
	p.sequencers = sequencers{
		music:    NewBeepSequencer(8, "music"),
		sfx:      NewBeepSequencer(8, "sfx"),
		narrator: NewBeepSequencer(32, "narrator"),
	}
	p.mixer = &beep.Mixer{}
	p.mixer.KeepAlive(true)
	p.mixer.Add(&p.sequencers.music, &p.sequencers.sfx, &p.sequencers.narrator)
	p.crossfadeNum = p.format.SampleRate.N(time.Second) / 2
	ui.Debug("crossfadeNum", p.crossfadeNum)
	p.volumed = wrapStreamerByVolumeEffect(p.mixer)
	ui.Debug("set initial volume 100")
	ui.Debug("format", p.format.SampleRate, "Hz")
	p.SetVolume(100)

	err := speaker.Init(p.format.SampleRate, p.format.SampleRate.N(time.Second/10))
	if err != nil {
		panic(err)
	}

	ui.Debug(fmt.Sprintf("beep player initialized, seqs: music=%p sfx=%p narrator=%p", &p.sequencers.music, &p.sequencers.sfx, &p.sequencers.narrator))
	p.sampleRate = p.format.SampleRate
	speaker.Play(p.volumed)
	p.initialized = true
	go p.watchMusicStream()
}

func (p *Beep) SetPlaylist(playlist *[]model.Music) {
	p.playlist = playlist
}

func getFileExtension(path string) (string, error) {
	base := filepath.Base(strings.ToLower(path))
	extIndex := strings.LastIndex(base, ".")
	if extIndex == -1 {
		return "", fmt.Errorf("path does not contain a dot: %s", path)
	}

	return base[extIndex:], nil
}

func openFile(path string) (beep.StreamSeekCloser, beep.Format, error) {
	ext, err := getFileExtension(path)

	if err != nil {
		return nil, beep.Format{}, err
	}

	file, err := os.Open(path)

	if err != nil {
		return nil, beep.Format{}, fmt.Errorf("cannot read file %v", path)
	}

	switch ext {
	case ".mp3":
		// fmt.Println("mp3")
		return mp3.Decode(file)
	case ".flac":
		// fmt.Println("flac")
		return flac.Decode(file)
	case ".ogg":
		// fmt.Println("ogg/vorbis")
		return vorbis.Decode(file)
		// case "mid":
		// return midi.Decode(file)
	}

	return nil, beep.Format{}, fmt.Errorf("cannot decode file type %v", ext)
}

func (p *Beep) playSfx(streamer beep.Streamer) {
	fmt.Println("appending sfx")
	p.sequencers.sfx.Append(streamer)
}

func (p *Beep) SetVolume(volume uint8) {
	p.volumePercent = volume
	setVolume(p.volumed, volume)
}

func (p *Beep) VolumeUp() {
	ui.Debug("volume up")
	if p.volumePercent == 100 {
		return
	}
	p.volumePercent += 2
	if p.volumePercent > 100 {
		p.volumePercent = 100
	}
	setVolume(p.volumed, p.volumePercent)
}

func (p *Beep) VolumeDown() {
	ui.Debug("volume down")
	if p.volumePercent == 0 {
		return
	}
	p.volumePercent -= 2
	setVolume(p.volumed, p.volumePercent)
}
func (p *Beep) Mute() {
	ui.Debug("mute")
	if p.muted {
		p.muted = false
		setVolume(p.volumed, p.volumePercent)
	} else {
		p.muted = true
		setVolume(p.volumed, 0)
	}
}
func (p *Beep) Pause() {
	ui.Debug("pause")
	if p.paused {
		p.paused = false
		speaker.Resume()
	} else {
		p.paused = true
		speaker.Suspend()
	}
}

func (p *Beep) Next() {
	ui.Debug("next")
	music := p.getNextMusic()
	ui.Debug("next music:", music)

	p.PlayMusic(music, p.currentConfig, false)
}

func (p *Beep) Prev() {
	ui.Debug("prev")
	music := p.getPrevMusic()
	ui.Debug("prev music:", music)

	p.PlayMusic(music, p.currentConfig, false)
}

// @todo shuffle
func (p *Beep) getNextMusic() *model.Music {
	index := 0
	for i, m := range *p.playlist {
		if m.Path == p.currentMusic.Path {
			index = i + 1
			break
		}
	}
	if index == len(*p.playlist) {
		index = 0
	}

	music := (*p.playlist)[index]
	if music.Volume == 0 {
		music.Volume = int(p.volumePercent)
	}
	return &music
}
func (p *Beep) getPrevMusic() *model.Music {
	index := 0
	for i, m := range *p.playlist {
		if m.Path == p.currentMusic.Path {
			index = i - 1
			break
		}
	}
	if index < 0 {
		index = len(*p.playlist) - 1
	}

	music := (*p.playlist)[index]
	if music.Volume == 0 {
		music.Volume = int(p.volumePercent)
	}
	return &music
}

func (p *Beep) SendControl(ctrl string) {
	switch ctrl {
	case "pause":
		p.Pause()
	case "volume-up":
		p.VolumeUp()
	case "volume-down":
		p.VolumeDown()
	case "mute":
		p.Mute()
	case "next":
		p.Next()
	case "prev":
		p.Prev()
	}
}

func (p *Beep) PlayMusic(music *model.Music, conf *model.Config, allowSame bool) {
	p.currentConfig = conf
	if p.currentMusic != nil && music.Path == p.currentMusic.Path && !allowSame {
		ui.Debug("PlayMusic, same song, continue")
		return
	}
	ui.Debug("PlayMusic", music, music.Volume)
	if p.musicStreamer == nil {
		p.beginMusicPlay(music, conf)
	} else {
		p.crossfadeCurrentMusic(music, conf)
	}

	// ui.Debug(fmt.Sprintf("playing music %v, duration: %vs", path, streamer.Len()/int(format.SampleRate)))
	// ui.SetCurrentMusicTitle(music.Path)
}

func (p *Beep) GetMusicEndedChan() chan (bool) {
	return p.musicEnded
}

func (p *Beep) GetCurrentMusic() *model.Music {
	return p.currentMusic
}

func (p *Beep) GetCurrentMusicProgress() float64 {
	if p.musicStreamer == nil {
		ui.Error("cm is nil")
		return 0.0
	}
	if p.musicStreamer.Len() == 0 {
		ui.Error("cm is empty")
		return 0.0
	}

	return float64(p.musicStreamer.Position()) / float64(p.musicStreamer.Len())
}

func (p *Beep) PlaySfx(sfx *model.Sfx, conf *model.Config) {
	path := conf.GetSfxPath(sfx)

	streamer, format, err := openFile(path)
	if err != nil {
		ui.Error(err)
		return
	}

	p.sfxStreamer = streamer

	ui.Debug(sfx)

	volumeSfxStreamer := wrapStreamerByVolumeEffect(streamer)
	setVolume(volumeSfxStreamer, sfx.Volume)

	var preparedStreamer beep.Streamer

	if format.SampleRate != p.format.SampleRate {
		preparedStreamer = beep.Resample(3, format.SampleRate, p.format.SampleRate, volumeSfxStreamer)
	} else {
		preparedStreamer = volumeSfxStreamer
	}

	p.playSfx(preparedStreamer)

	ui.Debug("playing sfx ", path, sfx.Volume)
}

func (p *Beep) Speak(sentence *[]model.Speech, c *model.Config) {
	var volumed *effects.Volume
	if len(*sentence) == 0 {
		return
	}
	for _, speech := range *sentence {
		err := speech.Prepare(c)
		if err != nil {
			ui.Error(speech.Meaning, "cannot prepare speech", err)
			continue
		}
		resampler := speech.Resample(p.sampleRate)
		volumed = wrapStreamerByVolumeEffect(resampler)
		setVolume(volumed, speech.Volume)
		speechCache[speech.Meaning] = volumed
		p.sequencers.narrator.Append(volumed)
		// p.mixers.narrator.Add(volumed)
	}
}

func (p *Beep) Play(s beep.Streamer) {
	speaker.Play(s)
}

func (p *Beep) Close() {
	ui.Debug("closing beep")
	if p.musicEnded != nil {
		close(p.musicEnded)
	}
	if p.stopWatchEnd != nil {
		p.stopWatchEnd <- true
		time.Sleep(time.Millisecond * 5)
		close(p.stopWatchEnd)
	}
	p.sequencers.music.Close()
	p.sequencers.sfx.Close()
	p.sequencers.narrator.Close()

	speaker.Close()
}

func (p *Beep) watchMusicStream() {
	sleepTime := time.Millisecond * 100
	ui.Debug("entering watchStreamEnds")
	for {
		if p.musicStreamer == nil {
			// ui.Debug("no stream yet")
			time.Sleep(sleepTime)
			continue
			// } else {
			// ui.Debug(fmt.Sprintf("music streamer %p", p.musicStreamer))
		}

		select {
		case <-p.stopWatchEnd:
			ui.Debug("exiting watchStreamEnds")
			return
		default:
			// ui.
			// ui.Debug(p.musicStreamer.Position()+p.crossfadeNum, p.musicStreamer.Len())
			if (p.musicStreamer.Position() + p.crossfadeNum) >= p.musicStreamer.Len() {
				// ui.Debug("music ", p.currentMusic.Path, " ending ", " mem ", &p.currentMusic)
				p.musicEnded <- true
			}
		}
		// ui.SetCurrentMusicProgress(p.GetCurrentMusicProgress())
		// ui.Debug(p.musicStreamer.Position(), " ", p.crossfadeNum, " p ", p.GetCurrentMusicProgress(), " l ", p.musicStreamer.Len())
		time.Sleep(sleepTime)
	}
}

func wrapStreamerByVolumeEffect(streamer beep.Streamer) *effects.Volume {
	return &effects.Volume{
		Base:     2,
		Silent:   false,
		Streamer: streamer,
	}
}

func setVolume(volumeEffect *effects.Volume, volumePercent uint8) {
	if volumePercent == 0 {
		volumeEffect.Silent = true
	} else {
		volumeEffect.Silent = false
		realVolume := float64(volumePercent)/20 - 5 + 0.5
		volumeEffect.Volume = realVolume
		ui.Debug("setVolume on", volumeEffect, "to", realVolume, volumePercent)
	}

	// ui.SetCurrentVolume(float64(volumePercent) / 100)
}

func (p *Beep) beginMusicPlay(music *model.Music, conf *model.Config) {
	if p.musicStreamer != nil {
		ui.Debug("already playing music")
		return
	}
	path := conf.GetMusicPath(music)

	p.currentMusic = music

	streamer, format, err := openFile(path)

	if err != nil {
		panic("cannot open file " + path)
	}

	var volumed *effects.Volume

	ui.Debug("Not music yet")
	stream, ok := streamer.(beep.Streamer)
	if !ok {
		panic("invalid streamer type")
	}
	if format.SampleRate != p.format.SampleRate {
		ui.Debug("resampling", music.Path, "from", format.SampleRate, "to", p.format.SampleRate)
		volumed = wrapStreamerByVolumeEffect(beep.Resample(3, format.SampleRate, p.format.SampleRate, stream))
	} else {
		volumed = wrapStreamerByVolumeEffect(stream)
	}
	setVolume(volumed, uint8(music.Volume))
	p.sequencers.music.Append(volumed)
	p.musicStreamer = streamer
}

func (p *Beep) crossfadeCurrentMusic(music *model.Music, conf *model.Config) {
	ui.Debug("crossfading")
	streamer, format, err := openFile(conf.GetMusicPath(music))
	ui.Debug("file opened")
	if err != nil {
		ui.Error(err)
		return
	}
	p.currentMusic = music
	ui.Debug("current music set", music)
	if format.SampleRate != p.format.SampleRate {
		ui.Debug("resampling", music.Path, "from", format.SampleRate, "to", p.format.SampleRate)
		streamer = resampleStream(3, format.SampleRate, p.format.SampleRate, streamer)
	}

	ui.Debug("crossfadeNum", p.crossfadeNum, "volume percent", music.Volume)
	ui.Debug(fmt.Sprintf("current: %p, new: %p", p.musicStreamer, streamer))

	current := p.sequencers.music.GetCurrentStreamer()

	ui.Debug("preparing current sample")
	currentSample := beep.Take(p.crossfadeNum, current)
	ui.Debug(currentSample)
	ui.Debug("volume effect wrapping streamer")
	newVolumed := wrapStreamerByVolumeEffect(streamer)
	ui.Debug("volume setting wrapped streamer")
	setVolume(newVolumed, uint8(music.Volume))
	ui.Debug("preparing new sample")
	newSample := beep.Take(p.crossfadeNum, newVolumed)
	ui.Debug(newSample)
	ui.Debug("preparing crossfade")
	mixed := crossfade(currentSample, newSample, p.crossfadeNum)
	ui.Debug(mixed)
	// ui.Debug("seeq on old to the end")
	// p.musicStreamer.Seek(p.musicStreamer.Len())
	// current.Close()
	ui.Debug("replacing current music streamer")
	p.musicStreamer = streamer

	// ui.Debug("locking")
	// p.sequencers.music.Lock()

	// ui.Debug("seek on new")
	// streamer.Seek(p.crossfadeNum)

	ui.Debug("replacing current with crossfade sample")
	p.sequencers.music.ReplaceCurrent(mixed)

	ui.Debug("append new to seq")
	p.sequencers.music.Append(newVolumed)
	// ui.Debug("unlocking")
	// p.sequencers.music.Unlock()
}

// func (b *beepState) prepareSpeech(speech model.Speech) *beep.Streamer {
// 	var streamer *beep.Streamer
// 	// speechSampelRate := beep.SampleRate(speech.SampleRate)
// 	// if b.sampleRate != speechSampelRate {
// 	// 	beep.Resample(3, speechSampelRate, b.sampleRate)
// 	// }

// 	return streamer
// }

func crossfade(stream1 beep.Streamer, stream2 beep.Streamer, length int) beep.Streamer {
	trans1 := effects.Transition(stream1, length, 1.0, 0.0, effects.TransitionEqualPower)
	trans2 := effects.Transition(stream2, length, 0.0, 1.0, effects.TransitionEqualPower)
	mixed := beep.Take(length, beep.Mix(trans1, trans2))

	return mixed
}
