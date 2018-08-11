package filesystem

import (
	"os"
	"errors"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-billy.v4"

	"github.com/chrootlogin/go-wiki/src/repo"
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
	Repository *git.Repository
	Filesystem billy.Filesystem
	Worktree *git.Worktree
	Error error
	FilePermissionMode os.FileMode
	Chroot string
}

func New(options ...Option) *filesystem {
	// init internal filesystem
	var fs = &filesystem{
		Repository: repo.GetRepo(),
		FilePermissionMode: 0644,
	}

	// init worktree
	wt, err := fs.Repository.Worktree()
	if err != nil {
		fs.Error = err;
	}

	fs.Worktree = wt
	fs.Filesystem = wt.Filesystem

	// run options
	for _, option := range options {
		err := option(fs)
		if err != nil {
			fs.Error = err
		}
	}

	return fs
}

func (fs *filesystem) Has(path string) (bool, error) {
	// check for error
	if fs.Error != nil {
		return false, fs.Error
	}

	if _, err := fs.Filesystem.Stat(path); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (fs *filesystem) Get(path string) (*File, error) {
	// check for error
	if fs.Error != nil {
		return nil, fs.Error
	}

	return readFile(fs, path)
}

func (fs *filesystem) Commit(path string, file File, commit repo.Commit) error {
	// check for error
	if fs.Error != nil {
		return fs.Error
	}

	return commitFile(fs, path, []byte(file.Content), commit)
}
