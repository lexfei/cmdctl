package cmd

import (
	"fmt"
	"io"

	"cmdctl/cmd/templates"
	cmdutil "cmdctl/cmd/util"
	"cmdctl/model"
	"cmdctl/pkg/i18n"

	"github.com/spf13/cobra"
)

type CreateOptions struct {
	username string
	password string
	email    string
}

var (
	addExample = templates.Examples(i18n.T(`
		# Add a new user lkong with password
		cmdctl add lkong lkongpasswd

		# Add a new user lkong with email
		cmdctl add lkong lkongpasswd -e 466701708@qq.com`))
)

func NewCmdAdd(f cmdutil.Factory, out io.Writer, cmdErr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add USERNAME PASSWORD",
		Short:   i18n.T("Add a user"),
		Long:    "Add a user",
		Example: addExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(validateCreateArgs(cmd, args))
			options := new(CreateOptions)
			cmdutil.CheckErr(options.Complete(cmd))
			if err := options.Validate(); err != nil {
				cmdutil.CheckErr(cmdutil.UsageErrorf(cmd, err.Error()))
			}
			cmdutil.CheckErr(options.RunAdd(f, out, cmdErr, args))
			return
		},
		Aliases: []string{},
	}

	cmd.Flags().StringP("email", "e", "", "Specify the user email.")
	return cmd
}

func validateCreateArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return cmdutil.UsageErrorf(cmd, "Unexpected args: %v", args)
	}

	if !isUsername(args[0]) {
		return fmt.Errorf("%s is not a legal user name", args[0])
	}

	return nil
}

func (o *CreateOptions) RunAdd(f cmdutil.Factory, out io.Writer, cmdErr io.Writer, args []string) error {
	user := model.UserModel{
		Username: args[0],
		Password: args[1],
		Email:    o.email,
	}

	db := model.GetSelfDB()
	defer db.Close()

	if err := db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (o *CreateOptions) Complete(cmd *cobra.Command) error {
	o.email = cmdutil.GetFlagString(cmd, "email")
	return nil
}

func (o *CreateOptions) Validate() error {
	if o.email != "" {
		if !isEmail(o.email) {
			return fmt.Errorf("%s is not a email format", o.email)
		}
	}

	return nil
}
