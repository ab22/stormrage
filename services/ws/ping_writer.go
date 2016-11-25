package ws

import (
	"fmt"
	"log"
	"os/exec"
)

// pingWriter is an implementation of io.Writer which makes it possible to intercept
// the data send by the os.Command execution output. Since our goal is to send the
// ping result to the client, then we intercept this data and send it directly to the
// websocket writer.
//
// Additionally, the client can set a flag to stop the ping and kill the ping process.
type pingWriter struct {
	client WebsocketClient
	cmd    *exec.Cmd
}

func (w *pingWriter) Write(p []byte) (n int, err error) {
	var (
		msg     = fmt.Sprintf("{\"payload\": \"%v\"}", string(p))
		success = w.client.WriteAndWait([]byte(msg))
	)

	if !success {
		return 0, fmt.Errorf("ping writer: write timeout!")
	}

	return len(p), nil
}

func (w *pingWriter) StartAndWait() {
	if err := w.cmd.Start(); err != nil {
		err = fmt.Errorf("error starting command: %v", err)
		w.client.LogError(err)
		return
	}

	if err := w.cmd.Wait(); err != nil {
		err = fmt.Errorf("error waiting for command: %v", err)
		w.client.LogError(err)
		return
	}

	log.Println("Ping command ended successfully!")
}

func (w *pingWriter) Kill() error {
	if w.cmd == nil {
		return nil
	}

	return w.cmd.Process.Kill()
}
