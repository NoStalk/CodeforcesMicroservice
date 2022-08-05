
# Microservice for handling Codeforces 

Built entirely on Golang, this microservice is a part of a larger cluster of services
that handle each of the platform data requests and retrievals of users signed up on NoStalk.
## Work-Flow

- This service handles the backend requests for users with a codeforces handle.

- Since, codeforces exposes a public API to query about a particular user, only the user handle is enough to retrieve
    the desired data (basic info, submissions, contests).

- The individual requests are REST compliant and made through the Get function of the HTTP module of GO.
    The JSON response is recieved and the data is read from the body of the response that is then unmarshalled into desired data structs and interfaces.

- The service itself exposes three [gRPC](https://grpc.io/) endpoints for backend-service communication, from which two support Unary Request-Response and the third one supports Bi-Directional streaming between the server(this service) and the client(backend or any other service).

- The communication itself is done through [Protocol Buffers](https://developers.google.com/protocol-buffers/docs/overview). Protocol buffers provide a language-neutral, platform-neutral, extensible mechanism for serializing structured data in a forward-compatible and backward-compatible way. It’s like JSON, except it's smaller and faster, and it generates native language bindings.

- The data is then written to the specific user's database entry through a [package](https://github.com/NoStalk/serviceUtilities) that has functions that facilitate opening/closing database connections and appending data to already existing entries.





## Acknowledgements

 - [Google Remote Procedure Call](https://grpc.io/)
 - [Protocol Buffers](https://developers.google.com/protocol-buffers/docs/overview)
 - [Using gRPC and Protocol Buffers](https://medium.com/aspnetrun/using-grpc-in-microservices-for-building-a-high-performance-interservice-communication-with-net-5-11f3e5fa0e9d)
 - [BloomRPC](https://github.com/bloomrpc/bloomrpc)
 - [QuickType](https://quicktype.io/)
