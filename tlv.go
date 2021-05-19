package redeye

import "log"


// TLV Type, Length, Value is just a slice of bytes
// where the TYPE of packet is identified by tlv[0].
// The length of the packet is in tlv[1] the remaining
// slice tlv[2:] is the value, to be used as needed.
type TLV struct {
	tlv []byte
}

var (
	tlvCallbacks map[byte]func(tlv TLV)
	cmdQ chan TLV
)

func init() {
	cmdQ = make(chan TLV)
	tlvCallbacks = map[byte]func(tlv TLV){
		CMDZero:  cmdZero,
		CMDTerm:  cmdTerm,
		CMDError: cmdError,

		CMDPlay:  cmdPlay,
		CMDPause: cmdPause,
		CMDSnap:  cmdSnap,
	}
}

const (

	// General purpose tlvs
	CMDZero  byte = 0x0
	CMDTerm  byte = 0x1
	CMDError byte = 0x2

	// For the Video Player
	CMDPlay  byte = 0x11
	CMDPause byte = 0x12
	CMDSnap  byte = 0x13

	// For MJPEG Server
	CMDFrame byte = 0x21
)

// NewTLV gets a new TLV ready to go
func NewTLV(typ, l byte) (t TLV) {
	if l < 2 {
		panic("TLV Len must be at least 2 bytes")
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
	if t == nil || t.tlv == nil {
		return 0
	}
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

func cmdZero(tlv TLV) {
	if Config.Debug {
		log.Println("ADD CMD ZERO")		
	}

}

func cmdTerm(tlv TLV) {
	if Config.Debug {
		log.Println("ADD CMD TERM")		
	}
}

func cmdError(tlv TLV) {
	if Config.Debug {
		log.Println("ADD CMD ERROR")		
	}
}

func cmdPlay(tlv TLV) {
	//video.Play()
}

func cmdPause(tlv TLV) {
	//video.Pause()
}

func cmdSnap(tlv TLV) {
	//video.Snap()
}
