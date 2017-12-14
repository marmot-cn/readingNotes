# prometheus

## 需要开发事项

1. 同步时间为中国时间, 否则时间序列监控会出现问题.
2. 监控主机`node_exporter`.
	* 网络流量监控?
	* 服务器监控
	* 运行结果统计分析(prometheus)
3. 监控容器`cadvisor`.
	* 应用服务监控
	* 运行结果统计分析(prometheus)
4. 运维服务审计`auditd`,https://linux.cn/article-4907-1.html
	* 运维服务记录