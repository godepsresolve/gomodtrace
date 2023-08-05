package gomodtrace

// copyMap does a shallow copy of provided map, because map is a reference type.
func copyMap(input map[string]bool) map[string]bool {
	cp := make(map[string]bool)
	for k, v := range input {
		cp[k] = v
	}
	return cp
}

// unique filters out non-unique items from a slice.
func unique[T comparable](input []T) []T {
	index := make(map[T]bool)
	var output []T
	for _, item := range input {
		if index[item] {
			continue
		}
		output = append(output, item)
		index[item] = true
	}
	return output
}

// toIndex makes an index for slice.
func toIndex[T comparable](input []T) map[T]bool {
	index := make(map[T]bool)
	for _, item := range input {
		if index[item] {
			continue
		}
		index[item] = true
	}
	return index
}
