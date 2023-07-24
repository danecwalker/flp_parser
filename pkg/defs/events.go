package defs

import (
	"fmt"
	"strings"
)

type EventKind = uint8

const (

	// TEXT EVENT
	EventKindProjectName EventKind = 0xCA
	EventKindFilePath    EventKind = 0xC4
	EventKindFLP_Version EventKind = 0xC7
)

type Event interface {
	Kind() EventKind
	Value() any
}

type TextEvent struct {
	kind      EventKind
	sizeCount uint8
	value     string
}

func (te *TextEvent) Kind() EventKind {
	return te.kind
}

func (te *TextEvent) Value() any {
	return te.value
}

type DWORDEvent struct {
	kind  EventKind
	value uint32
}

func (de *DWORDEvent) Kind() EventKind {
	return de.kind
}

func (de *DWORDEvent) Value() any {
	return de.value
}

type WORDEvent struct {
	kind  EventKind
	value uint16
}

func (we *WORDEvent) Kind() EventKind {
	return we.kind
}

func (we *WORDEvent) Value() any {
	return we.value
}

type BYTEEvent struct {
	kind  EventKind
	value uint8
}

func (be *BYTEEvent) Kind() EventKind {
	return be.kind
}

func (be *BYTEEvent) Value() any {
	return be.value
}

func NewWORDEvent(kind EventKind, value uint16) *WORDEvent {
	return &WORDEvent{kind: kind, value: value}
}

func NewDWORDEvent(kind EventKind, value uint32) *DWORDEvent {
	return &DWORDEvent{kind: kind, value: value}
}

func NewBYTEEvent(kind EventKind, value uint8) *BYTEEvent {
	return &BYTEEvent{kind: kind, value: value}
}

func NewTextEvent(kind EventKind, sizeCount uint8, value string) *TextEvent {
	return &TextEvent{kind: kind, sizeCount: sizeCount, value: value}
}

func ModTextEvent(ev *TextEvent, value string) {
	fmt.Println(ev.sizeCount)
	split := strings.Split(value, "")
	j := strings.Join(split, "\x00") + "\x00" + strings.Repeat("\x00", int(ev.sizeCount)+1)
	ev.value = j
}
