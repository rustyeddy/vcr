package main

import "github.com/rs/zerolog/log"

// TLV Type, Length and Value
type TLV struct {
	buffer []byte
}

type TLVCallbacks struct {
	Type     int
	Callback func(tlv *TLV)
}

const (
	// General purpose tlvs
	TLVZero = iota
	TLVTerm
	TLVError

	// For the Video Player
	TLVPlay
	TLVPause
)

// NewTLV gets a new TLV ready to go
func NewTLV(typ, l byte) (t TLV) {
	if l < 2 {
		log.Fatal().Int("len", int(l)).Msg("TLV Len must be at least 2 bytes")
	}
	t.buffer = make([]byte, l)
	t.buffer[0] = typ
	t.buffer[1] = l

	return t
}

// Type of TLV
func (tlv *TLV) Type() int {
	return int(tlv.buffer[0])
}

// Type of TLV
func (tlv *TLV) Len() int {
	return int(tlv.buffer[1])
}

// TypeLen of TLV
func (tlv *TLV) TypeLen() (t int, l int) {
	return int(tlv.buffer[0]), int(tlv.buffer[1])
}

// Value of the TLV
func (tlv *TLV) Value() []byte {
	return tlv.buffer[2:]
}

func (tlv *TLV) Str() string {
	return string(tlv.buffer)
}
