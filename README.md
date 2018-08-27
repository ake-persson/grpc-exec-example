## Create TLS Certificates

First generate TLS Certificates for localhost and RSA public/private key.

```bash
cd tls_setup
make preq ca req
```

## Setup Go

First install Go and then configure $GOHOME.

### Mac OS X

```bash
brew install go
```

### RedHat/CentOS/Fedora

```bash
yum install go
```

### Setup $GOHOME
```bash
mkdir -p ~/go/src
cat <<EOF>>~/.bash_profile
export GOPATH=~/go
EOF
source ~/.bash_profile
```

## Build auth-server

Build the auth-server.

```bash
cd auth-server
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
cd auth-server
go build
./info-server
```

## Build exec-server

```bash
cd exec-server
go build
./exec-server
```

## Build client

Build the client so you can interact with the services.

```bash
cd client
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
