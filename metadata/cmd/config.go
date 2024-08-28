package main

type Config struct {
	API ApiConfig `yaml:"api"`
}

type ApiConfig struct {
	Port int `yaml:"port"`
}
