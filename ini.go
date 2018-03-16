package ini

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var SectionStartKey = regexp.MustCompile("^\\[.*\\]")

type Entry struct {
	key   string
	value string
}

func (e *Entry) Key() string {
	return e.key
}

func (e *Entry) Value() string {
	return e.value
}

func (e *Entry) Update(s string) {
	e.value = s
}

type Section struct {
	name    string
	entries []*Entry
}

func (s *Section) Name() string {
	return s.name
}

func (s *Section) Entries() []*Entry {
	return s.entries
}

func (s *Section) Entry(k string) *Entry {
	for _, e := range s.entries {
		if e.Key() == k {
			return e
		}
	}
	return nil
}

func (s *Section) Add(k, v string) {
	s.entries = append(s.entries, &Entry{key: k, value: v})
}

func (s *Section) update(line string) {
	p := strings.Split(line, "=")
	if len(p) == 2 {
		s.Add(p[0], p[1])
	} else {
		s.Add(p[0], "")
	}
}

type File struct {
	header   []string
	sections []*Section
}

func ParseSectionName(line string) string {
	ret := string(line)
	ret = strings.Replace(ret, "[", "", -1)
	ret = strings.Replace(ret, "]", "", -1)
	ret = strings.Replace(ret, " ", "", -1)
	return ret
}

func (f *File) Header() []string {
	return f.header
}

func (f *File) Sections() []*Section {
	return f.sections
}

func (f *File) NumOfSections() int {
	return len(f.sections)
}

func (f *File) Section(name string) *Section {
	for _, s := range f.Sections() {
		if name == s.Name() {
			return s
		}
	}
	return nil
}

func (f *File) loadMain(scanner *bufio.Scanner) {
	var section *Section
	section = nil
	for scanner.Scan() {
		line := scanner.Text()

		if SectionStartKey.FindStringIndex(line) != nil {
			name := ParseSectionName(line)
			section = &Section{name: name}
			f.sections = append(f.sections, section)
			continue
		}

		if section != nil {
			section.update(line)
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

func (f *File) AddSection(name string) (*Section, error) {
	if f.Section(name) != nil {
		return nil, errors.New(fmt.Sprintf("section [%s] already exists", name))
	}
	ret := &Section{name: name}
	f.sections = append(f.sections, ret)
	return ret, nil
}

func (f *File) AddEntry(s, k, v string) error {
	dest := f.Section(s)
	if dest == nil {
		return errors.New(fmt.Sprintf("section [%s] not found", s))
	}
	dest.Add(k, v)
	return nil
}

func (f *File) ModifyEntry(s, k, v string) error {
	section := f.Section(s)
	if section == nil {
		return errors.New(fmt.Sprintf("section [%s] not found", s))
	}
	target := section.Entry(k)
	if target == nil {
		return errors.New(fmt.Sprintf("entry [%s]/%s not found", s, k))
	}
	target.Update(v)
	return nil
}
