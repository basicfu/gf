package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"os"
	"time"
)

type Cli struct {
	IP         string      //IP地址
	Username   string      //用户名
	Password   string      //密码
	Port       int         //端口号
	client     *ssh.Client //ssh客户端
}

func Ssh(ip string, username string, password string, port ...int) *Cli {
	cli := new(Cli)
	cli.IP = ip
	cli.Username = username
	cli.Password = password
	if len(port) <= 0 {
		cli.Port = 22
	} else {
		cli.Port = port[0]
	}
	return cli
}

func (c Cli) Run(shell string){
	if c.client == nil {
		if err := c.connect(); err != nil {
			panic(err)
		}
	}
	session,_:=c.client.NewSession()
	stdout, _ := session.StdoutPipe()
	stderr, _ := session.StderrPipe()
	_ = session.Start(shell)
	//_ = session.Run(shell)
	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)
	_ = session.Wait()
}

func (c *Cli) connect() error {
	config := ssh.ClientConfig{
		User: c.Username,
		Auth: []ssh.AuthMethod{ssh.Password(c.Password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 10 * time.Second,
	}
	addr := fmt.Sprintf("%s:%d", c.IP, c.Port)
	sshClient, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		return err
	}
	c.client = sshClient
	return nil
}
