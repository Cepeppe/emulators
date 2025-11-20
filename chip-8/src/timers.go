package main

import (
	"fmt"
	"strings"
	"time"
)

// Update frequency is 60 Hz
// We use millis, so delta is 1000 ms / 60
const update_delta_millis = int64(1000 / 60)

// Operations on this data structure are atomic
// and must be performed using exposed API only
type TIMERS struct {
	//The sound timer
	ST uint8

	//The delay timer
	DT uint8

	// must be set when setting up ST and after every decrement
	ST_last_decrement_ts int64

	// must be set when setting up DT and after every decrement
	DT_last_decrement_ts int64
}

func (t *TIMERS) Init() {
	t.ST = uint8(0)
	t.DT = uint8(0)
}

func (t *TIMERS) SetST(value uint8) {
	t.ST = value
	t.ST_last_decrement_ts = time.Now().UnixMilli()
}

func (t *TIMERS) SetDT(value uint8) {
	t.DT = value
	t.DT_last_decrement_ts = time.Now().UnixMilli()
}

func (t *TIMERS) Update() {
	now := time.Now().UnixMilli()

	if t.ST != 0 && now-t.ST_last_decrement_ts >= update_delta_millis {
		t.ST--
		t.ST_last_decrement_ts = now
	}

	if t.DT != 0 && now-t.DT_last_decrement_ts >= update_delta_millis {
		t.DT--
		t.DT_last_decrement_ts = now
	}
}

func (t *TIMERS) Dump() {
	const innerWidth = 76

	printSep := func() {
		fmt.Printf("\t+%s+\n", strings.Repeat("-", innerWidth+2))
	}
	printLine := func(text string) {
		fmt.Printf("\t| %-*s |\n", innerWidth, text)
	}

	fmt.Printf("\n\tTIMERS STATE DUMP\n")
	printSep()
	printLine("TIMERS [bit view] (dec)")
	printSep()

	// 8-bit timers, shown with padded bit field
	printLine(fmt.Sprintf("ST (sound timer): %s (%3d)", formatBits8(t.ST), t.ST))
	printLine(fmt.Sprintf("DT (delay timer): %s (%3d)", formatBits8(t.DT), t.DT))

	printSep()
	printLine(fmt.Sprintf("ST last decrement ts: %d", t.ST_last_decrement_ts))
	printLine(fmt.Sprintf("DT last decrement ts: %d", t.DT_last_decrement_ts))
	printSep()
}

func (t *TIMERS) shouldBeep() bool {
	if t.ST > 0 {
		return true
	}
	return false
}
