package share

import(
	"errors"
)

func StringContains(list []string, s string) bool {
	set := make(map[string]int)

	for k, v := range list {
		set[v] = k
	}

	if set[s] != 0 {
		return true
	}

	return false
}

func StringDeleteSlice(list []string, s string) ([]string, error) {
	set := make(map[string]int)

	for k, v := range list {
		set[v] = k
	}

	i, v := set[s]
	if v {
		list = append(list[:i], list[i+1:]...)
		return list, nil
	}

	return nil, errors.New("Slice not found")
}
