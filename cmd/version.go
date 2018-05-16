package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"cmdctl/cmd/templates"
	cmdutil "cmdctl/cmd/util"
	"cmdctl/pkg/version"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

type Version struct {
	ClientVersion *version.Info `json:"clientVersion,omitempty" yaml:"clientVersion,omitempty"`
	ServerVersion *version.Info `json:"serverVersion,omitempty" yaml:"serverVersion,omitempty"`
}

// VersionOptions: describe the options available to users of the "cmdctl
// version" command.
type VersionOptions struct {
	short  bool
	output string
}

var (
	versionExample = templates.Examples(`
		# Print the client and server versions for the current context
		cmdctl version`)
)

func NewCmdVersion(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Print the client and server version information",
		Long:    "Print the client and server version information for the current context",
		Example: versionExample,
		Run: func(cmd *cobra.Command, args []string) {
			options := new(VersionOptions)
			cmdutil.CheckErr(options.Complete(cmd))
			cmdutil.CheckErr(options.Validate())
			cmdutil.CheckErr(options.Run(out))
		},
	}
	cmd.Flags().BoolP("short", "", false, "Print just the version number.")
	cmd.Flags().StringP("output", "o", "", "One of 'yaml' or 'json'.")
	return cmd
}

func (o *VersionOptions) Run(out io.Writer) error {
	var (
		serverVersion *version.Info
		serverErr     error
		versionInfo   Version
	)

	//clientVersion := version.Info{Major:"1.0"}
	clientVersion := version.Get()
	versionInfo.ClientVersion = &clientVersion

	switch o.output {
	case "":
		if o.short {
			fmt.Fprintf(out, "Client Version: %s\n", clientVersion.GitTag)
			if serverVersion != nil {
				fmt.Fprintf(out, "Server Version: %s\n", serverVersion.GitTag)
			}
		} else {
			fmt.Fprintf(out, "Client Version: %s\n", fmt.Sprintf("%#v", clientVersion))
			if serverVersion != nil {
				fmt.Fprintf(out, "Server Version: %s\n", fmt.Sprintf("%#v", *serverVersion))
			}
		}
	case "yaml":
		marshalled, err := yaml.Marshal(&versionInfo)
		if err != nil {
			return err
		}
		fmt.Fprintln(out, string(marshalled))
	case "json":
		marshalled, err := json.MarshalIndent(&versionInfo, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintln(out, string(marshalled))
	default:
		// There is a bug in the program if we hit this case.
		// However, we follow a policy of never panicking.
		return fmt.Errorf("VersionOptions were not validated: --output=%q should have been rejected", o.output)
	}

	return serverErr
}

func (o *VersionOptions) Complete(cmd *cobra.Command) error {
	o.short = cmdutil.GetFlagBool(cmd, "short")
	o.output = cmdutil.GetFlagString(cmd, "output")
	return nil
}

func (o *VersionOptions) Validate() error {
	if o.output != "" && o.output != "yaml" && o.output != "json" {
		return errors.New(`--output must be 'yaml' or 'json'`)
	}

	return nil
}
