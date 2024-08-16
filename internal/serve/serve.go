package serve

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/dozyio/tls-proxy/internal/config"
	logging "github.com/ipfs/go-log/v2"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Server struct {
	config *config.Config
	logger *logging.ZapEventLogger
}

func New(cfg *config.Config, logger *logging.ZapEventLogger) *Server {
	return &Server{
		config: cfg,
		logger: logger,
	}
}

func (s *Server) Run() {
	s.logger.Info("Starting server")

	cert, err := os.ReadFile(s.config.Cert)
	if err != nil {
		s.logger.Error("Failed to read cert file", err)
		os.Exit(1)
	}

	key, err := os.ReadFile(s.config.Key)
	if err != nil {
		s.logger.Error("Failed to read key file", err)
		os.Exit(1)
	}

	tlsCert, err := tls.X509KeyPair(cert, key)
	if err != nil {
		s.logger.Error("Failed to create X509 key pair", err)
		os.Exit(1)
	}

	port := "443"

	if s.config.Target.Port == "" {
		if s.config.Target.Scheme == "http" {
			port = "80"
		}
	} else {
		port = s.config.Target.Port
	}

	target, err := net.ResolveTCPAddr("tcp", s.config.Target.Host+":"+port)
	if err != nil {
		s.logger.Error("Failed to resolve target address", err)
		os.Exit(1)
	}

	listen, err := net.Listen("tcp", s.config.Listen)
	if err != nil {
		s.logger.Error("Failed to listen on address", err)
		os.Exit(1)
	}

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = s.config.Target.Scheme
			req.URL.Host = target.String()
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Skip verification for upstream connections
			},
		},
	}

	tlsConfig := &tls.Config{
		Certificates:             []tls.Certificate{tlsCert},
		MinVersion:               tls.VersionTLS13,
		CurvePreferences:         []tls.CurveID{tls.CurveP256, tls.X25519},
		PreferServerCipherSuites: true,
		InsecureSkipVerify:       true, //nolint:gosec // ignore for self signed
	}

	h2cServer := &http2.Server{}
	server := &http.Server{
		Handler:           h2c.NewHandler(proxy, h2cServer),
		TLSConfig:         tlsConfig,
		ReadHeaderTimeout: 3 * time.Second,
	}

	tlsListener := tls.NewListener(listen, tlsConfig)

	err = server.Serve(tlsListener)
	if err != nil {
		s.logger.Error("Failed to start server", err)
		os.Exit(1)
	}
}
