package dateutil

import (
	"fmt"
	"time"
)

func DateDecrease(t *time.Time) {
	*t = t.AddDate(0, 0, -1)
}

func GetStringDate(t time.Time) string {
	y, m, d := t.Date()
	sy, sm, sd := fmt.Sprint(y), fmt.Sprint(int(m)), fmt.Sprint(d)

	if m < 10 {
		sm = "0" + sm
	}
	if d < 10 {
		sd = "0" + sd
	}
	return sd + "/" + sm + "/" + sy
}
