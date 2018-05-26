#!/bin/bash

me=$(basename $0)

msg() {
	echo >&2 $me: $@
}

die() {
	msg $@
	exit 1
}

do_install() {
	[ -f run/disbalance ] || die missing executable: run/disbalance
	[ -r run/disbalance ] || die missing executable: run/disbalance
	[ -x run/disbalance ] || die missing executable: run/disbalance

	[ -f /usr/local/sbin/disbalance ] && die already installed, use: $0 remove

	mkdir -p /var/run/disbalance

	cp -r run/console            /var/run/disbalance
	cp    run/disbalance         /usr/local/sbin/disbalance
	cp    run/disbalance.service /lib/systemd/system

	systemctl daemon-reload
	systemctl enable disbalance.service
	systemctl reload-or-restart disbalance.service
}

do_uninstall() {
	systemctl daemon-reload
	systemctl stop disbalance.service
	systemctl disable disbalance.service
	pkill -x -9 disbalance
	rm -f /usr/local/sbin/disbalance
}

case "$1" in
	remove)
		do_uninstall
	;;
	'')
		do_install
	;;
	*)
		msg invalid argument: "$1"
		echo >&2 usage: $(basename $0) [remove]
		exit 2
	;;
esac

