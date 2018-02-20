package main

//import "fmt"
import "os"
import "github.com/codegangsta/cli"
import "time"
import "math/rand"

func main() {
	rand.Seed(time.Now().UnixNano())
	app := cli.NewApp()
	app.Name = "gogo-gadget-kotlin"
	app.Usage = "change kotlin code from command line and avoid using an IDE "
	app.Version = "0.1.1"
	app.Commands = []cli.Command{
		{Name: "import statements", ShortName: "i",
			Usage: "fix all your imports", Action: ImportAction},
		{Name: "placeholder", ShortName: "p",
			Usage: "placeholder", Action: PlaceAction},
	}

	app.Run(os.Args)
}

func ImportAction(c *cli.Context) {
}
func PlaceAction(c *cli.Context) {

}
