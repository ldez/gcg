package label

import (
	"strings"

	"github.com/google/go-github/github"
)

type Predicate func(label github.Label) bool

func All(label github.Label) bool {
	return true
}

func Nothing(label github.Label) bool {
	return false
}

func Not(predicate Predicate) Predicate {
	return func(label github.Label) bool {
		return !predicate(label)
	}
}

func HasPrefix(prefix string) Predicate {
	return func(label github.Label) bool {
		return strings.HasPrefix(*label.Name, prefix)
	}
}

func HasSuffix(suffix string) Predicate {
	return func(label github.Label) bool {
		return strings.HasSuffix(*label.Name, suffix)
	}
}

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
