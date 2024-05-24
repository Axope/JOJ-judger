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
		"lang":   req.Lang,
	}
	data["testCases"] = req.TestCases
	if req.Lang == pb.LangSet_PYTHON {
		data["solution"] = "solution.py"
	}
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
	if err := utils.DeleteFilesByPrefix(ContainerDir, "solution"); err != nil {
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

func JudgeFromProtobuf(req *pb.Judge) (pb.StatusSet, int64, int64, error) {
	defer log.Logger.Sync()
	log.Logger.Info("judging req", log.Any("req", req))

	if err := downloadData(req.Pid); err != nil {
		return pb.StatusSet_UKE, -1, -1, err
	}
	if err := initContainer(req); err != nil {
		return pb.StatusSet_UKE, -1, -1, err
	}

	switch req.Lang {
	case pb.LangSet_CPP:
		log.Logger.Debug("cpp submission")
		result, executeTime, executeMemory, err := judgeSolutionByCppFromProtobuf(req)
		if err != nil {
			log.Logger.Debug("judge error", log.Any("err", err))
			return pb.StatusSet_UKE, -1, -1, err
		}
		log.Logger.Info("judge done",
			log.Any("req", req), log.Any("result", result))
		return result, executeTime, executeMemory, nil
	// TODO: other language
	// case model.JAVA:
	case pb.LangSet_PYTHON:
		log.Logger.Debug("python submission")
		result, executeTime, executeMemory, err := judgeSolutionByPythonFromProtobuf(req)
		if err != nil {
			log.Logger.Debug("judge error", log.Any("err", err))
			return pb.StatusSet_UKE, -1, -1, err
		}
		log.Logger.Info("judge done",
			log.Any("req", req), log.Any("result", result))
		return result, executeTime, executeMemory, nil
	// case model.GO:
	default:
		log.Logger.Warn("unsupported language")
		return pb.StatusSet_UKE, -1, -1, errors.New("unsupported language")
	}
}
