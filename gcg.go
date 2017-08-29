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
		BaseBranch:           types.DefaultBaseBranch,
		LabelEnhancement:     types.DefaultEnhancementLabel,
		LabelDocumentation:   types.DefaultDocumentationLabel,
		LabelBug:             types.DefaultBugLabel,
		OutputType:           types.DefaultOutputDestination,
		FileName:             types.DefaultFileName,
		DisplayLabel:         true,
		ThresholdPreviousRef: 1,
		ThresholdCurrentRef:  5,
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

			err := validateConfig(config)
			if err != nil {
				return err
			}

			core.Generate(config)
			return nil
		},
	}

	flag := flaeg.New(rootCmd, os.Args[1:])
	flag.AddParser(reflect.TypeOf(types.DisplayLabelOptions{}), &types.LabelDisplayOptionsParser{})
	flag.AddParser(reflect.TypeOf([]string{}), &types.SliceString{})
	flag.Run()
}

func validateConfig(config *types.Configuration) error {
	err := required(config.CurrentRef, "current-ref")
	if err != nil {
		return err
	}
	err = required(config.PreviousRef, "previous-ref")
	if err != nil {
		return err
	}
	err = required(config.Owner, "owner")
	if err != nil {
		return err
	}
	return required(config.RepositoryName, "repo-name")
}

func required(field string, fieldName string) error {
	if len(field) == 0 {
		log.Fatalf("%s is mandatory.", fieldName)
	}
	return nil
}
