package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/uuid"
)

const (
	repoURL = "git@github.com:remi-espie/Jambon.git"
	resourcePath = "apps/nginx/templates/pod.yml"
)

func cloneRepo() (*git.Repository, string) {
	uuid := uuid.New()
	dirName := fmt.Sprintf("jambon-%s", uuid.String())
	dirPath := path.Join(os.TempDir(), dirName)
	log.Print("Repository clone path: ", dirPath)

	repo, err := git.PlainClone(dirPath, false, &git.CloneOptions{
	    URL:      repoURL,
	    Progress: os.Stdout,
	})

	if err != nil {
		log.Fatal("Unable to clone the git repository:", err)
	}

	worktree, err := repo.Worktree()

	if err != nil {
		log.Fatal("Unable to get the worktree from the git repository:", err)
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/autofix/%s", uuid.String())),
		Create: true,
	})

	if err != nil {
		log.Fatal("Unable to switch branch in the git repository:", err)
	}

	return repo, dirPath
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

func pushAutofix(repo *git.Repository, commitMessage string) {
	worktree, err := repo.Worktree()

	if err != nil {
		log.Fatal("Unable to get the worktree from the git repository:", err)
	}

	_, err = worktree.Add(resourcePath)

	if err != nil {
		log.Fatal("Unable to add modified file to the staging area:", err)
	}

	commitSig := object.Signature{
		Name: "Qwen",
		Email: "qwen@jambon.bayonne",
		When: time.Now(),
	}

	_, err = worktree.Commit(commitMessage, &git.CommitOptions{
		Author: &commitSig,
		Committer: &commitSig,
	})

	if err != nil {
		log.Fatal("Unable to commit the autofixed change in the git repository:", err)
	}

	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
	})

	if err != nil {
		log.Fatal("Unable to push autofix to the remote:", err)
	}
}
