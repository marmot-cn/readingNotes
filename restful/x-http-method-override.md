#X-HTTP-Method-Override

---

While not normally an issue with thick clients, accessing full RESTful capabilities of available services via browsers often is problematic as many (if not all) browsers only allow a form to GET or POST. They don't allow for other HTTP methods, like HEAD, PUT, or DELETE. Google realized this and offers a solution, which is to add a header to the HTTP request, X-HTTP-Method-Override, that is supposed to be interpreted by the service and acted upon regardless of the actual HTTP method used.

简单说就是为了支持RESTFUL,一般服务器都只接收 `GET` 和 `POST`.所以Google提供一个解决方案,添加了`X-HTTP-Method-Override`头,也就是忽略了实际传输方法,而是使用`X-HTTP-Method-Override`的头信息来覆盖.

也就是实际使用的时`POST`方法, 但是这个头设置了:

		X-HTTP-Method-Override: PUT
		
服务端程序应该认为是PUT方法.