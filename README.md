# gobox
## upload

Upload a file in restricted linux environment
```shell
./upload https://bashupload.com file
```

## netstat

`netstat` in restricted linux environment
```shell
./netstat
========= tcp =========
Listening on 127.0.0.53:53
Listening on 0.0.0.0:443
Listening on 0.0.0.0:22
Listening on 0.0.0.0:81
Listening on 0.0.0.0:80
Listening on 0:::::::53
========= tcp =========

========= udp =========
Listening on 0:::::::53
Listening on :::0:323
========= udp =========

========= unix socket =========
....
Unix Socket Path: /run/docker.sock, RefCount: 00000000, Protocol: ....
========= unix socket =========
```
