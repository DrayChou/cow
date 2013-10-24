package main

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"
)

func Ssh2Running(socksServer string) bool {
	c, err := net.Dial("tcp", socksServer)
	if err != nil {
		return false
	}
	c.Close()
	return true
}

func runOneSSH2(server string) {
	// config parsing canonicalize sshServer config value
	arr := strings.SplitN(server, ":", 4)
	sshServer, passwd, localPort, sshPort := arr[0], arr[1], arr[2], arr[3]
	alreadyRunPrinted := false

	socksServer := "127.0.0.1:" + localPort
	for {
		if Ssh2Running(socksServer) {
			if !alreadyRunPrinted {
				debug.Println("plink socks server", socksServer, "maybe already running")
				alreadyRunPrinted = true
			}
			time.Sleep(30 * time.Second)
			continue
		}

		// -n redirects stdin from /dev/null
		// -N do not execute remote command
		debug.Println("connecting to ssh server", sshServer+":"+sshPort)
		fmt.Println("connecting to ssh server", sshServer+":"+sshPort)
		cmd := exec.Command("plink", "-C", "-N", "-D", localPort, "-pw", passwd, "-P", sshPort, sshServer)
		fmt.Println("plink", "-C", "-N", "-D", localPort, "-pw", passwd, "-P", sshPort, sshServer)
		fmt.Println(cmd.Path)
		if err := cmd.Run(); err != nil {
			debug.Println("plink:", err)
		}
		debug.Println("plink", sshServer+":"+sshPort, "exited, reconnect")
		fmt.Println("plink", sshServer+":"+sshPort, "exited, reconnect")
		time.Sleep(5 * time.Second)
		alreadyRunPrinted = false
	}
}

func runSSH2() {
	for _, server := range config.SshServer2 {
		go runOneSSH2(server)
	}
}
