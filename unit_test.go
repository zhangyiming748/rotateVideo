package rotateVideo

import (
	"testing"
)

func TestHelp(t *testing.T) {
	src := "/Users/zen/Movies"
	dst := "/Users/zen/Movies/bilibili"
	file := "電車.mp4"
	threads := "2"
	direction := "ToRight"
	rotate_help(src, dst, file, direction, threads, 1, 1)
}
