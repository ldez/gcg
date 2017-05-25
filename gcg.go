package main

import (
	"log"
	"os"
	"reflect"

	"github.com/containous/flaeg"
	"github.com/ldez/gcg/core"
	"github.com/ldez/gcg/types"
)

func main() {

	config := &types.Configuration{
		BaseBranch:         types.DefaultBaseBranch,
		LabelEnhancement:   types.DefaultEnhancementLabel,
		LabelDocumentation: types.DefaultDocumentationLabel,
		LabelBug:           types.DefaultBugLabel,
		OutputType:         types.DefaultOutputDestination,
		FileName:           types.DefaultFileName,
		DisplayLabel:       true,
	}

	defaultPointer := &types.Configuration{
		DisplayLabelOptions: &types.DisplayLabelOptions{},
	}

	rootCmd := &flaeg.Command{
		Name: "gcg",
		Description: `GCG is a GitHub Changelog Generator.
The generator use only Pull Requests.
		`,
		Config:                config,
		DefaultPointersConfig: defaultPointer,
		Run: func() error {
			if config.Debug {
				log.Printf("Run GCG command with config : %+v\n", config)
				log.Printf("Run GCG command with config : %+v\n", config.DisplayLabelOptions)
			}
			required(config.CurrentRef, "current-ref")
			required(config.PreviousRef, "previous-ref")
			required(config.Owner, "owner")
			required(config.RepositoryName, "repo-name")

			core.Generate(config)
			return nil
		},
	}

	flag := flaeg.New(rootCmd, os.Args[1:])
	flag.AddParser(reflect.TypeOf(types.DisplayLabelOptions{}), &types.LabelDisplayOptionsParser{})
	flag.AddParser(reflect.TypeOf([]string{}), &types.SliceString{})
	flag.Run()
}

func required(field string, fieldName string) error {
	if len(field) == 0 {
		log.Fatalf("%s is mandatory.", fieldName)
	}
	return nil
}
