package label

import "strings"

type NameTransform func(name string) string

func NameIdentity(name string) string {
	return name
}

func TrimPrefix(prefix string) NameTransform {
	return func(name string) string {
		return strings.TrimPrefix(name, prefix)
	}
}

func Chain(nts ...NameTransform) NameTransform {
	return func(name string) string {
		result := name
		for _, nt := range nts {
			result = nt(result)
		}
		return result
	}
}
