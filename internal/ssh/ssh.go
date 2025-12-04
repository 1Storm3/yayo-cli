package ssh

import (
	"bytes"
	"io"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

func RunSSHWithStdin(host, command, stdinData string) (string, error) {
	key, err := os.ReadFile(os.Getenv("HOME") + "/.ssh/id_rsa")
	if err != nil {
		return "", err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return "", err
	}

	config := &ssh.ClientConfig{
		User: host[:strings.Index(host, "@")],
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp",
		host[strings.Index(host, "@")+1:]+":22",
		config,
	)
	if err != nil {
		return "", err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	var stdout bytes.Buffer
	session.Stdout = &stdout

	stdin, err := session.StdinPipe()
	if err != nil {
		return "", err
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, stdinData)
	}()

	if err := session.Run(command); err != nil {
		return "", err
	}

	return stdout.String(), nil
}
