package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/rs/zerolog/log"

	"github.com/yoanm/github-tf/core"
)

var (
	errInputDirectoryDoesntExist = errors.New("input directory doesn't exist")
	errDuringWorkspaceLoading    = errors.New("error during configs loading")
	errDuringTemplateLoading     = errors.New("error during templates loading")
	errRepositoryAlreadyImported = errors.New("repository already imported")
)

func InputDirectoryDoesntExistError(path string) error {
	return fmt.Errorf("%w: %s", errInputDirectoryDoesntExist, path)
}

func ConfigDirectoryLoadingReadError(readErr error) error {
	return fmt.Errorf("%w\n\t - %w", errDuringWorkspaceLoading, readErr)
}

func ConfigDirectoryLoadingError(errList map[string]error) error {
	msgList := []string{}
	// sort file to always get a predictable output (for tests mostly)
	keys := make([]string, 0, len(errList))
	for k := range errList {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, file := range keys {
		msgList = append(msgList, fmt.Sprintf("\t - %s", errList[file].Error()))
	}

	return fmt.Errorf("%w\n%s", errDuringWorkspaceLoading, strings.Join(msgList, "\n"))
}

func AlreadyImportedRepositoryError(fName string, repoName string, firstFName string) error {
	return fmt.Errorf(
		"%w: repository %s imported by %s, but already imported by %s",
		errRepositoryAlreadyImported,
		repoName,
		fName,
		firstFName,
	)
}

func TemplateLoadingError(msgList []string) error {
	return fmt.Errorf("%w:\n\t - %s", errDuringTemplateLoading, strings.Join(msgList, "\n\t - "))
}

func readWorkspace(rootPath, configDir, templateDir, yamlAnchorDir string) (*core.Config, error) {
	var err error

	config := core.NewConfig()

	if _, err = os.Stat(rootPath); os.IsNotExist(err) {
		return nil, InputDirectoryDoesntExistError(rootPath)
	}

	configureYamlAnchorDirectory(rootPath, yamlAnchorDir)

	decoderOpts := []yaml.DecodeOption{yaml.UseOrderedMap()}

	confErr := readConfigDirectory(config, filepath.Join(rootPath, configDir), decoderOpts)
	tplErr := readTemplateDirectory(config, filepath.Join(rootPath, templateDir), decoderOpts)

	if confErr != nil || tplErr != nil {
		err2 := confErr
		if confErr != nil && tplErr != nil {
			err2 = fmt.Errorf("%w\n%w", confErr, tplErr)
		} else if tplErr != nil {
			err2 = tplErr
		}

		return nil, fmt.Errorf("error during workspace loading:\n%w", err2)
	}

	return config, nil
}

func configureYamlAnchorDirectory(path string, yamlAnchorDir string) {
	anchorDir := filepath.Join(path, yamlAnchorDir)
	fs, err := os.Stat(anchorDir)

	exists := !os.IsNotExist(err)
	isDir := exists && err == nil && fs.IsDir()

	if !exists || !isDir {
		return
	}

	core.YamlAnchorDirectory = &anchorDir
}

func readConfigDirectory(config *core.Config, rootPath string, decoderOpts []yaml.DecodeOption) error {
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		// Nothing to do
		return nil
	}

	var (
		filenames []string
		readErr   error
	)

	if filenames, readErr = readDirectory(rootPath); readErr != nil {
		return ConfigDirectoryLoadingReadError(readErr)
	}

	if err := loadConfigDirectoryFiles(config, rootPath, decoderOpts, filenames); err != nil {
		return err
	}

	return nil
}

func loadConfigDirectoryFiles(
	config *core.Config,
	rootPath string,
	decoderOpts []yaml.DecodeOption,
	filenames []string,
) error {
	visited := map[string]string{}
	errList := map[string]error{}

	for _, filename := range filenames {
		loadConfigDirectoryFile(config, filename, filepath.Join(rootPath, filename), decoderOpts, errList, visited)
	}

	uniqRepoList := map[string]string{}
	for fName, repoName := range visited {
		if firstFName, ok := uniqRepoList[repoName]; ok {
			errList[fName] = AlreadyImportedRepositoryError(fName, repoName, firstFName)
		} else {
			uniqRepoList[repoName] = fName
		}
	}

	if len(errList) > 0 {
		return ConfigDirectoryLoadingError(errList)
	}

	return nil
}

func loadConfigDirectoryFile(
	config *core.Config,
	filename string,
	path string,
	decoderOpts []yaml.DecodeOption,
	errList map[string]error,
	visited map[string]string,
) {
	switch {
	case filename == "repos.yaml" || filename == "repos.yml":
		loadReposConfigFile(config, filename, path, decoderOpts, errList, visited)
	case filename == "repos":
		loadReposConfigDirectory(config, path, decoderOpts, errList, visited)
	default:
		log.Debug().Msgf("%s is not a known file or directory => ignored", path)
	}
}

func loadReposConfigDirectory(
	config *core.Config,
	path string,
	decoderOpts []yaml.DecodeOption,
	errList map[string]error,
	visited map[string]string,
) {
	subVisited, loadErrList := readRepositoryDirectory(config, path, decoderOpts)
	for k, v := range loadErrList {
		errList[k] = v
	}

	for k, v := range subVisited {
		visited[k] = v
	}
}

func loadReposConfigFile(
	config *core.Config,
	filename string,
	path string,
	decoderOpts []yaml.DecodeOption,
	errList map[string]error,
	visited map[string]string,
) {
	repoConfigs, loadErr := core.LoadRepositoriesFromFile(path, decoderOpts...)
	if loadErr != nil {
		errList[filename] = loadErr
	} else {
		log.Debug().Msgf("Loaded '%s' as repositories config", path)
		for k, v := range repoConfigs {
			config.AppendRepo(v)
			visited[fmt.Sprintf("%s[%d]", path, k)] = *v.Name
		}
	}
}

func readRepositoryDirectory(
	config *core.Config,
	rootPath string,
	decoderOpts []yaml.DecodeOption,
) (map[string]string, map[string]error) {
	dirName := filepath.Base(rootPath)

	filenames, readErr := readDirectory(rootPath)
	if readErr != nil {
		return nil, map[string]error{dirName: readErr}
	}

	errList := map[string]error{}
	visited := map[string]string{}

	log.Debug().Msgf("Reading repository directory: %s", rootPath)

	for _, filename := range filenames {
		filePath := filepath.Join(rootPath, filename)

		ext := filepath.Ext(filename)
		if ext == ".yml" || ext == ".yaml" {
			repoConfig, loadErr := core.LoadRepositoryFromFile(filePath, decoderOpts...)
			if loadErr != nil {
				errList[filePath] = loadErr
			} else {
				log.Debug().Msgf("Loaded '%s' as repository config", filePath)
				config.AppendRepo(repoConfig)
				visited[filePath] = *repoConfig.Name
			}
		} else {
			log.Debug().Msgf("%s is not a YAML template => ignored", filePath)
		}
	}

	return visited, errList
}

func readTemplateDirectory(
	config *core.Config,
	rootPath string,
	decoderOpts []yaml.DecodeOption,
) error {
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		// Nothing to do
		return nil
	}

	errList := map[string]error{}

	dirContents, readErr := readDirectory(rootPath)
	if readErr == nil {
		log.Debug().Msgf("Reading template directory: %s", rootPath)

		for _, filename := range dirContents {
			path := filepath.Join(rootPath, filename)

			ext := filepath.Ext(filename)
			if ext == ".yml" || ext == ".yaml" {
				loadErr := loadTemplateFromFile(config, path, decoderOpts)
				if loadErr != nil {
					errList[path] = loadErr
				}
			} else {
				log.Debug().Msgf("%s is not a YAML template => ignored", path)
			}
		}
	}

	if len(errList) > 0 || readErr != nil {
		return TemplateLoadingError(generateTemplateLoadingErrorMessages(readErr, errList))
	}

	return nil
}

func generateTemplateLoadingErrorMessages(readErr error, errList map[string]error) []string {
	if readErr != nil {
		return []string{readErr.Error()}
	}

	msgList := []string{}
	// sort file to always get a predictable output (for tests mostly)
	keys := make([]string, 0, len(errList))
	for k := range errList {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, file := range keys {
		msgList = append(msgList, errList[file].Error())
	}

	return msgList
}

func loadTemplateFromFile(config *core.Config, filePath string, decoderOpts []yaml.DecodeOption) error {
	filename := filepath.Base(filePath)
	ext := filepath.Ext(filename)
	tplName := filename[:strings.LastIndex(filename, ext)]

	switch {
	case strings.HasSuffix(tplName, ".repo"):
		tplName = strings.TrimSuffix(tplName, ".repo")

		tpl, err := core.LoadRepositoryTemplateFromFile(filePath, decoderOpts...)
		if err != nil {
			//nolint:wrapcheck // Expecred to return error as is
			return err
		} else {
			config.Templates.Repos[tplName] = tpl
			log.Debug().Msgf("Loaded '%s' as repository template", filePath)
		}
	case strings.HasSuffix(tplName, ".branch-protection"):
		tplName = strings.TrimSuffix(tplName, ".branch-protection")

		tpl, err := core.LoadBranchProtectionTemplateFromFile(filePath, decoderOpts...)
		if err != nil {
			//nolint:wrapcheck // Expecred to return error as is
			return err
		} else {
			config.Templates.BranchProtections[tplName] = tpl
			log.Debug().Msgf("Loaded '%s' as branch protection template", filePath)
		}
	case strings.HasSuffix(tplName, ".branch"):
		tplName = strings.TrimSuffix(tplName, ".branch")

		tpl, err := core.LoadBranchTemplateFromFile(filePath, decoderOpts...)
		if err != nil {
			//nolint:wrapcheck // Expecred to return error as is
			return err
		} else {
			config.Templates.Branches[tplName] = tpl
			log.Debug().Msgf("Loaded '%s' as branch protection template", filePath)
		}
	default:
		log.Debug().Msgf("%s is not a known template type => ignored", filePath)
	}

	return nil
}
