# matrix-room-directory-server

A minimal implementation of a room directory server for Matrix

Support room: [#matrix-room-directory-server:t2bot.io](https://matrix.to/#/#matrix-room-directory-server:t2bot.io)

## Building and running

The key server will automatically generate itself a key to use on startup. The process is meant to be run 
only attached to a postgres instance and does not have any on-disk requirements other than the executable 
itself.

You will need to be running or otherwise have access to a [matrix-key-server](https://github.com/t2bot/matrix-key-server).

This project uses Go modules and requires Go 1.12 or higher. To enable modules, set `GO111MODULE=on`.

```bash
# Build
git clone https://github.com/t2bot/matrix-room-directory-server.git
cd matrix-room-directory-server
go build -v -o bin/matrix-room-directory-server

# Run
./bin/matrix-room-directory-server -keyserver="https://keys.t2host.io" -address="0.0.0.0" -port=8080 -postgres="postgres://username:password@localhost/dbname?sslmode=disable"
```

#### Docker

```bash
docker run -it --rm -e "ADDRESS=0.0.0.0" -e "PORT=8080" -e "KEYSERVER=https://keys.t2host.io" -e "POSTGRES=postgres://username:password@localhost/dbname?sslmode=disable" t2bot/matrix-room-directory-server
```

Build your own by checking out the repository and running `docker build -t t2bot/matrix-room-directory-server .`
