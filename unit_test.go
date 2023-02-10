package rotateVideo

import (
	"testing"
)

func TestHelp(t *testing.T) {
	src := "/Users/zen/Downloads/整理/dance/梓/Left/done"
	pattern := "mp4"
	threads := "10"
	direction := "ToRight"
	Rotate(src, pattern, direction, threads)
}
