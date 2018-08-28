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
export PATH=$PATH:$GOPATH/bin
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

## Download/update deps

```bash
cd $GOPATH/src/github.com/mickep76/grpc-exec-example
get get -u ./...
```

## Build auth-server

Build the auth-server.

```bash
cd $GOPATH/src/github.com/mickep76/grpc-exec-example/auth-server
go build
```

You can set the specific LDAP/AD settings in **~/.auth-server.toml**.

```toml
addr = "dc.example.com:389"
backend = "ad"
base = "DC=example,DC=com"
domain = "example"
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

```bash
./client login
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
