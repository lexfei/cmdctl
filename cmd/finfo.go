package cmd

import (
	"fmt"
	"io"

	"cmdctl/cmd/templates"
	cmdutil "cmdctl/cmd/util"

	"github.com/spf13/cobra"
)

var (
	finfoExample = templates.Examples(`
		# Get http server basic information(show how to send http request)
		cmdctl finfo`)
)

func NewCmdFinfo(f cmdutil.Factory, out io.Writer, cmdErr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "finfo",
		Short:   "Get http server basic information",
		Long:    "Get http server basic information",
		Example: finfoExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(RunFinfo(f, out, cmdErr, cmd, args))
			return
		},
		Aliases: []string{},
	}

	return cmd
}

func RunFinfo(f cmdutil.Factory, out io.Writer, cmdErr io.Writer, cmd *cobra.Command, args []string) error {
	request := f.Gorequest()
	resp, body, errs := request.Get("http://"+f.FileServer().Server+"/-/status").
		Set("Authorization", "Basic "+f.Auth()).
		Send(``).
		End()

	if err := cmdutil.CombineRequestErr(resp, body, errs); err != nil {
		return err
	}

	fmt.Printf("%s\n", body)
	return nil
}
