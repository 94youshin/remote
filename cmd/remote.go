package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"mcloud.chinaunicom.cn/remote/pkg/connect"
	"mcloud.chinaunicom.cn/remote/pkg/exec"
	"net/http"
)

func main() {
	http.HandleFunc("/exec", func(writer http.ResponseWriter, request *http.Request) {
		handle(writer, request, exec.Exec)
	})

	http.HandleFunc("/execPipe", func(writer http.ResponseWriter, request *http.Request) {
		handle(writer, request, exec.ExecPipe)
	})

	err := http.ListenAndServe(":9080", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

type ExecFunc func(string, *ssh.Session, http.ResponseWriter) error

func handle(writer http.ResponseWriter, request *http.Request, execFunc ExecFunc) {
	var (
		err     error
		conn    *connect.Connect
		script  []byte
		session *ssh.Session
	)
	script, err = ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
	}
	conn, err = connect.NewConnect(request)
	if err != nil {
		log.Println(err.Error())
	}
	session, err = conn.Session()
	if err != nil {
		log.Println(err.Error())
	}
	defer session.Close()
	err = execFunc(string(script), session, writer)
	if err != nil {
		writer.Write([]byte(err.Error()))
	}
}
