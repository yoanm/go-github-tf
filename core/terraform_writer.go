package core

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/rs/zerolog/log"

	"github.com/yoanm/go-gh2tf"
	"github.com/yoanm/go-tfsig"
)

/** Public **/

var errDuringWriteTerraformFiles = errors.New("error while writing terraform files")

func WriteTerraformFileError(errList []error) error {
	msgList := []string{}
	for _, err := range errList {
		msgList = append(msgList, fmt.Sprintf("\t - %s", err.Error()))
	}

	return fmt.Errorf("%w\n%s", errDuringWriteTerraformFiles, strings.Join(msgList, "\n"))
}

func GenerateHclRepoFiles(configList []*GhRepoConfig) (map[string]*hclwrite.File, error) {
	valueGenerator := gh2tf.NewValueGenerator()
	waitGroup := &sync.WaitGroup{}
	collector := make(fileCollector, len(configList))
	errCollector := make(errorCollector, len(configList))

	var errList []error

	for k, repoConfig := range configList {
		waitGroup.Add(1)

		if repoConfig.Name == nil {
			errList = append(errList, fmt.Errorf("config #%d: repository name is mandatory", k))
		} else {
			repoTfId := tfsig.ToTerraformIdentifier(*repoConfig.Name)
			go generateHclRepoFileAsync(repoConfig, valueGenerator, repoTfId, collector, waitGroup)
		}
	}

	waitGroup.Wait()
	close(collector)
	close(errCollector)

	if len(errCollector) > 0 || len(errList) > 0 {
		return nil, createFileGenerationError(errCollector, errList)
	}

	list := map[string]*hclwrite.File{}

	for fs := range collector {
		list[fs.name] = fs.file
	}

	return list, nil
}

func createFileGenerationError(errCollector errorCollector, errList []error) error {
	msgList := []string{"error while generating files:"}

	if len(errCollector) > 0 {
		subErrList, keys := internalSortFileErrors(errCollector)

		for _, file := range keys {
			generateErr := subErrList[file]
			msgList = append(msgList, fmt.Sprintf("\t - %s", generateErr))
		}
	}

	if len(errList) > 0 {
		for _, err := range errList {
			msgList = append(msgList, fmt.Sprintf("\t - %s", err))
		}
	}

	return fmt.Errorf(strings.Join(msgList, "\n"))
}

func internalSortFileErrors(errCollector errorCollector) (map[string]error, []string) {
	// sort file to always get a predictable output (for tests mostly)
	subErrList := map[string]error{}
	keys := []string{}

	for errItem := range errCollector {
		subErrList[errItem.File] = errItem.Err
		keys = append(keys, errItem.File)
	}

	sort.Strings(keys)

	return subErrList, keys
}

func WriteTerraformFiles(rootPath string, files map[string]*hclwrite.File) (err error) {
	if len(files) == 0 {
		return nil
	}

	waitGroup := &sync.WaitGroup{}
	errCollector := make(errorCollector, len(files))

	fs, statError := os.Stat(rootPath)
	exists := !os.IsNotExist(statError)
	isDir := exists && statError == nil && fs.IsDir()

	if exists && isDir {
		for fName, hclFile := range files {
			waitGroup.Add(1)

			go writeTerraformFileAsync(path.Join(rootPath, fName), hclFile, errCollector, waitGroup)
		}

		waitGroup.Wait()
	}

	close(errCollector)

	if len(errCollector) > 0 || statError != nil || !exists || !isDir {
		return WriteTerraformFileError(
			generateWritingFileErrors(rootPath, exists, statError, isDir, errCollector),
		)
	}

	return nil
}

func generateWritingFileErrors(
	rootPath string,
	exists bool,
	statError error,
	isDir bool,
	errCollector errorCollector,
) []error {
	switch {
	case !exists:
		return []error{fmt.Errorf("\t - open %s: no such file or directory", rootPath)}
	case statError != nil:
		return []error{fmt.Errorf("\t - %w", statError)}
	case !isDir:
		return []error{fmt.Errorf("\t - %s is not a directory", rootPath)}
	default:
		// sort file to always get a predictable output (for tests mostly)
		errList := map[string]error{}
		keys := make([]string, 0, len(errList))

		for errItem := range errCollector {
			errList[errItem.File] = errItem.Err
			keys = append(keys, errItem.File)
		}

		sort.Strings(keys)

		msgList := []error{}

		for _, file := range keys {
			msgList = append(msgList, fmt.Errorf("\t - %w", errList[file]))
		}

		return msgList
	}
}

/** Private **/

type fileCollectorItem struct {
	name string
	file *hclwrite.File
}
type errorCollectorItem struct {
	File string
	Err  error
}
type (
	fileCollector  chan fileCollectorItem
	errorCollector chan errorCollectorItem
)

func generateHclRepoFileAsync(
	repoConfig *GhRepoConfig,
	valGen tfsig.ValueGenerator,
	repoTfId string,
	collector fileCollector,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	fname := fmt.Sprintf("repo.%s.tf", repoTfId)

	collector <- fileCollectorItem{name: fname, file: NewHclRepository(repoTfId, repoConfig, valGen)}
}

func writeTerraformFileAsync(path string, hclFile *hclwrite.File, errCollector errorCollector, wg *sync.WaitGroup) {
	defer wg.Done()

	formatted := hclwrite.Format(hclFile.Bytes())

	var (
		file *os.File
		err  error
	)

	fName := filepath.Base(path)

	if file, err = os.Create(path); err != nil {
		errCollector <- errorCollectorItem{fName, err}
	} else {
		log.Debug().Msgf("Writing terraform file '%s'", file.Name())

		if _, err = file.Write(formatted); err != nil {
			errCollector <- errorCollectorItem{fName, err}
		}
	}
}
