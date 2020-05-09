package label

import "strings"

// NameTransform A name transformation.
type NameTransform func(name string) string

// NameIdentity Identity transformation.
func NameIdentity(name string) string {
	return name
}

// TrimPrefix transformation.
func TrimPrefix(prefix string) NameTransform {
	return func(name string) string {
		return strings.TrimPrefix(name, prefix)
	}
}

// Chain transformations.
func Chain(nts ...NameTransform) NameTransform {
	return func(name string) string {
		result := name
		for _, nt := range nts {
			result = nt(result)
		}
		return result
	}
}
