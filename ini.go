package ini

// to keep original data, comment and empty line are treated as "Invalid Entry"

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
	Key   string
	Value string
	Valid bool
}

func (e *Entry) Update(s string) {
	e.Value = s
}

type Section struct {
	Name    string
	Entries []*Entry
}

func (s *Section) Entry(k string) *Entry {
	for _, e := range s.Entries {
		if e.Key == k {
			return e
		}
	}
	return nil
}

func (s *Section) Add(k, v string) {
	s.Entries = append(s.Entries, &Entry{Key: k, Value: v, Valid: true})
}

func (s *Section) AddDummyEntry(k, v string) {
	s.Entries = append(s.Entries, &Entry{Key: k, Value: v, Valid: false})
}

func (s *Section) update(line string) {
	if ValidEntry(line) {
		p := strings.Split(line, "=")
		s.Add(p[0], strings.Join(p[1:], "="))
	} else {
		s.AddDummyEntry(line, "")
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

func ValidEntry(line string) bool {
	if len(line) == 0 {
		return false
	}
	if (line[0] == ';') || (line[0] == '#') {
		return false
	}
	if strings.Index(line, "=") == -1 {
		return false
	}
	return true
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
		if name == s.Name {
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
			section = &Section{Name: name}
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

func (f *File) WriteFile(path string) error {
	dest, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dest.Close()
	for _, h := range f.Header() {
		dest.WriteString(h + "\n")
	}
	for _, s := range f.Sections() {
		dest.WriteString("[" + s.Name + "]\n")
		for _, e := range s.Entries {
			if e.Valid {
				dest.WriteString(e.Key + "=" + e.Value + "\n")
			} else {
				dest.WriteString(e.Key + e.Value + "\n")
			}
		}
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
	ret := &Section{Name: name}
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

func (f *File) AppendEntry(s, k, v string) error {
	section := f.Section(s)
	if section == nil {
		return errors.New(fmt.Sprintf("section [%s] not found", s))
	}
	target := section.Entry(k)
	if target == nil {
		return errors.New(fmt.Sprintf("entry [%s]/%s not found", s, k))
	}
	target.Update(target.Value + "," + v)
	return nil
}

func (f *File) HasValue(s, k, v string) bool {
	section := f.Section(s)
	if section == nil {
		return false
	}
	entry := section.Entry(k)
	if entry == nil {
		return false
	}
	for _, p := range strings.Split(entry.Value, ",") {
		raw := strings.Replace(p, " ", "", -1)
		if v == raw {
			return true
		}
	}
	return false
}
