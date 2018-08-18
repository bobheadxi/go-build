package build

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	docker "github.com/docker/docker/client"
)

// NewDockerClient creates a new Docker Client from ENV values and negotiates
// the correct API version
func NewDockerClient() (*docker.Client, error) {
	c, err := docker.NewEnvClient()
	if err != nil {
		return nil, err
	}
	c.NegotiateAPIVersion(context.Background())
	return c, nil
}

// LogOptions is used to configure retrieved container logs
type LogOptions struct {
	Container    string
	Stream       bool
	Detailed     bool
	NoTimestamps bool
	Entries      int
}

// containerLogs get logs ;)
func containerLogs(docker *docker.Client, opts LogOptions) (io.ReadCloser, error) {
	ctx := context.Background()
	return docker.ContainerLogs(ctx, opts.Container, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     opts.Stream,
		Timestamps: !opts.NoTimestamps,
		Details:    opts.Detailed,
		Tail:       strconv.Itoa(opts.Entries),
	})
}

// wait blocks until given container ID stops
func wait(cli *docker.Client, id string, stop chan struct{}) (int64, error) {
	var status container.ContainerWaitOKBody
	statusCh, errCh := cli.ContainerWait(context.Background(), id, "")
	select {
	case err := <-errCh:
		if err != nil {
			return 0, err
		}
	case status = <-statusCh:
		// Exit log stream
		close(stop)
	}
	return status.StatusCode, nil
}

// startAndWait starts and waits for container to exit
func startAndWait(cli *docker.Client, containerID string, out io.Writer) error {
	ctx := context.Background()
	if err := cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	stop := make(chan struct{})
	go streamContainerLogs(cli, containerID, out, stop)
	exitCode, err := wait(cli, containerID, stop)
	if err != nil {
		return err
	}
	if exitCode != 0 {
		return fmt.Errorf("Container exited with non-zero status %d", exitCode)
	}

	return nil
}

// streamContainerLogs streams logs from given container ID. Best used as a
// goroutine.
func streamContainerLogs(client *docker.Client, id string, out io.Writer,
	stop chan struct{}) error {
	// Attach logs and report build progress until container exits
	reader, err := containerLogs(client, LogOptions{
		Container: id, Stream: true,
		NoTimestamps: true,
	})
	if err != nil {
		return err
	}
	defer reader.Close()
	flushRoutine(out, reader, stop)
	return nil
}
