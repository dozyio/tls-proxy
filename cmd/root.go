package cmd

import (
	"os"

	"github.com/dozyio/tls-proxy/internal/config"
	"github.com/dozyio/tls-proxy/internal/serve"
	logging "github.com/ipfs/go-log/v2"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tls-proxy",
	Short: "A TLS proxy",
	Long:  `A TLS proxy that can be used to forward traffic to a local or remote server over TLS.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New(cmd, args)

		logger := logging.Logger("tls-proxy")
	_:
		logging.SetLogLevel("*", cfg.LogLevel)

		s := serve.New(cfg, logger)
		s.Run()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("target", "t", "", "Where to forward traffic to. Can be a local address or a remote address.")
	rootCmd.Flags().StringP("listen", "l", ":443", "The address to listen on.")
	rootCmd.Flags().StringP("cert", "c", "", "The path to the certificate to use.")
	rootCmd.Flags().StringP("key", "k", "", "The path to the key to use.")
	rootCmd.Flags().StringP("loglevel", "v", "DEBUG", "Log level to use DEBUG,INFO,WARN,ERROR,PANIC,PANIC,FATAL")
}
