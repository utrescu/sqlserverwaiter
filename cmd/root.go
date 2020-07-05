package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version = "1.0.3"
)

var (

	// Debug shows info
	Debug bool

	// Server is the database server host
	Server string
	// Port is the database server port
	Port int
	// User is the database server username
	User string
	// Password is the database server password
	Password string
	// Database is the database name
	Database string
	// Timeout time to wait for database ready
	Timeout time.Duration
)

// RootCmd is the default command
var RootCmd = &cobra.Command{
	Use:     "sqlserverwaiter",
	Short:   "Waits until Sql Server is ready or timeout",
	Long:    "This application waits until a Microsoft SQL Server accepts connections for an specified time",
	Run:     func(cmd *cobra.Command, args []string) {},
	Version: version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Stop execution if help or version are requested
	var helpWanted = RootCmd.Flags().Lookup("help")
	if helpWanted.Changed {
		os.Exit(0)
	}

	if RootCmd.Flags().Lookup("version").Changed {
		os.Exit(0)
	}

}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().BoolVar(&Debug, "debug", false, "Show debug info")

	RootCmd.PersistentFlags().StringVar(&Server, "host", "localhost", "Database server host")
	RootCmd.PersistentFlags().IntVar(&Port, "port", 1433, "Numeric database port")
	RootCmd.PersistentFlags().StringVarP(&Database, "database", "d", "master", "Database name to connect: ex. BikeStores")

	RootCmd.PersistentFlags().StringVarP(&User, "user", "u", "sa", "Database user")
	RootCmd.PersistentFlags().StringVarP(&Password, "password", "p", "X1nGuXunG1", "Database user password")

	RootCmd.PersistentFlags().DurationVarP(&Timeout, "timeout", "t", 30*time.Second, "Time to wait: 30s, 2m, ...")

	viper.BindPFlag("host", RootCmd.Flags().Lookup("host"))
	viper.BindPFlag("port", RootCmd.Flags().Lookup("port"))
	viper.BindPFlag("database", RootCmd.Flags().Lookup("database"))

	viper.BindPFlag("user", RootCmd.Flags().Lookup("user"))
	viper.BindPFlag("password", RootCmd.Flags().Lookup("user"))

	viper.BindPFlag("timeout", RootCmd.Flags().Lookup("timeout"))

}

// initConfig reads ENV variables if set. (SQLW_*)
func initConfig() {
	viper.SetEnvPrefix("sqlw")
	viper.AutomaticEnv()
}
