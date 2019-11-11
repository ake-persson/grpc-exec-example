all: build

build:
	docker build . -t ldap-server

run:
	docker run -it -p 389:389 -p 636:636 \
	  -v $$PWD/../tls_setup/certs:/etc/ssl/certs \
	  -e TLS_CA="/etc/ssl/certs/ca.pem" \
	  -e TLS_CERT="/etc/ssl/certs/ldap.pem" \
	  -e TLS_KEY="/etc/ssl/certs/ldap.key" \
	  ldap-server

test:
	ldapsearch -x

.PHONY: build run test
