package judger

import (
	"context"
	"fmt"
	"os"

	"github.com/Axope/JOJ-Judger/common/log"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-units"
)

var (
	cli *client.Client
	ctx = context.Background()
)

const (
	dockerImage   = "axope/sandbox:judger-cpp-0.1"
	containerName = "JOJ-cpp"
)

func init() {
	var err error
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
}

func createAndRunContainer() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Logger.Error(err.Error())
	}

	// create container
	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:     dockerImage,
			OpenStdin: true,
			Tty:       true,
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: "/sys/fs/cgroup/memory/JOJ",
					Target: "/sys/fs/cgroup/memory",
				},
				{
					Type:   mount.TypeBind,
					Source: fmt.Sprintf("%s/JOJ-sandbox/container", currentDir),
					Target: "/root",
				},
			},
			NetworkMode: network.NetworkNone,
			Resources: container.Resources{
				Ulimits: []*units.Ulimit{ // 设置 ulimit
					&units.Ulimit{
						Name: "core",
						Soft: 0,
						Hard: 0,
					},
				},
			},
		},
		nil,
		nil,
		containerName)
	if err != nil {
		log.Logger.Error(err.Error())
	}
	containerID := resp.ID

	if err := cli.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		return "", err
	}

	// statusCh, errCh := cli.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)
	// select {
	// case err := <-errCh:
	// 	if err != nil {
	// 		return err
	// 	}
	// case <-statusCh:
	// }
	return containerID, nil
}

// func stopContainer(containerID string) {
// 	if containerID == "" {
// 		return
// 	}
// 	err := cli.ContainerStop(ctx, containerID, container.StopOptions{})
// 	if err != nil {
// 		log.Logger.Debug("stop container error", log.Any("err", err))
// 	}
// }

func execInContainer(containerID string, cmd []string) (int, error) {
	log.Logger.Debug("execInContainer", log.Any("cmd", cmd))
	// 容器执行命令的配置
	execConfig := types.ExecConfig{
		Cmd:          cmd,  // 要执行的命令及其参数
		AttachStdout: true, // 指定是否将标准输出附加到当前进程的stdout
		AttachStderr: true, // 指定是否将标准错误附加到当前进程的stderr
	}

	// 创建容器执行命令
	execID, err := cli.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return -1, err
	}

	// 执行容器命令并获取输出
	resp, err := cli.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		return -1, err
	}

	defer resp.Close()
	stdcopy.StdCopy(os.Stdout, os.Stderr, resp.Reader)
	inspect, err := cli.ContainerExecInspect(context.Background(), execID.ID)
	if err != nil {
		return -1, err
	}

	return inspect.ExitCode, nil
}

func removeContainer(containerID string) {
	if containerID == "" {
		return
	}
	err := cli.ContainerStop(ctx, containerID, container.StopOptions{})
	if err != nil {
		log.Logger.Debug("stop container error", log.Any("err", err))
	}
	err = cli.ContainerRemove(ctx, containerID, container.RemoveOptions{})
	if err != nil {
		log.Logger.Debug("remove container error", log.Any("err", err))
	}
}
