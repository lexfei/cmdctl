package cmd

import (
	"fmt"
	"io"
	"reflect"
	"strconv"

	"cmdctl/cmd/templates"
	cmdutil "cmdctl/cmd/util"
	"cmdctl/pkg/i18n"
	"cmdctl/util"

	"github.com/likexian/host-stat-go"
	"github.com/spf13/cobra"
)

type Info struct {
	HostName  string
	IPAddress string
	OSRelease string
	CPUCore   uint64
	MemTotal  string
	MemFree   string
}

var (
	infoExample = templates.Examples(i18n.T(`
		# Print the host information
		cmdctl info

		# Specify a server password
		cmdctl info -p newpass

		# Print details
		cmdctl info -d`))
)

func NewCmdInfo(f cmdutil.Factory, out io.Writer, cmdErr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "info",
		Short:   i18n.T("Print the host information"),
		Long:    "Print the host information",
		Example: infoExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(RunInfo(f, out, cmdErr, cmd, args))
		},
	}

	cmd.Flags().StringP("passwd", "p", "", "Specify the server password.")
	cmd.Flags().BoolP("detail", "d", false, "Print details.")

	return cmd
}

func RunInfo(f cmdutil.Factory, out io.Writer, cmdErr io.Writer, cmd *cobra.Command, args []string) error {
	passwd := cmdutil.GetFlagString(cmd, "passwd")
	detail := cmdutil.GetFlagBool(cmd, "detail")
	if detail {
		fmt.Printf("%12s %v\n", "OptionValue"+":", fmt.Sprintf("%v:%s", detail, passwd))
	}

	var info Info
	host_info, err := host_stat.GetHostInfo()
	if err != nil {
		return fmt.Errorf("get host info failed!error:%v", err)
	}

	info.HostName = host_info.HostName
	info.OSRelease = host_info.Release + " " + host_info.OSBit

	mem_stat, err := host_stat.GetMemStat()
	if err != nil {
		return fmt.Errorf("get mem stat failed!error:%v", err)
	}
	info.MemTotal = strconv.FormatUint(mem_stat.MemTotal, 10) + "M"
	info.MemFree = strconv.FormatUint(mem_stat.MemFree, 10) + "M"

	cpu_stat, err := host_stat.GetCPUInfo()
	if err != nil {
		return fmt.Errorf("get cpu stat failed!error:%v", err)
	}
	info.CPUCore = cpu_stat.CoreCount

	info.IPAddress = util.GetLocalAddress()

	s := reflect.ValueOf(&info).Elem()
	typeOfInfo := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		v := fmt.Sprintf("%v", f.Interface())
		if v != "" {
			fmt.Printf("%12s %v\n", typeOfInfo.Field(i).Name+":", f.Interface())
		}
	}

	return nil
}
