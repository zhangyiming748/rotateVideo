package rotateVideo

import (
	"testing"
)

func TestHelp(t *testing.T) {
	src := "/Volumes/T7/slacking/Telegram/未整理/dance/toRight/h265"
	pattern := "mp4"
	threads := "10"
	direction := "ToRight"
	Rotate(src, pattern, direction, threads)
}
