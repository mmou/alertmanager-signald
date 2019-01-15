package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"syscall"

	"github.com/cutelab/alertmanager-signald/alerts"
)

func main() {
	cmd := exec.Command("/signald/bin/signald")

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGTERM,
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	server := createServer()
	cmd.Run()
	log.Println("signald exited. Server shutting down")
	server.Shutdown(nil)
}

func createServer() *http.Server {
	server := &http.Server{Addr: ":8080"}
	http.HandleFunc("/", alerts.Handler)
	go func() {
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	return server
}
