package ini

// to keep original data, comment and empty line are treated as "Invalid Entry"

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

var SectionStartKey = regexp.MustCompile("^\\[.*\\]")
var SectionTitleTrim = regexp.MustCompile(`\[|\]| `)
var ValueSeparator = ","

func ParseSectionName(line string) string {
	ret := string(line)
	return SectionTitleTrim.ReplaceAllString(ret, "")
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

type Entry struct {
	Key   string
	Value string
	Valid bool
}

func (e *Entry) OverWrite(s string) {
	e.Value = s
}

func (e *Entry) Append(s string) {
	if e.Value == "" {
		e.Value = s
	} else {
		e.Value = e.Value + ValueSeparator + s
	}
}

func (e *Entry) String() string {
	if e.Valid {
		return fmt.Sprintf("%v=%v", e.Key, e.Value)
	} else {
		return fmt.Sprintf("%v%v", e.Key, e.Value)
	}
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
	Header   []string
	Sections []*Section
}

func NewFileScan(s *bufio.Scanner) (*File, error) {
	ret := &File{}
	ret.Load(s)
	if err := s.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

func NewFile(file string) (*File, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	scanner := bufio.NewScanner(fd)
	return NewFileScan(scanner)
}

func Load(file string) (*File, error) {
	return NewFile(file)
}

func LoadText(t string) (*File, error) {
	return NewFileScan(bufio.NewScanner(strings.NewReader(t)))
}

func (f *File) NumOfSections() int {
	return len(f.Sections)
}

func (f *File) Section(name string) *Section {
	for _, s := range f.Sections {
		if name == s.Name {
			return s
		}
	}
	return nil
}

func (f *File) Load(scanner *bufio.Scanner) {
	var section *Section
	section = nil
	for scanner.Scan() {
		line := scanner.Text()

		if SectionStartKey.FindStringIndex(line) != nil {
			name := ParseSectionName(line)
			section = &Section{Name: name}
			f.Sections = append(f.Sections, section)
			continue
		}

		if section != nil {
			section.update(line)
		}

		if f.NumOfSections() == 0 {
			f.Header = append(f.Header, line)
		}
	}
}

func (f *File) RawWrite(w io.Writer) error {
	for _, h := range f.Header {
		fmt.Fprintf(w, "%s\n", h)
	}
	for _, s := range f.Sections {
		fmt.Fprintf(w, "[%s]\n", s.Name)
		for _, e := range s.Entries {
			if e.Valid {
				fmt.Fprintf(w, "%s=%s\n", e.Key, e.Value)
			} else {
				fmt.Fprintf(w, "%s%s\n", e.Key, e.Value)
			}
		}
	}
	return nil
}

func (f *File) Write(w *bufio.Writer) error {
	defer w.Flush()
	return f.RawWrite(w)
}

func (f *File) WriteFile(path string) error {
	dest, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dest.Close()
	ret := f.Write(bufio.NewWriter(dest))
	return ret
}

func (f *File) AddSection(name string) (*Section, error) {
	if f.Section(name) != nil {
		return nil, errors.New(fmt.Sprintf("section [%s] already exists", name))
	}
	ret := &Section{Name: name}
	f.Sections = append(f.Sections, ret)
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

func (f *File) AppendEntry(s, k, v string) error {
	section := f.Section(s)
	if section == nil {
		return errors.New(fmt.Sprintf("section [%s] not found", s))
	}
	target := section.Entry(k)
	if target == nil {
		return errors.New(fmt.Sprintf("entry [%s]/%s not found", s, k))
	}
	target.Append(v)
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
	for _, p := range strings.Split(entry.Value, ValueSeparator) {
		raw := strings.Replace(p, " ", "", -1)
		if v == raw {
			return true
		}
	}
	return false
}
