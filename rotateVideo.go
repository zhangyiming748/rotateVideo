package rotateVideo

import (
	"github.com/zhangyiming748/GetFileInfo"
	"github.com/zhangyiming748/log"
	"github.com/zhangyiming748/replace"
	"github.com/zhangyiming748/voiceAlert"
	"os"
	"os/exec"
	"strings"
)

func Rotate(src, pattern, direction, threads string) {
	files := GetFileInfo.GetAllFileInfo(src, pattern)
	for index, file := range files {
		log.Debug.Printf("正在处理第 %d/%d 个文件\n", index+1, len(files))
		rotate(file, direction, threads)
		log.Debug.Printf("处理完成第 %d/%d 个文件\n", index+1, len(files))
		voiceAlert.Customize("done", voiceAlert.Samantha)
	}
	voiceAlert.Customize("complete", voiceAlert.Samantha)
}
func rotate(in GetFileInfo.Info, direction, threads string) {
	defer func() {
		if err := recover(); err != nil {
			voiceAlert.Customize("failed", voiceAlert.Samantha)
		}
	}()
	dst := strings.Join([]string{strings.Trim(in.FullPath, in.FullName), "rotate"}, "")
	os.Mkdir(dst, os.ModePerm)
	fname := strings.Join([]string{strings.Trim(in.FullName, in.ExtName), "mp4"}, ".")
	export := strings.Join([]string{dst, fname}, string(os.PathSeparator))
	var cmd *exec.Cmd
	var transport string
	switch direction {
	case "ToRight":
		transport = "transpose=1"
	case "ToLeft":
		transport = "transpose=2"
	default:
		return
	}
	//if info.Width > 1920 || info.Height > 1920 {
	//	cmd = exec.Command("ffmpeg", "-threads", threads, "-i", in.FullPath, "-strict", "2", "-vf", "scale=-1:1080", "-vf", transport, "-c:v", "libx265", "-threads", threads, export)
	//}
	cmd = exec.Command("ffmpeg", "-threads", threads, "-i", in.FullPath, "-vf", transport, "-threads", threads, export)
	log.Debug.Printf("开始处理文件%s\n生成的命令是:%s\n", in.FullPath, cmd)
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	if err != nil {
		log.Warn.Panicf("cmd.StdoutPipe产生的错误:%v", err)
	}
	if err = cmd.Start(); err != nil {
		log.Warn.Panicf("cmd.Run产生的错误:%v", err)
	}
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		//写成输出日志
		t := string(tmp)
		t = replace.Replace(t)
		log.TTY.Println(t)
		if err != nil {
			break
		}
	}
	if err = cmd.Wait(); err != nil {
		log.Warn.Panicf("命令执行中有错误产生:%v", err)
	}
	err = os.RemoveAll(in.FullPath)
	if err != nil {
		log.Warn.Panicf("删除文件%v出现错误%v\n", in.FullPath, err)
	} else {
		log.Debug.Printf("完成当前文件的处理:%s\n", dst)
	}
}
