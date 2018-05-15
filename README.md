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

API
===

Get rule list:

    curl -u admin:admin localhost:8080/api/rule/

Delete rule 'rule1':

    curl -u admin:admin -X DELETE localhost:8080/api/rule/rule1

Create/update inline rule 'rule1':

    curl -u admin:admin -X POST -d '{rule1: {protocol: tcp}}' localhost:8080/api/rule/

Replace inline rule 'rule1':

    curl -u admin:admin -X PUT -d '{listener: ":2000"}' localhost:8080/api/rule/rule1

Save rule list to file 'rules':

    curl -u admin:admin localhost:8080/api/rule/ > rules

Load rule list from file 'rules':

    curl -u admin:admin --data-binary @rules -X POST localhost:8080/api/rule/
