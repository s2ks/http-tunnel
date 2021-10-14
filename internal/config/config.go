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

type Config map[string]*Section

type Section struct {
	Accept string
	Connect string
}

type Line struct {
	body 	string
	no 	int
}

func ParseKeyVal(line *Line, s *Section) {
	key := ""
	val := ""

	/* key=val */

	/* ... TODO */
}

/* A section has the form of [section name] */
func ParseSection(line *Line) (string, error) {
	section_name := ""

	/* Assume the caller has made sure that line[0] == '[' */
	for i := 1; i < len(line.body); i++ {
		if line[i] == ']' {
			return section_name, nil
		}

		section_name += string(line.body[i])
	}

	return section_name, fmt.Errorf("Missing closing bracket ']' in section on line:
	\n\t%d: %s", line.no, line.body)
}

func Parse(cfg_file *os.File) (Config, error) {
	reader := bufio.NewReader(cfg_file)
	config_map := make(Config)

	section_name := "__global"
	config_map[section_name] = new(Section)

	var err error
	var line Line

	line.no = 0

	for {
		line.body, err := reader.ReadString('\n')
		line.no++

		if line.body == "" && err != nil {
			break
		}

		line.body = strings.Trim(line.body, " \t\r") /* Trim whitespace */

		switch line.body[0] {
			/* section */
			case '[':
				section_name, err = ParseSection(&line, reader)

				if err != nil {
					return nil, err
				}

				config_map[section_name] = new(Section)
				break
			/* comment */
			case ';':
			case '#':
				break /* ignore */
			default:
				/* Parse key=val */
				err = ParseKeyVal(&line, config_map[section_name])

				if err != nil {
					return nil, err
				}
		}
	}

	return config_map, nil
}

//func (c Config) GetSection(section string) *Section {}
