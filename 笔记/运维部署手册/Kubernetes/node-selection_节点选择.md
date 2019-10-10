#node-selection 节点选择

---

**选择`node`**

运行`kubectl get nodes`获取`node`名称.

**添加`label`**

运行`kubectl label nodes <node-name> <label-key>=<label-value>`.

可以使用参数实现`覆盖`,`删除`等操作.

**添加`nodeSelector`字段在配置文件**

在`pod`的配置文件添加:

		nodeSelector:
    		label-key: label-value
    		
**示例**

		[chloroplast@dev-server-2 docker-yaml]$ kubectl get nodes
		NAME            LABELS                                                     STATUS
		10.116.138.44   kubernetes.io/hostname=10.116.138.44,server=dev-server-1   Ready
		10.44.88.189    kubernetes.io/hostname=10.44.88.189,server=dev-server-3    Ready
		
我的node已经添加了lable: `server=dev-server-1`和`server=dev-server-3`

配置文件使用`nodeSelector`选择`node`部署:

		[chloroplast@dev-server-2 docker-yaml]$ cat nginx-test-rc.yaml
		apiVersion: v1
		kind: ReplicationController
		metadata:
		  name: nginx-test
		  labels:
		    name: nginx-test
		spec:
		  replicas: 1
		  selector:
		    name: nginx-test
		  template:
		    metadata:
		      labels:
		        name: nginx-test
		    spec:
		      containers:
		      - name: nginx-test
		        image: 120.25.87.35:5000/nginx:1.9
		        volumeMounts:
		        - name: data
		          mountPath: /usr/share/nginx/html
		        ports:
		        - containerPort: 80
		        hostPort: 80
		      nodeSelector:
		        server: dev-server-1
		      volumes:
		      - name: data
		        hostPath:
		          path: /data/html
		          
检查部署情况:

		[chloroplast@dev-server-2 docker-yaml]$ kubectl get pods -o=wide
		NAME               READY     STATUS    RESTARTS   AGE       NODE
		nginx-test-mr7vy   1/1       Running   0          2h        10.116.138.44
		
10.116.138.44 是我 `dev-server-1`外网的ip,现在确认已经部署在这个`node`上了.
