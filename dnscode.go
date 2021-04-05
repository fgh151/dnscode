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

	importCommand := flag.NewFlagSet("import", flag.ExitOnError)
	importTextPtr := importCommand.String("filename", "", "File name to save. If empty it will override "+FILENAME)
	importDirectiveImportPtr := importCommand.Bool("useImport", true, "Use import directive, default true")

	planCommand := flag.NewFlagSet("plan", flag.ExitOnError)
	planForceDeletePtr := planCommand.Bool("force", false, "Force delete")

	applyCommand := flag.NewFlagSet("apply", flag.ExitOnError)
	applyForceDeletePtr := applyCommand.Bool("force", false, "Force delete")

	if len(os.Args) < 2 {
		fmt.Println("list or count subcommand is required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "import":
		importCommand.Parse(os.Args[2:])

	case "plan":
		planCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if importCommand.Parsed() {

		filename := *importTextPtr
		if filename == "" {
			filename = FILENAME
		}

		fmt.Println("importing")
		commands.ImportDomains(FILENAME, filename, *importDirectiveImportPtr)
	}

	if planCommand.Parsed() {
		commands.Plan(FILENAME, planForceDeletePtr)
	}

	if applyCommand.Parsed() {
		commands.Apply(FILENAME, applyForceDeletePtr)
	}
}
