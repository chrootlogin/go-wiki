package filesystem

import (
	"gopkg.in/src-d/go-billy.v4"
	"github.com/chrootlogin/go-wiki/src/repo"
	"gopkg.in/src-d/go-git.v4"
	"os"
	"errors"
)

var (
	ErrIsDir = errors.New("Is a directory.")
)

type File struct {
	Content string
	Metadata map[string]string
	FileInfo os.FileInfo
}

type filesystem struct {
	r *git.Repository
	fs billy.Filesystem
	err error
	perms os.FileMode
}

func New(options ...Option) *filesystem {
	// init internal filesystem
	var fs = &filesystem{
		r: repo.GetRepo(),
		perms: 0644,
	}

	// init worktree
	wt, err := fs.r.Worktree()
	if err != nil {
		fs.err = err;
	}
	fs.fs = wt.Filesystem

	// run options
	for _, option := range options {
		err := option(fs)
		if err != nil {
			fs.err = err
		}
	}

	return fs
}

func (fs *filesystem) Has(path string) (bool, error) {
	// check for error
	if fs.err != nil {
		return false, fs.err
	}

	if _, err := fs.fs.Stat(path); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (fs *filesystem) Get(path string) (*File, error) {
	// check for error
	if fs.err != nil {
		return nil, fs.err
	}

	return readFile(fs, path)
}
