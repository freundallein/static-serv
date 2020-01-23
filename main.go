package main

import (
	"fmt"
	"github.com/freundallein/static-serv/server"
	"log"
	"os"
	"time"
)

const (
	timeFormat = "02.01.2006 15:04:05"

	portKey     = "PORT"
	defaultPort = "8000"

	rootDirKey     = "STATIC_ROOT"
	defaultRootDir = "/static"

	prefixKey     = "PREFIX"
	defaultPrefix = "/static"
)

type logWriter struct{}

// Write - custom logger formatting
func (writer logWriter) Write(bytes []byte) (int, error) {
	msg := fmt.Sprintf("%s | [staticserv] %s", time.Now().UTC().Format(timeFormat), string(bytes))
	return fmt.Print(msg)
}

func getEnv(key string, fallback string) (string, error) {
	if value := os.Getenv(key); value != "" {
		return value, nil
	}
	return fallback, nil
}

func main() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))

	port, err := getEnv(portKey, defaultPort)
	if err != nil {
		log.Fatalf("[config] %s\n", err.Error())
	}
	prefix, err := getEnv(prefixKey, defaultPrefix)
	if err != nil {
		log.Fatalf("[config] %s\n", err.Error())
	}
	rootDir, err := getEnv(rootDirKey, defaultRootDir)
	if err != nil {
		log.Fatalf("[config] %s\n", err.Error())
	}
	options := &server.Options{
		Port:    port,
		RootDir: rootDir,
		Prefix:  prefix,
	}
	log.Printf("[config] starting with prefix - %s, dir - %s\n", prefix, rootDir)
	srv, err := server.New(options)
	if err != nil {
		log.Fatalf("[config] %s\n", err.Error())
	}
	if err := srv.Run(); err != nil {
		log.Fatalf("[server] %s\n", err.Error())
	}
}
