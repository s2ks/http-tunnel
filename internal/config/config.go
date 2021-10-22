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
	//"log"
)

type Config map[string]Section
type Section map[string][]string

type Line struct {
	body 	string
	no 	int
}

func ParseKeyVal(line *Line, s Section) error {
	part := strings.SplitN(line.body, "=", 2)

	if len(part) < 2 {
		return fmt.Errorf("Malformed expression on line: %d" +
		"\n\t--> \"%s\"", line.no, line.body)
	}

	key := part[0]
	val := part[1]

	key = strings.Trim(key, " \t\r\n")
	val = strings.Trim(val, " \t\r\n")

	s[key] = append(s[key], val)

	return nil
}

/* A section has the form of [section name] */
func ParseSection(line *Line) (string, error) {
	section_name := ""

	/* Assume the caller has made sure that line[0] == '[' */
	for i := 1; i < len(line.body); i++ {
		if line.body[i] == ']' {
			return section_name, nil
		}

		section_name += string(line.body[i])
	}

	return section_name, fmt.Errorf("Missing closing bracket ']' " +
	"in section on line: %d\n\t--> \"%s\"", line.no, line.body)
}

func Parse(cfg_file *os.File) (Config, error) {
	reader := bufio.NewReader(cfg_file)
	config_map := make(Config)

	section_name := "__global"
	config_map[section_name] = make(Section)

	var err error
	var line Line

	line.no = 0

	for {
		line.body, err = reader.ReadString('\n')
		line.no++

		if line.body == "" && err != nil {
			break
		}

		line.body = strings.Trim(line.body, " \t\r\n") /* Trim whitespace */

		if len(line.body) == 0 {
			continue
		}

		/* Go switch statements don't 'fall thorugh' */
		switch line.body[0] {
			/* section */
			case '[':
				section_name, err = ParseSection(&line)

				if err != nil {
					return nil, err
				}

				config_map[section_name] = make(Section)
			/* ignore comments */
			case ';', '#', '\n':
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
