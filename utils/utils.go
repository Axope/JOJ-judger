package utils

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func WriteToFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

// rsync -a rsync://172.28.222.195/datas/$pid/ ./JOJ-sandbox/container/data
func DownloadTestCasesByRsync(pid string, dstPath string) error {
	source := fmt.Sprintf("rsync://172.28.222.195/datas/%s/", pid)
	rsyncCmd := exec.Command("rsync", "-a", source, dstPath)

	if err := rsyncCmd.Run(); err != nil {
		return err
	}
	return nil
}

func DeleteFiles(dirPath string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		filePath := filepath.Join(dirPath, file.Name())
		if err := os.Remove(filePath); err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		}
	}

	return nil
}

func CheckFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	} else if err != nil {
		return false
	} else {
		return true
	}
}

func DeleteFilesByPrefix(dirPath, prefix string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), prefix) {
			filePath := filepath.Join(dirPath, file.Name())
			if err := os.Remove(filePath); err != nil {
				return err
			}
		}
	}

	return nil
}

func CopyDirFile(srcDirPath, dstDirPath string) error {
	files, err := os.ReadDir(srcDirPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		sourcePath := filepath.Join(srcDirPath, file.Name())
		destinationPath := filepath.Join(dstDirPath, file.Name())

		// 打开源文件
		sourceFile, err := os.Open(sourcePath)
		if err != nil {
			return err
		}
		defer sourceFile.Close()

		// 创建目标文件
		destinationFile, err := os.Create(destinationPath)
		if err != nil {
			return err
		}
		defer destinationFile.Close()

		// 将源文件内容复制到目标文件
		_, err = io.Copy(destinationFile, sourceFile)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetNumber(path string) (int64, error) {
	file, err := os.Open(path)
	if err != nil {
		return -1, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return -1, err
	}
	strContent := strings.TrimSpace(string(content))
	value, err := strconv.ParseInt(strContent, 10, 64)
	if err != nil {
		return -1, err
	}
	return value, nil
}
