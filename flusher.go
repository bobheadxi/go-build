package build

import (
	"bufio"
	"io"
	"net/http"
)

// flushRoutine continuously writes everything in given ReadCloser
// to an io.Writer. Use this as a goroutine.
func flushRoutine(w io.Writer, rc io.Reader, stop chan struct{}) {
	reader := bufio.NewReader(rc)
ROUTINE:
	for {
		select {
		case <-stop:
			writeAndFlush(w, reader)
			break ROUTINE
		default:
			// Read from pipe then write to ResponseWriter and flush it,
			// sending the copied content to the client.
			err := writeAndFlush(w, reader)
			if err != nil {
				break ROUTINE
			}
		}
	}
}

// writeAndFlush reads from buffer, writes to writer, and flushes if possible
func writeAndFlush(w io.Writer, reader *bufio.Reader) error {
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return err
	}

	// Write to writer, and flush as well if it is a flusher
	w.Write(line)
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	return nil
}
