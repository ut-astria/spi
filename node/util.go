package node

import (
	"encoding/json"
	"fmt"
	"time"
)

func JSON(x interface{}) string {
	js, err := json.Marshal(&x)
	if err != nil {
		return fmt.Sprintf("%#v", x)
	}
	return string(js)
}

func RoundTime(t time.Time, resolution time.Duration) time.Time {
	ns := t.UnixNano()
	r := time.Duration(ns) / resolution
	r *= resolution
	return time.Unix(0, r.Nanoseconds()).UTC()
}
