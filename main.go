package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func initLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	serverIP := flag.String("server-ip", "0.0.0.0", "the IP to serve on")
	serverPort := flag.Int("server-port", 8080, "the Port to listen on")
	dirToServe := flag.String("serve-dir", "static", "the directory to serve")
	flag.Parse()

	staticDir := ""
	if *dirToServe == "static" {
		currentDir, _ := os.Getwd()
		staticDir = filepath.Join(currentDir, "static")
	} else {
		fileInfo, err := os.Stat(*dirToServe)
		if errors.Is(err, os.ErrNotExist) {
			log.WithError(err).WithField("target_dir", *dirToServe).Fatalf("the target directory does not exist")
		} else if !fileInfo.IsDir() {
			log.WithError(err).WithField("target_dir", *dirToServe).Fatalf("the target file is not a directory")
		} else {
			staticDir = *dirToServe
		}
	}

	http.Handle("/", http.FileServer(http.Dir(staticDir)))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	log.WithFields(log.Fields{
		"server_ip":   *serverIP,
		"server_port": *serverPort,
	}).Infof("the server is now listening")

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", *serverIP, *serverPort), nil)
	if err != nil {
		log.WithError(err).Fatal("an error occurred")
	}

}
