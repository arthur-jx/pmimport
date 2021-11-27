package media

import (
	"os/exec"
	"time"
)

type VideoMetadata struct {
	CreationTime time.Time
	Model        string
	Location     string
}

//调用命令 ffprobe -show_format FILE 获取视频信息
func GetVideoInfo(path string) (metadata *VideoMetadata, err error) {
	cmd := exec.Command("ffprobe", "-show_format", path)

	stdout, err := cmd.Output()

	if err != nil {
		println("exec command error:", err.Error())
		return
	}
	println("exec command output:", string(stdout))

	return
}
