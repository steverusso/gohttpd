# GoHTTPd

## Building

```shell
go build -o gohttpd cmd/gohttpd/main.go
```

## Install

### FreeBSD

To install GoHTTPd as a production-ready service on a FreeBSD system:

```shell
sudo cp gohttpd /usr/local/bin
sudo cp contrib/gohttpd.rc.d /usr/local/etc/rc.d/gohttpd
sudo sysrc gohttpd_enable=YES
sudo service gohttpd start
```
