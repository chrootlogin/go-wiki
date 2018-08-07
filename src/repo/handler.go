package repo

import (
	"path/filepath"
	"os"
	"log"
	"io/ioutil"
	"github.com/chrootlogin/go-wiki/src/common"
	"encoding/json"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"time"
	"fmt"
)

func HasRaw(path string) bool {
	path = filepath.Join(repositoryPath, path)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func HasFile(path string) bool {
	path = filepath.Join("pages", path)

	return HasRaw(path)
}

func GetRaw(path string) ([]byte, error) {
	path = filepath.Join(repositoryPath, path)

	// Open json file
	file, err := os.Open(path)
	if err != nil {
		log.Println("open: " + err.Error())
		return nil, err
	}

	// Read json file
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("read file: " + err.Error())
		return nil, err
	}

	return data, nil
}

func GetFile(path string) (*common.File, error) {
	path = filepath.Join("pages", path)

	commits := getCommitHistoryOfFile(path)
	fmt.Println(commits)

	data, err := GetRaw(path)
	if err != nil {
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

func SaveRaw(path string, data []byte, commit Commit) error {
	diskPath := filepath.Join(repositoryPath, path)

	// Write file
	err := ioutil.WriteFile(diskPath, data, permissions)
	if err != nil {
		log.Println("write file: " + err.Error())
		return err
	}

	// Open worktree
	wt, err := repo.Worktree()
	if err != nil {
		log.Println("worktree: " + err.Error())
		return err
	}

	// Add file
	_, err = wt.Add(path)
	if err != nil {
		log.Println("adding file: " + err.Error())
		return err
	}

	// Creating commit
	_, err = wt.Commit(commit.Message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  commit.Author.Username,
			Email: commit.Author.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		log.Println("commit: " + err.Error())
		return err
	}

	return nil
}

func SaveFile(path string, file *common.File, commit Commit) error {
	path = filepath.Join("pages", path)

	jsonBytes, err := json.Marshal(file)
	if err != nil {
		log.Println("marshal: " + err.Error())
		return err
	}

	return SaveRaw(path, jsonBytes, commit)
}

func MkdirPage(path string) error {
	path = filepath.Join(repositoryPath, "pages", path)

	return os.MkdirAll(path, os.ModePerm);
}

func getCommitHistoryOfFile(path string) []*object.Commit {
	objects := []*object.Commit{}

	ref, err := repo.Head()
	if err == nil {
		cIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
		if err == nil {
			cIter.ForEach(filterByChangesToPath(repo, path, func(c *object.Commit) error {
				objects = append(objects, c)
				return nil
			}))
		}
	}

	return objects
}