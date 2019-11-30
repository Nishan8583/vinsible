package SSH

import (
	"bufio"
	"bytes"
	"io"
	"log"

	"golang.org/x/crypto/ssh"
)

func CreateSSHConnection() {

	commands := []string{
		"sudo apt-get update",
		"ls",
		"sudo apt-get install radare2",
	}

	config := &ssh.ClientConfig{
		User: "nishan",
		Auth: []ssh.AuthMethod{
			ssh.Password("MeroSanoVai@7887"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Not a ssh key based authenticaiton for now
	}

	client, err := ssh.Dial("tcp", "192.168.100.143:22", config)
	if err != nil {
		log.Println("Error while creating a SSH conneciton with server -> ", err)
		return
	}

	session, err := client.NewSession()
	if err != nil {
		log.Println("Error could not create new session -> ", err)
		return
	}

	//session.Stdout = os.Stdout;
	//session.Stderr = os.Stderr;

	in, err := session.StdinPipe()
	out, err := session.StdoutPipe()

	ctx := make(chan int, 1)
	go func() {
		for {
			reader := bufio.NewReader(out)

			line, _, err := reader.ReadLine()

			if err != io.EOF {
				select {
				case <-ctx:
					break
				default:
					continue
				}
				log.Println("Error ", err)
				return
			}

			log.Println(string(line))
			if bytes.HasSuffix(line, []byte{58}) {

				_, err := in.Write([]byte("MeroSanoVai@7887" + "\n"))
				if err != nil {
					log.Println(err)
				}
			}
		}

	}()

	for _, cmd := range commands {
		err := session.Run(cmd)
		if err != nil {
			log.Println("Error while trying to run command -> ", err)
			return
		}
	}

	session.Close()
	ctx <- 1
}
