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

## Performance

The reason I have used gRPC is mainly because of performance concerns. **gRPC** is known to be faster than JSON-based REST frameworks because of its binary *serialization/deserialization* protocol - **Protocol Buffers**. However, I have used gRPC REST Gateway because there was a need of a REST API in the project description.

The REST API may have performance issues regarding the continuous need of binary deserialization/JSON serialization and binary serialization/JSON deserialization processes. I suggest using the `streaming gRPC functions` in the high-perfromance-needed tasks -such as retrieving all the Leak data-.

## Architectural Design Decisions

--