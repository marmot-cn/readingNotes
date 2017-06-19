### container-crontab 计划任务

该计划任务容器是属于官方镜像.

[github](https://github.com/rancher/container-crontab "crontab")

部署该容器以后,使用`--label`标签即可添加计划任务.

覆盖默认的`start`状态,可以设置`label` `cron.action`等于:

* `stop`
* `restart`

中的一种.

覆盖默认的 10s 重启时间,可以设置`label` `cron.restart_timeout=秒数` 覆盖默认的10秒.


