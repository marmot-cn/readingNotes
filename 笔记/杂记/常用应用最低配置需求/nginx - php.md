# nginx - php

---

服务器上使用`docker`虚拟化部署`nginx-phpfpm`. 因为是核心应用处理, 默认镜像`nginx`是按照4核心优化的. 所以这里暂时考虑使用**4核心8G内存**, 内存是因为需要部署`fluentd`和`memcached`