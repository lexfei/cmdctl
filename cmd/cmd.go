package cmd

import (
	"io"
	"path/filepath"
	"strings"

	"cmdctl/cmd/templates"
	cmdutil "cmdctl/cmd/util"
	"cmdctl/pkg/homedir"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	bashCompletionFunc = `# call cmdctl get $1,
__cmdctl_override_flag_list=(gconfig cluster user context namespace server)
__cmdctl_override_flags()
{
    local ${__cmdctl_override_flag_list[*]} two_word_of of
    for w in "${words[@]}"; do
        if [ -n "${two_word_of}" ]; then
            eval "${two_word_of}=\"--${two_word_of}=\${w}\""
            two_word_of=
            continue
        fi
        for of in "${__cmdctl_override_flag_list[@]}"; do
            case "${w}" in
                --${of}=*)
                    eval "${of}=\"${w}\""
                    ;;
                --${of})
                    two_word_of="${of}"
                    ;;
            esac
        done
        if [ "${w}" == "--all-namespaces" ]; then
            namespace="--all-namespaces"
        fi
    done
    for of in "${__cmdctl_override_flag_list[@]}"; do
        if eval "test -n \"\$${of}\""; then
            eval "echo \${${of}}"
        fi
    done
}

__cmdctl_get_namespaces()
{
    local template cmdctl_out
    template="{{ range .items  }}{{ .metadata.name }} {{ end }}"
    if cmdctl_out=$(cmdctl get -o template --template="${template}" namespace 2>/dev/null); then
        COMPREPLY=( $( compgen -W "${cmdctl_out[*]}" -- "$cur" ) )
    fi
}

__cmdctl_config_get_contexts()
{
    __cmdctl_parse_config "contexts"
}

__cmdctl_config_get_clusters()
{
    __cmdctl_parse_config "clusters"
}

__cmdctl_config_get_users()
{
    __cmdctl_parse_config "users"
}

# $1 has to be "contexts", "clusters" or "users"
__cmdctl_config_get()
{
    local template cmdctl_out
    template="{{ range .$1  }}{{ .name }} {{ end }}"
    if cmdctl_out=$(cmdctl config $(__cmdctl_override_flags) -o template --template="${template}" view 2>/dev/null); then
        COMPREPLY=( $( compgen -W "${cmdctl_out[*]}" -- "$cur" ) )
    fi
}

__cmdctl_parse_get()
{
    local template
    template="{{ range .items  }}{{ .metadata.name }} {{ end }}"
    local cmdctl_out
    if cmdctl_out=$(cmdctl get $(__cmdctl_override_flags) -o template --template="${template}" "$1" 2>/dev/null); then
        COMPREPLY=( $( compgen -W "${cmdctl_out[*]}" -- "$cur" ) )
    fi
}

__cmdctl_get_resource()
{
    if [[ ${#nouns[@]} -eq 0 ]]; then
        return 1
    fi
    __cmdctl_parse_get "${nouns[${#nouns[@]} -1]}"
}

__cmdctl_get_resource_pod()
{
    __cmdctl_parse_get "pod"
}

__cmdctl_get_resource_rc()
{
    __cmdctl_parse_get "rc"
}

__cmdctl_get_resource_node()
{
    __cmdctl_parse_get "node"
}

# $1 is the name of the pod we want to get the list of containers inside
__cmdctl_get_containers()
{
    local template
    template="{{ range .spec.containers  }}{{ .name }} {{ end }}"
    __debug "${FUNCNAME} nouns are ${nouns[*]}"

    local len="${#nouns[@]}"
    if [[ ${len} -ne 1 ]]; then
        return
    fi
    local last=${nouns[${len} -1]}
    local cmdctl_out
    if cmdctl_out=$(cmdctl get $(__cmdctl_override_flags) -o template --template="${template}" pods "${last}" 2>/dev/null); then
        COMPREPLY=( $( compgen -W "${cmdctl_out[*]}" -- "$cur" ) )
    fi
}

# Require both a pod and a container to be specified
__cmdctl_require_pod_and_container()
{
    if [[ ${#nouns[@]} -eq 0 ]]; then
        __cmdctl_parse_get pods
        return 0
    fi;
    __cmdctl_get_containers
    return 0
}

__custom_func() {
    case ${last_command} in
        cmdctl_get | cmdctl_describe | cmdctl_delete | cmdctl_label | cmdctl_stop | cmdctl_edit | cmdctl_patch |\
        cmdctl_annotate | cmdctl_expose | cmdctl_scale | cmdctl_autoscale | cmdctl_taint | cmdctl_rollout_*)
            __cmdctl_get_resource
            return
            ;;
        cmdctl_logs | cmdctl_attach)
            __cmdctl_require_pod_and_container
            return
            ;;
        cmdctl_exec | cmdctl_port-forward | cmdctl_top_pod)
            __cmdctl_get_resource_pod
            return
            ;;
        cmdctl_rolling-update)
            __cmdctl_get_resource_rc
            return
            ;;
        cmdctl_cordon | cmdctl_uncordon | cmdctl_drain | cmdctl_top_node)
            __cmdctl_get_resource_node
            return
            ;;
        cmdctl_config_use-context)
            __cmdctl_config_get_contexts
            return
            ;;
        *)
            ;;
    esac
}
`
)

const RecommendedHomeDir = ".cmdctl"

var cfgFile string

var (
	bash_completion_flags = map[string]string{
		"namespace": "__cmdctl_get_namespaces",
		"context":   "__cmdctl_config_get_contexts",
		"cluster":   "__cmdctl_config_get_clusters",
		"user":      "__cmdctl_config_get_users",
	}
)

func NewCommand(f cmdutil.Factory, in io.Reader, out, err io.Writer) *cobra.Command {
	// Parent command to which all subcommands are added.
	cmds := &cobra.Command{
		Use:   "cmdctl",
		Short: "A microservices toolkit",
		Long: templates.LongDesc(`
		Microctl is a toolkit for microservice development. It helps you build future-proof application platforms and services..`),
		Run: runHelp,
		BashCompletionFunction: bashCompletionFunc,
	}

	groups := templates.CommandGroups{
		{
			Message: "Basic Commands:",
			Commands: []*cobra.Command{
				NewCmdInfo(f, out, err),
				NewCmdNew(f, out, err),
			},
		},
		{
			Message: "User Control Commands:",
			Commands: []*cobra.Command{
				NewCmdInit(),
				NewCmdAdd(f, out, err),
				NewCmdList(f, out, err),
			},
		},
		{
			Message: "Troubleshooting and Debugging Commands:",
			Commands: []*cobra.Command{
				NewCmdTest(f, out, err),
			},
		},
		{
			Message: "File Server Commands:",
			Commands: []*cobra.Command{
				NewCmdFinfo(f, out, err),
			},
		},
		{
			Message: "Commands Have Sub Commands:",
			Commands: []*cobra.Command{
				NewCmdTemplate(f, out, err),
			},
		},
	}
	groups.Add(cmds)
	templates.ActsAsRootCommand(cmds, []string{}, groups...)

	cmds.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is ./sreconfig.yaml)")
	cmds.PersistentFlags().BoolP("debug", "", false, "enable the debug mode")
	cobra.OnInitialize(initConfig)

	cmds.AddCommand(NewCmdVersion(out))
	cmds.AddCommand(NewCmdCompletion(out, ""))
	cmds.AddCommand(NewCmdOptions(out))
	cmds.AddCommand(NewCmdValidate(f, out))

	return cmds
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath(filepath.Join(homedir.HomeDir(), RecommendedHomeDir))
		viper.SetConfigName("cmdctl")
	}

	viper.SetConfigType("yaml")
	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvPrefix("CMDCTL")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.ReadInConfig()
}
