package main

import (
    "encoding/json"
    "os"
    "io/ioutil"
)

type ircConfig struct {
    Server string
    Port float64
}

type config struct {
    IRC *ircConfig
    Nick string
    NickServPassword string
    Channels []string
}

var Config *config = new(config) 

func loadConfig(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }

    bytes, err := ioutil.ReadAll(file)
    if err != nil {
        return err
    }

    err = json.Unmarshal(bytes, Config)
    return err
}
