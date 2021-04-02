package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"openitstudio.ru/dnscode/commands"
	"os"
)

const FILENAME = "dnscode.json"

func main() {

	godotenv.Load()

	if len(os.Args) < 2 {
		fmt.Println("list or count subcommand is required")
		os.Exit(1)
	}

	importCommand := flag.NewFlagSet("import", flag.ExitOnError)

	importTextPtr := importCommand.String("filename", "", "File name to save. If empty it will override "+FILENAME)
	importDirectiveImportPtr := importCommand.Bool("useImport", true, "Use import directive, default true")

	switch os.Args[1] {
	case "import":
		importCommand.Parse(os.Args[2:])

	//case "count":
	//	countCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if importCommand.Parsed() {

		filename := *importTextPtr
		if filename == "" {
			filename = FILENAME
		}

		fmt.Println(filename)
		fmt.Println("importing")
		commands.ImportDomains(FILENAME, filename, *importDirectiveImportPtr)
	}
}
