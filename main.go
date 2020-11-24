package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"

	"net/http"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("could not start server", err)
	}

	// Get ENV variables
	err = godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	connPort, ok := os.LookupEnv("CONN_PORT")
	if !ok {
		log.Fatal("Set CONN_PORT variable in .env file")
	}

	upstreamServer, ok := os.LookupEnv("UPSTREAM_SERVER")
	if !ok {
		log.Println("Set UPSTREAM_SERVER variable in .env file")
	}

	upstreamServerPort, ok := os.LookupEnv("UPSTREAM_SERVER_PORT")
	if !ok {
		log.Println("Set UPSTREAM_SERVER_PORT variable in .env file")
	}

	bufferSize := 1024

	// create and handle TCP connections from client
	ln, err := net.Listen("tcp", ":"+connPort)
	if err != nil {
		log.Fatalf("Couldn't start server on port %s: %v", connPort, err)
	}
	log.Println("Application running on ", ln.Addr())
	for {
		clientConn, err := ln.Accept()
		if err != nil {
			log.Println("Unable to accept connections at this time, kindly try again later: ", err)
		}
		defer clientConn.Close()

		fmt.Println("client IP is: ", clientConn.LocalAddr())

		// Create buffer to hold client data
		buf := make([]byte, bufferSize)
		_, err = clientConn.Read(buf)
		if err != nil {
			log.Println("Unable to read data from client: ", err)
		}

		go handlRequest(clientConn, buf, upstreamServer, upstreamServerPort, bufferSize)
	}
}

func handlRequest(clientConn net.Conn, buf []byte, upstreamServer string, upstreamServerPort string, bufferSize int) {

	conf := &tls.Config{
		InsecureSkipVerify: false,
	}

	upstreamConn, err := tls.Dial("tcp", upstreamServer+":"+upstreamServerPort, conf)
	if err != nil {
		log.Println("Error connecting to upstream server", err)
	}
	fmt.Println("Connecting to upstream server:", upstreamConn.RemoteAddr())

	defer upstreamConn.Close()

	_, err = upstreamConn.Write(buf)
	if err != nil {
		log.Println("Error writing to upstream server: ", err)
	}

	// Read response from upstream server
	kBuf := make([]byte, bufferSize)
	_, _ = upstreamConn.Read(kBuf)

	// write response back to client
	_, err = clientConn.Write(kBuf)
	if err != nil {
		log.Println("Error writing to client, closing client connection")
		clientConn.Close()
	}
}
