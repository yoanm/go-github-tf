package main

import (
	"path"

	"github.com/goccy/go-yaml"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/yoanm/github-tf/core"
)

func loadYamlAndWriteTerraform(workspacePath, configDir, templateDir, terraformDir, yamlAnchorDir string) int {
	var err error

	var rawConfig *core.Config
	if rawConfig, err = readWorkspace(workspacePath, configDir, templateDir, yamlAnchorDir); err != nil {
		log.Error().Msgf("%s", err)
		return 1
	}

	if zerolog.GlobalLevel() == zerolog.TraceLevel {
		encoded, encodeError := yaml.Marshal(rawConfig)
		if encodeError != nil {
			log.Trace().Msgf("Decoded config: Error %s", encodeError)
		} else {
			log.Trace().Msgf("Decoded config: \n%s", string(encoded))
		}
	}

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
		return 2
	}

	if zerolog.GlobalLevel() == zerolog.TraceLevel {
		encoded, encodeError := yaml.Marshal(&config)
		if encodeError != nil {
			log.Trace().Msgf("Computed config: Error %s", encodeError)
		} else {
			log.Trace().Msgf("Computed config: \n%s", string(encoded))
		}
	}

	files, err := core.GenerateHclRepoFiles(config.Repos)
	if err != nil {
		log.Error().Msgf("%s", err)
		return 3
	}

	if err = core.WriteTerraformFiles(path.Join(workspacePath, terraformDir), files); err != nil {
		log.Error().Msgf("%s", err)
		return 4
	}

	return 0
}
