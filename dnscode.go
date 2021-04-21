package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"openitstudio.ru/dnscode/commands"
	localHttp "openitstudio.ru/dnscode/http"
	"os"
)

const FILENAME = "dnscode.json"

func main() {

	godotenv.Load()

	importCommand := flag.NewFlagSet("import", flag.ExitOnError)
	importTextPtr := importCommand.String("filename", "", "File name to save. If empty it will override "+FILENAME)
	importDirectiveImportPtr := importCommand.Bool("useImport", true, "Use import directive, default true")
	importProxyPtr := importCommand.String("proxy", "", "Proxy addr")

	planCommand := flag.NewFlagSet("plan", flag.ExitOnError)
	planForceDeletePtr := planCommand.Bool("force", false, "Force delete")
	planProxyPtr := importCommand.String("proxy", "", "Proxy addr")

	applyCommand := flag.NewFlagSet("apply", flag.ExitOnError)
	applyForceDeletePtr := applyCommand.Bool("force", false, "Force delete")
	applyProxyPtr := importCommand.String("proxy", "", "Proxy addr")

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
	case "apply":
		applyCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if importCommand.Parsed() {

		localHttp.SetProxy(*importProxyPtr)

		filename := *importTextPtr
		if filename == "" {
			filename = FILENAME
		}

		fmt.Println("importing")
		commands.ImportDomains(FILENAME, filename, *importDirectiveImportPtr)
	}

	if planCommand.Parsed() {
		localHttp.SetProxy(*planProxyPtr)
		commands.Plan(FILENAME, planForceDeletePtr)
	}

	if applyCommand.Parsed() {
		//localHttp.SetProxy(*applyProxyPtr)

		commands.Plan(FILENAME, planForceDeletePtr)
		if confirm("Apply?", 3) {
			commands.Apply(FILENAME, applyForceDeletePtr)
		}
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
