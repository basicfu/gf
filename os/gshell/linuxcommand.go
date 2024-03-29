package gshell

import (
	"github.com/basicfu/gf/text/gstr"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
)

// LinuxCommand结构体
type LinuxCommand struct {
}

// LinuxCommand的初始化函数
func NewLinuxCommand() *LinuxCommand {
	return &LinuxCommand{}
}
func shell() string {
	sh := os.Getenv("SHELL")
	if gstr.Trim(sh) == "" {
		sh = "/bin/sh"
	}
	return sh
}

// 执行命令行并返回结果
// args: 命令行参数
// return: 进程的pid, 命令行结果, 错误消息
func (lc *LinuxCommand) Exec(args ...string) (int, string, error) {
	args = append([]string{"-c"}, args...)
	cmd := exec.Command(shell(), args...)

	cmd.SysProcAttr = &syscall.SysProcAttr{}

	outpip, err := cmd.StdoutPipe()
	defer outpip.Close()

	if err != nil {
		return 0, "", err
	}

	err = cmd.Start()
	if err != nil {
		return 0, "", err
	}

	out, err := ioutil.ReadAll(outpip)
	if err != nil {
		return 0, "", err
	}

	return cmd.Process.Pid, string(out), nil
}
func (lc *LinuxCommand) ExecHidden(args ...string) (int, string, error) {
	args = append([]string{"-c"}, args...)
	cmd := exec.Command(shell(), args...)

	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	outpip, err := cmd.StdoutPipe()
	defer outpip.Close()

	if err != nil {
		return 0, "", err
	}

	err = cmd.Start()
	if err != nil {
		return 0, "", err
	}

	out, err := ioutil.ReadAll(outpip)
	if err != nil {
		return 0, "", err
	}

	return cmd.Process.Pid, string(out), nil
}

// 异步执行命令行并通过channel返回结果
// stdout: chan结果
// args: 命令行参数
// return: 进程的pid
// exception: 协程内的命令行发生错误时,会panic异常
func (lc *LinuxCommand) ExecAsync(stdout chan string, args ...string) int {
	var pidChan = make(chan int, 1)

	go func() {
		args = append([]string{"-c"}, args...)
		cmd := exec.Command(shell(), args...)

		cmd.SysProcAttr = &syscall.SysProcAttr{}

		outpip, err := cmd.StdoutPipe()
		defer outpip.Close()

		if err != nil {
			panic(err)
		}

		err = cmd.Start()
		if err != nil {
			panic(err)
		}

		pidChan <- cmd.Process.Pid

		out, err := ioutil.ReadAll(outpip)
		if err != nil {
			panic(err)
		}

		stdout <- string(out)
	}()

	return <-pidChan
}

// 执行命令行(忽略返回值)
// args: 命令行参数
// return: 错误消息
func (lc *LinuxCommand) ExecIgnoreResult(args ...string) error {

	args = append([]string{"-c"}, args...)
	cmd := exec.Command(shell(), args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{}

	err := cmd.Run()

	return err
}
