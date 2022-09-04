package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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
		err                 error
		script              string
		conn                *connect.Connect
		scriptByte, message []byte
		session             *ssh.Session

		entity struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}
	)

	defer func() {
		//writer.WriteHeader(http.StatusOK)
		writer.Header().Add("content-type", "application/json")
		data, _ := json.Marshal(entity)
		writer.Write(data)
	}()

	scriptByte, err = ioutil.ReadAll(request.Body)

	script = strings.ReplaceAll(string(scriptByte), "\r\n", "\n")

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

	// 执行命令
	message, err = session.CombinedOutput(script)
	if err != nil {
		msg := fmt.Sprintf("执行失败:%s%s", message, err.Error())
		log.Println(msg)
		entity.Message = msg
		return
	}

	entity.Message = string(message)
	entity.Success = true
}
