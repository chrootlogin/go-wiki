package repo

import (
	"log"
	"os"
	"time"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

)

func SaveRaw(path string, data []byte, commit Commit) error {
	// open workspace
	wt, err := repo.Worktree()
	if err != nil {
		log.Println("opening worktree: " + err.Error())
		return err
	}

	// Open file
	file, err := wt.Filesystem.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, R_PERMS)
	if err != nil {
		log.Println("open file: " + err.Error())
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		file.Close()
		log.Println("write to file: " + err.Error())
		return err
	}

	// close file
	err = file.Close()
	if err != nil {
		log.Println("close file: " + err.Error())
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