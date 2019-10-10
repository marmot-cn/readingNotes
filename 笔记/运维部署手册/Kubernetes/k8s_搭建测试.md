#k8s 搭建测试

---

###Deployments

相对于以前的RC,官方这次更推荐新出的`Deployments`.在做`roll update`和`roll back`的时候处理会更方便.

####测试deployments的创建

我们创建文档如下用于测试:

		apiVersion: extensions/v1beta1
		kind: Deployment
		metadata:
		  name: nginx-deployment
		spec:
		  template:
		    metadata:
		      labels:
		        app: nginx
		    spec:
		      containers:
		      - name: nginx
		        image: registry.aliyuncs.com/marmot/nginx:1.10
		        ports:
		        - containerPort: 80
		      
保存后,我们创建:

		[ansible@k8s-master-test k8s]$ kubectl create -f ./nginx-deployment.yml --record
		deployment "nginx-deployment" created
		
其中`--record`会记录我们发布的信息,是用于回滚使用的.在后面我们检查`history`时使用.

**检查发布情况**

		[ansible@k8s-master-test k8s]$ kubectl get deployments
		NAME               DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
		nginx-deployment   1         1         1            0           1m
		
因为我们的镜像之前没有使用,所以此时还在下载镜像.

		[ansible@k8s-master-test k8s]$ kubectl get deployment/nginx-deployment
		NAME               DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
		nginx-deployment   1         1         1            1           38s

**logs**

如果我们直接执行 kubectl logs xxx

我们会发现如下错误:

		xxx no such host 
		
这是因为k8s会直接把hostname 作为域名访问使用.所以解决方案是在 `/etc/hosts`里面添加节点对应的ip地址即可.

####测试更新和回滚

**更新配置文件**

我们更新上面的配置文件,把nginx版本修改为`1.11`

		apiVersion: extensions/v1beta1
		kind: Deployment
		metadata:
		  name: nginx-deployment
		spec:
		  template:
		    metadata:
		      labels:
		        app: nginx
		    spec:
		      containers:
		      - name: nginx
		        image: registry.aliyuncs.com/marmot/nginx:1.11
		        ports:
		        - containerPort: 80

**apply更新配置**

		[ansible@k8s-master-test k8s]$ kubectl apply -f ./nginx-deployment.yml
		deployment "nginx-deployment" configured
		
**确认滚动更新**

		kubectl get deployment
		...
		  32s		32s		1	{deployment-controller }			Normal		ScalingReplicaSet	Scaled up replica set nginx-deployment-907840193 to 1
		  32s		32s		1	{deployment-controller }			Normal		ScalingReplicaSet	Scaled down replica set nginx-deployment-829262528 to 0		

看最后的两条记录,就得被降低为0,新的被创建.

**查看历史**

检查历史

		[ansible@k8s-master-test k8s]$ kubectl rollout history deployment/nginx-deployment
		deployments "nginx-deployment":
		REVISION	CHANGE-CAUSE
		1		kubectl create -f ./nginx-deployment.yml --record
		2		kubectl apply -f ./nginx-deployment.yml

检查单条的历史:

		[ansible@k8s-master-test k8s]$ kubectl rollout history deployment/nginx-deployment --revision=2
		deployments "nginx-deployment" revision 2
		  Labels:	app=nginx,pod-template-hash=907840193
		  Annotations:	kubernetes.io/change-cause=kubectl apply -f ./nginx-deployment.yml
		  Containers:
		  nginx:
		    Image:	registry.aliyuncs.com/marmot/nginx:1.11
		    Port:	80/TCP
		    QoS Tier:
		      cpu:	BestEffort
		      memory:	BestEffort
		    Environment Variables:
		  No volumes.
		[ansible@k8s-master-test k8s]$ kubectl rollout history deployment/nginx-deployment --revision=1
		deployments "nginx-deployment" revision 1
		  Labels:	app=nginx,pod-template-hash=829262528
		  Annotations:	kubernetes.io/change-cause=kubectl create -f ./nginx-deployment.yml --record
		  Containers:
		  nginx:
		    Image:	registry.aliyuncs.com/marmot/nginx:1.10
		    Port:	80/TCP
		    QoS Tier:
		      cpu:	BestEffort
		      memory:	BestEffort
		    Environment Variables:
		  No volumes.
		
其中我们可以看见我们的镜像从1.10 变为1.11

**实现回滚操作**

		[ansible@k8s-master-test k8s]$ kubectl rollout undo deployment/nginx-deployment
		deployment "nginx-deployment" rolled back
		
**修改历史记录**

可以修改字段`revisionHistoryLimit`修改回滚历史,修改为5.

		apiVersion: extensions/v1beta1
		kind: Deployment
		metadata:
		  name: nginx-deployment
		spec:
		  revisionHistoryLimit: 5
		  template:
		    metadata:
		      labels:
		        app: nginx
		    spec:
		      containers:
		      - name: nginx
		        image: registry.aliyuncs.com/marmot/nginx:1.11
		        ports:
		        - containerPort: 80
		
####测试hpa

**创建hpa文件**

		apiVersion: extensions/v1beta1
		kind: HorizontalPodAutoscaler
		metadata:
		  name: nginx-deployment
		spec:
		  scaleRef:
		    kind: Deployment
		    name: nginx-deployment
		    subresource: scale
		  minReplicas: 1
		  maxReplicas: 2
		  cpuUtilization:
		    targetPercentage: 50
		    
**创建hpa**

		[ansible@k8s-master-test k8s]$ kubectl create -f ./nginx-hpa.yml
		horizontalpodautoscaler "nginx-deployment" created
		
		[ansible@k8s-master-test k8s]$ kubectl get hpa
		NAME               REFERENCE                           TARGET    CURRENT     MINPODS   MAXPODS   AGE
		nginx-deployment   Deployment/nginx-deployment/scale   50%       <waiting>   1         2         48s


###常见问题

####强制删除`Terminating`状态的`pod`

		kubectl delete pod pod名称 --grace-period=0 --namespace='命名空间'
