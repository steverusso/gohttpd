#!/bin/sh

# PROVIDE: gohttpd
# KEYWORD: shutdown

. /etc/rc.subr

name="gohttpd"
rcvar="${name}_enable"

load_rc_config ${name}
: ${gohttpd_enable:=NO}
: ${gohttpd_rootdir:="/var/www"}

pidfile="/var/run/${name}.pid"
logfile="/var/log/${name}.log"
procname="/usr/local/bin/${name}"

command="/usr/sbin/daemon"
command_args="-p ${pidfile} -o ${logfile} ${procname} ${gohttpd_rootdir} -tls"

run_rc_command "$1"
