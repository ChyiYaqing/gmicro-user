# Design Distribute System

## SideCar

```bash
chyiyaqing in gmicro-user/docker at HP-EliteDesk-800-G6-Desktop-Mini-PC on  main [!+?] via 🐹 v1.22.10 on 🐳 v27.3.1 runs 🐙 T 
➜ docker-compose up -d
WARN[0000] Found orphan containers ([topz]) for this project. If you removed or renamed this service in your compose file, you can run this command with the --remove-orphans flag to clean it up. 
[+] Running 2/2
 ✔ Network docker_default  Created                                                                                                                                                                                0.1s 
 ✔ Container user          Started                                                                                                                                                                                0.2s 

chyiyaqing in gmicro-user/docker at HP-EliteDesk-800-G6-Desktop-Mini-PC on  main [!+?] via 🐹 v1.22.10 on 🐳 v27.3.1 runs 🐙 TU 
➜ docker ps | grep user                                                                                      
e2086c21838e   chyiyaqing/user:v0.0.1-df3c3eb                  "./user"                 4 seconds ago   Up 3 seconds          0.0.0.0:8380-8381->8380-8381/tcp, :::8380-8381->8380-8381/tcp   user
821b5bbe96ed   7bd77774a775                                    "./user"                 4 days ago      Up 4 days                                                                             k8s_user_user-c67fb7c74-k2qll_default_69a72306-56b1-4dc3-a51a-eb960ba574b7_2
95f0e7483fcf   registry.k8s.io/pause:3.9                       "/pause"                 4 days ago      Up 4 days                                                                             k8s_POD_user-c67fb7c74-k2qll_default_69a72306-56b1-4dc3-a51a-eb960ba574b7_18

# 在应用容器相同PID命名空间中启动topz sidecar
chyiyaqing in gmicro-user/docker at HP-EliteDesk-800-G6-Desktop-Mini-PC on  main [!+?] via 🐹 v1.22.10 on 🐳 v27.3.1 runs 🐙 TU 
➜ docker run -d --pid=container:e2086c21838e -p 8080:8080 brendanburns/topz:db0fa58 /server --addr=0.0.0.0:8080
75ae42838cc7fb54bd59f1bc5d179b05a1b0068f7c826e48345b52cc5fe0040a

# 通过topz监控容器资源
chyiyaqing in gmicro-user/docker at HP-EliteDesk-800-G6-Desktop-Mini-PC on  main [!+?] via 🐹 v1.22.10 on 🐳 v27.3.1 runs 🐙 TU 
➜ curl http://192.168.100.16:8080/topz     
1  0 0.05365173  ./user
13 0 0.013289879 /server --addr=0.0.0.0:8080
```