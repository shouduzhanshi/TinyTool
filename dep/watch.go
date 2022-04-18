package dep

import (
	"bufio"
	"io"
	"os/exec"
	"tiny_tool/log"
)

func watch()  {

}

func baseCmd(shell string, mute bool, raw ...string) *exec.Cmd {
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
	if err := cmd.Start(); err != nil {
		log.E(err.Error())
		return nil
	}
	s := bufio.NewScanner(io.MultiReader(stdout, stderr))

	for s.Scan() {
		text := s.Text()
		if !mute {
			log.V(text)
		}
	}
	if err := cmd.Wait(); err != nil {
		log.E(err.Error())
	}
	return cmd
}

