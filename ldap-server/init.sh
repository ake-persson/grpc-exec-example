#!/bin/sh

sed -i "s!%SUFFIX%!$SUFFIX!g" /etc/openldap/slapd.conf
sed -i "s!%ADMIN%!$ADMIN!g" /etc/openldap/slapd.conf
sed -i "s!%PASS%!$PASS!g" /etc/openldap/slapd.conf

sed -i "s!%TLS_CA%!$TLS_CA!g" /etc/openldap/slapd.conf
sed -i "s!%TLS_KEY%!$TLS_KEY!g" /etc/openldap/slapd.conf
sed -i "s!%TLS_CERT%!$TLS_CERT!g" /etc/openldap/slapd.conf
sed -i "s!%TLS_VERIFY%!$TLS_VERIFY!g" /etc/openldap/slapd.conf

for f in /etc/openldap/ldif/*.ldif; do
	echo "Adding: $f"
	slapadd -l $f
done

exec /usr/sbin/slapd -h "ldap:/// ldaps:///" "$@"
