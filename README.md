# AmILeaked

**AmILeaked** is a simple leak/email data registry service with open REST API and gRPC service. This is given as a coding challenge by Binaryedge.

The project consists of two main parts, these are the **RPC server** and the **REST gateway**. The gateway is a HTTP reverse proxy which translates the binary gRPC messages to JSON messages.

## Installation

The project is implemented using `gRPC` framework with `Golang`. You need to install some packages first to compile the `.proto` files. Then you can compile the server microservices.

### Setting The Development Environment

First get the required packages to use `gRPC` environment to compile the **RPC server**.

```bash
go get google.golang.org/grpc
go get google.golang.org/grpc/reflection
go get github.com/grpc-ecosystem/go-grpc-middleware
go get github.com/grpc-ecosystem/go-grpc-middleware/recovery
```

Now get the packages required to compile the **REST gateway**.

```bash
go install \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
    github.com/golang/protobuf/protoc-gen-go
```

### Compiling The Source

Let's start with compiling the `.proto` files and generating the corresponding server/client codes.

---
**Note**

The `protocompile.sh` file is written to work in *Nix environments. In Win environments, the `-I/usr/local/include` flag should be changed to a proper one.

---

```bash
./protocompile.sh
```

This command generates **Golang** codes for gRPC server/client and **serializer/deserializer** codes in `rpc-server/` and `rest-server/` directories.

Now go to `rpc-server/` directory and compile the **RPC server**.

```bash
cd rpc-server/
go build
```

This command compiles the RPC server into the file `rpc-server`.

Now go to `rest-server/` directory and compile the **REST gateway**.

```bash
cd ../rest-server/
go build
```

This command compiles the RPC server into the file `rest-server`.

Now both the RPC service and API service are ready to be used.

## Usage

To use the system, you just need to run the `rpc-server` binary in the background, and then the `rest-server` after it.

```bash
./rpc-server &
./rest-server
```

---
**Note**

The RPC server is hardcoded to run on port `50051` on `localhost`.

The REST gateway is hardcoded to run on port `8081` on `localhost`.

---

### API Endpoints

#### ListLeaks

- **Endpoint**: `/v1-beta/leaks`
- **Method**: `GET`
- **Example response**: 
```json
{
    "result": {
        "id": "5efdb25c0479ba527b1788f6",
        "name": "fancy leak",
        "emails": [
            {
                "email": "ozan.sazak@ieee.metu.edu.tr",
                "domain": "ieee.metu.edu.tr",
                "first_occurance_ts": "1593684730000",
                "last_occurance_ts": "1593684730000"
            },
            {
                "email": "-----%%%@metu.edu.tr",
                "domain": "metu.edu.tr",
                "first_occurance_ts": "1593684730000",
                "last_occurance_ts": "1593684730000"
            }
        ],
        "email_count": "2"
    }
}
{
    "result": {
        "id": "5efdb2660479ba527b1788f7",
        "name": "moore fancy leak",
        "emails": [
            {
                "email": "ozan.sazak@ieee.metu.edu.tr",
                "domain": "ieee.metu.edu.tr",
                "first_occurance_ts": "1593684730000",
                "last_occurance_ts": "1593684730000"
            }
        ],
        "email_count": "1"
    }
}
{
    "result": {
        "id": "5efdb27c0479ba527b1788f8",
        "name": "the big turkey data leak",
        "emails": [
            {
                "email": "ozan.sazak@ieee.metu.edu.tr",
                "domain": "ieee.metu.edu.tr",
                "first_occurance_ts": "1593684730000",
                "last_occurance_ts": "1593684730000"
            },
            {
                "email": "ozan.sazak@ieee.metu.edu.tr",
                "domain": "ieee.metu.edu.tr",
                "first_occurance_ts": "1593684730000",
                "last_occurance_ts": "1593684730000"
            },
            {
                "email": "-----%%%@metu.edu.tr",
                "domain": "metu.edu.tr",
                "first_occurance_ts": "1593684730000",
                "last_occurance_ts": "1593684730000"
            },
            {
                "email": "AdM----iiiiNNN%%%++@metu.edu.tr",
                "domain": "metu.edu.tr",
                "first_occurance_ts": "1593684730000",
                "last_occurance_ts": "1593684730000"
            }
        ],
        "email_count": "4"
    }
}
```

---
**Note**

As this endpoint is the gateway-ed version of a **stream** RPC function, the response is actually **not** a JSON message, it is the concatenation of seperate JSONs each for a different `Leak` object.

---

#### GetLeaksByDomain

- **Endpoint**: `/v1-beta/leaks-by-domain`
- **Method**: `POST`
- **Example request**:
```json
{
	"domain": "ieee.metu.edu.tr"
}
```

- **Example response**: 
```json
{
    "leaks": [
        {
            "id": "5efdb25c0479ba527b1788f6",
            "name": "fancy leak",
            "emails": [
                {
                    "email": "ozan.sazak@ieee.metu.edu.tr",
                    "domain": "ieee.metu.edu.tr",
                    "first_occurance_ts": "1593684730000",
                    "last_occurance_ts": "1593684730000"
                }
            ],
            "email_count": "1"
        },
        {
            "id": "5efdb2660479ba527b1788f7",
            "name": "moore fancy leak",
            "emails": [
                {
                    "email": "ozan.sazak@ieee.metu.edu.tr",
                    "domain": "ieee.metu.edu.tr",
                    "first_occurance_ts": "1593684730000",
                    "last_occurance_ts": "1593684730000"
                }
            ],
            "email_count": "1"
        }
    ]
}
```

#### GetLeaksByEmail

- **Endpoint**: `/v1-beta/leaks-by-email`
- **Method**: `POST`
- **Example request**:
```json
{
	"domain": "ieee.metu.edu.tr"
}
```

- **Example response**: 
```json
{
    "leaks": [
        {
            "id": "5efdb25c0479ba527b1788f6",
            "name": "fancy leak"
        },
        {
            "id": "5efdb2660479ba527b1788f7",
            "name": "moore fancy leak"
        },
        {
            "id": "5efdb27c0479ba527b1788f8",
            "name": "the big turkey data leak"
        },
        {
            "id": "5efdb27c0479ba527b1788f8",
            "name": "the big turkey data leak"
        },
        {
            "id": "5efdb2940479ba527b1788f9",
            "name": "n00bZ l€@k"
        }
    ]
}
```

- **Example error response**:
```json
{
    "error": "Unable to get Email object ID",
    "code": 13,
    "message": "Unable to get Email object ID"
}
```

## Further Info

Further information about the system -such as API details and design decisions- can be found on the system documentations in `docs/`.