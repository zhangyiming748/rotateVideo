package rotateVideo

import (
	"fmt"
	"github.com/zhangyiming748/GetFileInfo"
	"github.com/zhangyiming748/replace"
	"github.com/zhangyiming748/voiceAlert"
	"golang.org/x/exp/slog"
	"io"
	"os"
	"os/exec"
	"strings"
)

var mylog *slog.Logger

func SetLog(level string) {
	var opt slog.HandlerOptions
	switch level {
	case "Debug":
		opt = slog.HandlerOptions{ // 自定义option
			AddSource: true,
			Level:     slog.LevelDebug, // slog 默认日志级别是 info
		}
	case "Info":
		opt = slog.HandlerOptions{ // 自定义option
			AddSource: true,
			Level:     slog.LevelInfo, // slog 默认日志级别是 info
		}
	case "Warn":
		opt = slog.HandlerOptions{ // 自定义option
			AddSource: true,
			Level:     slog.LevelWarn, // slog 默认日志级别是 info
		}
	case "Err":
		opt = slog.HandlerOptions{ // 自定义option
			AddSource: true,
			Level:     slog.LevelError, // slog 默认日志级别是 info
		}
	default:
		slog.Warn("需要正确设置环境变量 Debug,Info,Warn or Err")
		slog.Info("默认使用Debug等级")
		opt = slog.HandlerOptions{ // 自定义option
			AddSource: true,
			Level:     slog.LevelDebug, // slog 默认日志级别是 info
		}
	}
	file := "rotateVideo.log"
	logf, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}
	//defer logf.Close() //如果不关闭可能造成内存泄露
	mylog = slog.New(opt.NewJSONHandler(io.MultiWriter(logf, os.Stdout)))
}
func init() {
	l := os.Getenv("LEVEL")
	SetLog(l)
}
func Rotate(src, pattern, direction, threads string) {
	files := GetFileInfo.GetAllFileInfo(src, pattern)
	for index, file := range files {
		mylog.Info("正在处理第 %d/%d 个文件\n", index+1, len(files))
		rotate(file, direction, threads)
		mylog.Info("处理完成第 %d/%d 个文件\n", index+1, len(files))
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
	mylog.Info("开始处理文件", slog.Any("生成的命令", fmt.Sprint(cmd)))
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	if err != nil {
		mylog.Warn("cmd.StdoutPipe", slog.Any("错误", err))
		return
	}
	if err = cmd.Start(); err != nil {
		mylog.Warn("cmd.Run", slog.Any("错误", err))
		return
	}
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		//写成输出日志
		t := string(tmp)
		t = replace.Replace(t)
		fmt.Println(t)
		if err != nil {
			break
		}
	}
	if err = cmd.Wait(); err != nil {
		mylog.Warn("cmd.Wait", slog.Any("错误", err))
		return
	}
	err = os.RemoveAll(in.FullPath)
	if err != nil {
		mylog.Warn("删除失败", slog.Any("源文件", in.FullPath), slog.Any("错误", err))
	} else {
		mylog.Info("删除成功", slog.Any("源文件", in.FullPath))
	}
}
