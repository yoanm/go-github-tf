package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	flag "github.com/spf13/pflag"
)

// @TODO configTemplate: tpl => configTemplates: [tpl1, tpl2]
// @TODO
//  - Write test for core/terraform_loader.go
//  - Use the command to create repository for tfsig
//  - Push tfsig to github and check coverage miss from codecov.io
//  - Use the command to create repository for gh2tf
//  - Push gh2tf to github and check coverage miss from codecov.io
//  - Use the command to create repository for github-tf
//  - Push github-tf to github and check coverage miss from codecov.io
//  - then starts pushing others lib to github
// @TODO PUSH ON GITHUB + CREATE RELEASE + START USING IT
// LATER
// @TODO github action
// See https://docs.github.com/en/actions/security-guides/encrypted-secrets#storing-large-secrets
// See https://notes.nishkal.in/snippets/github-actions.html
// See https://itsfoss.com/gpg-encrypt-files-basic/
//  - Decrypt terraform.tfstate with a GPG key stored as secret (or symetric encryption ??) + password stored as secret
//  - terraform init
//  - terraform plan
//  - !! prevent updates on current repository
//  - For PR only: output as PR comment
//  - For manual trigger only *and only if my user (check if doable to use an id instead of a name)*:
//    - terraform apply
//    - encrypt terraform.tfstate to terraform.tfstate.gpg (exclude terraform.tfstate + terraform.tfstate* files and .terraform directory with .gitignore)
//    - git add terraform.tfstate.gpg + git commit -m "State update following XXXX commit" + git push
// @TODO re-read github_repository + github_branch_protection doc and manage missing field
// @TODO enhance json-schema (enum, min/max, etc)
// LATER LATER
// @TODO manage github_user_ssh_key and github_user_gpg_key
// @TODO manage github_repository_file (to add required common file on every repo)
// @TODO manage github action secrets
// LATER LATER MAYBE
// @TODO manage organization (user role, team + membership)

// Build time variables
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var (
	// Base flags
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
	// defaultSkipImportListFlag []string = nil

	// Logging flags
	verboseFlag     int
	quietFlag       bool
	disableAnsiFlag bool

	// Miscellaneous flags
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
