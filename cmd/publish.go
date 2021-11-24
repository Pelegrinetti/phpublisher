package cmd

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Pelegrinetti/phpublisher/lib/zip"
	"github.com/urfave/cli/v2"
)

func getWorkdir() string {
	wd, _ := os.Getwd()

	return wd
}

func getIgnoredFiles() []string {
	ignoreFilePath := getWorkdir() + "/.phpublisherignore"

	files := []string{
		".phpublisherignore",
		".git",
		".DS_Store",
	}

	ignoreFile, _ := os.Open(ignoreFilePath)
	defer ignoreFile.Close()

	scanner := bufio.NewScanner(ignoreFile)

	for scanner.Scan() {
		files = append(files, scanner.Text())
	}

	return files
}

func getFiles() []string {
	root := getWorkdir()
	files := []string{}

	// TODO: Ignorar os paths que são ignorados antes de fazer a leitura inteira do diretório
	// 			 para que não seja necessário caminhar por todos os arquivos de um diretório interno.
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		pathWithoutRoot := strings.Replace(path, root+"/", "", -1)

		if !d.IsDir() {
			files = append(files, pathWithoutRoot)
		}

		return nil
	})

	newFiles := []string{}
	ignoredFiles := getIgnoredFiles()

	for _, file := range files {
		count := 0

		for _, ignoredFile := range ignoredFiles {
			regex := regexp.MustCompile(ignoredFile)
			matches := regex.MatchString(file)

			if !matches {
				count++
			}
		}

		if count == len(ignoredFiles) {
			newFiles = append(newFiles, file)
		}
	}

	if err != nil {
		panic(err)
	}

	return newFiles
}

func createTemporaryFile(version string) (*os.File, func() error) {
	file, _ := os.CreateTemp("", version+"-")

	return file, func() error {
		return os.Remove(file.Name())
	}
}

func PublishCmd() *cli.Command {
	var version, registry string

	return &cli.Command{
		Name:  "publish",
		Usage: "Publish a package in Nexus.",
		Action: func(c *cli.Context) error {
			files := getFiles()

			zipFile, removeTemporaryFile := createTemporaryFile(version)
			defer removeTemporaryFile()

			if err := zip.ZipFiles(zipFile, files); err != nil {
				panic(err)
			}

			fmt.Println("Zipped File:", zipFile.Name())

			// Tomorrow

			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "registry", Aliases: []string{"r"}, Required: true, Destination: &registry},
			&cli.StringFlag{Name: "version", Aliases: []string{"v"}, Required: true, Destination: &version},
		},
	}
}
