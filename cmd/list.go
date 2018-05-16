package cmd

import (
	"io"

	"cmdctl/cmd/templates"
	cmdutil "cmdctl/cmd/util"
	"cmdctl/model"
	"cmdctl/pkg/i18n"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	listExample = templates.Examples(i18n.T(`
	# List existing users
	cmdctl list`))
)

func NewCmdList(f cmdutil.Factory, out io.Writer, cmdErr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   i18n.T("List existing users"),
		Long:    "List existing users",
		Example: listExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(RunList(f, out, cmdErr, cmd, args))
			return
		},
		Aliases: []string{"li"},
	}

	return cmd
}

func RunList(f cmdutil.Factory, out io.Writer, cmdErr io.Writer, cmd *cobra.Command, args []string) error {
	model.DB.Init()
	defer model.DB.Close()

	users, _, err := model.ListUser("", 0, 0)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(out)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColWidth(TABLE_WIDTH)
	table.SetHeader([]string{"Username", "Password", "Email"})
	for _, user := range users {
		table.Append([]string{color.RedString(user.Username), user.Password, user.Email})
	}
	table.Render()
	return nil
}
