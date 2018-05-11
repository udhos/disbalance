# disbalance
disbalance - automatic load balancer

If you want to use TLS, you will need a certificate:

    $ openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem
