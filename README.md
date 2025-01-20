# gmicro-user

user serivce

## The Clean Architecture

![The Clean Architrecture](misc/images/TheCleanArchitrecture.png)

同心圆代表软件的不同领域，一般来说，越往里走，软件的级别就越高

这种架构发挥作用的首要规则是依赖规则，改规则规定依赖关系只能指向内部，内圈中的任何事物都无法知道外圈中的任何事物。特别是，外圈中声明的事物的名称不得在内圈的代码中提及，包括函数、类、变量或任何其他命名的软件实体。

同样的道理，外层循环使用的数据格式不应该被内层循环使用，尤其是那些由外层循环中的框架生成的格式，我们不希望外层循环中的任何东西影响内层循环。

## gRPC-Gateway

> The gRPC-Gateway is a plugin of the Google protocol buffers compiler protoc. It reads protobuf service definitions and generates a reverse-proxy server which translates a RESTful HTTP API into gRPC. This server is generated according to the google.api.http annotations in your service definitions.

![gRPC-Gateway](/misc/images/grpc-gateway.png)

## SetUp

```bash
SQLITE_DB=data/sqlite.db APPLICATION_GRPC_PORT=8380 APPLICATION_HTTP_PORT=8381 ENV=development JWT_SECRET=bWFjaW50b3NoCg JWT_TOKEN_DURATION=5  go run cmd/main.go
```

```bash

➜ grpcurl -d '{"username": "admin", "password": "admin123"}' -plaintext 192.168.100.16:8380 user.v1.User/Login
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzcwMjg1OTMsImlhdCI6MTczNzAyNjc5Mywic3ViIjoiXHUwMDAxIiwicm9sZSI6ImFkbWluIn0.5GQiDIYYFf9cEs6WVgUQk7kDPemqCZMEdvfKl5II3sE"
}

"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzcwMjU5NTYsImlhdCI6MTczNzAyNTY1Niwic3ViIjoiXHUwMDAxIiwicm9sZSI6ImFkbWluIn0.LhA96aK_bJNBSCWHS-CX41p9xIw7yKh9oL88115-GjQ"
"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzcwMjU5NTYsImlhdCI6MTczNzAyNTY1Niwic3ViIjoiXHUwMDAxIiwicm9sZSI6ImFkbWluIn0.LhA96aK_bJNBSCWHS-CX41p9xIw7yKh9oL88115-GjQ"

chyiyaqing in ~ at HP-EliteDesk-800-G6-Desktop-Mini-PC
➜ grpcurl -H 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzcwMjg1OTMsImlhdCI6MTczNzAyNjc5Mywic3ViIjoiXHUwMDAxIiwicm9sZSI6ImFkbWluIn0.5GQiDIYYFf9cEs6WVgUQk7kDPemqCZMEdvfKl5II3sE' -d '{"username":"user1","email":"user1@gmail.com","phone":"101010", "address": "中国", "password": "user123", "role": "user"}' -plaintext 192.168.100.16:8380 user.v1.User/Create        
{
  "userId": "2"
}


➜ grpcurl -d '{"username": "user1", "password": "user123"}' -plaintext 192.168.100.16:8380 user.v1.User/Login
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzcwMjEyNDksImlhdCI6MTczNzAyMDk0OSwic3ViIjoiXHUwMDAxIn0.1n-BNWWjltr_4sS6W4yvAB1absSUG3rmgIBNCc9nxk0"
}

chyiyaqing in ~ at HP-EliteDesk-800-G6-Desktop-Mini-PC took 2.6s
➜ grpcurl -H 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzcwMjg2NTksImlhdCI6MTczNzAyNjg1OSwic3ViIjoiXHUwMDAyIiwicm9sZSI6InVzZXIifQ.r7FnGV342ky9rDHx-z-GypulBGrL-aW7LawXYlZdpvE' -d '{"userId": 1}' -plaintext 192.168.100.16:8380 user.v1.User/Get
{
  "userId": "2",
  "name": "user1",
  "email": "user1@gmail.com",
  "phone": "101010",
  "address": "中国"
}
```

可以使用`[jwt.io](https://jwt.io/)` 解析token内容

## Interceptor 拦截器

* 服务端拦截器

> gRPC请求在到达实际RPC方法之前将调用的函数,它可用于多种用途, 日志记录、跟踪、速率限制、身份验证和授权

* 客户端拦截器

> 在实际调用RPC之前由gRPC客户端调用的函数




## JWT

JWT 结构分为三个base64url编码部分, [Header: [令牌类型,算法标头]、 Payload: [数据的有效负载]、 Signature: [签名]]
![jwt](/misc/images/jwt.png)

```
Header: {
  "alg": "HS256",
  "typ": "JWT"
}

Payload: {
  "sub": "userID1234", // 标识令牌的主题,通常是用户ID
  "iat": 1516239022, // 令牌的签发时间
  "exp": 1516249022 // 令牌到期时间
}

Signature: {
  Base64UrlSafe (
    HMACSH256(<header>.<payload>,<secret key>)
  )
}
```

__gRPC元数据__ 允许将其他信息附加到请求或响应.类似于HTTP标头，此元数据可用于身份验证、跟踪和其他的目的. gRPC支持两种类型的元数据: `标头`和`尾部`,他们在RPC生命周期的不同阶段发送,标头在消息数据之前发送，而尾部在响应消息数据发送到客户端之后发送.

![metadata](/misc/images/jwt-metadata.png)

* 客户端请求添加元数据

```go
md := metadata.Pairs("authorization", "Bearer your_jwt_token")
// 将元数据对象包含到RPC调用的上下文中
ctx := metadata.NewOutgoingContext(context.Background(), md)
```

* 服务端 从传入的上下文中提取元数据验证令牌或检索必要信息

```go
md, ok := metadata.FromIncomingContext(ctx)
if !ok {
  return nil, status.Error(codes.Unauthenticated, "missing metadata")
}
token := md["authorization"]
```

