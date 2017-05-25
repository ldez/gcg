package label

import (
	"github.com/google/go-github/github"
)

func FilterAndTransform(labels []github.Label, filter Predicate, transform NameTransform) []string {

	var results []string
	for _, label := range labels {
		if filter(label) {
			results = append(results, transform(*label.Name))
		}
	}
	return results
}

func FilteredBy(namePredicate func(string) Predicate, values []string) Predicate {
	predicate := All
	if len(values) != 0 {
		predicate = Nothing
		for _, value := range values {
			if len(value) != 0 {
				predicate = AnyMatch(predicate, namePredicate(value))
			} else {
				predicate = AnyMatch(predicate, All)
			}
		}
	}
	return predicate
}

func ExcludedBy(namePredicate func(string) Predicate, values []string) Predicate {
	predicate := All
	if len(values) != 0 {
		for _, value := range values {
			if len(value) != 0 {
				predicate = AllMatch(predicate, Not(namePredicate(value)))
			} else {
				predicate = AllMatch(predicate, All)
			}
		}
	}
	return predicate
}

func TrimAllPrefix(values []string) NameTransform {
	transform := NameIdentity
	for _, rp := range values {
		transform = Chain(transform, TrimPrefix(rp))
	}
	return transform
}
