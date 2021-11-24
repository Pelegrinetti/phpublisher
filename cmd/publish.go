package cmd

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/fs"
	"net/http"
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
	var version, registry, vendor, project, user, password string

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

			zipFile.Close()

			zipFileName := zipFile.Name()

			fmt.Println("Zipped File:", zipFileName)

			fmt.Println("Uploading file...")

			repositoryUrl := fmt.Sprintf("%s/%s/%s/%s", registry, vendor, project, version)

			basicToken := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, password)))

			client := &http.Client{}

			file, _ := os.Open(zipFileName)
			defer file.Close()

			req, err := http.NewRequest(http.MethodPut, repositoryUrl, file)
			if err != nil {
				panic(err)
			}
			defer req.Body.Close()

			req.Header.Set("Authorization", fmt.Sprintf("Basic %s", basicToken))

			response, responseError := client.Do(req)
			if responseError != nil {
				fmt.Println("ERROR: ", responseError)
			}
			defer response.Body.Close()

			fmt.Println("Package published! Request status: ", response.StatusCode)

			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "registry", Required: true, Destination: &registry},
			&cli.StringFlag{Name: "version", Required: true, Destination: &version},
			&cli.StringFlag{Name: "vendor", Required: true, Destination: &vendor},
			&cli.StringFlag{Name: "project", Required: true, Destination: &project},
			&cli.StringFlag{Name: "user", Required: true, Destination: &user},
			&cli.StringFlag{Name: "password", Required: true, Destination: &password},
		},
	}
}
