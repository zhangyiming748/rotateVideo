package rotateVideo

import (
	"github.com/zhangyiming748/rotateVideo/rotate"
	"testing"
)

func TestUnit(t *testing.T) {
	src := "/Users/zen/Movies"
	dst := "/Users/zen/Movies/bilibili"
	pattern := "mp4"
	threads := "2"
	direction := "ToRight"
	rotate.Rotate(src, pattern, direction, dst, threads)
}
