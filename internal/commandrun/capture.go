package commandrun

import (
	"bytes"
	"io"
	"os"
)

func CaptureStdout(run func() error) (string, error) {
	originalStdout := os.Stdout

	reader, writer, err := os.Pipe()
	if err != nil {
		return "", err
	}

	os.Stdout = writer

	var buffer bytes.Buffer
	done := make(chan error, 1)
	go func() {
		_, copyErr := io.Copy(&buffer, reader)
		done <- copyErr
	}()

	runErr := run()

	os.Stdout = originalStdout
	closeErr := writer.Close()
	copyErr := <-done
	readCloseErr := reader.Close()

	if runErr != nil {
		return buffer.String(), runErr
	}
	if closeErr != nil {
		return buffer.String(), closeErr
	}
	if copyErr != nil {
		return buffer.String(), copyErr
	}
	if readCloseErr != nil {
		return buffer.String(), readCloseErr
	}

	return buffer.String(), nil
}
