package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/google/uuid"
)

const (
	repoURL = "git@github.com:remi-espie/Jambon.git"
	resourcePath = "apps/nginx/templates/pod.yml"
)

func cloneRepo() string {
	uuid := uuid.New()
	dirName := fmt.Sprintf("jambon-%s", uuid.String())
	dirPath := path.Join(os.TempDir(), dirName)
	log.Print("Repository clone path: ", dirPath)

	_, err := git.PlainClone(dirPath, false, &git.CloneOptions{
	    URL:      repoURL,
	    Progress: os.Stdout,
	})

	if err != nil {
		log.Fatal("Unable to clone the git repository:", err)
	}

	return dirPath
}

func getResourceContents(repoPath string) string {
	filePath := path.Join(repoPath, resourcePath)
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("Unable to read the k8s resource file:", err)
	}

	return string(content)
}

func setResourceContents(repoPath string, content string) {
	filePath := path.Join(repoPath, resourcePath)
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		log.Fatal("Unable to write the k8s resource file:", err)
	}
}
