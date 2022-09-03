package cmd

import (
	"fmt"
	"os"

	"github.com/hertg/egpu-switcher/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var gendocsCommand = &cobra.Command{
	Use:    "gendocs",
	Short:  "Generate docs for egpu-switcher in ${pwd}/docs",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := genMarkdown(); err != nil {
			return err
		}
		if err := genManpage(); err != nil {
			return err
		}
		logger.Success("generated all docs in ./docs")
		return nil
	},
}

func genMarkdown() error {
	if err := os.MkdirAll("docs/md", os.ModePerm); err != nil {
		return fmt.Errorf("unable to create 'md' docs folder: %s", err)
	}
	return doc.GenMarkdownTree(rootCmd, "./docs/md")
}

func genManpage() error {
	if err := os.MkdirAll("docs/man", os.ModePerm); err != nil {
		return fmt.Errorf("unable to create 'man' docs folder: %s", err)
	}
	header := &doc.GenManHeader{
		Title:   "egpu-switcher",
		Section: "1",
	}
	return doc.GenManTree(rootCmd, header, "./docs/man")
}

func init() {
	rootCmd.AddCommand(gendocsCommand)
}
