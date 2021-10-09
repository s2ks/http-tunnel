package config

/* NOTE: This is a super simple parser, don't expect any
fancy features */

import (
	"os"
	"bufio"
	"strings"
	"fmt"
)

/*
#comment
;comment
[tunnel name]
accept=addr
connect=addr
*/

/* lines starting with either of these characters are treated
as comments */
var (
	comment_char = []byte{';', '#'}
)

type Config map[string]*Section

type Section struct {
	Accept string
	Connect string
}

func IsComment(line string) (bool) {
	for _, c := range comment_char {
		if line[0] == c {
			return true
		}
	}

	return false
}

func IsSection(line string) (bool) {
	if line[0] == '[' {
		return true
	} else {
		return false
	}
}

func ParseSection(line string, r *bufio.Reader) (string, error) {
	if line[0] != '[' {
		return "", fmt.Errorf("Not a section")
	}

	for i := 0; i < len(line); i++ {
		err := r.UnreadByte()

		if err != nil {
			return "", err
		}
	}

	section_name := ""

	if strings.Contains(line, "]") == false {
		section_name = line
	}

	section_part, err := r.ReadString(']')

	if err == nil {
		return "", err
	}

	section_name = strings.Join([]string{section_name, section_part}, "")

	return section_name, nil
}

func Parse(cfg_file *os.File) (Config, error) {
	reader := bufio.NewReader(cfg_file)
	config_map := make(Config)

	section_name := "__global"
	config_map[section_name] = new(Section)

	for line, err := reader.ReadString('\n'); (err == nil && line != ""); {
		line = strings.Trim(line, " \t") /* Trim whitespace */

		if len(line) == 0 {
			continue
		}
		if IsComment(line) {
			continue /* skip */
		}
		if IsSection(line) {
			section_name, err = ParseSection(line, reader)
			/* TODO Append a new Section to the Config map */
		}
	}

	return nil, nil
}

//func (c Config) GetSection(section string) *Section {}
