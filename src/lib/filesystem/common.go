package filesystem

import (
	"os"
	"errors"
	"log"
	"path/filepath"

	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"github.com/chrootlogin/go-wiki/src/lib/common"
	"github.com/gin-gonic/gin/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"time"
)

const DEFAULT_FILE_PERMISSIONS = 0644

var (
	ErrIsDir = errors.New("is a directory")
	ErrIsFile = errors.New("is a file")
	dataPath = ""
	filesystemCache *cache.Cache
)

type File struct {
	Content string
	Metadata Metadata
	FileInfo os.FileInfo
}

type Metadata struct {
	Permissions map[string][]string `json:"permissions"`
}

type Filesystem struct {
	Filesystem billy.Filesystem
	Error error
	FilePermissionMode os.FileMode
	ChrootDirectory string
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

	filesystemCache = cache.New(30*time.Minute, 10*time.Minute)
}

func initDataDir() {
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		log.Println("Creating new repository...")

		os.Mkdir(dataPath, os.ModePerm)
		fs := osfs.New(dataPath)
		for filename, defaultFile := range common.DefaultFiles {
			jsonData, err := json.Marshal(defaultFile)
			if err != nil {
				log.Fatal(fmt.Sprintf("Couldn't marshal default file %s: %s", filename, err.Error()))
			}

			err = saveFile(fs, DEFAULT_FILE_PERMISSIONS, filename, jsonData)
			if err != nil {
				log.Fatal(fmt.Sprintf("Couldn't save default file %s: %s", filename, err.Error()))
			}
		}
	}
}

func New(options ...Option) *Filesystem {
	// set default values
	var fs = &Filesystem{
		FilePermissionMode: DEFAULT_FILE_PERMISSIONS,
		ChrootDirectory: "",
		Filesystem: osfs.New(dataPath),
	}

	// run options
	for _, option := range options {
		err := option(fs)
		if err != nil {
			fs.Error = err
		}
	}

	return fs
}

func (fs *Filesystem) Stat(path string) (os.FileInfo, error) {
	// check for error
	if fs.Error != nil {
		return nil, fs.Error
	}

	return fs.Filesystem.Stat(path)
}

func (fs *Filesystem) Has(path string) (bool, error) {
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

func (fs *Filesystem) List(path string) ([]os.FileInfo, error) {
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

func (fs *Filesystem) Get(path string) (*File, error) {
	// check for error
	if fs.Error != nil {
		return nil, fs.Error
	}

	// set cache path
	cachePath := path
	if len(fs.ChrootDirectory) > 0 {
		cachePath = filepath.Join(fs.ChrootDirectory, cachePath)
	}

	// check if file is cached
	f, exists := filesystemCache.Get(cachePath)
	if exists {
		file, ok := f.(File)
		if !ok {
			return nil, errors.New("cached file corrupt")
		}

		return &file, nil
	}

	data, fileinfo, err := readFile(fs, path)
	if err != nil {
		return nil, err
	}

	file := &File{
		Content: string(data),
		FileInfo: fileinfo,
	}

	// write to cache
	filesystemCache.Set(cachePath, file, cache.DefaultExpiration)

	return file, nil
}

func (fs *Filesystem) Save(path string, file File) error {
	// check for error
	if fs.Error != nil {
		return fs.Error
	}

	err := saveFile(fs.Filesystem, fs.FilePermissionMode, path, []byte(file.Content))
	if err != nil {
		return err
	}

	// set cache path
	cachePath := path
	if len(fs.ChrootDirectory) > 0 {
		cachePath = filepath.Join(fs.ChrootDirectory, cachePath)
	}

	// Delete from cache path
	filesystemCache.Delete(cachePath)

	return nil
}
