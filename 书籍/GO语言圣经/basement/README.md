# README

### 1. 获取go镜像

```
docker pull golang:1.19.0-buster
```

稳定的Debian发行版是10.4，它的代号是“buster”。 “stretch”是所有版本9变种的代号，“jessie”是所有版本8变种的代号。 正在开发的未来版本是“bullseye ”和“bookworm”，但还不稳定。

### 2. 制作`docker-compose`

```
version: "3"

services:
  golang:
    image: "golang:1.19.0-buster"
    ports:
      - "8080:8080"
    volumes:
      - ./:/go
    tty: true
    container_name: go-env
```