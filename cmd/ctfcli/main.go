package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

func LoadCondig(path string) ([]string) {
	type Config struct {
		Genres []string `yaml:"genre"`
	}
	var config Config = Config{
		Genres: []string{"web", "misc", "rev", "pwn"},
	}
	yml, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("failed to load config file: %s\n", err.Error())
		return config.Genres
	}
	err = yaml.Unmarshal(yml, &config)
	if err != nil {
		fmt.Printf("failed to unmarshal config file: %s\n", err.Error())
	}
	return config.Genres
}

func execCmd(cmd *exec.Cmd) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return cmd.String(), stdout.String(), err
}

func main() {
	genres := LoadCondig("config.yaml")

	for _, genre := range genres {
		if _, err := os.Stat(genre); os.IsNotExist(err) {
			fmt.Printf("No genre found: %s\n", genre)
			continue
		}

		var sorted_challs []string
		filepath.Walk(genre, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() && strings.HasSuffix(info.Name(), "challenge.yaml") {
				sorted_challs = append(sorted_challs, path)
			}
			return nil
		})

		sort.Strings(sorted_challs)

		for _, c := range sorted_challs {
			fp, err := filepath.Abs(c)
			if err != nil {
				fmt.Println("Error getting absolute path: ", err.Error())
				continue
			}
			fmt.Printf("Syncing challenge: %s\n", fp)
			cmd := exec.Command("python", "-m", "ctfcli", "challenge", "sync", fp)
			stdout, stderr, err := execCmd(cmd)
			fmt.Println(stdout)
			if err != nil {
				if(strings.Contains(stderr, "Perhaps you meant install instead of sync")) {
					fmt.Println("Installing challenge instead of syncing...")
					cmd := exec.Command("python", "-m", "ctfcli", "challenge", "install", fp)
					stdout, stderr, err := execCmd(cmd)
					fmt.Println(stdout)
					if err != nil {
						fmt.Println("Error installing challenge: ", err.Error())
						fmt.Printf("Stderr: %s\n", stderr)
					}
				} else {
					fmt.Println("Error syncing challenge: ", err.Error())
					fmt.Printf("Stderr: %s\n", stderr)
				}
			}

			time.Sleep(1 * time.Second)
		}
	}
}
