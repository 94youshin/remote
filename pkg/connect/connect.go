package connect

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"net/http"
	"time"
)

type Connect struct {
	User     string
	Password string
	Address  string
}

// NewConnect 构造NewConnect用于连接远程机器
func NewConnect(request *http.Request) (*Connect, error) {
	var (
		ip, port, address, username, password string
	)
	// 解析认证相关参数
	query := request.URL.Query()
	ip = query.Get("ip")
	port = query.Get("port")
	address = fmt.Sprintf("%s:%s", ip, port)
	username = query.Get("username")
	password = query.Get("password")

	if ip == "" || port == "" || username == "" || password == "" {
		return nil, errors.New("ip,port,username,password 均不能为空")
	}
	return &Connect{
		User:     username,
		Password: password,
		Address:  address,
	}, nil
}

// Session 构造一个ssh Session
func (c *Connect) Session() (session *ssh.Session, err error) {
	var (
		client       *ssh.Client
		clientConfig *ssh.ClientConfig
	)

	clientConfig = &ssh.ClientConfig{
		User:            c.User,
		Auth:            []ssh.AuthMethod{ssh.Password(c.Password)},
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	if client, err = ssh.Dial("tcp", c.Address, clientConfig); err != nil {
		return nil, err
	}

	if session, err = client.NewSession(); err != nil {
		return nil, err
	}

	return session, nil
}
