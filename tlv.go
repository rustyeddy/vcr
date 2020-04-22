package main

import "github.com/rs/zerolog/log"

// TLV Type, Length, Value is just a slice of bytes
// where the TYPE of packet is identified by tlv[0].
// The length of the packet is in tlv[1] the remaining
// slice tlv[2:] is the value, to be used as needed.
type TLV struct {
	tlv []byte
}

type TLVCallbacks struct {
	Type     int
	Callback func(tlv *TLV)
}

const (
	// General purpose tlvs
	TLVZero  byte = 0x0
	TLVTerm  byte = 0x1
	TLVError byte = 0x2

	// For the Video Player
	TLVPlay  byte = 0x11
	TLVPause byte = 0x12
	TLVSnap  byte = 0x13

	// For MJPEG Server
	TLVFrame byte = 0x21
)

// NewTLV gets a new TLV ready to go
func NewTLV(typ, l byte) (t TLV) {
	if l < 2 {
		log.Fatal().Int("len", int(l)).Msg("TLV Len must be at least 2 bytes")
	}
	t.tlv = make([]byte, l)
	t.tlv[0] = typ
	t.tlv[1] = l

	return t
}

// Type of TLV
func (t *TLV) Type() byte {
	return t.tlv[0]
}

// Type of TLV
func (t *TLV) Len() int {
	return int(t.tlv[1])
}

// TypeLen of TLV
func (t *TLV) TypeLen() (ty int, l int) {
	return int(t.tlv[0]), int(t.tlv[1])
}

// Value of the TLV
func (t *TLV) Value() []byte {
	return t.tlv[2:]
}

func (t *TLV) Str() string {
	return string(t.tlv)
}
