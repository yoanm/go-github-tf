package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	flag "github.com/spf13/pflag"
)

// Build time variables.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var (
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

func init() {
	// init command line flags
	flag.StringVarP(&workspacePathFlag, "workspace", "w", defaultWorkspacePathFlag, `Workspace directory`)
	flag.StringVarP(&configDirFlag, "config", "c", defaultConfigDirFlag, `Config directory`)
	flag.StringVarP(&templateDirFlag, "templates", "t", defaultTemplateDirFlag, `Template directory`)
	flag.StringVar(&yamlAnchorDirFlag, "yaml-anchors", defaultYamlAnchorDirFlag, `YAML anchors directory`)

	// flag.BoolVar(&printImportsFlag, "print-imports", false, "Read terraform files and print related terraform import commands")
	// flag.StringSliceVar(&skipImportListFlag, "skip", defaultSkipImportListFlag, "Skip provided import from the list")

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
	logLevel := zerolog.WarnLevel
	if quietFlag {
		logLevel = zerolog.Disabled
	} else if verboseFlag == 1 {
		logLevel = zerolog.InfoLevel
	} else if verboseFlag == 2 {
		logLevel = zerolog.DebugLevel
	} else if verboseFlag > 2 {
		logLevel = zerolog.TraceLevel
	}

	setupLogOutput(logLevel, disableAnsiFlag)

	terraformDir := "terraform"
	log.Debug().Msgf("Workspace: %s", workspacePathFlag)
	log.Debug().Msgf("Config directory: %s", configDirFlag)
	log.Debug().Msgf("Template directory: %s", templateDirFlag)
	log.Debug().Msgf("YAML anchor directory: %s", yamlAnchorDirFlag)
	exitCode := 0
	if helpFlag {
		flag.PrintDefaults()
	} else if versionFlag {
		fmt.Printf("github-tf version: %s (commit %s from %s)\n", version, commit, date)

		/*}  else if printImportsFlag {
		exitCode = printTerraformImports(filepath.Join(workspacePathFlag, terraformDir), &skipImportListFlag)
		*/
	} else {
		exitCode = loadYamlAndWriteTerraform(workspacePathFlag, configDirFlag, templateDirFlag, terraformDir, yamlAnchorDirFlag)
	}

	return exitCode
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
