package main

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
)

const (
    sampleRate = 44100
    beepHz     = 440
    beepFrames = sampleRate / 60
)

var (
	audioContext *audio.Context
	beepPlayer   *audio.Player
)

func InitAudio() error {
	audioContext = audio.NewContext(sampleRate)

	// 16-bit PCM, stereo: 4 bytes per frame
	buf := make([]byte, beepFrames*4)

	period := sampleRate / beepHz

	for i := 0; i < beepFrames; i++ {
		// square wave
		var v int16
		if i%period < period/2 {
			v = 0x3FFF
		} else {
			v = -0x3FFF
		}

		off := i * 4

		buf[off+0] = byte(v)
		buf[off+1] = byte(v >> 8)

		buf[off+2] = byte(v)
		buf[off+3] = byte(v >> 8)
	}

	// NewPlayerFromBytes 
	beepPlayer = audioContext.NewPlayerFromBytes(buf)
	return nil
}

func PlayBeep() {
	if beepPlayer == nil {
		return
	}

	_ = beepPlayer.Rewind()
	beepPlayer.Play()
}
