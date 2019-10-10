/*
 Author:chloroplast
 */
#Json-LD

---

###传统接口的缺点

JSON-LD 让你可以将一个被称为`上下文(context)`的机器可以读的文档与普通的JSON文档结合起来.

我们先看一个传统的API锁提供的未加修饰的JSON表述:

		HTTP/1.1 200 ok
		Content-Type: application/json
		
		{"n":"Jenny Gallegos",
		 "photo_link":"http://api.example.com/img/xxx"}
		 
从上面这个接口我们能得到如下信息:

1. 一个语意描述符`n`
2. `photo_link` 和一个**貌似**是连接url的连接`http://api.example.com/img/xxx`

这里为什么说貌似是url的连接,因为`application/json`媒体类型并没有超媒体控件,这个连接仅仅是一个正好看上去很像URL的字符串.

所以我们得到如下结论:

1. `n` 是字符串
2. `photo_link`是字符串(尽管张的很像连接)

###传统接口的解决方式

我们可以通过将API所提供的每个文档连接到它的人类可读的profile来稍微改善下这一状况

		HTTP/1.1 200 ok
		Content-Type: application/json
		Link: <http://help.example.com/api/>;rel="profile"
		...

但是这并不能带来太多的帮助.

###JSON-LD的优势

JSON-LD会怎么来处理这些问题. 提供一个链向JSON-LD上下文的连接. 

		HTTP/1.1 200 ok
		Content-Type: application/json
		Link: <http://help.example.com/person.jsonld/>;rel="http://www.w3.org/ns/json-ld#context"

通过向help.example.com/person.jsonld发起第二个 HTTP GET 请求, 你将会找到对应的上下文.

		HTTP/1.1 200 ok
		Content-Type: application/ld+json
		{
		 "@context":
		  {
		   "n": "http://api.example.org/docs/Person#name",
		   
		   "photo_link":
		   {
		    "@id":"http://api.example.org/docs/Person#photo_link"           		    "@type":"@id"
		    }
		   }
		 }
* `@context`的JSON对象都可以是JSON-LD上下文  
* `@id`通常意味着"超媒体连接".对象的`@id`属性是一个链向该词条的`应用语义说明`的连接.  
* `@type`属性同样也是`@id`.这就是将JSON文档转化成超媒体文档的魔法.将phtot_link的`@type`设置为`@id`表明了JSON文档中的photo_link无论在何时出现,客户端都可以将它作为一个超媒体链接来处理,而并非当做一个只是看上去像URL的字符串.


**如何理解上述例子**

通过JSON-LD对最初的JSON文档中的n属性的应用语义进行说明.

		"n": "http://api.example.org/docs/Person#name"
		"n": "Jenny Gallegos"

可以访问`http://api.example.org/docs/Person#name`并阅读相应的说明.

**常用标签**

* @context   	
* @id
* @value
* @language
* @type
* @container
* @list
* @set
* @reverse
* @index
* @base
* @vocab
* @graph


###总结

计算机可以通过将JSON文档和JSON-LD上下文结合起来从而识别出超媒体链接.


更多参考 [官方手册][id]
[id]:http://json-ld.org/