package utils

import (
	"fmt"
	"os"
	"os/exec"
)

func WriteToFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

// rsync -a rsync://172.28.222.195/datas/$pid/ ./JOJ-sandbox/container/data
func DownloadTestCasesByRsync(pid string) error {
	source := fmt.Sprintf("rsync://172.28.222.195/datas/%s/", pid)
	dst := "./JOJ-sandbox/container/data"
	rsyncCmd := exec.Command("rsync", "-a", source, dst)

	if err := rsyncCmd.Run(); err != nil {
		return err
	}
	return nil
}
