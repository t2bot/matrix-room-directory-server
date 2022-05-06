# matrix-room-directory-server

A minimal implementation of a room directory server for Matrix, using a Space as a backend.

Configure `matrix-room-directory-server` to look at a space and it will expose all the rooms in that space under the [federation `/publicRooms`](https://spec.matrix.org/v1.1/server-server-api/#public-room-directory) endpoint. This can then be used by clients by signing into your actual homeserver, asking for the room directory for the domain this project is hosted on, and it will return the rooms in the configured space. This project is used to host the `t2bot.io` room directory for example.

Support room: [#matrix-room-directory-server:t2bot.io](https://matrix.to/#/#matrix-room-directory-server:t2bot.io)

**Caution**: Although this claims to be a room directory server, it is not yet recommended for full-featured deployment. 
Check the GitHub issues before deploying.

A room directory server is simply a service that resolves aliases to room IDs, where it then identifies one or more resident servers for the calling (usually joining) server to talk to. This further extends into offering the [`/publicRooms`](https://spec.matrix.org/v1.1/server-server-api/#public-room-directory) endpoint as a directory of rooms, usually supplied by a directory server.

## Building and running

This project does not provide any guidelines on how to run this in your infrastructure. It is up to you to determine
how best to deploy this, and how much of it actually gets deployed.

The process is meant to be run only attached to a postgres instance and does not have any on-disk requirements other 
than the executable itself.

You will need to be running or otherwise have access to a [matrix-key-server](https://github.com/t2bot/matrix-key-server).
This project also expects that you have extensive knowledge on how to set up an application service for
your server, as demonstrated by the program arguments.

This project uses Go modules and requires Go 1.17 or higher. To enable modules, set `GO111MODULE=on`.

```bash
# Build
git clone https://github.com/t2bot/matrix-room-directory-server.git
cd matrix-room-directory-server
go build -v -o bin/matrix-room-directory-server

# Run
./bin/matrix-room-directory-server \
    -keyserver="https://keys.t2host.io" \
    -address="0.0.0.0" \
    -port=8080 \
    -space="#directory:example.org" \
    -accesstoken="syt_randomstringfromserver" \
    -hsurl="https://t2bot.io"
```

#### Docker

```bash
docker run -it --rm \
    -e "ADDRESS=0.0.0.0" \
    -e "PORT=8080" \
    -e "KEYSERVER=https://keys.t2host.io" \
    -e "HSURL=https://t2bot.io" \
    -e "SPACE=#directory:example.org" \
    -e "ACCESSTOKEN=syt_randomstringfromserver" \
    t2bot/matrix-room-directory-server
```

Build your own by checking out the repository and running `docker build -t t2bot/matrix-room-directory-server .`
