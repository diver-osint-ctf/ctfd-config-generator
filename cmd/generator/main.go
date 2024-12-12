package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"text/template"

	_ "embed"

	"gopkg.in/yaml.v3"
	"github.com/manifoldco/promptui"
)

var (
	//go:embed templates/challenge.yml.tmpl
	challengeTemplate string

	//go:embed templates/writeup.md.tmpl
	writeupTemplate string
	
	defaultFlagPrefix = "HogeCTF"
	defaultGenres = []string{"web", "misc", "rev", "pwn"}
	flagPrefix, genres = LoadCondig("config.yaml")

	challengeFormat = "^[A-Za-z0-9_!?]+$"
	challengeRegExp = regexp.MustCompile(challengeFormat)

	flagFormat = fmt.Sprintf("^%v{[^{}]+}$", flagPrefix)
	flagRegExp = regexp.MustCompile(flagFormat)
)

type challengeInfo struct {
	FlagPrefix    string
	ChallengeName string
	Author        string
	Genre         string
	Flag          string
}

func main() {
	// get genre
	promptForSelect := promptui.Select{
		Label: "genre",
		Items: genres,
	}
	_, genre, err := promptForSelect.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed promptForSelect.Run(): %s", err.Error())
		os.Exit(1)
	}

	// get challenge name
	prompt := promptui.Prompt{
		Label: "challenge name",
		Validate: func(input string) error {
			if !challengeRegExp.MatchString(input) {
				return fmt.Errorf("challenge name and author name should meet %s", challengeFormat)
			}
			return nil
		},
	}
	challengeName, err := prompt.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed prompt.Run() for challengeName: %s", err.Error())
		os.Exit(1)
	}

	// get author name
	prompt.Label = "author name"
	author, err := prompt.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed prompt.Run() for author: %s", err.Error())
		os.Exit(1)
	}

	// get flag
	prompt = promptui.Prompt{
		Label: "flag",
		Validate: func(input string) error {
			if !flagRegExp.MatchString(input) {
				return fmt.Errorf("flag should meet %s", flagFormat)
			}
			return nil
		},
	}
	flag, err := prompt.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed prompt.Run(): %s", err.Error())
		os.Exit(1)
	}

	info := challengeInfo{
		FlagPrefix:    flagPrefix,
		ChallengeName: challengeName,
		Author:        author,
		Genre:         genre,
		Flag:          flag,
	}

	// ready a directory structure
	// - make directory(./genre/challengeName)
	//   - directory: build, files, solver
	//   - file: flag.txt, challenge.yml, writeup/README.md
	err = os.MkdirAll(filepath.Join(genre, challengeName), os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed os.MkdirAll(genre/challengeName): %s", err.Error())
		os.Exit(1)
	}

	dirs := []string{"build", "public", "solver", "writeup"}
	for _, dirName := range dirs {
		err = os.MkdirAll(filepath.Join(genre, challengeName, dirName), os.ModePerm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed os.MkdirAll(genre/challengeName/%s): %s", dirName, err.Error())
			os.Exit(1)
		}
	}

	// write default description for each file
	files := []string{"flag.txt", "challenge.yml", "writeup/README.md"}
	for _, fileName := range files {
		if err := readyFile(fileName, info); err != nil {
			fmt.Fprintf(os.Stderr, "failed readyFile: %s", err.Error())
			os.Exit(1)
		}
	}
}

func LoadCondig(path string) (string, []string) {
	type Config struct {
		FlagPrefix string `yaml:"flag_prefix"`
		Genres []string `yaml:"genre"`
	}
	var config Config = Config{
		FlagPrefix: defaultFlagPrefix,
		Genres: defaultGenres,
	}
	yml, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("failed to load config file: %s\n", err.Error())
		return config.FlagPrefix, config.Genres
	}
	err = yaml.Unmarshal(yml, &config)
	if err != nil {
		fmt.Printf("failed to unmarshal config file: %s\n", err.Error())
	}
	return config.FlagPrefix, config.Genres
}

func generateMarkdown(templateName string, templateStr string, info challengeInfo) (string, error) {
	tpl, err := template.New(templateName).Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("failed template.New: %w", err)
	}
	writer := &bytes.Buffer{}
	err = tpl.Execute(writer, info)
	return writer.String(), err
}

func readyFile(fileName string, info challengeInfo) error {
	fp, err := os.Create(filepath.Join(info.Genre, info.ChallengeName, fileName))
	if err != nil {
		return fmt.Errorf("failed os.Create(genre/challengeName/%s): %w", fileName, err)
	}
	defer fp.Close()

	// write template message
	switch fileName {
	case "flag.txt":
		fmt.Fprintln(fp, info.Flag)
		break
	case "challenge.yml":
		challengeYaml, err := generateMarkdown("challenge", challengeTemplate, info)
		if err != nil {
			return fmt.Errorf("failed generateMarkdown for challenge.yml: %w", err)
		}
		fmt.Fprint(fp, challengeYaml)
		break
	case "writeup/README.md":
		writeupMd, err := generateMarkdown("writeup", writeupTemplate, info)
		if err != nil {
			return fmt.Errorf("failed generateMarkdown for challenge.yml: %w", err)
		}
		fmt.Fprint(fp, writeupMd)
		break
	}

	return nil
}
