package ini

import (
	"bufio"
	//"errors"
	"fmt"
	"os"
)

type Entry struct {
	key   string
	value string
}

type Section struct {
	name    string
	entries []Entry
}

type File struct {
	header   []string
	sections []Section
}

func (f *File) Load(s string) error {
	fd, err := os.Open(s)
	if err != nil {
		return err
	}
	defer fd.Close()
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func NewFile(s string) (*File, error) {
	ret := &File{}
	err := ret.Load(s)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func Load(s string) (*File, error) {
	return NewFile(s)
}
