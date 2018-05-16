package cmd

import (
	"fmt"
	"io"
	"net"
	"os"

	"cmdctl/cmd/templates"
	cmdutil "cmdctl/cmd/util"
	"cmdctl/pkg/i18n"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	validateExample = templates.Examples(i18n.T(`
		# Validate the basic environment for cmdctl to run
		cmdctl validate`))
)

type ValidateInfo struct {
	ItemName string
	Status   string
	Message  string
}

func NewCmdValidate(f cmdutil.Factory, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "validate",
		Short:   i18n.T("Validate the basic environment for cmdctl to run"),
		Long:    "Validate the basic environment for cmdctl to run",
		Example: validateExample,
		Run: func(cmd *cobra.Command, args []string) {
			err := RunValidate(f, out, cmd)
			cmdutil.CheckErr(err)
		},
		Aliases: []string{"va", ""},
	}
	return cmd
}

func RunValidate(f cmdutil.Factory, out io.Writer, cmd *cobra.Command) error {
	data := [][]string{}
	FAIL := color.RedString("FAIL")
	PASS := color.GreenString("PASS")
	validateInfo := ValidateInfo{}

	// check if can access db
	validateInfo.ItemName = "db connection"
	_, err := net.Dial("tcp", f.FileServer().Server)
	//defer client.Close()
	if err != nil {
		validateInfo.Status = FAIL
		validateInfo.Message = fmt.Sprintf("%v", err)
	} else {
		validateInfo.Status = PASS
		validateInfo.Message = ""
	}
	data = append(data, []string{validateInfo.ItemName, validateInfo.Status, validateInfo.Message})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColWidth(TABLE_WIDTH)
	table.SetHeader([]string{"ValidateItem", "Result", "Message"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output
	return nil
}
