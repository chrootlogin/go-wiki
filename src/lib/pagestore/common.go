package pagestore

import (
	"fmt"
	"log"

	"gopkg.in/src-d/go-git.v4"

	"github.com/chrootlogin/go-wiki/src/lib/filesystem"
)

const PAGES_DIR = "pages"

type pagestore struct {
	filesystem *filesystem.Filesystem
	repository *git.Repository
}

func New() *pagestore {
	fs := filesystem.New(filesystem.WithChroot(PAGES_DIR))
	repo, err := getRepository(fs)
	if err != nil {
		log.Fatal(fmt.Printf("Error opening pages repository: %s", err.Error()))
	}

	return &pagestore{
		repository: repo,
		filesystem: fs,
	}
}

func (ps *pagestore) Get(path string) (*filesystem.File, error) {
	return ps.filesystem.Get(path)
}

func (ps *pagestore) Has(path string) (bool, error) {
	return ps.filesystem.Has(path)
}

func (ps *pagestore) Commit(path string, file filesystem.File, commit Commit) error {
	err := commitFile(ps, path, file, commit)
	if err != nil {
		return err
	}

	return nil
}
