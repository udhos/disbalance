# disbalance
disbalance - automagic load balancer

Quick Start
===========

	go get github.com/udhos/disbalance
	cd ~/go/src/github.com/udhos/disbalance
	./build.sh
	disbalance

Then open http://localhost:8080/console

HTTPS
=====

If you want to use TLS, you will need a certificate:

    openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout run/key.pem -out run/cert.pem

Then start disbalance and open https://localhost:8080/console

