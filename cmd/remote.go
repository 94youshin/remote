package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/crypto/ssh"
	"mcloud.chinaunicom.cn/remote/pkg/connect"
	"mcloud.chinaunicom.cn/remote/pkg/exec"
)

func main() {
	http.HandleFunc("/exec", func(writer http.ResponseWriter, request *http.Request) {
		err := handle(writer, request, exec.Exec)
		if err != nil {
			errHandle(err, writer)
		}
	})

	http.HandleFunc("/execPipe", func(writer http.ResponseWriter, request *http.Request) {
		err := handle(writer, request, exec.ExecPipe)
		if err != nil {
			errHandle(err, writer)
		}
	})

	go func() {
		fmt.Println("Start to listening the incoming requests on http address: 0.0.0.0:9080")
		if err := http.ListenAndServe(":9080", nil); err != nil {
			log.Println(err.Error())
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

type ExecFunc func(string, *ssh.Session, http.ResponseWriter) error

func handle(writer http.ResponseWriter, request *http.Request, execFunc ExecFunc) error {
	var (
		err     error
		conn    *connect.Connect
		script  []byte
		session *ssh.Session
	)
	script, err = ioutil.ReadAll(request.Body)
	if err != nil {
		return err
	}
	conn, err = connect.NewConnect(request)
	if err != nil {
		return err
	}
	session, err = conn.Session()
	if err != nil {
		return err
	}
	defer session.Close()
	err = execFunc(string(script), session, writer)
	if err != nil {
		return err
	}

	return nil
}

func errHandle(err error, w http.ResponseWriter) {
	log.Println(err.Error())
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}
