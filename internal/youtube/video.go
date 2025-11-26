package youtube

import "fmt"

type Video struct {
	ID           string  `json:"id"`
	Title        string  `json:"title"`
	Duration     float64 `json:"duration"`
	Channel      string  `json:"uploader"`
	ThumbnailURL string
}

func (v *Video) GetThumbnailURL() string {
	if v.ThumbnailURL != "" {
		return v.ThumbnailURL
	}

	return fmt.Sprintf("https://img.youtube.com/vi/%s/hqdefault.jpg", v.ID)
}

func (v *Video) GetFallbackThumbnailURL() string {
	return fmt.Sprintf("https://img.youtube.com/vi/%s/hqdefault.jpg", v.ID)
}

func (v *Video) FormatDuration() string {
	return FormatDuration(v.Duration)
}

func FormatDuration(seconds float64) string {
	secInt := int(seconds)
	hours := secInt / 3600
	minutes := (secInt % 3600) / 60
	secs := secInt % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, secs)
	}
	return fmt.Sprintf("%d:%02d", minutes, secs)
}
