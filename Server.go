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
	"d-haven.org/spa-server/handlers"
	"d-haven.org/spa-server/version"
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/heptiolabs/healthcheck"
	"golang.org/x/net/context"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const appName = "SPA server"

type Config struct {
	Server struct {
		// Host is the server host name
		Host string `yaml:"host"`

		// Port is the local machine TCP Port to bind the HTTP Server to
		Port string `yaml:"port"`

		// SitePath string is the location to serve up static files
		SitePath string `yaml:"sitePath"`

		// TLS provides the TLS tuning configuration
		TLS struct {
			CertFile string `yaml:"certificate"`
			KeyFile  string `yaml:"key"`
		} `yaml:"tls"`
	} `yaml:"server"`
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ReadConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	if err := ValidateConfigPath(configPath); err != nil {
		log.Println("No config file, will use default settings")

		config.Server.Port = "8080"
		config.Server.SitePath = "./www"

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

	CheckError(ValidateOptionalFile(config.Server.TLS.KeyFile))
	CheckError(ValidateOptionalFile(config.Server.TLS.CertFile))

	return nil
}

func ValidateOptionalFile(path string) error {
	if len(path) > 0 {
		fileInfo, err := os.Stat(path)
		CheckError(err)

		if fileInfo.IsDir() {
			return fmt.Errorf("'%s' is directory, not a file", path)
		}
	}

	return nil
}

func ShowLogo() {
	logo := figure.NewFigure(appName, "trek", true)
	figure.Write(log.Writer(), logo)
	_, err := io.WriteString(log.Writer(), "\n")
	CheckError(err)
}

func main() {
	ShowLogo()
	log.Printf("Version: %s Commit: %s Timestamp: %s\n", version.Release, version.Commit, version.BuildTime)

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "--hostname lookup error: " + err.Error()
	}

	log.Printf("Server host name: %s", hostname)
	log.Printf("Configuring %s...", appName)

	config, err := ReadConfig("./config.yaml")
	CheckError(err)
	CheckError(ValidateConfig(config))

	fileServer := http.FileServer(http.Dir(config.Server.SitePath))
	redirectDefault := handlers.NotFoundRedirectHandler(fileServer)
	handler := redirectDefault

	health := healthcheck.NewHandler()
	health.AddLivenessCheck("go-routinethreshold", healthcheck.GoroutineCountCheck(100))

	multiplexHandler := http.NewServeMux()
	multiplexHandler.Handle("/ready", health)
	multiplexHandler.Handle("/live", health)
	multiplexHandler.Handle("/", handler)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	server := &http.Server{
		Addr:    config.Server.Host + ":" + config.Server.Port,
		Handler: multiplexHandler,
	}

	go func() {
		useTLS := len(config.Server.TLS.CertFile) != 0 && len(config.Server.TLS.KeyFile) != 0

		if useTLS {
			CheckError(server.ListenAndServeTLS(config.Server.TLS.CertFile, config.Server.TLS.KeyFile))
		} else {
			CheckError(server.ListenAndServe())
		}
	}()

	log.Printf("%s started, listening on %s", appName, server.Addr)
	log.Print("K8s Readiness Check: /ready")
	log.Print("K8s Liveness Check: /live")

	killSignal := <-interrupt
	switch killSignal {
	case os.Interrupt:
		log.Print("Received OS Interrupt")
	case syscall.SIGTERM:
		log.Print("Received Termination Signal")
	}

	log.Printf("%s is shutting down...", appName)
	CheckError(server.Shutdown(context.Background()))
}
