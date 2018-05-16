package util

import (
	"encoding/base64"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type FileServer struct {
	Server   string
	Timeout  int
	Username string
	Password string
}

type Factory struct {
	flags *pflag.FlagSet
}

func NewFactory() Factory {
	flags := pflag.NewFlagSet("", pflag.ContinueOnError)
	f := Factory{
		flags: flags,
	}

	return f
}

func (f *Factory) FlagSet() *pflag.FlagSet {
	return f.flags
}

// TODO: We need to filter out stuff like secrets.
func (f *Factory) Command() string {
	if len(os.Args) == 0 {
		return ""
	}
	base := filepath.Base(os.Args[0])
	args := append([]string{base}, os.Args[1:]...)
	return strings.Join(args, " ")
}

func (f *Factory) BindFlags(flags *pflag.FlagSet) {
	// Merge factory's flags
	flags.AddFlagSet(f.flags)
}

func (f *Factory) BindExternalFlags(flags *pflag.FlagSet) {
	// any flags defined by external projects (not part of pflags)
	flags.AddGoFlagSet(flag.CommandLine)
}

func (f *Factory) Auth() string {
	return base64.StdEncoding.EncodeToString([]byte(f.FileServer().Username + ":" + f.FileServer().Password))
}

func (f *Factory) Gorequest() *gorequest.SuperAgent {
	request := gorequest.New()
	request.DoNotClearSuperAgent = true
	return request.Timeout(time.Duration(f.FileServer().Timeout)*time.Second).
		Set("Content-Type", "application/json").
		Set("Username", "cc")
}

func (f *Factory) FileServer() *FileServer {
	return &FileServer{
		Server:   viper.GetString("fileserver.server"),
		Timeout:  viper.GetInt("fileserver.timeout"),
		Username: viper.GetString("fileserver.username"),
		Password: viper.GetString("fileserver.password"),
	}
}
