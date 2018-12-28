package share

import(
	"errors"
)

func StringContains(a []string, s string) bool {
	for _, n := range a {
		if s == n {
			return true
		}
	}

	return false
}

func StringDeleteSlice(a []string, s string) ([]string, error) {
	for i, n := range a {
		if s == n {
			a = append(a[:i], a[i+1:]...)
			return a, nil
		}
	}

	return nil, errors.New("slice not found")
}
