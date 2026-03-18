// Package commands contains all cobra commands that are used from the main
package commands

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Fabianexe/gocoverageplus/pkg/cleaner"
	"github.com/Fabianexe/gocoverageplus/pkg/complexity"
	"github.com/Fabianexe/gocoverageplus/pkg/config"
	"github.com/Fabianexe/gocoverageplus/pkg/coverage"
	"github.com/Fabianexe/gocoverageplus/pkg/source"
	"github.com/Fabianexe/gocoverageplus/pkg/writer"
)

func RootCommand() {
	var rootCmd = &cobra.Command{ //nolint:gochecknoglobals
		Use:   "gocoverageplus",
		Short: "gocoverageplus optimise the coverage report in go text format with source code.",
		Run: func(cmd *cobra.Command, _ []string) {
			initLogger()
			slog.Info("Start flag parsing")
			configPath, err := cmd.Flags().GetString("config")
			if err != nil {
				slog.Error(fmt.Sprintf("%+v", err))
				os.Exit(1)
			}

			inputPath, err := cmd.Flags().GetString("input")
			if err != nil {
				slog.Error(fmt.Sprintf("%+v", err))
				os.Exit(1)
			}

			outputPath, err := cmd.Flags().GetString("output")
			if err != nil {
				slog.Error(fmt.Sprintf("%+v", err))
				os.Exit(1)
			}

			slog.Info("Read Config")
			conf, err := config.ReadConfig(configPath)
			if err != nil {
				slog.Error(fmt.Sprintf("%+v", err))
				os.Exit(1)
			}
			if err := conf.Validate(); err != nil {
				slog.Error(fmt.Sprintf("%+v", err))
				os.Exit(1)
			}

			sourcePath, err := filepath.Abs(conf.SourcePath)
			if err != nil {
				slog.Error(fmt.Sprintf("%+v", err))
				os.Exit(1)
			}

			slog.Info("Load sources")
			project, err := source.LoadSources(sourcePath, conf.ExcludePaths)
			if err != nil {
				slog.Error(fmt.Sprintf("%+v", err))
				os.Exit(1)
			}

			slog.Info("Clean data")
			project = cleaner.CleanData(
				project,
				conf.Cleaner.Generated,
				conf.Cleaner.NoneCodeLines,
				conf.Cleaner.ErrorIf,
				conf.Cleaner.CustomIf,
			)

			if conf.Complexity.Active {
				slog.Info("Add complexity")
				cyclomatic := false
				if conf.Complexity.Type == "cyclomatic" {
					cyclomatic = true
				}
				project = complexity.AddComplexity(project, cyclomatic)
			}

			if inputPath != "-" {
				slog.Info("Load coverage")
				project, err = coverage.LoadCoverage(project, inputPath)

				if err != nil {
					slog.Error(fmt.Sprintf("%+v", err))
					os.Exit(1)
				}
			}

			slog.Info("Write output")
			if conf.OutputFormat == "cobertura" {
				err = writer.WriteXML(sourcePath, project, outputPath)
				if err != nil {
					slog.Error(fmt.Sprintf("%+v", err))
					os.Exit(1)
				}
			} else if conf.OutputFormat == "textfmt" {

				err = writer.WriteTextFMT(project, outputPath)
				if err != nil {
					slog.Error(fmt.Sprintf("%+v", err))
					os.Exit(1)
				}
			} else {
				slog.Error("Unknown output format")
				os.Exit(1)
			}
		},
	}

	rootCmd.PersistentFlags().StringP(
		"config",
		"c",
		".cov.json",
		"The config file path",
	)

	rootCmd.PersistentFlags().StringP(
		"input",
		"i",
		"coverage.cov",
		"The input file path",
	)

	rootCmd.PersistentFlags().StringP(
		"output",
		"o",
		"coverage.xml",
		"The output file path",
	)

	verboseFlag := rootCmd.PersistentFlags().VarPF(
		&verbose,
		"verbose",
		"v",
		"Add verbose output. Multiple -v options increase the verbosity.",
	)
	verboseFlag.NoOptDefVal = "1"

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
