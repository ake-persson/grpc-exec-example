FROM alpine

RUN apk add --no-cache openldap

EXPOSE 389 636

ENTRYPOINT ["/usr/sbin/slapd", "-h", "ldap:/// ldapi:/// ldaps:///"]
