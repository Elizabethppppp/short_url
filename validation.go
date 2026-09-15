package main

import (
	"errors"
	"fmt"
	"strings"
)

const maxLenURL = 2048

func splitProtocol(url string) (protocol, domen string, ok bool) {

	protocol, domen, ok = strings.Cut(url, "://")
	if !ok {
		return "", "", false
	}
	return strings.ToLower(protocol), domen, true
}

func splitDomain(url string) (domain, path string) {
	i := strings.IndexAny(url, "/?#")
	if i == -1 {
		return url, ""
	}
	domain, path = url[:i], url[i:]
	return domain, path
}

func splitPort(url string) (domain, port string, ok bool) {
	domain, port, ok = strings.Cut(url, ":")
	if !ok {
		return url, "", false
	}
	return domain, port, true
}

func validateUrl(url string) (string, error) {
	url = strings.TrimSpace(url)
	if url == "" {
		return "", errors.New("empty url")
	}
	if len(url) > maxLenURL {
		return "", errors.New("url too long")
	}
	if strings.ContainsAny(url, " \t\r\n\"'<>\\") {
		return "", errors.New("url contains whitespace")
	}

	protocol, domen, ok := splitProtocol(url)
	if !ok {
		return "", errors.New("invalid protocol")
	}
	if protocol != "http" && protocol != "https" {
		return "", errors.New("invalid protocol")
	}

	domainURL, path := splitDomain(domen)
	if domainURL == "" {
		return "", errors.New("invalid domain")
	}
	for _, ch := range path {
		if ch < 0x21 || ch > 0x7E {
			return "", fmt.Errorf("invalid character %q in path", ch)
		}
		switch ch {
		case '"', '\'', '<', '>', '\\':
			return "", fmt.Errorf("invalid character %q in path", ch)
		}
	}

	domain, port, hasPort := splitPort(domainURL)
	if domain == "" {
		return "", errors.New("invalid domain")
	}
	if domain != "localhost" && !strings.Contains(domain, ".") {
		return "", errors.New("invalid domain")
	}
	if hasPort {
		if port == "" {
			return "", errors.New("invalid port")
		}
		n := 0
		for _, ch := range port {
			if ch < '0' || ch > '9' {
				return "", fmt.Errorf("invalid character %q in port", ch)
			}
			n = n*10 + int(ch-'0')
			if n > 65535 {
				return "", fmt.Errorf("port out of range: %q", port)
			}
		}
		if n == 0 {
			return "", errors.New("port must be > 0")
		}
	}

	for _, ch := range domain {
		valid := false
		switch {
		case ch >= 'a' && ch <= 'z',
			ch >= 'A' && ch <= 'Z',
			ch >= '0' && ch <= '9',
			ch == '-', ch == '.':
			valid = true
		}
		if !valid {
			return "", fmt.Errorf("invalid character %q in host", ch)
		}

	}

	return url, nil
}
