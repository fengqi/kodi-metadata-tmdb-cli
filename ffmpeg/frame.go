package ffmpeg

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

func Frame(fileName, outfile string, options ...string) error {
	return FrameWithTimeout(fileName, outfile, 0, options...)
}

func FrameWithTimeout(fileName, outfile string, timeout time.Duration, options ...string) error {
	options = append([]string{
		"-vframes", "1",
		"-format", "image2",
		"-vcodec", "mjpeg",
	}, options...)

	return FrameWithTimeoutExec(fileName, outfile, timeout, options...)
}

func FrameWithTimeoutExec(filename, outfile string, timeout time.Duration, options ...string) error {
	args := append([]string{
		"-i", filename,
	}, options...)
	args = append(args, "-y")
	args = append(args, outfile)

	ctx := context.Background()
	if timeout > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
		defer cancel()
	}

	var outputBuf bytes.Buffer
	var stdErr bytes.Buffer

	cmd := exec.CommandContext(ctx, ffmpeg, args...)
	cmd.Stdout = &outputBuf
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		return errors.New(fmt.Sprintf("%s\n %s", err.Error(), stdErr.String()))
	}

	return nil
}
