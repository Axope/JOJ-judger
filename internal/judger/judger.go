package judger

import (
	"encoding/json"
	"errors"
	"path/filepath"

	"github.com/Axope/JOJ-Judger/common/log"
	pb "github.com/Axope/JOJ-Judger/protocol"
	"github.com/Axope/JOJ-Judger/utils"
)

const (
	DatasDir = "./datas"

	RunConfigPath = "./JOJ-sandbox/container/config/run.json"
	DataDir       = "./JOJ-sandbox/container/data"
	OutputDir     = "./JOJ-sandbox/container/output"
	ContainerDir  = "./JOJ-sandbox/container"
)

func writeRunJsonFromProtobuf(req *pb.Judge) error {
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

	if err = utils.WriteToFile(RunConfigPath, jsonData); err != nil {
		return err
	}
	return nil
}

func checkRequiredFiles() bool {
	return utils.CheckFileExist("./JOJ-sandbox/container/config/compile.json") &&
		utils.CheckFileExist("./JOJ-sandbox/container/config/run.json") &&
		utils.CheckFileExist("./JOJ-sandbox/container/script/compile.sh") &&
		utils.CheckFileExist("./JOJ-sandbox/container/script/run.sh") &&
		utils.CheckFileExist("./JOJ-sandbox/container/sandbox") &&
		utils.CheckFileExist("./JOJ-sandbox/container/solution.cpp")
}

func judgeSolutionByCPPFromProtobuf(req *pb.Judge) (pb.StatusSet, error) {
	// 额外的写文件
	if err := utils.WriteToFile("./JOJ-sandbox/container/solution.cpp", []byte(req.SubmitCode)); err != nil {
		return pb.StatusSet_UKE, err
	}

	if !checkRequiredFiles() {
		return pb.StatusSet_UKE, nil
	}
	containerID, err := createAndRunContainer()
	if err != nil {
		log.Logger.Debug("run container error", log.Any("err", err))
		return pb.StatusSet_UKE, err
	}
	defer removeContainer(containerID)
	log.Logger.Debug("create container success")

	cmd := []string{"sh", "-c", "cd /root && ./sandbox -type=compile"}
	exitCode, err := execInContainer(containerID, cmd)
	if err != nil {
		return pb.StatusSet_UKE, err
	}
	if exitCode != 0 || !utils.CheckFileExist("./JOJ-sandbox/container/solution") {
		return pb.StatusSet_CE, nil
	}

	cmd = []string{"sh", "-c", "cd /root && ./sandbox -type=run"}
	exitCode, err = execInContainer(containerID, cmd)
	if err != nil {
		return pb.StatusSet_UKE, err
	}
	switch exitCode {
	case 0:
		return pb.StatusSet_AC, nil
	case 2:
		return pb.StatusSet_WA, nil
	case 3:
		return pb.StatusSet_RE, nil
	case 4:
		return pb.StatusSet_TLE, nil
	default:
		log.Logger.Sugar().Debugf("Unknown error, exit code = %v", exitCode)
		return pb.StatusSet_UKE, nil
	}

}

func initContainer(req *pb.Judge) error {
	// config
	log.Logger.Info("judging req", log.Any("req", req))
	if err := writeRunJsonFromProtobuf(req); err != nil {
		return err
	}
	log.Logger.Debug("write run.json done")
	// data
	if err := utils.DeleteFiles(DataDir); err != nil {
		return err
	}
	dataDirPath := filepath.Join(DatasDir, req.Pid)
	if err := utils.CopyDirFile(dataDirPath, DataDir); err != nil {
		return err
	}
	// output
	if err := utils.DeleteFiles(OutputDir); err != nil {
		return err
	}
	// script (nothing to do)

	// solution
	if err :=utils.DeleteFilesByPrefix(ContainerDir, "solution");err != nil {
		return err
	}

	return nil
}

func downloadData(pid string) error {
	path := filepath.Join(DatasDir, pid)
	if utils.CheckFileExist(path) {
		return nil
	}
	err := utils.DownloadTestCasesByRsync(pid, path)
	return err
}

func JudgeFromProtobuf(req *pb.Judge) (pb.StatusSet, error) {
	defer log.Logger.Sync()
	log.Logger.Info("judging req", log.Any("req", req))

	if err := downloadData(req.Pid); err != nil {
		return pb.StatusSet_UKE, err
	}
	if err := initContainer(req); err != nil {
		return pb.StatusSet_UKE, err
	}

	switch req.Lang {
	case pb.LangSet_CPP:
		log.Logger.Debug("cpp submission")
		result, err := judgeSolutionByCPPFromProtobuf(req)
		if err != nil {
			log.Logger.Debug("judge error", log.Any("err", err))
			return pb.StatusSet_UKE, err
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
		return pb.StatusSet_UKE, errors.New("unsupported language")
	}
}
