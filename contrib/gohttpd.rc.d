#!/bin/sh

# PROVIDE: gohttpd
# KEYWORD: shutdown

. /etc/rc.subr

name="gohttpd"
rcvar="${name}_enable"

load_rc_config ${name}
: ${gohttpd_enable:=NO}
: ${gohttpd_logfile:="/var/log/${name}.log"}
: ${gohttpd_rootdir:="/var/www/root"}
: ${gohttpd_domains:=""}

pidfile="/var/run/${name}.pid"
procname="/usr/local/bin/${name}"

command="/usr/sbin/daemon"
command_args="-p ${pidfile} -o ${gohttpd_logfile} ${procname} ${gohttpd_rootdir} -domains=${gohttpd_domains}"

run_rc_command "$1"
