package services

import (
	"bytes"
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func CreateAndRunDockerContainer(command string) (string, error) {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.43"))
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "ubuntu",
		Cmd:   []string{"sh", "-c", command},
		Tty:   true,
	}, nil, nil, nil, "")
	if err != nil {
		return "", err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}

	exitCode, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case status := <-exitCode:
		if status.StatusCode != 0 {
			return "", fmt.Errorf("Container exited with non-zero status code: %d", status.StatusCode)
		}
	case err := <-errCh:
		return "", err
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return "", err
	}
	defer cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{})

	buf := new(bytes.Buffer)
	buf.ReadFrom(out)
	output := buf.String()

	if err := cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{}); err != nil {
		return "", err
	}

	return output, nil
}
