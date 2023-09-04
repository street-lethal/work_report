package model

type Work struct {
	Contents []string
	Hours    float64
}

func (w *Work) AddContent(content string) {
	duplicates := false
	for _, c := range w.Contents {
		if c != content {
			continue
		}
		duplicates = true
		break
	}
	if !duplicates {
		w.Contents = append(w.Contents, content)
	}
}

func (w *Work) AddHour(hour float64) {
	w.Hours += hour
}

func (w Work) HourMin() (hour, minute int) {
	hour = int(w.Hours)
	minute = int(w.Hours*60) - hour*60

	return
}
