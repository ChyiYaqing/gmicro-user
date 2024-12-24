# gmicro-user
user serivce

## SetUp

```bash

chyiyaqing in ~ at HP-EliteDesk-800-G6-Desktop-Mini-PC
➜ grpcurl -d '{"name":"chyiyaqing","email":"chyiyaqing@gmail.com","phone":"101010", "address": "中国"}' -plaintext 192.168.100.16:8380 User/Create
{
  "userId": "1"
}

chyiyaqing in ~ at HP-EliteDesk-800-G6-Desktop-Mini-PC took 2.6s
➜ grpcurl -d '{"userId": 1}' -plaintext 192.168.100.16:8380 User/Get
{
  "userId": "1",
  "name": "chyiyaqing",
  "email": "chyiyaqing@gmail.com",
  "phone": "101010",
  "address": "中国"
}
```