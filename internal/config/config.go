package config

import (
	"fmt"
	"os"

	"github.com/dozyio/tls-proxy/internal/utils"
	"github.com/spf13/cobra"
)

type Config struct {
	Target   *utils.DomainInfo
	Listen   string
	Cert     string
	Key      string
	LogLevel string
}

func New(cmd *cobra.Command, _ []string) *Config {
	// target
	if cmd.Flag("target").Value.String() == "" {
		fmt.Println("target is required")
		cmd.Help()
		os.Exit(1)
	}

	domainParts, err := utils.ParseDomainWithScheme(cmd.Flag("target").Value.String())
	if err != nil {
		fmt.Printf("target is invalid - %s\n", err.Error())
		cmd.Help()
		os.Exit(1)
	}

	// listen
	if cmd.Flag("listen").Value.String() == "" {
		fmt.Println("listen is required")
		cmd.Help()
		os.Exit(1)
	}

	if !utils.IsValidIPPort(cmd.Flag("listen").Value.String()) {
		fmt.Println("listen is not valid - ip:port")
		cmd.Help()
		os.Exit(1)
	}

	// cert path
	if cmd.Flag("cert").Value.String() == "" {
		fmt.Println("cert is required")
		cmd.Help()
		os.Exit(1)
	}

	if !utils.IsReadableFile(cmd.Flag("cert").Value.String()) {
		fmt.Println("key is not readable")
		os.Exit(1)
	}

	// key path
	if cmd.Flag("key").Value.String() == "" {
		fmt.Println("key is required")
		cmd.Help()
		os.Exit(1)
	}

	if !utils.IsReadableFile(cmd.Flag("key").Value.String()) {
		fmt.Println("key is not readable")
		os.Exit(1)
	}

	return &Config{
		Target: domainParts,
		Listen: cmd.Flag("listen").Value.String(),
		Cert:   cmd.Flag("cert").Value.String(),
		Key:    cmd.Flag("key").Value.String(),
	}
}
