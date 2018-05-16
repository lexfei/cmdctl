package cmd

import (
	"fmt"
	"io"

	"cmdctl/cmd/templates"
	cmdutil "cmdctl/cmd/util"
	"cmdctl/pkg/i18n"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type ImportOptions struct {
	user        bool
	appId       int
	description string
}

var (
	importExample = templates.Examples(i18n.T(`
	# Import template
	cmdctl template import template.tar.gz

	# Import template with options
	cmdctl template import -a 3xx -u lkong template.tar.gz`))
)

func NewCmdTemplateImport(f cmdutil.Factory, out io.Writer, cmdErr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "import",
		Short:   i18n.T("Import template from tar file"),
		Long:    "Import template from tar file",
		Example: importExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(validateArgs(cmd, args))

			options := new(ImportOptions)
			cmdutil.CheckErr(options.Complete(cmd))
			if err := options.Validate(); err != nil {
				cmdutil.CheckErr(cmdutil.UsageErrorf(cmd, err.Error()))
			}
			cmdutil.CheckErr(options.Run(f, out, cmdErr, args))
			return
		},
		Aliases: []string{"imp"},
	}

	cmd.Flags().StringP("description", "", "", "Specify the description.")
	cmd.Flags().BoolP("user", "", false, "Specify the create username.")
	cmd.Flags().IntP("appId", "", 0, "AppId to use")

	return cmd
}

func validateArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return cmdutil.UsageErrorf(cmd, "Unexpected args: %v", args)
	}

	if len(args) > 1 {
		color.Yellow("only import %s\n", args[0])
	}

	return nil
}

func (o *ImportOptions) Run(f cmdutil.Factory, out io.Writer, cmdErr io.Writer, args []string) error {
	fmt.Printf("import succ!")
	return nil
}

func (o *ImportOptions) Complete(cmd *cobra.Command) error {
	o.user = cmdutil.GetFlagBool(cmd, "user")
	o.appId = cmdutil.GetFlagInt(cmd, "appId")
	o.description = cmdutil.GetFlagString(cmd, "description")
	return nil
}

func (o *ImportOptions) Validate() error {
	return nil
}
