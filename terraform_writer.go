package main

import (
	"path"

	"github.com/rs/zerolog/log"

	"github.com/yoanm/github-tf/core"
)

const (
	noErrorExitCode                     = 0
	readWorkspaceErrorExitCode          = 1
	computeConfigErrorExitCode          = 2
	generateTerraformFilesErrorExitCode = 3
	writeTerraformFilesErrorExitCode    = 4
)

func loadYamlAndWriteTerraform(workspacePath, configDir, templateDir, terraformDir, yamlAnchorDir string) int {
	var err error

	var rawConfig *core.Config

	if rawConfig, err = readWorkspace(workspacePath, configDir, templateDir, yamlAnchorDir); err != nil {
		log.Error().Msgf("%s", err)

		return readWorkspaceErrorExitCode
	}

	core.ConfigTrace("Decoded config", rawConfig)

	log.Info().Msgf(
		"Found: %d repos / %d repo templates / %d branch templates / %d branch protection templates",
		len(rawConfig.Repos),
		len(rawConfig.Templates.Repos),
		len(rawConfig.Templates.Branches),
		len(rawConfig.Templates.BranchProtections),
	)

	var config *core.Config

	if config, err = core.ComputeConfig(rawConfig); err != nil {
		log.Error().Msgf("%s", err)

		return computeConfigErrorExitCode
	}

	core.ConfigTrace("Computed config", config)

	files, err := core.GenerateHclRepoFiles(config.Repos)
	if err != nil {
		log.Error().Msgf("%s", err)

		return generateTerraformFilesErrorExitCode
	}

	if err = core.WriteTerraformFiles(path.Join(workspacePath, terraformDir), files); err != nil {
		log.Error().Msgf("%s", err)

		return writeTerraformFilesErrorExitCode
	}

	return noErrorExitCode
}
