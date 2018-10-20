package main

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/containous/flaeg"
	"github.com/containous/staert"
	"github.com/ldez/gcg/core"
	"github.com/ldez/gcg/types"
	"github.com/ogier/pflag"
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
The generator use only Pull Requests.`,
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
		Metadata: map[string]string{
			"parseAllSources": "true",
		},
	}

	flag := flaeg.New(rootCmd, os.Args[1:])
	flag.AddParser(reflect.TypeOf(types.DisplayLabelOptions{}), &types.LabelDisplayOptionsParser{})
	flag.AddParser(reflect.TypeOf([]string{}), &types.SliceString{})

	if _, err := flag.Parse(rootCmd); err != nil {
		if err == pflag.ErrHelp {
			os.Exit(0)
		}
		log.Fatalf("Error parsing command: %v\n", err)
	}

	s := staert.NewStaert(rootCmd)

	// init TOML source
	toml := staert.NewTomlSource("gcg", []string{config.ConfigFile, "."})

	// add sources to staert
	s.AddSource(toml)
	s.AddSource(flag)

	if _, err := s.LoadConfig(); err != nil {
		if err != pflag.ErrHelp {
			os.Exit(0)
		}
		log.Fatalf("Error reading TOML config file %s : %v\n", toml.ConfigFileUsed(), err)
	}

	if err := s.Run(); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
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
		return fmt.Errorf("%s is mandatory", fieldName)
	}
	return nil
}
