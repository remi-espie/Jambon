package main

import (
	"context"
	"fmt"
	ssh2 "golang.org/x/crypto/ssh"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/google/go-github/v72/github"
	"github.com/google/uuid"
)

const (
	repoURL      = "git@github.com:remi-espie/Jambon.git"
	resourcePath = "apps/nginx/templates/deployment.yml"

	githubRepoOwner = "remi-espie"
	githubRepoName  = "Jambon"
)

func cloneRepo(sshKey string) (*git.Repository, string) {
	uuid := uuid.New()
	dirName := fmt.Sprintf("jambon-%s", uuid.String())
	dirPath := path.Join(os.TempDir(), dirName)
	log.Print("Repository clone path: ", dirPath)

	publicKey, err := ssh.NewPublicKeys("git", []byte(sshKey), "")

	if err != nil {
		log.Fatal("Unable to create the SSH public key for git:", err)
	}

	publicKey.HostKeyCallback = ssh2.InsecureIgnoreHostKey()

	repo, err := git.PlainClone(dirPath, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
		Auth:     publicKey,
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

func pushAutofix(repo *git.Repository, sshKey string, commitMessage string) {
	worktree, err := repo.Worktree()

	if err != nil {
		log.Fatal("Unable to get the worktree from the git repository:", err)
	}

	_, err = worktree.Add(resourcePath)

	if err != nil {
		log.Fatal("Unable to add modified file to the staging area:", err)
	}

	commitSig := object.Signature{
		Name:  "Qwen",
		Email: "qwen@jambon.bayonne",
		When:  time.Now(),
	}

	_, err = worktree.Commit(commitMessage, &git.CommitOptions{
		Author:    &commitSig,
		Committer: &commitSig,
	})

	if err != nil {
		log.Fatal("Unable to commit the autofixed change in the git repository:", err)
	}

	publicKey, err := ssh.NewPublicKeys("git", []byte(sshKey), "")

	if err != nil {
		log.Fatal("Unable to create the SSH public key for git:", err)
	}

	publicKey.HostKeyCallback = ssh2.InsecureIgnoreHostKey()

	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       publicKey,
	})

	if err != nil {
		log.Fatal("Unable to push autofix to the remote:", err)
	}
}

func mergeAutofix(repo *git.Repository, token string) {
	client := github.NewClient(nil).WithAuthToken(token)

	head, err := repo.Head()

	if err != nil {
		log.Fatal("Unable to get the git repository head:", err)
	}

	cIter, err := repo.Log(&git.LogOptions{From: head.Hash()})

	if err != nil {
		log.Fatal("Unable to get the git repository log iter:", err)
	}

	commit, err := cIter.Next()

	if err != nil {
		log.Fatal("Unable to get the first commit:", err)
	}

	cIter.Close()

	commitTitle := strings.Split(commit.Message, "\n")[0]
	headName := head.Name().Short()
	baseName := "main"

	pr, _, err := client.PullRequests.Create(context.TODO(), githubRepoOwner, githubRepoName, &github.NewPullRequest{
		Title: &commitTitle,
		Head:  &headName,
		Base:  &baseName,
	})

	if err != nil {
		log.Fatal("Unable to create the pull request:", err)
	}

	_, _, err = client.PullRequests.Merge(context.TODO(), githubRepoOwner, githubRepoName, *pr.Number, "", &github.PullRequestOptions{MergeMethod: "squash"})

	if err != nil {
		log.Fatal("Unable to merge the pull request:", err)
	}

	_, err = client.Git.DeleteRef(context.TODO(), githubRepoOwner, githubRepoName, head.Name().String())

	if err != nil {
		log.Fatal("Unable to delete the remote git branch:", err)
	}
}
