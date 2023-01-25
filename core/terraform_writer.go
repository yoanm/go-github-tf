package core

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/rs/zerolog/log"

	"github.com/yoanm/go-tfsig"

	"github.com/yoanm/go-gh2tf"
)

/** Public **/

func GenerateHclRepoFiles(configList []*GhRepoConfig) (map[string]*hclwrite.File, error) {
	valueGenerator := gh2tf.NewValueGenerator()
	wg := &sync.WaitGroup{}
	collector := make(fileCollector, len(configList))
	errCollector := make(errorCollector, len(configList))
	var errList []error

	for k, c := range configList {
		wg.Add(1)
		if c.Name == nil {
			errList = append(errList, fmt.Errorf("config #%d: repository name is mandatory", k))
		} else {
			repoTfId := tfsig.ToTerraformIdentifier(*c.Name)
			go generateHclRepoFileAsync(c, valueGenerator, repoTfId, collector, wg)
		}
	}

	wg.Wait()
	close(collector)
	close(errCollector)
	if len(errCollector) > 0 || len(errList) > 0 {
		msgList := []string{"error while generating files:"}
		if len(errCollector) > 0 {
			// sort file to always get a predictable output (for tests mostly)
			subErrList := map[string]error{}
			for errItem := range errCollector {
				subErrList[errItem.File] = errItem.Err
			}
			var keys []string
			for k := range subErrList {
				keys = append(keys, k)
			}
			sort.Strings(keys)

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

		return nil, fmt.Errorf(strings.Join(msgList, "\n"))
	}

	list := map[string]*hclwrite.File{}
	for fs := range collector {
		list[fs.name] = fs.file
	}

	return list, nil
}

func WriteTerraformFiles(rootPath string, files map[string]*hclwrite.File) (err error) {
	if len(files) == 0 {
		return nil
	}

	wg := &sync.WaitGroup{}
	errCollector := make(errorCollector, len(files))

	fs, statError := os.Stat(rootPath)
	exists := !os.IsNotExist(statError)
	isDir := exists && statError == nil && fs.IsDir()
	if exists && isDir {
		for fName, hclFile := range files {
			wg.Add(1)
			go writeTerraformFileAsync(path.Join(rootPath, fName), hclFile, errCollector, wg)
		}

		wg.Wait()
	}

	close(errCollector)
	if len(errCollector) > 0 || statError != nil || !exists || !isDir {
		msgList := []string{"error while writing files:"}
		if !exists {
			msgList = append(msgList, fmt.Sprintf("\t - open %s: no such file or directory", rootPath))
		} else if statError != nil {
			msgList = append(msgList, fmt.Sprintf("\t - %s", statError))
		} else if !isDir {
			msgList = append(msgList, fmt.Sprintf("\t - %s is not a directory", rootPath))
		} else {
			// sort file to always get a predictable output (for tests mostly)
			errList := map[string]error{}
			for errItem := range errCollector {
				errList[errItem.File] = errItem.Err
			}
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

/** Private **/

type fileCollectorItem struct {
	name string
	file *hclwrite.File
}
type errorCollectorItem struct {
	File string
	Err  error
}
type fileCollector chan fileCollectorItem
type errorCollector chan errorCollectorItem

func generateHclRepoFileAsync(c *GhRepoConfig, valGen tfsig.ValueGenerator, repoTfId string, collector fileCollector, wg *sync.WaitGroup) {
	defer wg.Done()

	fname := fmt.Sprintf("repo.%s.tf", repoTfId)

	collector <- fileCollectorItem{name: fname, file: NewHclRepository(repoTfId, c, valGen)}
}

func writeTerraformFileAsync(path string, hclFile *hclwrite.File, errCollector errorCollector, wg *sync.WaitGroup) {
	defer wg.Done()

	formatted := hclwrite.Format(hclFile.Bytes())
	// _, diags := hclwrite.ParseConfig(formatted, "", hcl.InitialPos)
	// if diags.HasErrors() {
	//	var errorMessage string
	//	for _, err := range diags.Errs() {
	//		errorMessage = fmt.Sprintf("%s\n%s", errorMessage, err)
	//	}
	//	log.Error().Msgf("errors with terraform config: \n‰s", errorMessage)
	// } else {
	var f *os.File
	var err error
	fName := filepath.Base(path)
	if f, err = os.Create(path); err != nil {
		errCollector <- errorCollectorItem{fName, err}
	} else {
		log.Debug().Msgf("Writing terraform file '%s'", f.Name())
		/*if zerolog.GlobalLevel() == zerolog.TraceLevel {
			log.Trace().Msgf("Terraform content: \n%s", formatted)
		}*/
		if _, err = f.Write(formatted); err != nil {
			errCollector <- errorCollectorItem{fName, err}
		}
	}
	// }
}
