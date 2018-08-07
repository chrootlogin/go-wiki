package repo

import (
	"os"
	"log"
	"time"
	"path/filepath"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

const permissions = 0644

var repositoryPath = ""
var repo *git.Repository

func init() {
	repositoryPath = os.Getenv("REPOSITORY_PATH")

	if len(repositoryPath) == 0 {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		repositoryPath = filepath.Join(dir, "data")
		log.Println("Environment variable REPOSITORY_PATH is empty. Using: " + repositoryPath)
	} else {
		log.Println("Using repository path: " + repositoryPath)
	}

	repo = initRepository()
}

func initRepository() *git.Repository {
	if _, err := os.Stat(repositoryPath); os.IsNotExist(err) {
		log.Println("Creating new repository...")

		os.Mkdir(repositoryPath, os.ModePerm)
	}

	repo, err := git.PlainOpen(repositoryPath)
	if err == git.ErrRepositoryNotExists {
		log.Println("Initializing new git repository...")
		repo, err = git.PlainInit(repositoryPath, false)
		if err != nil {
			log.Fatal(err)
		}

		// Restore all assets from default.go
		err := RestoreAssets(repositoryPath, "")
		if err != nil {
			log.Fatal(err)
		}

		wt, err := repo.Worktree()
		if err != nil {
			log.Fatal(err)
		}

		// Add all files
		_, err = wt.Add(".")
		if err != nil {
			log.Fatal(err)
		}

		// Creating initial commit
		_, err = wt.Commit("Initial commit...", &git.CommitOptions{
			Author: &object.Signature{
				Name:  "Go Wiki",
				Email: "go-wiki@example.org",
				When:  time.Now(),
			},
		})
	} else if err != nil {
		log.Fatal(err)
	}

	return repo
}