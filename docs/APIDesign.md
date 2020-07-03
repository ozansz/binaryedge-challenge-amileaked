# API Design Documentation

This document describes the detailed design process of the APIs -both RPC and REST-.

## Overall Information

The system uses microservice approach using a gRPC service. The REST API is handled by the **HTTP reverse proxy** (REST Gateway) in front of the RPC service. Both the RPC and REST services are accessible over the internet.

## Endpoints

### Listing Leaks

- **RPC Function**: `ListLeaks`
- **RPC Request**: `google.protobuf.Empty`
- **RPC Response**: `stream Leak`
- **REST Endpoint**: `/v1-beta/leaks`
- **Method**: `GET`
- **REST Response**: 
```json
{
    "result": {
        "id": <ObjectID>,
        "name": <String>,
        "emails": [
            {
                "email": <String>,
                "domain": <String>,
                "first_occurance_ts": <Int64|String>,
                "last_occurance_ts": <Int64|String>
            },
            ...
        ],
        "email_count": <Int64|String>
    }
}
...
```

### Get Leaks by the Domain Specified

- **RPC Function**: `GetLeaksByDomain`
- **RPC Request**: `GetLeaksByDomainRequest`
- **RPC Response**: `GetLeaksByDomainResponse`
- **REST Endpoint**: `/v1-beta/leaks-by-domain`
- **Method**: `POST`
- **REST Request**:
```json
{
	"domain": <String>
}
```

- **REST Response**: 
```json
{
    "leaks": [
        {
            "id": <ObjectID>,
            "name": <String>,
            "emails": [
                {
                    "email": <String>,
                    "domain": <String>,
                    "first_occurance_ts": <Int64|String>,
                    "last_occurance_ts": <Int64|String>
                },
                ...
            ],
            "email_count": <Int64|String>
        },
        ...
    ]
}
```

### Get Leaks by the Email Specified

- **RPC Function**: `GetLeaksByEmail`
- **RPC Request**: `GetLeaksByEmailRequest`
- **RPC Response**: `GetLeaksByEmailResponse`
- **Endpoint**: `/v1-beta/leaks-by-email`
- **Method**: `POST`
- **REST Request**:
```json
{
	"domain": <String>
}
```

- **REST Response**: 
```json
{
    "leaks": [
        {
            "id": <ObjectID>,
            "name": <String>
        },
        ...
    ]
}
```

- **REST Error Response**:
```json
{
    "error": <String>,
    "code": <Int64>,
    "message": <String>
}
```

### Get Leaks by the Domain Specified (Streamed)

- **RPC Function**: `GetLeaksByDomainStreamed`
- **RPC Request**: `GetLeaksByDomainRequest`
- **RPC Response**: `stream Leak`
- **REST Endpoint**: -

### Get Leaks by the Email Specified (Streamed)

- **RPC Function**: `GetLeaksByEmailStreamed`
- **RPC Request**: `GetLeaksByEmailRequest`
- **RPC Response**: `stream Leak`
- **REST Endpoint**: -

## Security

Since the project was not descibed to have a closed -internal- API, I did not add any authorization mechanisms into the system.

Although, a more secure gRPC service which uses SSL/TLS certificates can be configured, just need to change 1-2 lines of code.

## Performance

The reason I have used gRPC is mainly because of performance concerns. **gRPC** is known to be faster than JSON-based REST frameworks because of its binary *serialization/deserialization* protocol - **Protocol Buffers**. However, I have used gRPC REST Gateway because there was a need of a REST API in the project description.

The REST API may have performance issues regarding the continuous need of binary deserialization/JSON serialization and binary serialization/JSON deserialization processes. I suggest using the `streaming gRPC functions` in the high-perfromance-needed tasks -such as retrieving all the Leak data-.

## Architectural Design Decisions

The biggest decision was for the framework and language choice. I was choosing between Django and gRPC + REST Gateway. I choosed gRPC over Django after several considerations. These are,

* ***Performance**: A gRPC (with `Golang`) service is way better than Django (with `Python3`) in the means of performance. This is mainly because -as I have explained before-, gRPC's **binary data serialization policy**. Also, gRPC uses multithread approach in code level. In Django, one can only use this approach in system level, with multiple workers on Gunicorn for example.

* **Design Patterns**: While Django is a monolithic backend framework, gRPC is closer to be a microservice. Following **microservice patterns** both ease the deployment process and lets the DevOps manager use **lookahead load balancing** with the backend service, this makes it easier to use a load balancer in both client and server side.

* **Production Readiness**: Addition to the upper one, microservices -especially gRPC- are more production-ready espacially when the deal is deploying the service to **serverless -cloud- services** with **multiple instances**. This readiness also benefits the service to have more performance `(Multiple lambda function instances, lookahead load balancing and binary serialization all together)`.

Lastly, although the project needs to serve a **RESTful API**, using a HTTP reverse proxy to perform binary-to-JSON translation does **not** cause a remarkable fall in performance. Mathematically, let the time to perform the JSON serialization of a very big data structure be `t_JSONser`, the time to perform binary serialization be `t_Binser`. Then **barely** the performance ratio of Django to gRPC + REST gateway is `(t_JSONser + t_Binser) / t_JSONser`. This ratio is very close to `1` for **very big data structures**, like in the project description. Now if you add the load balancing, built-in-multithreading and multiple micro-instances on deployment, it is obvious that this **tradeoff** is not a hard choice to make, the gRPC one is way more better.

## Extra: Data Imports

The challenge of **mass** data placement to the system can be solved basically: Implementing a new endpoint on the gRPC service (let's say `CreateLeaks`) with a proper endpoint for the REST gateway.

The critical point here is, as the system deals with **hundreds of millions** of documents, this endpoint should be **streamful** (idk if such a word exists), that should be used to stream into thousands of documents continuously. 

Also using **load balancers** and **multiple microservice instances** are very beneficial here in the means of performance.