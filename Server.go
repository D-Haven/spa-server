/*
 * Copyright (c) 2020. D-Haven
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
)

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Config struct {
	Server struct {
		// Host is the server host name
		Host string `yaml:"host"`

		// Port is the local machine TCP Port to bind the HTTP Server to
		Port string `yaml:"port"`

		// Compress is a flag for if results should be GZip compressed
		Compress bool `yaml:"compress"`

		// SitePath string is the location to serve up static files
		SitePath string `yaml:"sitePath"`

		// TLS provides the TLS tuning configuration
		TLS struct {
			CertFile string `yaml:"certificate"`
			KeyFile  string `yaml:"key"`
		} `yaml:"tls"`
	} `yaml:"server"`
}

func ReadConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	if err := ValidateConfigPath(configPath); err != nil {
		log.Println("No config file, will use default settings")

		config.Server.Port = "8443"
		config.Server.Compress = false
		config.Server.SitePath = "./www"
		config.Server.TLS.CertFile = "./server.crt"
		config.Server.TLS.KeyFile = "./server.key"

		return config, nil
	}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	defer func() { CheckError(file.Close()) }()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func ValidateConfigPath(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		// File doesn't exist will use defaults
		return err
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}

	return nil
}

func ValidateConfig(config *Config) error {
	fileInfo, err := os.Stat(config.Server.SitePath)
	CheckError(err)

	if !fileInfo.IsDir() {
		return fmt.Errorf("'%s' is not a directory, cannot start server", config.Server.SitePath)
	}

	return nil
}

func ShowLogo() {
	logo := figure.NewFigure("SPA server", "trek", true)
	logo.Print()
	fmt.Println()
}

func main() {
	config, err := ReadConfig("./config.yaml")
	CheckError(err)
	CheckError(ValidateConfig(config))

	ShowLogo()
	fmt.Printf("    Starting server at port %s\n", config.Server.Port)

	fileServer := http.FileServer(http.Dir(config.Server.SitePath))
	redirectDefault := NotFoundRedirectHandler("/", fileServer)
	handler := redirectDefault

	if config.Server.Compress {
		handler = GzipHandler(redirectDefault)
	}

	server := &http.Server{
		Addr:    config.Server.Host + ":" + config.Server.Port,
		Handler: handler,
	}

	CheckError(server.ListenAndServeTLS(config.Server.TLS.CertFile, config.Server.TLS.KeyFile))
}
