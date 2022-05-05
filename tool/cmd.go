package tool

import (
	"bufio"
	"io"
	"os/exec"
	"strings"
	"tiny_tool/log"
)

func Adb(raw ...string) int {
	cmd, _ := BaseCmd("adb", false, raw...)
	return cmd
}

func ExecCmd(shell string, raw ...string) int {
	cmd, _ := BaseCmd(shell, false, raw...)
	return cmd
}

func BaseCmd(shell string, mute bool, raw ...string) (int, []string) {
	cmd := exec.Command(shell, raw...)
	defer func() {
		if cmd != nil && cmd.Process != nil {
			cmd.Process.Release()
			cmd.Process.Kill()
		}
	}()
	result := make([]string, 0)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.E(err.Error())
		return -1000, result
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.E(err.Error())
		return -1000, result
	}
	if err := cmd.Start(); err != nil {
		log.E(err.Error())
		return -1000, result
	}
	s := bufio.NewScanner(io.MultiReader(stdout, stderr))

	for s.Scan() {
		text := s.Text()
		result = append(result, text)
		if !mute {
			log.V(text)
		}
	}

	if err := cmd.Wait(); err != nil {
		log.E(err.Error())
	}
	return cmd.ProcessState.ExitCode(), result
}

func CmdWatch(init func(), shell string, raw ...string) *exec.Cmd {
	cmd := exec.Command(shell, raw...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.E(err.Error())
		return nil
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.E(err.Error())
		return nil
	}
	go func(cmd *exec.Cmd) {
		defer func() {
			cmd.Process.Kill()
		}()
		if err := cmd.Start(); err != nil {
			log.E(err.Error())
			return
		}
		s := bufio.NewScanner(io.MultiReader(stdout, stderr))
		for s.Scan() {
			text := s.Text()
			log.V(text)
			if strings.Contains(text,"webpack")&& strings.Contains(text,"compiled") && init != nil {
				init()
			}
		}
		if err := cmd.Wait(); err != nil {
			log.E(err.Error())
		}
	}(cmd)
	return cmd
}
