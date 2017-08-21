# 容器单机收集nginx容器日志并logrotate

---

### Nginx 接收`USR1`信号

NGINX will re-open its logs in response to the `USR1` signal.

### 更新`docker-compose.yml`

把`nginx`的日志文件夹挂载到宿主机上.`- /xxxx:/var/log/nginx`

我挂载单独一个文件夹下, 并把用户和用户组都给为`www-data:www-data`;

### 更新宿主机 `logrotate.d/nginx`

因我们需要回滚日志.

```shell
root@5896e93a5e9e:/var/www/html# cat /etc/logrotate.d/nginx
/日志路径(你容器日志也就是你挂载的日志目录的路径)/*.log {
        daily
        missingok
        rotate 52
        compress
        delaycompress
        notifempty
        create 640 www-data www-data
        sharedscripts
        postrotate
        		#检查容器是否存在
				docker ps | grep demo-nginx > /dev/null 2>&1
				#如果容器存在
				if [ $? -eq 0 ]; then
					docker kill -s USR1 demo-nginx
                fi
        endscript
}
```

### 测试



