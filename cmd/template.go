package cmd

import (
	"io"

	cmdutil "cmdctl/cmd/util"
	"cmdctl/pkg/i18n"

	"github.com/spf13/cobra"
)

func NewCmdTemplate(f cmdutil.Factory, out io.Writer, cmdErr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template SUBCOMMAND",
		Short: i18n.T("Import and Export template"),
		Long:  "Import and Export template",
		Run: func(cmd *cobra.Command, args []string) {
			// run sub command
			defaultRunFunc := cmdutil.DefaultSubCommandRun(out)
			defaultRunFunc(cmd, args)
			return
		},
		Aliases: []string{"tp"},
	}

	// sub command
	cmd.AddCommand(NewCmdTemplateImport(f, out, cmdErr))
	cmd.AddCommand(NewCmdTemplateExport(f, out, cmdErr))

	//cmdutil.AddCleanFlags(cmd)
	//cmd.Flags().StringP("user", "u", "", "Specify the user name.")
	//cmd.Flags().BoolP("erase", "c", false, "Erase the records from the db")

	return cmd
}
