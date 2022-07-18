package ffmpeg

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

func Probe(filename string, options ...string) (*ProbeData, error) {
	return ProbeWithTimeout(filename, 0, options...)
}

func ProbeWithTimeout(filename string, timeout time.Duration, options ...string) (*ProbeData, error) {
	options = append([]string{
		"-loglevel", "fatal",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
	}, options...)

	return ProbeWithTimeoutExec(filename, timeout, options...)
}

func ProbeWithTimeoutExec(filename string, timeout time.Duration, options ...string) (*ProbeData, error) {
	options = append(options, filename)

	ctx := context.Background()
	if timeout > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
		defer cancel()
	}

	var outputBuf bytes.Buffer
	var stdErr bytes.Buffer

	cmd := exec.CommandContext(ctx, ffprobe, options...)
	cmd.Stdout = &outputBuf
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s\n %s", err.Error(), stdErr.String()))
	}

	data := &ProbeData{}
	err = json.Unmarshal(outputBuf.Bytes(), data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
