package cmd

import (
	"database/sql"
	"fmt"

	"cmdctl/cmd/templates"
	cmdutil "cmdctl/cmd/util"
	"cmdctl/model"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	initExample = templates.Examples(`
		# Init db
		cmdctl init
		
		# Drop db first && init
		cmdctl init -f
		
		`)
)

func NewCmdInit() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Short:   "Init database",
		Long:    "Init database",
		Example: initExample,
		Run: func(cmd *cobra.Command, args []string) {
			//cmdutil.CheckErr(validateArgs(cmd, args))
			cmdutil.CheckErr(RunInit(cmd, args))
			return
		},
		Aliases: []string{},
	}

	cmd.Flags().BoolP("force", "f", false, "Drop table if exists")

	return cmd
}

func RunInit(cmd *cobra.Command, args []string) error {
	force := cmdutil.GetFlagBool(cmd, "force")

	if err := Createdb(force); err != nil {
		return err
	}

	if err := Createtb(); err != nil {
		return err
	}

	fmt.Println("sync db end, please reopen app again")
	return nil
}

//创建数据库
func Createdb(force bool) error {
	db_type := "mysql" // current only support mysql
	db_addr := viper.GetString("db.addr")
	db_user := viper.GetString("db.username")
	db_pass := viper.GetString("db.password")
	db_name := viper.GetString("db.name")

	var dns string
	var sqlstring, sql1string string

	dns = fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8", db_user, db_pass, db_addr)
	sql1string = fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", db_name)
	sqlstring = fmt.Sprintf("CREATE DATABASE if not exists `%s` CHARSET utf8 COLLATE utf8_general_ci", db_name)

	db, err := sql.Open(db_type, dns)
	if err != nil {
		return err
	}
	defer db.Close()

	if force {
		fmt.Println(sql1string)
		if _, err := db.Exec(sql1string); err != nil {
			return err
		}
	}

	if _, err := db.Exec(sqlstring); err != nil {
		return err
	}

	fmt.Printf("database %s created\n", db_name)
	return nil
}

func Createtb() error {
	db := model.GetSelfDB()
	defer db.Close()
	if err := db.AutoMigrate(&model.UserModel{}).Error; err != nil {
		return err
	}

	return nil
}
