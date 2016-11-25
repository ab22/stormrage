package ws

import "os/exec"

func newPingWriter(client WebsocketClient, ip string) *pingWriter {
	var (
		cmd    = exec.Command("ping", "-t", ip)
		writer = &pingWriter{
			client: client,
			cmd:    cmd,
		}
	)

	cmd.Stdout = writer
	cmd.Stderr = writer

	return writer
}
