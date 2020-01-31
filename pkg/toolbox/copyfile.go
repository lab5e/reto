package toolbox

import (
	"errors"
	"io/ioutil"
	"os"
)

const DefaultFilePerm = 0700

// CopyFile copies the a file from A to B. The file B must not exist prior to this
func CopyFile(from, to string) error {
	_, err := os.Stat(to)
	if !os.IsNotExist(err) {
		return errors.New("file already exists")
	}
	_, err = os.Stat(from)
	if os.IsNotExist(err) {
		return errors.New("file does not exist")
	}
	input, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(to, input, DefaultFilePerm); err != nil {
		return err
	}
	return nil
}
