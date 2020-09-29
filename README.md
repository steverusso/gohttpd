# GoHTTPd

```
GoHTTPd is a specific, stubbornly simple web server.

Usage:

    gohttpd <directory> [options]

Options:

    -mem         Cache files in memory instead of using disk.
    -port int    The port to use for the local server. (default 8080)
    -tls         Use TLS.
```

If the given `directory` contains a GoHTTPd configuration file, then it will
treat all child directories as individual sites. Otherwise, it will just treat
the given directory itself as a site.

## Install

### FreeBSD

To install GoHTTPd as a production-ready service on a FreeBSD system:

```shell
sudo cp gohttpd /usr/local/bin
sudo cp contrib/gohttpd.rc.d /usr/local/etc/rc.d/gohttpd
sudo sysrc gohttpd_enable=YES
sudo service gohttpd start
```
