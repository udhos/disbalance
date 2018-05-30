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

If you want to use TLS for the console, you will need a certificate:

    openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout run/key.pem -out run/cert.pem

Then start disbalance and open https://localhost:8080/console

How to test
===========

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

3. Visit http://localhost:8000. Your web traffic should be distributed between localhost:8001 and localhost:8002.

API
===

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

-x-

