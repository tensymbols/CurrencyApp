package dateutil

import (
	"fmt"
	"time"
)

func DateDecrease(t *time.Time) {
	*t = t.AddDate(0, 0, -1)
}

func GetStringDate(t time.Time) string {
	y, m, d := t.Date()                                            // год месяц день
	sy, sm, sd := fmt.Sprint(y), fmt.Sprint(int(m)), fmt.Sprint(d) // год месяц день в виде строки
	// добавляем к месяцу 0 если он меньше 10 для корректного формата
	if m < 10 {
		sm = "0" + sm
	}
	// добавляем к дню 0 если он меньше 10 для корректного формата
	if d < 10 {
		sd = "0" + sd
	}
	return sd + "/" + sm + "/" + sy
}
