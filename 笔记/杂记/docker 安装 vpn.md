# docker 安装 vpn

---

## setup

```shell
106.15.56.49: 是我服务器的外网ip
$CONNECTION_NAME: 是对外vpn的名字

docker run --name dvpn-data -v /etc/openvpn busybox
docker run --volumes-from dvpn-data --rm kylemanna/openvpn ovpn_genconfig -u udp://106.15.56.49:1194
docker run --volumes-from dvpn-data --rm -it kylemanna/openvpn ovpn_initpki

docker run --volumes-from dvpn-data --rm -it kylemanna/openvpn easyrsa build-client-full $CONNECTION_NAME nopass
docker run --volumes-from dvpn-data --rm kylemanna/openvpn ovpn_getclient $CONNECTION_NAME > demo-vpn.ovpn
```

## 制作服务

```shell
vim /etc/systemd/system/dvpn.service
[Unit]
Description=OpenVPN Docker Container
Requires=docker.service
After=docker.service

[Service]
Restart=always
ExecStart=/usr/bin/docker run --name vpn --volumes-from dvpn-data --rm -p 1194:1194/udp --cap-add=NET_ADMIN kylemanna/openvpn
ExecStop=/usr/bin/docker stop vpn

[Install]
WantedBy=local.target

systemctl start dvpn.service
systemctl enable dvpn.service
```

## 客户端

下载`demo-vpn.ovpn`证书即可, 配合`open-vpn`使用.

`Windows`的客户端, 需要使用`管理员模式`启动.