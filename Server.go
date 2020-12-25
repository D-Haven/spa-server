package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"path"
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

		// Default is the default path to redirect responses to
		Default string `yaml:"default"`

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
		config.Server.Default = "/index.html"
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

	defer CheckError(file.Close())

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	if config.Server.Default[:1] != "/" {
		config.Server.Default = "/" + config.Server.Default
	}

	return config, nil
}

func ValidateConfigPath(path string) error {
	fileInfo, err := os.Stat(path)
	CheckError(err)

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

	if len(config.Server.Default) <= 1 {
		config.Server.Default = "/index.html"
		log.Println("No default path provided, using '/index.html'")
	}

	if config.Server.Default[:1] != "/" {
		config.Server.Default = "/" + config.Server.Default
	}

	defaultFile := path.Join(config.Server.SitePath, config.Server.Default[1:len(config.Server.Default)])
	fileInfo, err = os.Stat(defaultFile)
	CheckError(err)

	if fileInfo.IsDir() {
		return fmt.Errorf("no default file exists, invalid configuration: '%s'", defaultFile)
	}

	return nil
}

func main() {
	config, err := ReadConfig("./config.yaml")
	CheckError(err)

	CheckError(ValidateConfig(config))

	fmt.Printf("Starting server at port %s\n", config.Server.Port)
	fmt.Printf("    default path: %s\n", config.Server.Default)

	fileServer := http.FileServer(http.Dir(config.Server.SitePath))
	redirectDefault := NotFoundRedirectHandler(config.Server.Default, fileServer)
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
