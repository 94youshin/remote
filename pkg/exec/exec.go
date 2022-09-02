package exec

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"sync"

	"golang.org/x/crypto/ssh"
)

// ExecPipe 执行脚本，实时返回
func ExecPipe(script string, session *ssh.Session, writer http.ResponseWriter) error {
	var (
		err    error
		stdout io.Reader
		wg     sync.WaitGroup
	)

	stdout, err = session.StdoutPipe()
	if err != nil {
		return err
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(stdout)
		flusher, ok := writer.(http.Flusher)
		if !ok {
			// err 不会发生
			panic("expected http.ResponseWriter to be an http.Flusher")
		}
		for {
			read, err := reader.ReadByte()
			if err != nil || err == io.EOF {
				return
			}
			fmt.Fprintf(writer, "%s", string(read))
			flusher.Flush()
		}
	}()
	err = session.Run(script)
	if err != nil {
		writer.Write([]byte(err.Error()))
	}
	wg.Wait()
	return nil
}

// Exec 执行脚本，等待执行结束后统一返回
func Exec(script string, session *ssh.Session, writer http.ResponseWriter) error {
	session.Stdout = writer
	session.Stderr = writer
	return session.Run(script)
}
