[![license](http://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/udhos/disbalance/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/udhos/disbalance)](https://goreportcard.com/report/github.com/udhos/disbalance)

# disbalance
disbalance - automagic load balancer

* [Features](#features)
* [Quick Start](#quick-start)
* [HTTPS](#https)
* [How to test](#how-to-test)
* [API](#api)
* [Install](#install)

Created by [gh-md-toc](https://github.com/ekalinin/github-markdown-toc.go)

# Features

- Minimal required configuration. You are supposed to fire up disbalance and start using it.
- Configuration automatically kept as YAML file. You are not required to edit it by hand.
- Integrated web console. Use the web interface to quickly define load balancing rules.
- REST API. Use the API to dynamically combine the load balancer with other services.

# Quick Start

    git clone https://github.com/udhos/disbalance
    cd disbalance
    ./build.sh
    disbalance

Then open http://localhost:8080/console

# HTTPS

If you want to use TLS for the console, you will need a certificate:

    openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout run/key.pem -out run/cert.pem

Then start disbalance and open https://localhost:8080/console

# How to test

1. Run webservers on ports 8001 and 8002.

Example:

    go get github.com/udhos/gowebhello

    gowebhello -addr :8001 -disableKeepalive &
    gowebhello -addr :8002 -disableKeepalive &

2. Then open the disbalance console http://localhost:8080/console and create a rule with targets for ports 8001 and 8002.

Create the rule with this information:

    rule name: rule-1000
    listener:  :8000
    target:    localhost:8001
    target:    localhost:8002

3. Visit http://localhost:8000.

Example:

    curl http://localhost:8000

Your web traffic should be distributed between localhost:8001 and localhost:8002.

# API

Get rule list:

    curl -u admin:admin localhost:8080/api/rule/

Delete rule 'rule1':

    curl -u admin:admin -X DELETE localhost:8080/api/rule/rule1

Create/update inline rule 'rule1':

    curl -u admin:admin -X POST -d '{rule1: {listener: ":2000", protocol: tcp, targets: {2.2.2.2:80: {}}}}' localhost:8080/api/rule/

Replace inline rule 'rule1':

    curl -u admin:admin -X PUT -d '{listener: ":3000", targets: {3.3.3.3:80: {}}}' localhost:8080/api/rule/rule1

Save rule list to file 'rules':

    curl -u admin:admin localhost:8080/api/rule/ > rules

Load rule list from file 'rules':

    curl -u admin:admin --data-binary @rules -X POST localhost:8080/api/rule/

# Install

You can install disbalance as a systemd service by running the install.sh script:

    sudo ./install.sh

    sudo ./install.sh remove ;# this removes disbalance service from systemd

-x-

