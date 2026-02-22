package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "amctx",
	Short:        "amtool context manager",
	Args:         cobra.MaximumNArgs(1),
	RunE:         runRoot,
	SilenceUsage: true,
}

const noConfigFileMsg = "no config file found, file created"

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	red := color.New(color.FgRed, color.Bold)
	rootCmd.SetErrPrefix(red.Sprint("error:"))

	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println(`USAGE:
  amctx                         : list the aliases
  amctx <ALIAS>                 : switch to context <ALIAS>
  amctx <ALIAS>=<URL>           : create or update context
  amctx -h, --help              : show this message`)
	})
}

func runRoot(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return printAliases()
	}
	return parseAlias(args[0])
}

func printAliases() error {
	aliases, created, err := ListAliases()
	if err != nil {
		return err
	}

	if created {
		fmt.Println(noConfigFileMsg)
		return nil
	}

	if len(aliases) == 0 {
		fmt.Println("no aliases found")
		return nil
	}

	for _, alias := range aliases {
		fmt.Println(alias)
	}

	return nil
}

func addAlias(alias, url string) error {
	if alias == "" || url == "" {
		return fmt.Errorf("invalid argument format, expected <alias>=<url>")
	}
	created, err := CreateOrUpdateAlertmanagerAlias(alias, url)
	if err != nil {
		return err
	}
	if created {
		fmt.Println(noConfigFileMsg)
	}
	fmt.Printf("alias '%s' set with url: %s\n", alias, url)
	return nil
}

func switchContext(alias string) error {
	created, err := SwitchContext(alias)
	if err != nil {
		return err
	}
	if created {
		fmt.Println(noConfigFileMsg)
	}
	fmt.Printf("switched to context '%s'\n", alias)
	return nil
}

func parseAlias(arg string) error {
	if alias, url, ok := strings.Cut(arg, "="); ok {
		return addAlias(alias, url)
	}
	return switchContext(arg)
}
