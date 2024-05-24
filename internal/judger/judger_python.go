package judger

import (
	"github.com/Axope/JOJ-Judger/common/log"
	pb "github.com/Axope/JOJ-Judger/protocol"
	"github.com/Axope/JOJ-Judger/utils"
)

const (
	dockerImagePython   = "axope/sandbox:judger-python-0.1"
	containerNamePython = "JOJ-python"
)

func checkRequiredFilesPython() bool {
	return utils.CheckFileExist("./JOJ-sandbox/container/config/run.json") &&
		utils.CheckFileExist("./JOJ-sandbox/container/script/run.sh") &&
		utils.CheckFileExist("./JOJ-sandbox/container/sandbox") &&
		utils.CheckFileExist("./JOJ-sandbox/container/solution.py")
}

func judgeSolutionByPythonFromProtobuf(req *pb.Judge) (pb.StatusSet, int64, int64, error) {
	// 额外的写文件
	if err := utils.WriteToFile("./JOJ-sandbox/container/solution.py", []byte(req.SubmitCode)); err != nil {
		return pb.StatusSet_UKE, -1, -1, err
	}

	if !checkRequiredFilesPython() {
		return pb.StatusSet_UKE, -1, -1, nil
	}

	containerID, err := createAndRunContainer(dockerImagePython, containerNamePython)
	if err != nil {
		log.Logger.Debug("run container error", log.Any("err", err))
		return pb.StatusSet_UKE, -1, -1, err
	}
	defer removeContainer(containerID)
	log.Logger.Debug("create container success")

	// python是解释性语言
	cmd := []string{"sh", "-c", "cd /root && ./sandbox -type=run"}
	exitCode, err := execInContainer(containerID, cmd)
	if err != nil {
		return pb.StatusSet_UKE, -1, -1, err
	}
	switch exitCode {
	case 0:
		executeTime, err := utils.GetNumber("./JOJ-sandbox/container/output/executeTime")
		if err != nil {
			return pb.StatusSet_UKE, -1, -1, nil
		}
		executeMemory, err := utils.GetNumber("/sys/fs/cgroup/memory/JOJ-judger/memory.usage_in_bytes")
		if err != nil {
			return pb.StatusSet_UKE, -1, -1, nil
		}
		return pb.StatusSet_AC, executeTime, executeMemory, nil
	case 2:
		return pb.StatusSet_WA, -1, -1, nil
	case 3:
		return pb.StatusSet_RE, -1, -1, nil
	case 4:
		return pb.StatusSet_TLE, -1, -1, nil
	default:
		log.Logger.Sugar().Debugf("Unknown error, exit code = %v", exitCode)
		return pb.StatusSet_UKE, -1, -1, nil
	}

}
