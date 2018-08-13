package gox

import "os"

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
