package redeye

import "time"

type TimeVals struct {
	Year  int        `json:"year"`
	Month time.Month `json:"month"`
	Day   int        `json:"day"`

	Hour   int `json:"hour"`
	Minute int `json:"minute"`
	Second int `json:"second"`

	Action string `json:"action"`
}

type TimeMsg struct {
	TimeVals `json"time"`
}

func NewTimeMsg(now time.Time) (tm *TimeMsg) {
	tv := TimeVals{
		Year:   now.Year(),
		Month:  now.Month(),
		Day:    now.Day(),
		Hour:   now.Hour(),
		Minute: now.Minute(),
		Second: now.Second(),
		Action: "setTime",
	}
	tm = &TimeMsg{tv}
	return tm
}
