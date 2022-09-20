package rotate

import "testing"

func TestHelp(t *testing.T) {
	src := "/Users/zen/Github/rotateVideo/file"
	dst := "/Users/zen/Github/rotateVideo/file/done"
	file := "電車.mp4"
	threads := "2"
	direction := "ToRight"
	rotate_help(src, dst, file, direction, threads)
}
