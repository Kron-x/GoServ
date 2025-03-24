package main

import (
    "encoding/json"
    "log"
    "io/ioutil"
)


type Config struct {
    Port      string `json:"port"`
    ImagesDir string `json:"images_dir"`
    LogFile   string `json:"log_file"`
}

func LoadConfig() Config {
    file, err := ioutil.ReadFile("config.json")
    if err != nil {
        log.Fatalf("Failed to read config file: %v", err)
    }
    var config Config
    err = json.Unmarshal(file, &config)
    if err != nil {
        log.Fatalf("Failed to parse config file: %v", err)
    }
    return config
}