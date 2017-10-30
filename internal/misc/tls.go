package misc

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
)

func LoadTLS(certFile string, keyFile string, caFile string, insecureSkipTLSVerify bool) (*tls.Config, error) {
	certificate, err := tls.LoadX509KeyPair(
		certFile,
		keyFile,
	)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	if caFile != "" {
		bs, err := ioutil.ReadFile(caFile)
		if err != nil {
			return nil, err
		}

		if !certPool.AppendCertsFromPEM(bs) {
			return nil, errors.New("failed to append client certs")
		}
	}
	TLS := &tls.Config{
		ClientAuth:         tls.RequireAndVerifyClientCert,
		InsecureSkipVerify: insecureSkipTLSVerify,
		Certificates:       []tls.Certificate{certificate},
		ClientCAs:          certPool,
	}
	return TLS, nil
}

/*
	tlscfg := &tls.Config{
		MinVersion:         tls.VersionTLS10,
		InsecureSkipVerify: yc.InsecureSkipTLSVerify,
		RootCAs:            cp,
	}
	if cert != nil {
		tlscfg.Certificates = []tls.Certificate{*cert}
	}
	cfg.TLS = tlscfg

*/
