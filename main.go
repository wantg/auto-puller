package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

type config struct {
	Path                string `json:"path"`
	AdditionalInstructs []struct {
		Path      string `json:"path"`
		Instructs string `json:"instructs"`
	} `json:"additional-instructs"`
}

func AppPath(subPath string) string {
	rootPath, _ := os.Executable()
	return path.Join(path.Dir(rootPath), subPath)
}

func loadConfig() *config {
	c := config{}
	configBytes, _ := ioutil.ReadFile(AppPath("/config.json"))
	json.Unmarshal(configBytes, &c)
	return &c
}

func runInstruct(path string, instruct string) string {

	cmd := exec.Command("sh", "-c", instruct)
	cmd.Dir = path
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(path, instruct)
	log.Println(string(stdoutStderr))
	log.Println()
	return string(stdoutStderr)
}

func main() {
	for {
		config := loadConfig()
		additionalInstructs := config.AdditionalInstructs
		stdoutStderr := runInstruct(config.Path, "git remote update;git status -uno")
		upToDate := strings.Index(string(stdoutStderr), "up-to-date") > -1
		if !upToDate {
			for _, additionalInstruct := range additionalInstructs {
				path := additionalInstruct.Path
				if path == "." {
					path = config.Path
				}
				instructs := additionalInstruct.Instructs
				runInstruct(path, instructs)
			}
		}
		time.Sleep(1 * time.Second)
	}
}
