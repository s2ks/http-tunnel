/* Super simple INI-style parser */

/*
#comment
;comment
[tunnel name]
accept=addr
connect=addr
*/

package config

import (
	"os"
	"bufio"
	"strings"
	"fmt"
	"log"
)

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

/* A section has the form of [section name] */
func ParseSection(line string, r *bufio.Reader) (string, error) {
	if line[0] != '[' {
		return "", fmt.Errorf("Not a section")
	}

	section_name := ""

	for i := 1; i < len(line); i++ {
		if line[i] == ']' {
			log.Print(section_name)
			return section_name, nil
		}

		section_name += string(line[i])
	}

	section_part, err := r.ReadString(']')
	section_part = strings.TrimSuffix(section_part, "]")

	if err != nil {
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

	var err error

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			break
		}

		line = strings.Trim(line, " \t") /* Trim whitespace */
		log.Print(line)

		switch line[0] {
			/* section */
			case '[':
				section_name, err = ParseSection(line, reader)

				if err != nil {
					return nil, err
				}

				config_map[section_name] = new(Section)
				break
			/* comment */
			case ';':
			case '#':
				break
			default:
				/* Parse key=val */

		}
	}

	log.Print(config_map)

	return config_map, err
}

//func (c Config) GetSection(section string) *Section {}
