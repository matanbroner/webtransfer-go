package main

import (
	"context"
	"crypto/tls"
	"log"

	"github.com/lucas-clemente/quic-go"
)

// Configuration for the server
type Config struct {
	Host                 string
	Port                 string
	CertificatePath      string
	KeyPath              string
	AllowedAccessOrigins []string
}

type QuicServer struct {
	config Config
}

func (server *QuicServer) TLSConfig() *tls.Config {
	cert, err := tls.LoadX509KeyPair(server.config.CertificatePath, server.config.KeyPath)
	if err != nil {
		log.Fatal(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		// We identify ourselves as a HTTP/3 based client
		NextProtos: []string{"h3"},
	}
}

func (server *QuicServer) handleSession(session quic.Session) {
	// TODO: Implement
	log.Println("handle session")
}

func (server *QuicServer) Start() error {
	var addr string = server.config.Host + ":" + server.config.Port
	// A Listener for incoming QUIC connections
	listener, err := quic.ListenAddr(addr, server.TLSConfig(), nil)
	log.Printf("Listening on %s", addr)
	if err != nil {
		return err
	}
	for {
		session, err := listener.Accept(context.Background())
		if err != nil {
			return err
		}
		go func() {
			defer func() {
				_ = session.CloseWithError(0, "Session Closed")
				log.Println("close session")
			}()
			server.handleSession(session)
		}()
	}
}

func main() {
	config := Config{
		Host:                 "0.0.0.0",
		Port:                 "4433",
		CertificatePath:      "quic_cert.pem",
		KeyPath:              "quic_key.pem",
		AllowedAccessOrigins: []string{"localhost", "googlechrome.github.io"},
	}
	server := QuicServer{config: config}
	server.Start()
}
