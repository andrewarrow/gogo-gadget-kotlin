package main

import "fmt"
import "os"
import "github.com/codegangsta/cli"
import "time"
import "math/rand"
import "path/filepath"

func main() {

	rand.Seed(time.Now().UnixNano())
	app := cli.NewApp()
	app.Name = "gogo-gadget-kotlin"
	app.Usage = "change kotlin code from command line and avoid using an IDE "
	app.Version = "0.1.1"
	app.Commands = []cli.Command{
		{Name: "import", ShortName: "i",
			Usage: "fix all your imports", Action: ImportAction},
		{Name: "vim", ShortName: "v", Usage: "open in vim",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "query", Value: "Resource", Usage: "match"},
			},
			Action: VimAction},
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
	fileList := AllSrcFiles()
	for _, file := range fileList {
		fmt.Println(file)
	}
}

func VimAction(c *cli.Context) {
	query := c.String("query")
	fmt.Println(query)
	fileList := AllSrcFiles()
	for _, file := range fileList {
		fmt.Println(file)
	}

}
