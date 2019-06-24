#!/bin/bash

docker run -t -p 389:389 -p 636:636 \
	-v $PWD/../tls_setup/certs:/etc/ssl/certs \
	-e ORGANISATION_NAME="Example" \
	-e SUFFIX="dc=example,dc=com" \
	-e TLS_VERIFY_CLIENT="allow" \
	-e CA_FILE="/etc/ssl/certs/ca.pem" \
	-e KEY_FILE="/etc/ssl/certs/ldap.key" \
	-e CERT_FILE="/etc/ssl/certs/ldap.pem" \
	ldap-server
