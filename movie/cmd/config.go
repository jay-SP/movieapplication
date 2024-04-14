package main

type Config struct {
	API apiConfig `yaml:"api"`
}

type apiConfig struct {
	Port int `yaml:"port"`
}
