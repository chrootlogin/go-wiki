package filesystem

import (
	"os"
	"errors"
	"encoding/json"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-billy.v4"

	"github.com/chrootlogin/go-wiki/src/lib/repo"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"log"
	"path/filepath"
)

var (
	ErrIsDir = errors.New("is a directory")
	ErrIsFile = errors.New("is a file")
	dataPath = ""
)

type File struct {
	Content string
	Metadata Metadata
	FileInfo os.FileInfo
}

type Metadata struct {
	Permissions map[string][]string `json:"permissions"`
}

type filesystem struct {
	Repository *git.Repository
	Filesystem billy.Filesystem
	Worktree *git.Worktree
	Error error
	FilePermissionMode os.FileMode
	ChrootDirectory string
	WithMetadata bool
}

func init() {
	dataPath = os.Getenv("DATA_DIR")

	if len(dataPath) == 0 {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		dataPath = filepath.Join(dir, "data")
		log.Println("Environment variable DATA_DIR is empty. Using: " + dataPath)
	} else {
		log.Println("Using data directory: " + dataPath)
	}

	initDataDir()
}

func initDataDir() {
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		log.Println("Creating new repository...")

		os.Mkdir(dataPath, os.ModePerm)
	}
}

func New(options ...Option) *filesystem {
	// set default values
	var fs = &filesystem{
		Repository: repo.GetRepo(),
		FilePermissionMode: 0644,
		ChrootDirectory: "",
		WithMetadata: false,
		Filesystem: osfs.New(dataPath),
	}

	// init worktree
	/* wt, err := fs.Repository.Worktree()
	if err != nil {
		fs.Error = err;
		// return broken fs, error will be catched in chained function
		return fs
	} */

	//fs.Worktree = wt
	//fs.Filesystem = wt.Filesystem

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

func (fs *filesystem) List(path string) ([]os.FileInfo, error) {
	// check for error
	if fs.Error != nil {
		return []os.FileInfo{}, fs.Error
	}

	// check if is a dir
	fileinfo, err := fs.Filesystem.Stat(path)
	if err != nil {
		return []os.FileInfo{}, err
	}

	if !fileinfo.IsDir() {
		return []os.FileInfo{}, ErrIsFile
	}

	return fs.Filesystem.ReadDir(path)
}

func (fs *filesystem) Get(path string) (*File, error) {
	// check for error
	if fs.Error != nil {
		return nil, fs.Error
	}

	data, fileinfo, err := readFile(fs, path)
	if err != nil {
		return nil, err
	}

	file := &File{
		Content: string(data),
		FileInfo: fileinfo,
	}

	// if metadata enabled, get metadata file if available
	if fs.WithMetadata {
		metajson, _, err := readFile(fs, path + ".meta")
		if err != nil {
			// if meta file is not available, return empty metadata
			if os.IsNotExist(err) {
				file.Metadata = Metadata{}
				return file, nil
			}

			// otherwise return error
			return nil, err
		}

		// unmarshal and add metadata
		var metadata Metadata
		err = json.Unmarshal(metajson, &metadata)
		if err != nil {
			return nil, err
		}

		file.Metadata = metadata
	}

	return file, nil
}

func (fs *filesystem) Save(path string, file File) error {
	// check for error
	if fs.Error != nil {
		return fs.Error
	}

	err := saveFile(fs, path, []byte(file.Content))
	if err != nil {
		return err
	}

	return nil
}

func (fs *filesystem) Commit(path string, file File, commit Commit) error {
	/*// check for error
	if fs.Error != nil {
		return fs.Error
	}

	err := commitFile(fs, path, []byte(file.Content), commit)
	if err != nil {
		return err
	}

	if fs.WithMetadata {
		metadata, err := json.Marshal(file.Metadata)
		if err != nil {
			return err
		}

		return commitFile(fs, path + ".meta", metadata, commit)
	}*/

	return nil
}
