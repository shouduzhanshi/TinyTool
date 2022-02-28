package tool

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

func Log(msg string) {
	fmt.Println(msg)
}

func Adb(raw ...string)(int, error) {
	return ExecCmd("adb",raw...)
}

func ExecCmd(shell string, raw ...string) (int, error) {
	cmd := exec.Command(shell, raw...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		Log(err.Error())
		return 0, nil
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		Log(err.Error())
		return 0, nil
	}
	if err := cmd.Start(); err != nil {
		Log(err.Error())
		return 0, nil
	}
	s := bufio.NewScanner(io.MultiReader(stdout, stderr))
	for s.Scan() {
		text := s.Text()
		Log(text)
	}
	if err := cmd.Wait(); err != nil {
		Log(err.Error())
	}

	return cmd.ProcessState.ExitCode(), nil
}


func DeCodeAppJson(appJson string) map[string]interface{} {
	if file, err := ioutil.ReadFile(appJson); err == nil {
		decoder := json.NewDecoder(bytes.NewBuffer(file))
		buildConfig := make(map[string]interface{})
		decoder.Decode(&buildConfig)
		return buildConfig
	}else{
		panic(err)
	}
	return nil
}


func GetApplicationId(buildConfig map[string]interface{}) string {
	build := buildConfig["build"].(map[string]interface{})

	applicationId := build["applicationId"].(string)
	return applicationId
}


func GetCurrentPath() string {
	if getwd, err := os.Getwd(); err == nil {
		return getwd
	} else {
		panic(err)
	}
}
