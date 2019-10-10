/*
 Author:chloroplast
 */
#Collection+Json

---

[Collection+Json手册][CollectionJsonUrl]

[CollectionJsonUrl]:http://amundsen.com/media-types/collection/

格式:

vnd.collection+json
		

一种注册为Collection + Json的媒体类型.

		{ "collection" :
  		  {
    		"version" : "1.0",
    		"href" : "http://example.org/friends/",
    
    		"links" : [
      		{"rel" : "feed", "href" : "http://example.org/friends/rss"}
    		],
    
    		"items" : [
      		{
       		 "href" : "http://example.org/friends/jdoe",
        		"data" : [
          		{"name" : "full-name", "value" : "J. Doe", "prompt" : "Full Name"},
          		{"name" : "email", "value" : "jdoe@example.org", "prompt" : "Email"}
       		 ],
        	 "links" : [
          		{"rel" : "blog", "href" : "http://examples.org/blogs/jdoe", "prompt" : "Blog"},
         		 {"rel" : "avatar", "href" : "http://examples.org/images/jdoe", "prompt" : "Avatar", "render" : "image"}
       		 ]
     		 },
      
     	 	{
        		"href" : "http://example.org/friends/msmith",
        		"data" : [
          		{"name" : "full-name", "value" : "M. Smith", "prompt" : "Full Name"},
          		{"name" : "email", "value" : "msmith@example.org", "prompt" : "Email"}
       			 ],
        		"links" : [
          		{"rel" : "blog", "href" : "http://examples.org/blogs/msmith", "prompt" : "Blog"},
          		{"rel" : "avatar", "href" : "http://examples.org/images/msmith", "prompt" : "Avatar", "render" : "image"}
        		]
      		},
      
      		{
        		"href" : "http://example.org/friends/rwilliams",
        		"data" : [
          		{"name" : "full-name", "value" : "R. Williams", "prompt" : "Full Name"},
          		{"name" : "email", "value" : "rwilliams@example.org", "prompt" : "Email"}
        		],
        		"links" : [
          		{"rel" : "blog", "href" : "http://examples.org/blogs/rwilliams", "prompt" : "Blog"},
          		{"rel" : "avatar", "href" : "http://examples.org/images/rwilliams", "prompt" : "Avatar", "render" : "image"}
        		]
      		}      
    		],
    
    		"queries" : [
      		 {"rel" : "search", "href" : "http://example.org/friends/search", "prompt" : "Search",
        	"data" : [ {"name" : "search", "value" : ""}]
      		 }
    		],
    
    		"template" : {
      		 "data" : [
       		  {"name" : "full-name", "value" : "", "prompt" : "Full Name"},
       		  {"name" : "email", "value" : "", "prompt" : "Email"},
       		  {"name" : "blog", "value" : "", "prompt" : "Blog"},
       		  {"name" : "avatar", "value" : "", "prompt" : "Avatar"}
      		 ]
    		}
  		} 
		}

**href**

一个指向集合本身的永久连接.

**items**

包含指向集合成员的连接,以及他们的部分表述.

**links**

指向其他与集合相关资源的连接.

**queries**

用于搜索集合的超媒体控件.

**template**

用于向集合添加子项的超媒体控件.

###子项的表示

**href**

一个指向子项的永久链接,这个子项被作为独立的资源看待.如果你向href属性中的URL发起GET连接请求,服务器将会向你发送该单个子项的Collection+JSON表述.

**links**

指向子项相关的其他资源的超媒体连接.`rel`属性是一个为连接关系准备的数据槽.制定从源文档到目标文档的关系.`prompt`属性是一个用于放置人类可读的描述信息的地方.


**data**

任何其他区信息, 这是子项表述的一个重要部分.`name`和`value`属性描述为单个键值对.`name`属性是该键值对的键,而`value`是值.`prompt`值是一段人类可读的描述信息.

###写入模板

		"template" : {
         "data" : [
          {"name" : "full-name", "value" : "", "prompt" : "Full Name"},
          {"name" : "email", "value" : "", "prompt" : "Email"},
          {"name" : "blog", "value" : "", "prompt" : "Blog"},
          {"name" : "avatar", "value" : "", "prompt" : "Avatar"}
         ]
        }
向该集合href属性中的url发送一个post请求才能向该集合添加一个子项.

等价于
		
		<form action="http://example.org/friends" method="post">
		<p>Search</p>
		<label for="name">name</label>
		<input lable="Full Name" id="full-name" name="full-name" value="" />
		<input lable="Email" id="email" name="email" value="" />
		<input lable="Blog" id="blog" name="blog" value="" />
		<input lable="Avatar" id="avatar" name="avatar" value="" />
		</form>	
			
###搜索模板
	
	"queries" : [
         {"rel" : "search", "href" : "http://example.org/friends/search", "prompt" : "Search",
        "data" : [ {"name" : "search", "value" : "", "prompt":"search input"}]
         }
    ],

向该`href`发起`GET`请求,使用`data`里面键值.

		<form action="http://example.org/friends/search" method="get">
		<p>Search</p>
		<label for="name">name</label>
		<input lable="search input" id="search" name="search" value="" />
		</form>

###错误模板

	{
  		"error" :
  		{
    		"title" : "xxxx",
    		"code" : "203",
   			"message" : "xxxx"  
  		}
	}

###分页

用于在分页列表中导航的通用链接关系:

* next
* previous
* first
* last
* prev

这些连接关系原本是为HTML定义的,但是现在它们都在IANA上注册过,所以你可以在任何媒体类型中使用它们.

		"links" : [ {
		  "name" : "next_page",
		  "prompt" : "Next",
		  "rel" : "next",
		  "href" : "/collection/page/3",
		  "render" : "link"
		 },
		 {
		  "name" : "previous_page",
		  "prompt" : "Back",
		  "rel" : "previous",
		  "href" : "/collection/page/1",
		  "render" : "link"
		 }
		]


###对象

**1.collection**
	
	// sample collection object
	{
 	 	"collection" :
  		{
    		"version" : "1.0",
   	 		"href" : URI,
    		"links" : [ARRAY],
    		"items" : [ARRAY],
    		"queries" : [ARRAY],
    		"template" : {OBJECT},
    		"error" : {OBJECT}
  		}
	}

* `REQUIRED` object  
* `MUST NOT` be more than one collection object

**2.error**

	// sample error object
	{
  	 	"error" :
  		{
    		"title" : STRING,
    		"code" : STRING,
    		"message" : STRING  
  		}
	}
	
* `OPTIONAL` object  
* `MUST NOT` be more than one error object

**3.template**

	// sample template object
	{
  		"template" :
  		{
    		"data" : [ARRAY]  
 		 }
	}

* `OPTIONAL` object  
* `MUST NOT` be more than one template  object


###数组

`items`,`data`,`queries` and `links` 这些类型用数组表示

**1.items**

	// sample items array
	{
  		"collection" :
  		{
   		 	"version" : "1.0",
    		"href" : URI,
   			"items" :
    		[
     		 {
       			"href" : URI,
        		"data" : [ARRAY],
        		"links" : [ARRAY]
      		 },
     		 ...
     		 {
        		"href" : URI,
       			"data" : [ARRAY],
        		"links" : [ARRAY]
      		 }
    		]
  		}
	}

**2.data**

	// sample data array
	{
 	 "template" :
  		{
   	 		"data" :
    		[
      			{"prompt" : STRING, "name" : STRING, "value" : VALUE},
      			{"prompt" : STRING, "name" : STRING, "value" : VALUE},
      			...
      			{"prompt" : STRING, "name" : STRING, "value" : VALUE}
    		]
  		}
	}

* `REQUIRED` name
* `OPTIONAL` value
* `OPTIONAL` prompt

**3.queries**

	// sample queries array
	{
	 	"queries" :
  		[
    		{"href" : URI, "rel" : STRING, "prompt" : STRING, "name" : STRING},
    		{"href" : URI, "rel" : STRING, "prompt" : STRING, "name" : STRING,
      			"data" :
      			[
       	 		{"name" : STRING, "value" : VALUE}
      			]
    		},
   		   ...
    		{"href" : URI, "rel" : STRING, "prompt" : STRING, "name" : STRING}
  		]
	 }
	 
**4.links**

`links` 是 `OPTIONAL`

	// sample links array
	{
	  "collection" :
	  {
	    "version" : "1.0",
	    "href" : URI,
	    "items" :
	    [
	      {
	        "href" : URI,
	        "data" : [ARRAY],
	        "links" :
	        [
	          {"href" : URI, "rel" : STRING, "prompt" : STRING, "name" : STRING, "render" : "image"},
	          {"href" : URI, "rel" : STRING, "prompt" : STRING, "name" : STRING, "render" : "link"},
	          ...
	          {"href" : URI, "rel" : STRING, "prompt" : STRING, "name" : STRING}
	        ]
	      }
	    ]
	  }
	}
		
* `REQUIRED` href
* `OPTIONAL` rel
* `OPTIONAL` name
* `OPTIONAL` render
* `OPTIONAL` prompt


###属性

**1.code**

code 是属于 error object. 

`SHOULD` be a `STRING` type.

**2.href**

`MUST` contain a valid `URI`.

**3.message**

message 属于 error object.

`SHOULD` be a `STRING` type. 

**4.name**

`SHOULD` be a `STRING` type. 

**5.prompt**

`SHOULD` be a `STRING` type. 

**6.rel**

`SHOULD` be a `STRING` type. 

**7.render**

`SHOULD` be a `STRING` type. 

`MUST` be either `image` or `link`.

**8.title**

`SHOULD` be a `STRING` type. 

**9.value**

`MAY` contain one of the following data types: `STRING` `NUMBER` `true` `false` or `null`

**10.version**

`SHOULD` be a `STRING` type. 

###数据类型

**1.OBJECT**

An OBJECT structure is represented as a pair of curly brackets surrounding zero or more name/value pairs (or members). A name is a string. A single colon comes after each name, separating the name from the value. A single comma separates a value from a following name. The names within an object SHOULD be unique. 

**2.ARRAY**

An ARRAY structure is represented as square brackets surrounding zero or more values (or elements). Elements are separated by commas.

**3.NUMBER**

 A NUMBER contains an integer component that may be prefixed with an optional minus sign, which may be followed by a fraction part and/or an exponent part.

Octal and hex forms are not allowed. Leading zeros are not allowed.

A fraction part is a decimal point followed by one or more digits.

An exponent part begins with the letter E in upper or lowercase, which may be followed by a plus or minus sign. The E and optional sign are followed by one or more digits. 

**4.STRING**

A STRING begins and ends with quotation marks. All Unicode characters may be placed within the quotation marks except for the characters that must be escaped: quotation mark, reverse solidus, and the control characters (U+0000 through U+001F). 

**5.URI**

A URI is defined by RFC 3986

**6.VALUE**

A VALUE data type MUST be a NUMBER, STRING, or one of the following literal names: false, null, or true. 

##通用方法示例

---

###Reading Collections

		*** REQUEST ***
		GET /my-collection/ HTTP/1.1
		Host: www.example.org
		Accept: application/vnd.collection+json

		*** RESPONSE ***
		200 OK HTTP/1.1
		Content-Type: application/vnd.collection+json
		Content-Length: xxx

		{ "collection" : {...}, ... }

###Adding an Item

		*** REQUEST ***
		POST /my-collection/ HTTP/1.1
		Host: www.example.org
		Content-Type: application/vnd.collection+json

		{ "template" : { "data" : [ ...] } }
 
		*** RESPONSE ***
		201 Created HTTP/1.1
		Location: http://www.example.org/my-collection/1

###Readning an Item

		*** REQUEST ***
		GET /my-collection/1 HTTP/1.1
		Host: www.example.org
		Accept: application/vnd.collection+json

		*** RESPONSE ***
		200 OK HTTP/1.1
		Content-Type: application/vnd.collection+json
		Content-Length: xxx

		{ "collection" : { "href" : "...", "items" : [ { "href" : "...", "data" : [...] } } }

###Updating an Item
		
		*** REQUEST ***
		PUT /my-collection/1 HTTP/1.1
		Host: www.example.org
		Content-Type: application/vnd.collection+json

		{ "template" : { "data" : [ ...] } }
 
		*** RESPONSE ***
		200 OK HTTP/1.1 

###Deleting an Item

		*** REQUEST ***
		DELETE /my-collection/1 HTTP/1.1
		Host: www.example.org
 
		*** RESPONSE ***
		204 No Content HTTP/1.1 
		

