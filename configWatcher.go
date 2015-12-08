package main

import(
  "github.com/spf13/viper"
  "fmt"

)


func ReadConfig() {
  viper.SetConfigType("json");
  viper.SetConfigName("config") // name of config file (without extension)
  viper.AddConfigPath("$HOME/whalee-core")
  viper.AddConfigPath(".")
  err := viper.ReadInConfig();
  if err != nil {
    panic(fmt.Errorf("Fatal error config file: %s \n", err));
  }
}
