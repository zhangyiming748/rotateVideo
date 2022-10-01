package rotateVideo

import (
	"testing"
)

func TestUnit(t *testing.T) {
	src := "/Users/zen/Movies"
	dst := "/Users/zen/Movies/bilibili"
	pattern := "mp4"
	threads := "2"
	direction := "ToRight"
	Rotate(src, pattern, direction, dst, threads)
}
