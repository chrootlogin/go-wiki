package pagestore

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	gitfs "gopkg.in/src-d/go-git.v4/storage/filesystem"

	"github.com/chrootlogin/go-wiki/src/lib/filesystem"
)

func getRepository(fs *filesystem.Filesystem) (*git.Repository, error) {
	exists, err := fs.Has(".git")
	if err != nil {
		log.Fatal(fmt.Printf("Error checking .git repository: %s", err.Error()))
	}

	dot, _ := fs.Filesystem.Chroot(".git")

	s, err := gitfs.NewStorage(dot)
	if err != nil {
		log.Fatal(fmt.Printf("Adding storage error: %s", err.Error()))
	}

	if !exists {
		repo, err := git.Init(s, fs.Filesystem)
		if err != nil {
			log.Fatal(fmt.Sprintf("Can't create git repository: %s", err.Error()))
		}

		for _, assetName := range AssetNames() {
			fileContent, err := Asset(assetName)
			if err != nil {
				log.Fatal(fmt.Sprintf("Getting asset error: %s", err.Error()))
			}

			err = fs.Save(assetName, filesystem.File{
				Content: string(fileContent),
			})
			if err != nil {
				log.Fatal(fmt.Sprintf("Saving asset error: %s", err.Error()))
			}
		}

		wt, err := repo.Worktree()
		if err != nil {
			log.Fatal(fmt.Sprintf("Failed opening worktree: %s", err.Error()))
		}

		_, err = wt.Add(".")
		if err != nil {
			log.Fatal(err)
		}

		// Creating initial commit
		_, err = wt.Commit("Initial commit...", &git.CommitOptions{
			Author: &object.Signature{
				Name:  "system",
				Email: "go-wiki@example.org",
				When:  time.Now(),
			},
		})

		return repo, nil
	}

	return git.Open(s, fs.Filesystem)
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
