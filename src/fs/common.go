package fs

import (
	"gopkg.in/src-d/go-billy.v4"
	"github.com/chrootlogin/go-wiki/src/repo"
	"gopkg.in/src-d/go-git.v4"
	"os"
)

type filesystem struct {
	r *git.Repository
	fs billy.Filesystem
	err error
	perms os.FileMode
}

func New(options ...Option) *filesystem {
	// init internal fs
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

func (fs *filesystem) Get(path string) ([]byte, error) {
	// check for error
	if fs.err != nil {
		return nil, fs.err
	}

	return readFile(fs, path)
}
