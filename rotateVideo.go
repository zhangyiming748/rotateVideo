package rotateVideo

import (
	"fmt"
	"github.com/zhangyiming748/GetFileInfo"
	"github.com/zhangyiming748/voiceAlert"
	"golang.org/x/exp/slog"

	"os"
	"os/exec"
	"strings"
)

func Rotate(src, pattern, direction, threads string) {
	files := GetFileInfo.GetAllFileInfo(src, pattern)
	for index, file := range files {
		slog.Info(fmt.Sprintf("开始处理第 %d/%d 个文件", index+1, len(files)))
		rotate(file, direction, threads)
		slog.Info(fmt.Sprintf("处理完成第 %d/%d 个文件", index+1, len(files)))
		voiceAlert.Customize("done", voiceAlert.Ava)
	}
	voiceAlert.Customize("complete", voiceAlert.Ava)
}
func rotate(in GetFileInfo.Info, direction, threads string) {
	defer func() {
		if err := recover(); err != nil {
			voiceAlert.Customize("failed", voiceAlert.Ava)
		}
	}()
	dst := strings.Join([]string{strings.Trim(in.FullPath, in.FullName), "rotate"}, "")
	os.Mkdir(dst, os.ModePerm)
	fname := strings.Join([]string{strings.Trim(in.FullName, in.ExtName), "mp4"}, "")
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
	slog.Info("开始处理文件", slog.Any("生成的命令", fmt.Sprint(cmd)))
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	if err != nil {
		slog.Warn("cmd.StdoutPipe", slog.Any("错误", err))
		return
	}
	if err = cmd.Start(); err != nil {
		slog.Warn("cmd.Run", slog.Any("错误", err))
		return
	}
	for {
		tmp := make([]byte, 1024)
		_, err = stdout.Read(tmp)
		//写成输出日志
		t := string(tmp)
		t = strings.Replace(t, "\u0000", "", -1)
		fmt.Println(t)
		if err != nil {
			break
		}
	}
	if err = cmd.Wait(); err != nil {
		slog.Warn("cmd.Wait", slog.Any("错误", err))
		return
	}
	err = os.RemoveAll(in.FullPath)
	if err != nil {
		slog.Warn("删除失败", slog.Any("源文件", in.FullPath), slog.Any("错误", err))
	} else {
		slog.Info("删除成功", slog.Any("源文件", in.FullPath))
	}
}
