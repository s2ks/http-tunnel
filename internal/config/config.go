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
        "regexp"
        //"log"
)

var (
        whitespace      = regexp.MustCompile(`^[[:space:]]+|[[:space:]]+$`)

        section         = regexp.MustCompile(`^\[(.+)\]$`)
        comment         = regexp.MustCompile(`^;.*|^#.*`)
        keyval          = regexp.MustCompile(`(^[^[:space:]]+)=(.*)`)
        blankline       = regexp.MustCompile(`^[[:space:]]*$`)
)

type Config map[string]Section
type Section map[string][]string

type Line struct {
        body    string
        no      int
}

func ParseKeyVal(line *Line, s Section) error {
        part := strings.SplitN(line.body, "=", 2)

        if len(part) < 2 {
                return fmt.Errorf("Malformed expression on line: %d" +
                "\n\t--> \"%s\"", line.no, line.body)
        }

        key := part[0]
        val := part[1]

        key = whitespace.ReplaceAllString(key, "")
        val = whitespace.ReplaceAllString(val, "")

        s[key] = append(s[key], val)

        return nil
}

/* A section has the form of [section name] */
func ParseSection(line *Line) (string, error) {
        section_name := section.ReplaceAllString(line.body, "$1")

        if section_name == "" {
                return "", fmt.Errorf("Unable to parse section name on line %d: \"%s\"",
                                line.no, line.body)
        } else {
                return section_name, nil
        }
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

                /* Trim whitespace at start and end of line */
                line.body = whitespace.ReplaceAllString(line.body, "")

                switch true {
                        case blankline.MatchString(line.body):
                        case comment.MatchString(line.body):
                        case section.MatchString(line.body):
                                section_name, err = ParseSection(&line)

                                if err != nil {
                                        return nil, err
                                }
                                config_map[section_name] = make(Section)
                        case keyval.MatchString(line.body):
                                err = ParseKeyVal(&line, config_map[section_name])

                                if err != nil {
                                        return nil, err
                                }
                        default:
                                return nil, fmt.Errorf("Malformed expression on line " +
                                                "%d: \"%s\"", line.no, line.body)
                }
        }

        return config_map, nil
}
