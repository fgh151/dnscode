package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"openitstudio.ru/dnscode/commands"
	localHttp "openitstudio.ru/dnscode/http"
	"os"
	"strings"
)

const FILENAME = "dnscode.json"

func main() {

	godotenv.Load()

	var forcePtr = flag.Bool("force", true, "Force delete")
	var proxyPtr = flag.String("proxy", "", "Proxy server")
	var importPtr = flag.Bool("useImport", true, "Use import directive, default true")
	var importTextPtr = flag.String("filename", "", "File name to save. If empty it will override "+FILENAME)

	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("list or count subcommand is required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "import":
		filename := *importTextPtr
		if filename == "" {
			filename = FILENAME
		}
		commands.ImportDomains(FILENAME, filename, *importPtr)

	case "plan":
		localHttp.SetProxy(*proxyPtr)
		commands.Plan(FILENAME, forcePtr)
	case "apply":
		localHttp.SetProxy(*proxyPtr)
		commands.Plan(FILENAME, forcePtr)
		if confirm("Apply?", 3) {
			commands.Apply(FILENAME, forcePtr)
		}
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func confirm(s string, tries int) bool {
	r := bufio.NewReader(os.Stdin)

	for ; tries > 0; tries-- {
		fmt.Printf("%s [y/n]: ", s)

		res, err := r.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		// Empty input (i.e. "\n")
		if len(res) < 2 {
			continue
		}

		return strings.ToLower(strings.TrimSpace(res))[0] == 'y'
	}

	return false
}
