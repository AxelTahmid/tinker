package main

import (
	"context"
	"crypto/tls"
	"log"

	"github.com/AxelTahmid/golang-starter/config"
	"github.com/AxelTahmid/golang-starter/db"
	"github.com/AxelTahmid/golang-starter/internal/server"
)

func main() {
	ctx := context.Background()

	conf := config.New()

	serverTLSCert, err := tls.LoadX509KeyPair(conf.Server.TLSCertPath, conf.Server.TLSKeyPath)
	if err != nil {
		log.Fatalf("Error loading certificate and key file: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverTLSCert},
	}

	dbconn, err := db.CreatePool(ctx, conf.Database)
	if err != nil {
		log.Fatalf("Db Connection Failed: %v", err)
	}

	err = dbconn.Ping(ctx)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	server := server.NewServer(conf, dbconn, tlsConfig)
	server.Start(ctx)
}
