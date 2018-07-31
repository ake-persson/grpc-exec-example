## Create TLS Certificates

First generate Certificates for localhost.

```bash
cd tls_setup
make
```

## Build auth-server

Start by compiling the auth-server.

```bash
cd auth-server
go build
./auth-server
```

You can set the specific LDAP/AD settings in **~/.auth-server.toml**.

```toml
addr = "dc.example.com:389"
backend = "ad"
base = "DC=example,DC=com"
domain = "example"
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

Generate a JWT token by loggin in.

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
