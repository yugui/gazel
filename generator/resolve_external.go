package generator

import (
	"strings"
)

type externalResolver struct{}

func (r externalResolver) resolve(importpath, dir string) (label, error) {
	components := strings.Split(importpath, "/")

	labels := strings.Split(components[0], ".")
	var reversed []string
	for i := range labels {
		l := labels[len(labels)-i-1]
		reversed = append(reversed, l)
	}
	return label{
		repo: strings.Join(reversed, "_"),
		pkg:  strings.Join(components[1:], "/"),
		name: "go_default_library",
	}, nil
}
