package rotate

import (
	"fmt"
	"github.com/zhangyiming748/rotateVideo/log"
	"os"
	"os/exec"
	"path"
	"strings"
)

func Rotate(src, pattern, direction, dst, threads string) {
	files := getFiles(src, pattern)
	for index, file := range files {
		log.Info.Printf("正在处理第 %d/%d 个文件\n", index+1, len(files))
		rotate_help(src, dst, file, direction, threads)
	}
}
func rotate_help(src, dst, file, direction, threads string) {
	var errorReport string
	in := strings.Join([]string{src, file}, "/")
	extname := path.Ext(file) //.txt
	filename := strings.Trim(file, extname)
	filename = replace(filename)
	export := strings.Join([]string{dst, strings.Join([]string{filename, "mp4"}, ".")}, "/")
	var cmd *exec.Cmd
	switch direction {
	case "ToRight":
		cmd = exec.Command("ffmpeg", "-threads", threads, "-i", in, "-vf", "transpose=1", "-threads", threads, export)
	case "ToLeft":
		cmd = exec.Command("ffmpeg", "-threads", threads, "-i", in, "-vf", "transpose=2", "-threads", threads, export)
	default:
		return
	}
	log.Debug.Printf("开始处理文件%s\t生成的命令是:%s", file, cmd)

	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	if err != nil {
		errorReport = strings.Join([]string{errorReport, fmt.Sprintf("cmd.StdoutPipe产生的错误:%v", err)}, "")
	}
	if err = cmd.Start(); err != nil {
		errorReport = strings.Join([]string{errorReport, fmt.Sprintf("cmd.Run产生的错误:%v", err)}, "")
	}
	// 从管道中实时获取输出并打印到终端
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		//写成输出日志
		t := string(tmp)
		t = strings.Replace(t, "\u0000", "", -1)
		log.Info.Println(t)
		if err != nil {
			break
		}
	}
	if err = cmd.Wait(); err != nil {
		errorReport = strings.Join([]string{errorReport, fmt.Sprintf("命令执行中有错误产生:%v", err)}, "")
	}
	log.Debug.Printf("完成当前文件的处理:dst是%s\tfile是%s\n", dst, file)
}

func replace(str string) string {
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "《", "", -1)
	str = strings.Replace(str, "》", "", -1)
	str = strings.Replace(str, "【", "", -1)
	str = strings.Replace(str, "】", "", -1)
	str = strings.Replace(str, "(", "", -1)
	str = strings.Replace(str, ")", "", -1)
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\u00A0", "", -1)
	str = strings.Replace(str, "\u0000", "", -1)
	return str
}
func getFiles(dir, pattern string) []string {
	files, _ := os.ReadDir(dir)
	var aim []string
	types := strings.Split(pattern, ";") //"wmv;rm"
	for _, f := range files {
		//fmt.Println(f.Name())
		if l := strings.Split(f.Name(), ".")[0]; len(l) != 0 {
			//log.Info.Printf("有效的文件:%v\n", f.Name())
			for _, v := range types {
				if strings.HasSuffix(f.Name(), v) {
					log.Debug.Printf("有效的目标文件:%v\n", f.Name())
					//absPath := strings.Join([]string{dir, f.Name()}, "/")
					//log.Printf("目标文件的绝对路径:%v\n", absPath)
					aim = append(aim, f.Name())
				}
			}
		}
	}
	return aim
}
