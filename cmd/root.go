/*
Copyright © 2020 Shiun

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
  "bytes"
  "demo_project/config"
  "fmt"
  "io/ioutil"
  "os"
  "github.com/spf13/cobra"
  "reflect"
  "strings"

  "github.com/spf13/viper"
  log "github.com/sirupsen/logrus"
)


var cfgFile string
var version string


// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
  Use:   "demoPro",
  Short: "Demo Project",
  Long: ``,
  RunE: run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func init() {
  cobra.OnInitialize(initConfig)

  // Here you will define your flags and configuration settings.
  // Cobra supports persistent flags, which, if defined here,
  // will be global for your application.

  rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.demoPro.yaml)")


  // Cobra also supports local flags, which will only run
  // when this action is called directly.
  rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


// initConfig reads in config file and ENV variables if set.
func initConfig() {
  if cfgFile != "" {
    b, err := ioutil.ReadFile(cfgFile)
    if err != nil {
      log.WithError(err).WithField("config", cfgFile).Fatal("error loading config file")
    }
    viper.SetConfigType("toml")
    if err := viper.ReadConfig(bytes.NewBuffer(b)); err != nil {
      log.WithError(err).WithField("config", cfgFile).Fatal("error loading config file")
    }
  } else {
    viper.SetConfigName("demo-project")
    viper.AddConfigPath(".")
    viper.AddConfigPath("./config")
    viper.AddConfigPath("$HOME/.config/demo-project")
    viper.AddConfigPath("/etc/demo-project")
    if err := viper.ReadInConfig(); err != nil {
      switch err.(type) {
      case viper.ConfigFileNotFoundError:
        log.Warning("No configuration file found, using defaults.")
      default:
        log.WithError(err).Fatal("read configuration file error")
      }
    }
  }

  viperBindEnvs(config.Conf)

  if err := viper.Unmarshal(&config.Conf); err != nil {
    log.WithError(err).Fatal("unmarshal config error")
  }
}

func viperBindEnvs(iface interface{}, parts ...string) {
  ifv := reflect.ValueOf(iface)
  ift := reflect.TypeOf(iface)
  for i := 0; i < ift.NumField(); i++ {
    v := ifv.Field(i)
    t := ift.Field(i)
    tv, ok := t.Tag.Lookup("mapstructure")
    if !ok {
      tv = strings.ToLower(t.Name)
    }
    if tv == "-" {
      continue
    }

    switch v.Kind() {
    case reflect.Struct:
      viperBindEnvs(v.Interface(), append(parts, tv)...)
    default:
      key := strings.Join(append(parts, tv), ".")
      if err := viper.BindEnv(key); err != nil {
        log.WithError(err).Error("BindEnv error")
      }
    }
  }
}