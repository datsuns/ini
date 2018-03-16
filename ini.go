package ini

import (
	"bufio"
	//"errors"
	//	"fmt"
	"os"
	"regexp"
	"strings"
)

var SectionStartKey = regexp.MustCompile("^\\[.*\\]")

type Entry struct {
	key   string
	value string
}

type Section struct {
	name    string
	entries []Entry
}

func (s *Section) Name() string {
	return s.name
}

type File struct {
	header   []string
	sections []Section
}

func ParseSectionName(line string) string {
	ret := string(line)
	strings.Replace(ret, "[", "", -1)
	strings.Replace(ret, "]", "", -1)
	strings.Replace(ret, " ", "", -1)
	return ret
}

func (f *File) Header() []string {
	return f.header
}

func (f *File) Sections() []Section {
	return f.sections
}

func (f *File) NumOfSections() int {
	return len(f.sections)
}

func (f *File) loadMain(scanner *bufio.Scanner) {
	var section *Section
	section = nil
	for scanner.Scan() {
		line := scanner.Text()

		if SectionStartKey.FindStringIndex(line) != nil {
			section = &Section{name: ParseSectionName(line)}
			f.sections = append(f.sections, *section)
		}

		if f.NumOfSections() == 0 {
			f.header = append(f.header, line)
		}
	}
}

func (f *File) Load(s string) error {
	fd, err := os.Open(s)
	if err != nil {
		return err
	}
	defer fd.Close()
	scanner := bufio.NewScanner(fd)
	f.loadMain(scanner)
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
