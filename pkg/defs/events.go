package defs

type EventKind = uint8

const (

	// TEXT EVENT
	EventKindFilePath    EventKind = 0xC4
	EventKindFLP_Version EventKind = 0xC7
)

type Event interface {
	Kind() EventKind
	Content() any
}

type TextEvent struct {
	kind    EventKind
	content string
}

func (te *TextEvent) Kind() EventKind {
	return te.kind
}

func (te *TextEvent) Content() any {
	return te.content
}

type DWORDEvent struct {
	kind    EventKind
	content uint32
}

func (de *DWORDEvent) Kind() EventKind {
	return de.kind
}

func (de *DWORDEvent) Content() any {
	return de.content
}

type WORDEvent struct {
	kind    EventKind
	content uint16
}

func (we *WORDEvent) Kind() EventKind {
	return we.kind
}

func (we *WORDEvent) Content() any {
	return we.content
}

type BYTEEvent struct {
	kind    EventKind
	content uint8
}

func (be *BYTEEvent) Kind() EventKind {
	return be.kind
}

func (be *BYTEEvent) Content() any {
	return be.content
}

func NewWORDEvent(kind EventKind, content uint16) *WORDEvent {
	return &WORDEvent{kind: kind, content: content}
}

func NewDWORDEvent(kind EventKind, content uint32) *DWORDEvent {
	return &DWORDEvent{kind: kind, content: content}
}

func NewBYTEEvent(kind EventKind, content uint8) *BYTEEvent {
	return &BYTEEvent{kind: kind, content: content}
}

func NewTextEvent(kind EventKind, content string) *TextEvent {
	return &TextEvent{kind: kind, content: content}
}
