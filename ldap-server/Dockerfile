FROM alpine

ENV SUFFIX "DC=example,DC=com"
ENV ORG "Example Ltd"
ENV ADMIN "admin"
ENV PASS "admin"

ENV TLS_CA "/etc/ssl/certs/ca.pem"
ENV TLS_CERT "/etc/ssl/certs/server.pem"
ENV TLS_KEY "/etc/ssl/certs/server.key"
ENV TLS_VERIFY "allow"

RUN set -e ;\
    apk add --no-cache openldap openldap-back-mdb ;\
    mkdir -p /var/lib/openldap/openldap-data ;\
    chmod 700 /var/lib/openldap/openldap-data ;\
    mkdir -p /run/openldap ;\
    chmod 700 /run/openldap

COPY slapd.conf /etc/openldap/slapd.conf
COPY ldif /etc/openldap/ldif
COPY init.sh /init.sh

RUN set -e ;\
    chmod 600 /etc/openldap/slapd.conf ;\
    chmod 700 /etc/openldap/ldif ;\
    chmod 755 /init.sh

EXPOSE 389
EXPOSE 636

ENTRYPOINT ["/init.sh", "-d", "254" ]
