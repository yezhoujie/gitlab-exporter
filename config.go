package main

import (
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Gitlab struct {
		Url   string `yaml:"url"`
		Token string `yaml:"token"`
	} `yaml:"gitlab"`
	Oss struct {
		AccessKeyId  string `yaml:"accessKeyId"`
		AccessSecret string `yaml:"accessSecret"`
		Endpoint     string `yaml:"endpoint"`
		BucketName   string `yaml:"bucketName"`
	} `yaml:"oss"`
	KeepLocalBackup bool `yaml:"keepLocalBackup"`
}

func loadConfig(path string) Config {
	var config Config

	configFile, err := os.Open(path)
	if err != nil {
		log.Fatal("Cannot open config file: ", err)
	}
	defer configFile.Close()

	byteValue, _ := io.ReadAll(configFile)

	err = yaml.Unmarshal(byteValue, &config)
	if err != nil {
		log.Fatal("Cannot decode config YAML: ", err)
	}
	return config
}
