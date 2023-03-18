package core

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
		return nil, MaxTemplateDepthReachedError(tplType, path)
	}

	var (
		err     error
		tplList []*T
		tpl     *T
	)

	if tpl = loaderFn(tplName); tpl == nil {
		return nil, UnknownTemplateError(tplType, tplName)
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
			return nil, MaxTemplateCountReachedError(tplType, path)
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
