#docker 私有仓库

---

**下载`registry`镜像**

		docker pull registry
		
会将仓库存放于容器内的`/tmp/registry`目录下,这样如果容器被删除,则存放于容器中的镜像也会丢失,所以我们一般情况下会指定本地一个目录挂载到容器内的`/tmp/registry`下.

**编写`dokcer-compose.yml`**

		registry:
 		 image: "registry:latest"
 		 ports:
  		  - "5000:5000"
 		 volumes:
          - "/data/docker-registry:/tmp/registry"
         container_name: self-registry
         
 
**`push`镜像**

`tag`本地镜像

  		docker images
  		mysql                                 5.6                 1a69fcdb790d        13 days ago         324.3 MB
  		
  		docker tag mysql:5.6 120.25.161.1:5000/mysql:5.6
  		docker push 120.25.161.1:5000/mysql
  		
  		Pushing repository 120.25.161.1:5000/mysql (1 tags)
		6d1ae97ee388: Image successfully pushed
		8b9a99209d5c: Image successfully pushed
		410c2fae2283: Image successfully pushed
		e3a6552a83c2: Image successfully pushed
		1b0e180fd8fa: Image successfully pushed
		0d5f060b62c4: Image successfully pushed
		2e8a186e254e: Image successfully pushed
		dc1434565071: Image successfully pushed
		75c8c65fc91d: Image successfully pushed
		f9c4e0df39dc: Image successfully pushed
		0fbaa1d94bf3: Image successfully pushed
		5287fd6e217a: Image successfully pushed
		eaeca35b15e4: Image successfully pushed
		7aab961bc74d: Image successfully pushed
		edc3302a6b2e: Image successfully pushed
		1a69fcdb790d: Image successfully pushed
		Pushing tag for rev [1a69fcdb790d] on {http://120.25.161.1:5000/v1/repositories/mysql/tags/5.6}
		
**查看仓库**

		curl http://120.25.161.1:5000/v1/search
		
		{"num_results": 3, "query": "", "results": [{"description": "", "name": "library/gitlab"}, {"description": "", "name": "library/gitlab-postgresql"}, {"description": "", "name": "library/gitlab-redis"}]}
		
		
####docker配置修改在Centos7上

**修改配置**

修改`/etc/sysconfig/docker`:

		INSECURE_REGISTRY='--insecure-registry 120.25.87.35:5000'
		
这个IP地址是仓库的ip地址

**重启docker服务**

		service docker restart
		
**上传镜像**

重新`tag`镜像:

		docker tag chloroplast/php-helloworld 120.25.87.35:5000/php-test
		
`push`镜像:

		docker push 120.25.87.35:5000/php-test

测试另外一台服务器`pull`镜像:

		root@iZ94xwu3is8Z ~]# docker pull 120.25.87.35:5000/php-test
		Using default tag: latest
		970c5a38b506: Pulling dependent layers
		cb6fb082434e: Downloading [=>                                                 ] 1.032 MB/51.26 MB
