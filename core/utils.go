package core

import "sort"

func MapToSortedList[T any](list map[string]T) []T {
	_, newList := MapToSortedListWithKeys(list)

	return newList
}

func MapToSortedListWithKeys[T any](list map[string]T) ([]string, []T) {
	// sort file to always get a predictable output (for tests mostly)
	newList := []T{}
	keys := make([]string, 0, len(list))

	for k := range list {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, file := range keys {
		newList = append(newList, list[file])
	}

	return keys, newList
}
