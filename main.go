package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	flag "github.com/spf13/pflag"
)

// Build time variables.
//
//nolint:gochecknoglobals // Expected to be global
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	// Base flags.
	workspacePathFlag        string
	defaultWorkspacePathFlag = "."
	configDirFlag            string
	defaultConfigDirFlag     = "config"
	templateDirFlag          string
	defaultTemplateDirFlag   = "templates"
	yamlAnchorDirFlag        string
	defaultYamlAnchorDirFlag = "yaml-anchors"

	// printImportsFlag          bool
	// skipImportListFlag        []string
	// defaultSkipImportListFlag []string = nil.

	// Logging flags.
	verboseFlag     int
	quietFlag       bool
	disableAnsiFlag bool

	// Miscellaneous flags.
	helpFlag    bool
	versionFlag bool
)

//nolint:gochecknoinits // Normal CLI flag declaration
func init() {
	// init command line flags
	flag.StringVarP(&workspacePathFlag, "workspace", "w", defaultWorkspacePathFlag, `Workspace directory`)
	flag.StringVarP(&configDirFlag, "config", "c", defaultConfigDirFlag, `Config directory`)
	flag.StringVarP(&templateDirFlag, "templates", "t", defaultTemplateDirFlag, `Template directory`)
	flag.StringVar(&yamlAnchorDirFlag, "yaml-anchors", defaultYamlAnchorDirFlag, `YAML anchors directory`)

	// flag.BoolVar(
	// 	&printImportsFlag,
	// 	"print-imports",
	// 	false,
	// 	"Read terraform files and print related terraform import commands"
	// )
	// flag.StringSliceVar(
	// 	&skipImportListFlag,
	// 	"skip",
	// 	defaultSkipImportListFlag,
	// 	"Skip provided import from the list"
	// )

	flag.BoolVarP(&quietFlag, "quiet", "q", false, "Disable output")
	flag.CountVarP(&verboseFlag, "verbose", "v", "Enable verbose output. -v for Info, -vv for Debug and -vvv for Trace")
	flag.BoolVar(&disableAnsiFlag, "no-ansi", false, "Disable ANSI output")

	flag.BoolVarP(&versionFlag, "version", "V", false, `Print current version`)
	flag.BoolVarP(&helpFlag, "help", "h", false, "Display this help")
}

func main() {
	os.Exit(run())
}

func run() int {
	parseFlags()
	setupLogOutput(computeLogLevel(), disableAnsiFlag)

	terraformDir := "terraform"

	log.Debug().Msgf("Workspace: %s", workspacePathFlag)
	log.Debug().Msgf("Config directory: %s", configDirFlag)
	log.Debug().Msgf("Template directory: %s", templateDirFlag)
	log.Debug().Msgf("YAML anchor directory: %s", yamlAnchorDirFlag)

	exitCode := 0

	switch {
	case helpFlag:
		flag.PrintDefaults()
	case versionFlag:
		//nolint:forbidigo // Expected output
		fmt.Printf("github-tf version: %s (commit %s from %s)\n", version, commit, date)

		/*}  else if printImportsFlag {
		exitCode = printTerraformImports(filepath.Join(workspacePathFlag, terraformDir), &skipImportListFlag)
		*/
	default:
		exitCode = loadYamlAndWriteTerraform(
			workspacePathFlag,
			configDirFlag,
			templateDirFlag,
			terraformDir,
			yamlAnchorDirFlag,
		)
	}

	return exitCode
}

func computeLogLevel() zerolog.Level {
	switch {
	case quietFlag:
		return zerolog.Disabled
	case verboseFlag == 1:
		return zerolog.InfoLevel
	case verboseFlag == 2: //nolint:gomnd // Doesn't make sense here to wrap 2
		return zerolog.DebugLevel
	case verboseFlag > 2: //nolint:gomnd // Doesn't make sense here to wrap 2
		return zerolog.TraceLevel
	}

	return zerolog.WarnLevel
}

func parseFlags() {
	// Reset input (useful only for tests execution)
	workspacePathFlag = defaultWorkspacePathFlag
	configDirFlag = defaultConfigDirFlag
	templateDirFlag = defaultTemplateDirFlag
	yamlAnchorDirFlag = defaultYamlAnchorDirFlag
	// printImportsFlag = false
	// skipImportListFlag = defaultSkipImportListFlag
	helpFlag = false
	verboseFlag = 0
	quietFlag = false
	disableAnsiFlag = false
	versionFlag = false

	flag.Parse()
}
