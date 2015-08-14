package utils

import (
	"io/ioutil"
	"os"
)

/*
	Return the file content in a string.
*/
func FileReadAll(filepath string) (string, error) {
	contents, err := ioutil.ReadFile(filepath)
	return string(contents), err
}

/*
	Return the file content in a []byte.
*/
func FileReadAllBytes(filepath string) ([]byte, error) {
	contents, err := ioutil.ReadFile(filepath)
	return contents, err
}

/*
	Write a string in a file.
*/
func FileWriteAll(filepath, content string) error {
	return ioutil.WriteFile(filepath, []byte(content), 0600)
}

/*
	Write a []byte in a file.
*/
func FileWriteAllBytes(filepath string, content []byte) error {
	return ioutil.WriteFile(filepath, content, 0600)
}

/*
	Test if a file exists.
	(if the target is a dir, the function returns false)
*/
func FileExists(path string) bool {
	if info, err := os.Stat(path); err == nil {
		if !info.IsDir() {
			return true
		}
	}
	return false
}

/*
	Test if a dir exists.
	(if the target is a file, the function returns false)
*/
func DirExists(path string) bool {
	if info, err := os.Stat(path); err == nil {
		if info.IsDir() {
			return true
		}
	}
	return false
}
