First generate Certificates for localhost.

```bash
cd tls_setup
make
```

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

Now build and start  **info-server** and **exec-server**.

Once all the servers are running change to the client directory.

```bash
cd client
go build
./client login
./client verify
./client info localhost
./client info localhost,localhost
./client exec localhost /bin/ls -la /
./client exec localhost,localhost /bin/ls -la /
```
