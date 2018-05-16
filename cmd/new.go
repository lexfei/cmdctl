package cmd

import (
	"fmt"
	"io"
	"os"
	"text/template"

	"cmdctl/cmd/templates"
	cmdutil "cmdctl/cmd/util"
	"cmdctl/pkg/i18n"

	"github.com/spf13/cobra"
)

type replace struct {
	Cmd     string
	Cmdfunc string
	Desc    string
	Dot     string
}

var optionText string = `package cmd

import (
	"fmt"
	"io"

	"cmdctl/cmd/templates"
	cmdutil "cmdctl/cmd/util"

	"github.com/spf13/cobra"
)

var (
	{{.Cmd}}Example = templates.Examples({{.Dot}}
		# Command desc
		cmdctl {{.Cmd}}

		# Command desc
		cmdctl {{.Cmd}} -u lkong update

		# Command desc
		cmdctl {{.Cmd}} -c download{{.Dot}})
)

func NewCmd{{.Cmdfunc}}() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "{{.Cmd}}",
		Short:   "Command desc",
		Long:    "Command desc",
		Example: {{.Cmd}}Example,
		Run: func(cmd *cobra.Command, args []string) {
			//cmdutil.CheckErr(validateArgs(cmd, args))
			cmdutil.CheckErr(Run{{.Cmdfunc}}(cmd, args))
			return
		},
		Aliases: []string{},
	}

	//cmdutil.Add{{.Cmdfunc}}Flags(cmd)
	cmd.Flags().BoolP("create", "", false, "create the template")
	cmd.Flags().StringP("appId", "a", "", "Specify the user appId.")
	cmd.Flags().StringP("format", "", "yaml", "Specify the file format: 'json' or 'yaml'.")

	return cmd
}

/*
func validateArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return cmdutil.UsageErrorf(cmd, "Unexpected args: %v", args)
	}

	if len(args) > 1 {
		color.Yellow("only import %s\n", args[0])
	}

	return nil
}
*/

func Run{{.Cmdfunc}}(f cmdutil.Factory, out io.Writer, cmdErr io.Writer, cmd *cobra.Command, args []string) error {
	create := cmdutil.GetFlagBool(cmd, "create")
	appId := cmdutil.GetFlagString(cmd, "appId")
	format := cmdutil.GetFlagString(cmd, "format")
	fmt.Printf("appId: %v, create: %v, format: %v\n", appId, create, format)

	// do some thing
	return nil
}`

var groupText string = `package cmd

import (
	"io"

	"cmdctl/cmd/templates"
	cmdutil "cmdctl/cmd/util"

	"github.com/spf13/cobra"
)

type {{.Cmdfunc}}Options struct {
	appId  string
	format string
	create bool
}

var (
	{{.Cmd}}Example = templates.Examples({{.Dot}}
		# Command desc
		cmdctl {{.Cmd}}

		# Command desc
		cmdctl {{.Cmd}} -u lkong update

		# Command desc
		cmdctl {{.Cmd}} -c download{{.Dot}})
)

func NewCmd{{.Cmdfunc}}() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "{{.Cmd}}",
		Short:   "Command desc",
		Long:    "Command desc",
		Example: {{.Cmd}}Example,
		Run: func(cmd *cobra.Command, args []string) {
			//cmdutil.CheckErr(validateArgs(cmd, args))
			defaultRunFunc := cmdutil.DefaultSubCommandRun(out)
			defaultRunFunc(cmd, args)
			return
		},
		Aliases: []string{},
	}

	//cmdutil.Add{{.Cmdfunc}}Flags(cmd)
	cmd.Flags().BoolP("create", "", false, "create the template")
	cmd.Flags().StringP("appId", "a", "", "Specify the user appId.")
	cmd.Flags().StringP("format", "", "yaml", "Specify the file format: 'json' or 'yaml'.")

	// {{.Cmd}} subcommands
	cmd.AddCommand(NewCmd{{.Cmdfunc}}(f, out, cmdErr))
	cmd.AddCommand(NewCmd{{.Cmdfunc}}(f, out, cmdErr))

	return cmd
}

// some common functions
func ValidateArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return cmdutil.UsageErrorf(cmd, "Unexpected args: %v", args)
	} 
	return nil
}`

var subText string = `package cmd

import (
	"fmt"
	"io"

	"cmdctl/cmd/templates"
	cmdutil "cmdctl/cmd/util"

	"github.com/spf13/cobra"
)

var (
	{{.Cmd}}Example = templates.Examples({{.Dot}}
		# Command desc
		gctl {{.Cmd}}

		# Command desc
		gctl {{.Cmd}} -u lkong update

		# Command desc
		gctl {{.Cmd}} -c download{{.Dot}})
)

func NewCmd{{.Cmdfunc}}() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "{{.Cmd}}",
		Short:   "Command desc",
		Long:    "Command desc",
		Example: {{.Cmd}}Example,
		Run: func(cmd *cobra.Command, args []string) {
			//cmdutil.CheckErr(validateArgs(cmd, args))
			cmdutil.CheckErr(Run{{.Cmdfunc}}(cmd, args))
			return
		},
		Aliases: []string{},
	}

	//cmdutil.Add{{.Cmdfunc}}Flags(cmd)
	cmd.Flags().BoolP("create", "", false, "create the template")
	cmd.Flags().StringP("appId", "a", "", "Specify the user appId.")
	cmd.Flags().StringP("format", "", "yaml", "Specify the file format: 'json' or 'yaml'.")

	return cmd
}

/*
func validateArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return cmdutil.UsageErrorf(cmd, "Unexpected args: %v", args)
	}

	if len(args) > 1 {
		color.Yellow("only import %s\n", args[0])
	}

	return nil
}
*/

func Run{{.Cmdfunc}}(f cmdutil.Factory, out io.Writer, cmdErr io.Writer, cmd *cobra.Command, args []string) error {
	create := cmdutil.GetFlagBool(cmd, "create")
	appId := cmdutil.GetFlagString(cmd, "appId")
	format := cmdutil.GetFlagString(cmd, "format")
	fmt.Printf("appId: %v, create: %v, format: %v\n", appId, create, format)

	// do some thing
	return nil
}`

var (
	newExample = templates.Examples(i18n.T(`
		# Create a default cmd file(cmdctl $cmd_name $cmd_function_name $cmd_description)
		cmdctl test Test "This is a test command"

		# Create a group command
		cmdctl -g test Test "This is a test command"

		# Create a command wiht options filled
		cmdctl -o test Test "This is a test command"`))
)

func NewCmdNew(f cmdutil.Factory, out io.Writer, cmdErr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "new CMDNAME CMDFUNCNAME CMDDESCRIPTION",
		Short:   i18n.T("New cmd format go source file"),
		Long:    "New cmd format go source file",
		Example: newExample,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmdutil.CheckErr(cmdutil.UsageErrorf(cmd, "Unexpected args: %v", args))
			}
			cmdutil.CheckErr(RunNew(f, out, cmdErr, cmd, args))
			return
		},
		Aliases: []string{},
	}

	cmd.Flags().BoolP("group", "g", false, "If the command have subcommands")
	cmd.Flags().BoolP("option", "o", false, "Build with options")

	return cmd
}

func RunNew(f cmdutil.Factory, out io.Writer, cmdErr io.Writer, cmd *cobra.Command, args []string) error {
	group := cmdutil.GetFlagBool(cmd, "group")
	option := cmdutil.GetFlagBool(cmd, "option")

	var text string
	if group {
		text = groupText
	} else {
		if option {
			text = optionText
		} else {
			text = subText
		}
	}

	var desc string = "Description of the command."
	if len(args) > 2 {
		desc = args[2]
	}

	r := replace{
		Cmd:     args[0],
		Cmdfunc: args[1],
		Desc:    desc,
		Dot:     "`",
	}

	tmpl, err := template.New("cmd").Parse(text)
	if err != nil {
		return err
	}

	file, err := os.Create(fmt.Sprintf("%s.go", r.Cmd))
	if err != nil {
		return err
	}

	err = tmpl.Execute(file, r)
	if err != nil {
		return err
	}

	fmt.Printf("New cmd format file generated: %s.go\n", r.Cmd)
	return nil
}
