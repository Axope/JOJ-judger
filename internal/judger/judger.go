package judger

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/Axope/JOJ-Judger/common/log"
	"github.com/Axope/JOJ-Judger/common/request"
	"github.com/Axope/JOJ-Judger/internal/model"
	pb "github.com/Axope/JOJ-Judger/protocol"
	"github.com/Axope/JOJ-Judger/utils"
)

const (
	configPath   = "./JOJ-sandbox/container/config/run.json"
	dataDir      = "./JOJ-sandbox/container/data"
	containerDir = "./JOJ-sandbox/container"
)

func writeRunJson(req request.JudgeRequest) error {
	if err := utils.DownloadTestCasesByRsync(req.PID); err != nil {
		return err
	}
	// make run.json
	data := map[string]interface{}{
		"memLimit":  req.MemoryLimit,
		"timeLimit": req.TimeLimit,
		"solution":  "solution",
	}
	data["testCases"] = req.TestCases
	log.Logger.Debug("write run.json", log.Any("data", data))

	// marshal
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err = utils.WriteToFile(configPath, jsonData); err != nil {
		return err
	}
	return nil
}

func writeRunJsonFromProtobuf(req *pb.Judge) error {
	if err := utils.DownloadTestCasesByRsync(req.Pid); err != nil {
		return err
	}
	// make run.json
	data := map[string]interface{}{
		"memLimit":  req.MemoryLimit,
		"timeLimit": req.TimeLimit,
		"solution":  "solution",
	}
	data["testCases"] = req.TestCases
	log.Logger.Debug("write run.json", log.Any("data", data))

	// marshal
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err = utils.WriteToFile(configPath, jsonData); err != nil {
		return err
	}
	return nil
}

func deleteFiles(dirPath string) error {
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

func clean() error {
	dir := "./JOJ-sandbox/container/output"
	if err := deleteFiles(dir); err != nil {
		return err
	}
	if err := os.Remove("./JOJ-sandbox/container/solution"); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func checkFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	} else if err != nil {
		return false
	} else {
		return true
	}
}

func checkFiles() bool {
	return checkFileExist("./JOJ-sandbox/container/config/compile.json") &&
		checkFileExist("./JOJ-sandbox/container/config/run.json") &&
		checkFileExist("./JOJ-sandbox/container/script/compile.sh") &&
		checkFileExist("./JOJ-sandbox/container/script/run.sh") &&
		checkFileExist("./JOJ-sandbox/container/sandbox") &&
		checkFileExist("./JOJ-sandbox/container/solution.cpp")
}

func judgeSolutionByCPP(req request.JudgeRequest) (model.StatusSet, error) {
	if err := utils.WriteToFile("./JOJ-sandbox/container/solution.cpp", []byte(req.SubmitCode)); err != nil {
		return model.UKE, err
	}
	if err := clean(); err != nil {
		return model.UKE, err
	}
	log.Logger.Debug("clean done")
	if !checkFiles() {
		return model.UKE, nil
	}
	containerID, err := createAndRunContainer()
	if err != nil {
		log.Logger.Debug("run container error", log.Any("err", err))
		return model.UKE, err
	}
	defer removeContainer(containerID)
	log.Logger.Debug("create container success")

	cmd := []string{"sh", "-c", "cd /root && ./sandbox -type=compile"}
	exitCode, err := execInContainer(containerID, cmd)
	if err != nil {
		return model.UKE, err
	}
	if exitCode != 0 || !checkFileExist("./JOJ-sandbox/container/solution") {
		return model.CE, nil
	}

	cmd = []string{"sh", "-c", "cd /root && ./sandbox -type=run"}
	exitCode, err = execInContainer(containerID, cmd)
	if err != nil {
		return model.UKE, err
	}
	switch exitCode {
	case 0:
		return model.AC, nil
	case 2:
		return model.WA, nil
	case 3:
		return model.RE, nil
	case 4:
		return model.TLE, nil
	default:
		log.Logger.Sugar().Debugf("Unknown error, exit code = %v", exitCode)
		return model.UKE, nil
	}

}

func judgeSolutionByCPPFromProtobuf(req *pb.Judge) (model.StatusSet, error) {
	if err := utils.WriteToFile("./JOJ-sandbox/container/solution.cpp", []byte(req.SubmitCode)); err != nil {
		return model.UKE, err
	}
	if err := clean(); err != nil {
		return model.UKE, err
	}
	log.Logger.Debug("clean done")
	if !checkFiles() {
		return model.UKE, nil
	}
	containerID, err := createAndRunContainer()
	if err != nil {
		log.Logger.Debug("run container error", log.Any("err", err))
		return model.UKE, err
	}
	defer removeContainer(containerID)
	log.Logger.Debug("create container success")

	cmd := []string{"sh", "-c", "cd /root && ./sandbox -type=compile"}
	exitCode, err := execInContainer(containerID, cmd)
	if err != nil {
		return model.UKE, err
	}
	if exitCode != 0 || !checkFileExist("./JOJ-sandbox/container/solution") {
		return model.CE, nil
	}

	cmd = []string{"sh", "-c", "cd /root && ./sandbox -type=run"}
	exitCode, err = execInContainer(containerID, cmd)
	if err != nil {
		return model.UKE, err
	}
	switch exitCode {
	case 0:
		return model.AC, nil
	case 2:
		return model.WA, nil
	case 3:
		return model.RE, nil
	case 4:
		return model.TLE, nil
	default:
		log.Logger.Sugar().Debugf("Unknown error, exit code = %v", exitCode)
		return model.UKE, nil
	}

}

func Judge(req request.JudgeRequest) (model.StatusSet, error) {
	log.Logger.Info("judging req", log.Any("req", req))
	if err := writeRunJson(req); err != nil {
		return model.UKE, err
	}
	log.Logger.Debug("write run.json done")

	switch req.Lang {
	case model.CPP:
		log.Logger.Debug("cpp submission")
		result, err := judgeSolutionByCPP(req)
		if err != nil {
			log.Logger.Debug("judge error", log.Any("err", err))
			return model.UKE, err
		}
		log.Logger.Info("judge done",
			log.Any("req", req), log.Any("result", result))
		return result, nil
	// TODO: other language
	// case model.JAVA:
	// case model.PYTHON:
	// case model.GO:
	default:
		log.Logger.Warn("unsupported language")
		return model.UKE, errors.New("unsupported language")
	}

}

func JudgeFromProtobuf(req *pb.Judge) (model.StatusSet, error) {
	log.Logger.Info("judging req", log.Any("req", req))
	if err := writeRunJsonFromProtobuf(req); err != nil {
		return model.UKE, err
	}
	log.Logger.Debug("write run.json done")

	switch req.Lang {
	case pb.LangSet_CPP:
		log.Logger.Debug("cpp submission")
		result, err := judgeSolutionByCPPFromProtobuf(req)
		if err != nil {
			log.Logger.Debug("judge error", log.Any("err", err))
			return model.UKE, err
		}
		log.Logger.Info("judge done",
			log.Any("req", req), log.Any("result", result))
		return result, nil
	// TODO: other language
	// case model.JAVA:
	// case model.PYTHON:
	// case model.GO:
	default:
		log.Logger.Warn("unsupported language")
		return model.UKE, errors.New("unsupported language")
	}
}
