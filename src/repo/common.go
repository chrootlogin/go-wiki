package repo

import (
	"os"
	"log"
	"time"
	"io/ioutil"
	"path/filepath"
	"encoding/json"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/chrootlogin/go-wiki/src/common"
)

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
		log.Println("Environment variable REPOSITORY is empty. Using: " + repositoryPath)
	} else {
		log.Println("Using repository path: " + repositoryPath)
	}

	repo = initRepository()
}

func GetFile(path string) (*common.File, error) {
	path = filepath.Join(repositoryPath, "pages", path)

	// Open json file
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Println("open: " + err.Error())
		return nil, err
	}
	defer jsonFile.Close()

	// Read json file
	data, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Println("read file: " + err.Error())
		return nil, err
	}

	// Convert json to object
	var file = &common.File{}
	err = json.Unmarshal(data, file)
	if err != nil {
		log.Println("unmarshal: " + err.Error())
		return nil, err
	}

	return file, nil
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