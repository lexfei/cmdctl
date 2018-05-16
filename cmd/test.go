package cmd

import (
	"fmt"
	"io"

	"cmdctl/cmd/templates"
	cmdutil "cmdctl/cmd/util"
	"cmdctl/pkg/i18n"

	"github.com/spf13/cobra"
)

var (
	testExample = templates.Examples(i18n.T(`
		# Run simple test command
		cmdctl test

		# Run command with option
		cmdctl test -a 8888`))
)

func NewCmdTest(f cmdutil.Factory, out io.Writer, cmdErr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "test",
		Short:   i18n.T("Hello world command"),
		Long:    "Hello world command",
		Example: testExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(RunTest(f, out, cmdErr, cmd, args))
			return
		},
		Aliases: []string{},
	}

	//cmdutil.AddTestFlags(cmd)
	cmd.Flags().BoolP("create", "", false, "Create the template")
	cmd.Flags().StringP("appId", "a", "0949", "Specify the user appId.")
	cmd.Flags().StringP("format", "", "yaml", "Specify the file format: 'json' or 'yaml'.")

	return cmd
}

func RunTest(f cmdutil.Factory, out io.Writer, cmdErr io.Writer, cmd *cobra.Command, args []string) error {
	create := cmdutil.GetFlagBool(cmd, "create")
	format := cmdutil.GetFlagString(cmd, "format")
	appId := cmdutil.GetFlagString(cmd, "appId")

	fmt.Fprintf(out, "Hello world, params is:\n")
	fmt.Fprintf(out, "==> create: %v\n==> format: %v\n==> appId: %s\n", create, format, appId)
	return nil
}
