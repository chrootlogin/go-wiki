package pagestore

import (
	"time"
	"log"
	"fmt"

	"gopkg.in/src-d/go-git.v4"
	gitfs "gopkg.in/src-d/go-git.v4/storage/filesystem"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/chrootlogin/go-wiki/src/lib/filesystem"
)

func getRepository(fs *filesystem.Filesystem) (*git.Repository, error) {
	exists, err := fs.Has(".git")
	if err != nil {
		log.Fatal(fmt.Printf("Error checking .git repository: %s", err.Error()))
	}

	wt := fs.Filesystem
	dot, _ := wt.Chroot(".git")

	s, err := gitfs.NewStorage(dot)
	if err != nil {
		log.Fatal(fmt.Printf("Adding storage error: %s", err.Error()))
	}

	if !exists {
		return git.Init(s, wt)
	}

	return git.Open(s, wt)
}

func commitFile(ps *pagestore, path string, file filesystem.File, commit Commit) error {
	err := ps.filesystem.Save(path, file)
	if err != nil {
		return err
	}

	wt, err := ps.repository.Worktree()
	if err != nil {
		return err
	}

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