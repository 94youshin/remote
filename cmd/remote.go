package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"golang.org/x/crypto/ssh"
	"mcloud.chinaunicom.cn/remote/pkg/connect"
)

func main() {
	http.HandleFunc("/exec", handle)
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

func handle(writer http.ResponseWriter, request *http.Request) {
	var (
		err     error
		message string
		conn    *connect.Connect
		script  []byte
		session *ssh.Session

		entity struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}
	)

	defer func() {
		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("content-type", "application/json")
		data, _ := json.Marshal(entity)
		writer.Write(data)
	}()

	script, err = ioutil.ReadAll(request.Body)
	if err != nil {
		msg := fmt.Sprintf("读取脚本数据失败:%s", err.Error())
		log.Println(msg)
		entity.Message = msg
		return
	}
	conn, err = connect.NewConnect(request)
	if err != nil {
		msg := fmt.Sprintf("创建远程连接失败:%s", err.Error())
		log.Println(msg)
		entity.Message = msg
		return
	}
	session, err = conn.Session()
	if err != nil {
		msg := fmt.Sprintf("创建远程Session失败:%s", err.Error())
		log.Println(msg)
		entity.Message = msg
		return
	}
	defer session.Close()

	message, err = exec(string(script), session)
	if err != nil {
		msg := fmt.Sprintf("执行失败:%s", err.Error())
		log.Println(msg)
		entity.Message = msg
		return
	}

	entity.Message = message
	entity.Success = true
}

// Exec 执行脚本，执行结束后统一返回
func exec(script string, session *ssh.Session) (string, error) {
	var (
		err     error
		message string
		stdout  io.Reader
		wg      sync.WaitGroup
	)

	stdout, err = session.StdoutPipe()
	if err != nil {
		return "", err
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(stdout)
		for {
			read, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				return
			}
			message = fmt.Sprintf("%s%s", message, read)
		}
	}()
	err = session.Run(script)
	if err != nil {
		return "", err
	}
	wg.Wait()
	return message, nil
}
