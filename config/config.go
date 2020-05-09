package config

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

type EmpConfig struct{
	Start int8
	CreateIndex bool
	EthIp string
	MongoDBIp string
	DatabaseName string
}

var (
	EmpApp = NewEmpApp()
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "emp",
	Short: "A test emp",
	Long:  `Configure basic setting`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if EmpApp.Start == 0 {
			cmd.Help()
			os.Exit(-1)
		}
		if EmpApp.CreateIndex == true {
			fmt.Printf("Create index\n")
		}
		if EmpApp.DatabaseName != ""{
			fmt.Printf("mongodb database name %s\n", EmpApp.DatabaseName)
		}
	},
}

func NewEmpApp() *EmpConfig{
	return &EmpConfig{

	}
}

func (e *EmpConfig)EmpSetting() {
	v := viper.New()
	v.AddConfigPath("./config")
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	// viper解析配置文件
	if err := v.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	if e.EthIp == "" {
		e.EthIp = fmt.Sprintf("http://%s:%s",v.Get("localETH.host"),v.Get("localETH.port"))
		fmt.Println(e.EthIp)
	}
	if e.DatabaseName == "" {
		e.DatabaseName = fmt.Sprintf("%s",v.Get("database.mongodb.dbName"))
		fmt.Println(e.DatabaseName)
	}
	if e.MongoDBIp == "" {
		e.MongoDBIp = fmt.Sprintf("mongodb://%s:%s",v.Get("database.labMongodb.host"),v.Get("database.labMongodb.port"))
		fmt.Println(e.MongoDBIp)
	}
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		panic(err)
	}
}

func init() {
	RootCmd.Flags().Int8VarP(&EmpApp.Start, "start", "s",0, "start mod")
	RootCmd.Flags().BoolVarP(&EmpApp.CreateIndex, "index", "i",false, "create index default:false")
	RootCmd.Flags().StringVarP(&EmpApp.EthIp, "ethIp", "e","", "ethereum node ip")
	RootCmd.Flags().StringVarP(&EmpApp.DatabaseName, "database", "d","", "mongodb database name")
	RootCmd.Flags().StringVarP(&EmpApp.MongoDBIp, "mongodbIp", "m","", "mongodb ip")
}






