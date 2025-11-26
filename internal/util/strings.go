package util

func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

func ClearKittyGraphics() string {
	return "\x1b_Ga=d,d=A\x1b\\"
}
