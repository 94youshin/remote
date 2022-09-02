# remote
## 部署方式
> docker 部署
```bash
$ docker build -t remote:v1 .
$ docker run -dit --name remote-agent -p 9080:9080 remote:v1
```
> 二进制部署
```bash
$ make build
$ ./remote
```

## 调用方式
```bash
$ curl --location --request POST '10.92.119.242:9080/exec?ip=127.0.0.1&username=roott&password=123456&port=22' \
--header 'Content-Type: text/plain' \
--data-raw 'systemctl status sshd'
```
