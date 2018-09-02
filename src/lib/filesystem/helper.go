package filesystem

import (
	"os"
	"log"
	"io/ioutil"
)

func readFile(fs *Filesystem, path string) ([]byte, os.FileInfo, error) {
	// Get FileInfo
	fileinfo, err := fs.Filesystem.Stat(path)
	if err != nil {
		return nil, nil, err
	}

	if fileinfo.IsDir() {
		return nil, nil, ErrIsDir
	}

	// Open file
	file, err := fs.Filesystem.OpenFile(path, os.O_RDONLY, fs.FilePermissionMode)
	if err != nil {
		log.Println("open file: " + err.Error())
		return nil, nil, err
	}

	// Read file
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("read file: " + err.Error())
		return nil, nil, err
	}

	// Close file
	err = file.Close()
	if err != nil {
		log.Println("close file: " + err.Error())
		return nil, nil, err
	}

	return data, fileinfo, nil
}

func saveFile(fs *Filesystem, path string, data []byte) error {
	// Open file
	file, err := fs.Filesystem.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fs.FilePermissionMode)
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

	return nil
}