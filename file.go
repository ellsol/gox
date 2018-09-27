package gox

import (
	"os"
	"strings"
	"path/filepath"
	"fmt"
)

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func CreateFileIfNotExistent(filename string) error {
	if !FileExists(filename) {
		f, err := os.Create(filename)

		if err != nil {
			return err
		}

		err = f.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

func OpenFile(filename string) (*os.File, error) {
	return os.Create(filename)
}

func OpenFileSafe(filename string) (*os.File, error) {
	file, err := os.Open(filename)

	if err != nil && strings.Contains(err.Error(), "no such file or directory") {
		return nil, nil
	}

	return file, err
}

func RemoveDirectoryContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}

	return nil
}

func GetOnlyFileOfDirectory(dir string) (string, error){

	d, err := os.Open(dir)
	if err != nil {
		return "", err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return "", err
	}

	if len(append(names)) == 0 {
		return "", fmt.Errorf("no file found")
	}


	if len(append(names)) > 1 {
		return "", fmt.Errorf("many files found")
	}

	return names[0], nil
}