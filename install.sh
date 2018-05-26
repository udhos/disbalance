#!/bin/bash

die() {
	echo >&2 $(basename $0): $@
	exit 1
}

[ -f run/disbalance ] || die missing executable: run/disbalance
[ -r run/disbalance ] || die missing executable: run/disbalance
[ -x run/disbalance ] || die missing executable: run/disbalance

mkdir -p /var/run/disbalance

cp -p -r run/console            /var/run/disbalance
cp -p    run/disbalance.service /var/run/disbalance
cp -p    run/disbalance         /usr/local/sbin/disbalance

systemctl reenable /var/run/disbalance/disbalance.service
systemctl reload disbalance.service
systemctl restart disbalance.service
