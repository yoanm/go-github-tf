package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/rs/zerolog/log"

	"github.com/yoanm/github-tf/core"
)

func readWorkspace(rootPath, configDir, templateDir, yamlAnchorDir string) (config *core.Config, err error) {
	config = core.NewConfig()

	if _, err = os.Stat(rootPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("input directory '%s' doesn't exist", rootPath)
	}

	configureYamlAnchorDirectory(rootPath, yamlAnchorDir)

	decoderOpts := []yaml.DecodeOption{yaml.UseOrderedMap()}

	confErr := readConfigDirectory(config, filepath.Join(rootPath, configDir), decoderOpts)
	tplErr := readTemplateDirectory(config, filepath.Join(rootPath, templateDir), decoderOpts)

	if confErr != nil || tplErr != nil {
		e := confErr
		if confErr != nil && tplErr != nil {
			e = fmt.Errorf("%s\n%s", confErr, tplErr)
		} else if tplErr != nil {
			e = tplErr
		}
		return nil, fmt.Errorf("Error during workspace loading:\n%v", e)
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

	errList := map[string]error{}
	filenames, readErr := readDirectory(rootPath)
	visited := map[string]string{}

	if readErr == nil {
		for _, filename := range filenames {
			path := filepath.Join(rootPath, filename)
			if filename == "repos.yaml" || filename == "repos.yml" {
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
			} else if filename == "repos" {
				subVisited, loadErrList := readRepositoryDirectory(config, path, decoderOpts)
				if loadErrList != nil {
					for k, v := range loadErrList {
						errList[k] = v
					}
				}
				if subVisited != nil {
					for k, v := range subVisited {
						visited[k] = v
					}
				}
			} else {
				log.Debug().Msgf("%s is not a known file or directory => ignored", path)
			}
		}
	}

	uniqRepoList := map[string]string{}
	for fName, repoName := range visited {
		if firstFName, ok := uniqRepoList[repoName]; ok {
			errList[fName] = fmt.Errorf("file %s imports %s, but already imported by %s", fName, repoName, firstFName)
		} else {
			uniqRepoList[repoName] = fName
		}
	}

	if len(errList) > 0 || readErr != nil {
		msgList := []string{"Error during configs loading:"}
		if readErr != nil {
			msgList = append(msgList, fmt.Sprintf("\t - %s", readErr))
		} else {
			// sort file to always get a predictable output (for tests mostly)
			keys := make([]string, 0, len(errList))
			for k := range errList {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, file := range keys {
				msgList = append(msgList, fmt.Sprintf("\t - %s", errList[file]))
			}
		}

		return fmt.Errorf(strings.Join(msgList, "\n"))
	}

	return nil
}
func readRepositoryDirectory(config *core.Config, rootPath string, decoderOpts []yaml.DecodeOption) (visited map[string]string, errList map[string]error) {
	dirName := filepath.Base(rootPath)
	filenames, readErr := readDirectory(rootPath)
	if readErr != nil {
		return nil, map[string]error{dirName: readErr}
	}

	errList = map[string]error{}
	visited = map[string]string{}

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
func readTemplateDirectory(config *core.Config, rootPath string, decoderOpts []yaml.DecodeOption) error {
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
		msgList := []string{"Error during templates loading:"}
		if readErr != nil {
			msgList = append(msgList, fmt.Sprintf("\t - %s", readErr))
		} else {
			// sort file to always get a predictable output (for tests mostly)
			keys := make([]string, 0, len(errList))
			for k := range errList {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, file := range keys {
				decodeErr := errList[file]
				msgList = append(msgList, fmt.Sprintf("\t - %s", decodeErr))
			}
		}
		return fmt.Errorf(strings.Join(msgList, "\n"))
	}

	return nil
}

func loadTemplateFromFile(config *core.Config, filePath string, decoderOpts []yaml.DecodeOption) error {
	filename := filepath.Base(filePath)
	ext := filepath.Ext(filename)
	tplName := filename[:strings.LastIndex(filename, ext)]
	if strings.HasSuffix(tplName, ".repo") {
		tplName = strings.TrimSuffix(tplName, ".repo")
		tpl, err := core.LoadRepositoryTemplateFromFile(filePath, decoderOpts...)
		if err != nil {
			return err
		} else {
			config.Templates.Repos[tplName] = tpl
			log.Debug().Msgf("Loaded '%s' as repository template", filePath)
		}
	} else if strings.HasSuffix(tplName, ".branch-protection") {
		tplName = strings.TrimSuffix(tplName, ".branch-protection")
		tpl, err := core.LoadBranchProtectionTemplateFromFile(filePath, decoderOpts...)
		if err != nil {
			return err
		} else {
			config.Templates.BranchProtections[tplName] = tpl
			log.Debug().Msgf("Loaded '%s' as branch protection template", filePath)
		}
	} else if strings.HasSuffix(tplName, ".branch") {
		tplName = strings.TrimSuffix(tplName, ".branch")
		tpl, err := core.LoadBranchTemplateFromFile(filePath, decoderOpts...)
		if err != nil {
			return err
		} else {
			config.Templates.Branches[tplName] = tpl
			log.Debug().Msgf("Loaded '%s' as branch protection template", filePath)
		}
	} else {
		log.Debug().Msgf("%s is not a known template type => ignored", filePath)
	}

	return nil
}
