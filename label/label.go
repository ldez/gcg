package label

import (
	"github.com/google/go-github/v30/github"
)

// FilterAndTransform Filter and transform labels
func FilterAndTransform(labels []*github.Label, filter Predicate, transform NameTransform) []string {
	var results []string
	for _, lbl := range labels {
		if filter(lbl) {
			results = append(results, transform(*lbl.Name))
		}
	}
	return results
}

// FilteredBy Filter labels by a Predicate
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

// ExcludedBy Exclude labels by a Predicate
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

// TrimAllPrefix Trim all prefix
func TrimAllPrefix(values []string) NameTransform {
	transform := NameIdentity
	for _, rp := range values {
		transform = Chain(transform, TrimPrefix(rp))
	}
	return transform
}
