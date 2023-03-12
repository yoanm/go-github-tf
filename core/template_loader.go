package core

import (
	"fmt"
	"strings"
)

const (
	TemplateMaxDepth = 10
	TemplateMaxCount = 10
)

func loadTemplate[T any](
	tplName string,
	loaderFn func(s string) *T,
	finderFn func(c *T) *[]string,
	tplType string,
	path ...string,
) ([]*T, error) {
	if len(path) > TemplateMaxDepth {
		return nil, fmt.Errorf(
			"more than %d levels of %s template detected for %s",
			TemplateMaxDepth,
			tplType,
			strings.Join(path, "->"),
		)
	}

	var (
		err     error
		tplList []*T
		tpl     *T
	)

	if tpl = loaderFn(tplName); tpl == nil {
		return nil, fmt.Errorf("unknown %s template %s", tplType, tplName)
	}

	if tplList, err = loadTemplateList(finderFn(tpl), loaderFn, finderFn, tplType, append(path, tplName)...); err != nil {
		return nil, err
	}

	tplList = append(tplList, tpl)

	return tplList, nil
}

func loadTemplateList[T any](
	tplNameList *[]string,
	loaderFn func(s string) *T,
	finderFn func(c *T) *[]string,
	tplType string,
	path ...string,
) ([]*T, error) {
	var (
		tplList []*T
		err     error
	)

	if tplNameList != nil {
		if len(*tplNameList) > TemplateMaxCount {
			pathString := "ROOT"

			if len(path) > 0 {
				pathString = strings.Join(path, "->")
			}

			return nil, fmt.Errorf("more than %d %s template detected for %s", TemplateMaxCount, tplType, pathString)
		}

		for _, tplName := range *tplNameList {
			var subTplList []*T

			if subTplList, err = loadTemplate[T](tplName, loaderFn, finderFn, tplType, path...); err != nil {
				return nil, err
			}

			tplList = append(tplList, subTplList...)
		}
	}

	return tplList, nil
}
