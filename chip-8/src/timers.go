package main

import (
	"fmt"
	"strings"
	"time"
)

// Update frequency is 60 Hz
// We use millis, so delta is 1000 ms / 60
const update_delta_millis = uint64((time.Second / 60) / time.Millisecond)

// Operations on this data structure are atomic
// and must be performed using exposed API only
type TIMERS struct {
	//The sound timer
	ST uint8

	//The delay timer
	DT uint8

	// must be set when setting up ST and after every decrement
	ST_last_decrement_ts uint64
	ST_next_decrement_ts uint64

	// must be set when setting up DT and after every decrement
	DT_last_decrement_ts uint64
	DT_next_decrement_ts uint64
}

func (t *TIMERS) Init() {
	t.ST = uint8(0)
	t.DT = uint8(0)
}

func (t *TIMERS) SetST(value uint8) {
	t.ST = value
	t.ST_last_decrement_ts = uint64(time.Now().UnixMilli())
	t.ST_next_decrement_ts = t.ST_last_decrement_ts + update_delta_millis
}

func (t *TIMERS) SetDT(value uint8) {
	t.DT = value
	t.DT_last_decrement_ts = uint64(time.Now().UnixMilli())
	t.DT_next_decrement_ts = t.DT_last_decrement_ts + update_delta_millis
}

// Check current time and decide if decrement is needed
func (t *TIMERS) Update() {
	now := uint64(time.Now().UnixMilli())

	for t.ST > 0 && now >= t.ST_next_decrement_ts {
		t.ST--
		t.ST_last_decrement_ts = t.ST_next_decrement_ts
		t.ST_next_decrement_ts += update_delta_millis
	}

	for t.DT > 0 && now >= t.DT_next_decrement_ts {
		t.DT--
		t.DT_last_decrement_ts = t.DT_next_decrement_ts
		t.DT_next_decrement_ts += update_delta_millis
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
	return t.ST > 0
}
