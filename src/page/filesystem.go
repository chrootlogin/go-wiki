package page

import (
	"log"
	"os"
	"io/ioutil"

	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-git.v4"

	"github.com/chrootlogin/go-wiki/src/repo"
)

const (
	CHROOT_DIR = "pages"
)

type filesystem struct {
	r *git.Repository
}

func PageFS() (*filesystem) {
	var fs = filesystem{
		r: repo.GetRepo(),
	}

	return &fs
}

func (fs *filesystem) Has(path string) (bool, error) {
	_, chroot, err := getWorktreeAndChroot(fs.r)
	if err != nil {
		return false, err
	}

	_, err = chroot.Stat(path);
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (fs *filesystem) Get(path string) ([]byte, error) {
	_, chroot, err := getWorktreeAndChroot(fs.r)
	if err != nil {
		return nil, err
	}

	return readFile(chroot, path)
}

func getWorktreeAndChroot(r *git.Repository) (*git.Worktree, billy.Filesystem, error){
	// open workspace
	wt, err := r.Worktree()
	if err != nil {
		log.Println("opening worktree: " + err.Error())
		return nil, nil, err
	}

	fs, err := wt.Filesystem.Chroot(CHROOT_DIR)
	if err != nil {
		log.Println("open file: " + err.Error())
		return nil, nil, err
	}

	return wt, fs, nil
}

func readFile(fs billy.Filesystem, path string) ([]byte, error) {
	// Open file
	file, err := fs.OpenFile(path, os.O_RDONLY, repo.R_PERMS)
	if err != nil {
		log.Println("open file: " + err.Error())
		return nil, err
	}

	// Read file
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("read file: " + err.Error())
		return nil, err
	}

	err = file.Close()
	if err != nil {
		log.Println("close file: " + err.Error())
		return nil, err
	}

	return data, nil
}