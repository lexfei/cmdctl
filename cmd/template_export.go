package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"cmdctl/cmd/templates"
	cmdutil "cmdctl/cmd/util"
	"cmdctl/pkg/i18n"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type exportOptions struct {
	application bool
}

var (
	exportExample = templates.Examples(i18n.T(`
	# Export template
	cmdctl template export templateName

	# Export template with option
	cmdctl template export templateName -a app-afnbdef`))
)

func NewCmdTemplateExport(f cmdutil.Factory, out io.Writer, cmdErr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "export",
		Short:   i18n.T("Export template"),
		Long:    "Export template",
		Example: exportExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(validateExportArgs(cmd, args))

			options := new(exportOptions)
			cmdutil.CheckErr(options.Complete(cmd))
			if err := options.Validate(); err != nil {
				cmdutil.CheckErr(cmdutil.UsageErrorf(cmd, err.Error()))
			}
			cmdutil.CheckErr(options.Run(f, out, cmdErr, cmd, args))
			return
		},
		Aliases: []string{"imp"},
	}

	cmd.Flags().BoolP("application", "a", false, "Export from")
	cmd.Flags().StringP("format", "f", "yaml", "Specify the file format: 'json' or 'yaml'.")
	return cmd
}

func validateExportArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return cmdutil.UsageErrorf(cmd, "Unexpected args: %v", args)
	}

	if len(args) > 1 {
		color.Yellow("only export %s\n", args[0])
	}

	return nil
}

func (o *exportOptions) Run(f cmdutil.Factory, out io.Writer, cmdErr io.Writer, cmd *cobra.Command, args []string) error {
	fmt.Printf("template name: %s\n", args[0])
	return nil
	id := args[0]
	suffix := cmdutil.GetFlagString(cmd, "format")

	// create dir
	workDir := filepath.Join("/tmp", id)
	os.MkdirAll(workDir, 0755)
	defer func() {
		os.RemoveAll(workDir)
	}()

	rsList := []string{"nginx", "busybox"}
	for _, rs := range rsList {
		f, err := os.Create(filepath.Join(workDir, rs) + "." + suffix)
		if err != nil {
			return err
		}
		defer f.Close()

		f.WriteString(rs)
		f.Sync()
	}

	cmdArgs := []string{"-czf", id + ".tar.gz", "-C", "/tmp/", id}
	command := exec.Command("tar", cmdArgs...)
	cmdOut, err := command.CombinedOutput()
	if err != nil {
		os.Remove(id + ".tar.gz")
		return fmt.Errorf("run tar command, failed:%v, arguments:%v", string(cmdOut), cmdArgs)
	}

	return nil
}

func (o *exportOptions) Complete(cmd *cobra.Command) error {
	o.application = cmdutil.GetFlagBool(cmd, "application")
	return nil
}

func (o *exportOptions) Validate() error {
	return nil
}
