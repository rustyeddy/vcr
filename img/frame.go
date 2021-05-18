package img

import "log"

type Frame struct {
	Buffer interface{}
}

func (f *Frame) Save(fname string) (err error) {
	log.Println("TODO: Must implement Save!!!")
	return err
}
