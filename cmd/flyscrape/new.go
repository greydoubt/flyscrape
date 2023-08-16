package main

import (
	"flag"
	"fmt"
	"os"

	"flyscrape/js"
)

type NewCommand struct{}

func (c *NewCommand) Run(args []string) error {
	fs := flag.NewFlagSet("flyscrape-new", flag.ContinueOnError)
	fs.Usage = c.Usage

	if err := fs.Parse(args); err != nil {
		return err
	} else if fs.NArg() == 0 || fs.Arg(0) == "" {
		return fmt.Errorf("script path required")
	} else if fs.NArg() > 1 {
		return fmt.Errorf("too many arguments")
	}

	script := fs.Arg(0)
	if _, err := os.Stat(script); err == nil {
		return fmt.Errorf("script already exists")
	}

	if err := os.WriteFile(script, js.Template, 0o644); err != nil {
		return fmt.Errorf("failed to create script %q: %w", script, err)
	}

	fmt.Printf("Scraping script %v created.\n", script)
	return nil
}

func (c *NewCommand) Usage() {
	fmt.Println(`
The new command creates a new scraping script.

Usage:

    flyscrape new SCRIPT


Examples:

    # Create a new scraping script.
    $ flyscrape new example.js
`[1:])
}