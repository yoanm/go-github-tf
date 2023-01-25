package core

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func LoadTerraformFiles(root string, pathList []string) (list []*hclwrite.File, err error) {
	errList := map[string]error{}
	for _, path := range pathList {
		filePath := filepath.Join(root, path)

		conf, loadErr := LoadTerraformFile(filePath)
		if loadErr != nil {
			errList[filePath] = fmt.Errorf("\n%s", loadErr)
		} else {
			list = append(list, conf)
		}
	}

	if len(errList) > 0 {
		var msgList []string
		for k, v := range errList {
			msgList = append(msgList, fmt.Sprintf("%s: %s", k, v))
		}

		return nil, fmt.Errorf(strings.Join(msgList, "\n"))
	}

	return list, nil
}

func LoadTerraformFile(filePath string) (f *hclwrite.File, err error) {
	var content []byte
	if content, err = ioutil.ReadFile(filePath); err != nil {
		return nil, err
	}

	f, diags := hclwrite.ParseConfig(content, filePath, hcl.InitialPos)
	if diags != nil && diags.HasErrors() {
		var msgList []string
		for _, err = range diags.Errs() {
			msgList = append(msgList, fmt.Sprintf("\t - %s", err))
		}

		return nil, fmt.Errorf(strings.Join(msgList, "\n"))
	}

	return f, nil
}

func FindResourcesToImport(hclFile *hclwrite.File) map[string]string {
	l := map[string]string{}

	for _, block := range hclFile.Body().Blocks() {
		if block.Type() == "resource" {
			labels := block.Labels()
			if len(labels) < 2 {
				continue
			}
			label := labels[0]
			id := labels[1]
			if label == "github_repository" {
				l[fmt.Sprintf("%s.%s", label, id)] = getCleanAttributeValue(block.Body().Attributes()["name"])
			} else if label == "github_branch_default" {
				l[fmt.Sprintf("%s.%s", label, id)] = getCleanAttributeValue(block.Body().Attributes()["repository"])
			} else if label == "github_branch_protection" {
				l[fmt.Sprintf("%s.%s", label, id)] = fmt.Sprintf(
					"%s:%s",
					getCleanAttributeValue(block.Body().Attributes()["repository_id"]),
					getCleanAttributeValue(block.Body().Attributes()["pattern"]),
				)
			}
		}
	}

	return l
}

func getCleanAttributeValue(attr *hclwrite.Attribute) string {
	if attr == nil {
		return ""
	}

	// Convert attribute value to tokens, then to bytes, trim spaces and double quotes
	return strings.Trim(string(attr.Expr().BuildTokens(hclwrite.Tokens{}).Bytes()), " \"")
}
