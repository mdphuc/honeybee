/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
  "fmt"
  "os"
  "os/exec"
  "github.com/spf13/cobra"
  "log"
  homedir "github.com/mitchellh/go-homedir"
  "github.com/spf13/viper"
  "strings"
  // "bytes"
)


var cfgFile string


// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
  Use:   "beaver",
  Short: "Set up remote development environment",
  Long: `biever v1.0
    
Set up remote development environment in isolated
enivronment like docker or proxy server to use
remote machine or cloud machine as development environment
  `,
  CompletionOptions: cobra.CompletionOptions{
    DisableDefaultCmd: true,
  },
  Version: "1.0",
  // Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
  if err := rootCmd.Execute(); err != nil {
    initLog(fmt.Sprintf("%v", err))
    os.Exit(1)
  }
}

func init() {
  cobra.OnInitialize(initConfig)
}

func initLog(err string) {
  mkdir_run := exec.Command("mkdir", "./log")

  _ = mkdir_run.Run()
  err1 := strings.Split(err, "\n")[0]

  logFile, err_log := os.OpenFile("./log/app.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
  log.SetOutput(logFile)
  log.SetFlags(log.Lshortfile | log.LstdFlags | log.LUTC) 

  log.Println("UTC: " + fmt.Sprintf("%v",err1))  


  if err_log != nil{
    log.Println("UTC: " + fmt.Sprintf("%v",err_log))
  }
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
  if cfgFile != "" {
    // Use config file from the flag.
    viper.SetConfigFile(cfgFile)
  } else {
    // Find home directory.
    home, err := homedir.Dir()
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

    // Search config in home directory with name ".beaver" (without extension).
    viper.AddConfigPath(home)
    viper.SetConfigName(".beaver")
  }

  viper.AutomaticEnv() // read in environment variables that match

  // If a config file is found, read it in.
  if err := viper.ReadInConfig(); err == nil {
    fmt.Println("Using config file:", viper.ConfigFileUsed())
  }
}

