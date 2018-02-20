package main

import "fmt"
import "os"
import "github.com/codegangsta/cli"
import "time"
import "strings"
import "io/ioutil"
import "math/rand"
import "sort"
import "path/filepath"
import "os/exec"

func main() {

	rand.Seed(time.Now().UnixNano())
	app := cli.NewApp()
	app.Name = "gogo-gadget-kotlin"
	app.Usage = "change kotlin code from command line and avoid using an IDE "
	app.Version = "0.1.1"
	app.Commands = []cli.Command{
		{Name: "import", ShortName: "i", Usage: "import query", Action: ImportAction},
		{Name: "vim", ShortName: "v", Usage: "vim query", Action: VimAction},
	}

	app.Run(os.Args)
}

func AllSrcFiles() []string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(dir)

	fileList := []string{}
	err = filepath.Walk(dir+"/src", func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	return fileList
}

func ImportAction(c *cli.Context) {
	query := c.Args().Get(0)
	fileList := AllSrcFiles()
	hash := make(map[string]bool)
	for _, file := range fileList {
		b, err := ioutil.ReadFile(file)
		if err == nil {
			for _, line := range strings.Split(string(b), "\n") {
				trimmed := strings.TrimSpace(line)
				if strings.Contains(trimmed, "import ") {
					if strings.HasSuffix(trimmed, ";") {
						trimmed = trimmed[0 : len(trimmed)-1]
					}
					hash[trimmed] = true
				}
			}
		}
	}
	lasts := make(map[string][]string)
	imports := make(map[string]bool)
	for k, _ := range hash {
		tokens := strings.Split(k, ".")
		last := tokens[len(tokens)-1]
		if last != "*" {
			lasts[last] = append(lasts[last], k)
		}
	}

	path := FindJustOne(strings.ToLower(query))
	if path == "" {
		return
	}
	b, err := ioutil.ReadFile(path)
	if err == nil {
		for _, line := range strings.Split(string(b), "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.Contains(trimmed, "import ") {
				continue
			}
			if strings.Contains(trimmed, "package ") {
				continue
			}
			replaced := strings.Replace(line, "@", " ", -1)
			replaced = strings.Replace(replaced, "(", " ", -1)
			replaced = strings.Replace(replaced, ".", " ", -1)
			replaced = strings.Replace(replaced, ":", " ", -1)
			replaced = strings.Replace(replaced, "?", " ", -1)
			replaced = strings.Replace(replaced, ",", " ", -1)
			replaced = strings.Replace(replaced, "<", " ", -1)
			replaced = strings.Replace(replaced, ">", " ", -1)
			for last, v := range lasts {
				tokens := strings.Split(replaced, " ")
				trimmedToken0 := strings.TrimSpace(tokens[0])
				trimmedToken1 := ""
				if len(tokens) > 1 {
					trimmedToken1 = strings.TrimSpace(tokens[1])
				}
				if trimmedToken0 != "class" && trimmedToken1 != "class" {
					for _, t := range tokens {
						if t == last {
							for _, l := range v {
								imports[l] = true
							}
						}
					}
				}
			}
		}
	}
	var keys []string
	for k, _ := range imports {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Println(k)
	}

}

func FindJustOne(query string) string {
	fileList := AllSrcFiles()
	fileListMatch := []string{}
	for _, file := range fileList {
		lower := strings.ToLower(file)
		if strings.Contains(lower, query) {
			fileListMatch = append(fileListMatch, file)
		}
	}
	if len(fileListMatch) == 1 {
		return fileListMatch[0]
	} else {
		for _, file := range fileListMatch {
			fmt.Println(file)
		}
	}
	return ""
}

func VimAction(c *cli.Context) {
	query := c.Args().Get(0)
	path := FindJustOne(strings.ToLower(query))
	if path == "" {
		return
	}
	cmd := exec.Command("vim", path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

}
