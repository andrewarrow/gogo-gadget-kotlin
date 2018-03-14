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
		{Name: "add_item_to_array", ShortName: "a", Usage: "add item to array", Action: AddAction},
		{Name: "new_table", ShortName: "t", Usage: "add new table", Action: TableAction},
		{Name: "vim", ShortName: "v", Usage: "vim query", Action: VimAction},
	}

	app.Run(os.Args)
}

func WriteLine() {
}
func AddAction(c *cli.Context) {
	//role := c.Args().Get(0)
	fileList := AllSrcFiles()
	for _, file := range fileList {
		b, err := ioutil.ReadFile(file)
		if err != nil {
			continue
		}
		buffer := []string{}
		for _, line := range strings.Split(string(b), "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.Contains(trimmed, "@RolesAllowed(\"") {
				tokens := strings.Split(trimmed, "RolesAllowed")
				more := tokens[1]
				evenMore := strings.Split(more[1:len(more)-1], ",")
				roles := map[string]int{}
				fixed := []string{}
				for _, em := range evenMore {
					key := strings.TrimSpace(em)
					name := key[1 : len(key)-1]
					fixed = append(fixed, "Role."+strings.ToUpper(name))
					roles[name] = 1
				}
				if true {
					//@RolesAllowed(Role.TECHNICIAN, Role.MANAGER, Role.LAUNCHER, Role.SUPPORT)
					//  @RolesAllowed("superuser","admin","support","manager","test")
					newline := fmt.Sprintf("  @RolesAllowed(%s)", strings.Join(fixed, ","))
					fmt.Println(newline)
					buffer = append(buffer, newline)
				} else {
					buffer = append(buffer, line)
				}
			} else {
				buffer = append(buffer, line)
			}
		}
		newcon := strings.Join(buffer, "\n")
		d1 := []byte(newcon)
		ioutil.WriteFile(file, d1, 0644)
	}
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
	imports := make(map[string]string)
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
	blackList := make(map[string]bool)
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
			replaced := strings.Replace(trimmed, "@", " ", -1)
			replaced = strings.Replace(replaced, "(", " ", -1)
			replaced = strings.Replace(replaced, ".", " ", -1)
			replaced = strings.Replace(replaced, ":", " ", -1)
			replaced = strings.Replace(replaced, "?", " ", -1)
			replaced = strings.Replace(replaced, ",", " ", -1)
			replaced = strings.Replace(replaced, "<", " ", -1)
			replaced = strings.Replace(replaced, ">", " ", -1)
			tokens := strings.Split(replaced, " ")
			t0 := strings.TrimSpace(tokens[0])
			t1 := ""
			if len(tokens) > 1 {
				t1 = strings.TrimSpace(tokens[1])
			}
			if t0 == "class" {
				blackList[t1] = true
				continue
			}
			if t1 == "class" {
				blackList[strings.TrimSpace(tokens[2])] = true
				continue
			}

			for last, v := range lasts {
				for _, t := range tokens {
					if t == last {
						for _, l := range v {
							imports[l] = last
						}
					}
				}
			}
		}
	}
	var keys []string
	for k, _ := range imports {
		tokens := strings.Split(k, ".")
		last := tokens[len(tokens)-1]
		if blackList[last] == false {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	buffer := []string{}
	b, _ = ioutil.ReadFile(path)
	for _, line := range strings.Split(string(b), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "package ") {
			buffer = append(buffer, trimmed)
		}
		if strings.Contains(trimmed, "import ") {
			break
		}
	}
	buffer = append(buffer, "")
	for _, k := range keys {
		buffer = append(buffer, k)
	}
	for _, line := range strings.Split(string(b), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "package ") {
			continue
		}
		if strings.Contains(trimmed, "import ") {
			continue
		}
		buffer = append(buffer, trimmed)
	}
	file := strings.Join(buffer, "\n")

	d1 := []byte(file)
	ioutil.WriteFile(path, d1, 0644)

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
