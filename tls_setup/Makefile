all: ca req

clean:
	rm -rf certs

preq:
	go get -u -tags nopkcs11 github.com/cloudflare/cfssl/cmd/cfssl
	go get -u github.com/cloudflare/cfssl/cmd/cfssljson

ca:
	mkdir -p certs
	cfssl gencert -initca config/ca-csr.json | cfssljson -bare certs/ca
	mv certs/ca-key.pem certs/ca.key

req:
	cfssl gencert \
	  -ca certs/ca.pem \
	  -ca-key certs/ca.key \
	  -config config/ca.json config/auth.json | cfssljson -bare certs/auth
	mv certs/auth-key.pem certs/auth.key

	cfssl gencert \
	  -ca certs/ca.pem \
	  -ca-key certs/ca.key \
	  -config config/ca.json config/info.json | cfssljson -bare certs/info
	mv certs/info-key.pem certs/info.key

	cfssl gencert \
	  -ca certs/ca.pem \
	  -ca-key certs/ca.key \
	  -config config/ca.json config/exec.json | cfssljson -bare certs/exec
	mv certs/exec-key.pem certs/exec.key

	cfssl gencert \
	  -ca certs/ca.pem \
	  -ca-key certs/ca.key \
	  -config config/ca.json config/catalog.json | cfssljson -bare certs/catalog
	mv certs/catalog-key.pem certs/catalog.key

	openssl genrsa -out certs/private.rsa 2048
	openssl rsa -in certs/private.rsa -outform PEM -pubout -out certs/public.rsa

.PHONY: clean preq ca req
