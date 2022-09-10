package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hertg/egpu-switcher/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docsPath string

var gendocsCommand = &cobra.Command{
	Use:    "gendocs",
	Short:  "Generate docs for egpu-switcher",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := genMarkdown(); err != nil {
			return err
		}
		if err := genManpage(); err != nil {
			return err
		}
		if verbose {
			logger.Success("generated all docs in %s", docsPath)
		}
		return nil
	},
}

func genMarkdown() error {
	dir := filepath.Join(docsPath, "md")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("unable to create 'md' docs folder: %s", err)
	}
	return doc.GenMarkdownTree(rootCmd, dir)
}

func genManpage() error {
	dir := filepath.Join(docsPath, "man")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("unable to create 'man' docs folder: %s", err)
	}
	header := &doc.GenManHeader{
		Title:   "egpu-switcher",
		Section: "1",
	}
	return doc.GenManTree(rootCmd, header, dir)
}

func init() {
	rootCmd.AddCommand(gendocsCommand)
	gendocsCommand.PersistentFlags().StringVarP(&docsPath, "out", "o", "./docs", "directory used for output")
}
