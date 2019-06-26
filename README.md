## Install Go

First install Go and then configure Go environment.

### Mac OS X

```bash
brew install go
```

### RedHat/CentOS/Fedora

```bash
yum install golang
```

### Setup Go environment

```bash
mkdir -p ~/go/{src,bin}
cat << EOF >>~/.bash_profile
export GOPATH=~/go
export PATH=\$PATH:\$GOPATH/bin
EOF
source ~/.bash_profile
```

## Clone code

```bash
mkdir -p $GOPATH/src/github.com/mickep76
cd $GOPATH/src/github.com/mickep76
git clone https://github.com/mickep76/grpc-exec-example.git
```

## Create TLS Certificates
  
First generate TLS Certificates for localhost and RSA public/private key.

```bash
cd $GOPATH/src/github.com/mickep76/grpc-exec-example/tls_setup
make preq ca req
```

## Setup OpenLDAP server inside Docker

```bash
cd $GOPATH/src/github.com/mickep76/grpc-exec-example/ldap-server
make build run
```

## Download/update deps

```bash
cd $GOPATH/src/github.com/mickep76/grpc-exec-example
go get -u ./...
```

## Build auth-server

Build the auth-server.

```bash
cd $GOPATH/src/github.com/mickep76/grpc-exec-example/auth-server
go build
```

You can set the specific LDAP/AD settings in **~/.auth-server.toml**.

```bash
cat << EOF >~/.auth-server.toml
addr = "localhost:389"
base = "dc=example,dc=com"
ou = "ou=users"
ca = "../tls_setup/certs/ca.pem"
verify = true
EOF
```

Start auth-server.

```bash
./auth-server
```

## Build info-server

```bash
cd $GOPATH/src/github.com/mickep76/grpc-exec-example/info-server
go build
./info-server
```

## Build exec-server

```bash
cd $GOPATH/src/github.com/mickep76/grpc-exec-example/exec-server
go build
./exec-server
```

## Build client

Build the client so you can interact with the services.

```bash
cd $GOPATH/src/github.com/mickep76/grpc-exec-example/client
go build
```

Generate a JWT token by logging in.

For user "jdoe" use password "secret".

```bash
./client login -user jdoe
./client verify
```

Query host info from one or more hosts.

```bash
./client info localhost
./client info localhost,localhost
```

Execute command on one or more hosts.

```bash
./client exec localhost /bin/ls -la /
./client exec localhost,localhost /bin/ls -la /
```

## Build catalog-server

Catalog server will allow each info-server to register and send a keep-alive.

```bash
cd $GOPATH/src/github.com/mickep76/grpc-exec-example/catalog-server
go build
./catalog-server
```

## Restart info-server

Stop info-server CTRL+C then start it with the option **_-register_**. This requires that you have created a token.

```bash
cd $GOPATH/src/github.com/mickep76/grpc-exec-example/info-server
./info-server -register
```

## List hosts registered in catalog-server

```bash
cd $GOPATH/src/github.com/mickep76/grpc-exec-example/client
./client list
```
