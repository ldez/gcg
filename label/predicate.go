package label

import (
	"strings"

	"github.com/google/go-github/v27/github"
)

// Predicate A label predicate
type Predicate func(label github.Label) bool

// All Keep all labels
func All(_ github.Label) bool {
	return true
}

// Nothing Exclude all labels
func Nothing(_ github.Label) bool {
	return false
}

// Not Negate a predicate
func Not(predicate Predicate) Predicate {
	return func(label github.Label) bool {
		return !predicate(label)
	}
}

// HasPrefix label predicate
func HasPrefix(prefix string) Predicate {
	return func(label github.Label) bool {
		return strings.HasPrefix(*label.Name, prefix)
	}
}

// HasSuffix label predicate
func HasSuffix(suffix string) Predicate {
	return func(label github.Label) bool {
		return strings.HasSuffix(*label.Name, suffix)
	}
}

// AllMatch label predicate
func AllMatch(predicates ...Predicate) Predicate {
	return func(label github.Label) bool {
		for _, predicate := range predicates {
			if !predicate(label) {
				return false
			}
		}
		return true
	}
}

// AnyMatch label predicate
func AnyMatch(predicates ...Predicate) Predicate {
	return func(label github.Label) bool {
		for _, predicate := range predicates {
			if predicate(label) {
				return true
			}
		}
		return false
	}
}
