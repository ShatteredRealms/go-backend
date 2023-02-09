package main

import (
    "errors"
    "fmt"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "os"
)

func SetupConfig() error {
    setupDefaults()
    err := setupConfig()
    if err != nil {
        return err
    }

    return viper.ReadInConfig()
}

func setupDefaults() {
    viper.SetDefault("dir", ".")
}

func setupConfig() error {
    home, err := os.UserHomeDir()
    cobra.CheckErr(err)

    name := "benchmark"
    configType := "yaml"
    location := fmt.Sprintf("%s/.sro", home)
    configFile := fmt.Sprintf("%s/%s.%s", location, name, configType)
    viper.SetConfigName(name)
    viper.SetConfigType(configType)
    viper.AddConfigPath(location)

    if _, err = os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
        err = os.MkdirAll(location, 0760)
        _, err = os.Create(configFile)
        if err != nil {
            return err
        }
    }

    return nil
}
