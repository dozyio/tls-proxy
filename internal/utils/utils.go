package utils

import (
	"errors"
	"net"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type DomainInfo struct {
	Scheme string
	Host   string
	Port   string
}

func ParseDomainWithScheme(s string) (*DomainInfo, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, errors.New("invalid URL format")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, errors.New("invalid scheme, must be http or https")
	}

	domainRegexp := `^(([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}|(\d{1,3}\.){3}\d{1,3})$`

	domain := u.Hostname()
	port := u.Port()

	matchedDomain, err := regexp.MatchString(domainRegexp, domain)
	if err != nil || !matchedDomain {
		return nil, errors.New("invalid domain")
	}

	return &DomainInfo{
		Scheme: u.Scheme,
		Host:   domain,
		Port:   port,
	}, nil
}

func IsValidIPPort(s string) bool {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return false
	}

	ip := parts[0]
	port := parts[1]

	if net.ParseIP(ip) == nil {
		return false
	}

	if _, err := net.LookupPort("tcp", port); err != nil {
		return false
	}

	return true
}

func IsReadableFile(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	_, err = file.Stat()

	return err == nil
}
