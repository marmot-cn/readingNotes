#JSON API

[http://jsonapi.org](http://http://jsonapi.org/ "http://jsonapi.org")

---

###Content Negotiation

###内容协商

####Client Responsibilities

Clients **MUST** send all JSON API data in request documents with the header `Content-Type: application/vnd.api+json` without any media type parameters.

Clients that include the JSON API media type in their `Accept` header **MUST** specify the media type there at least once without any media type parameters.

Clients **MUST** ignore any parameters for the application/vnd.api+json media type received in the Content-Type header of response documents.
Server Responsibilities

Servers **MUST** send all JSON API data in response documents with the header `Content-Type: application/vnd.api+json` without any media type parameters.

Servers **MUST** respond with a `415 Unsupported Media Type` status code if a request specifies the header `Content-Type: application/vnd.api+json` with any media type parameters.

Servers **MUST** respond with a `406 Not Acceptable` status code if a request's Accept header contains the JSON API media type and all instances of that media type are modified with media type parameters.

	Note: The content negotiation requirements exist to allow future versions of this specification to use media type parameters for extension negotiation and versioning.
		

####客户端响应

客户端**必须**在传输JSON API格式的请求文档时附带一个`Content-Type: application/vnd.api+json`头且没有任何其他媒体类型(参数).
 
客户在引用JSON API媒体类型在其`Accept`头内**必须**至少标注一次媒体类型且不带任何其他媒体类型(参数).

服务端**必须**在传输JSON API格式的相应文档时附带一个`Content-Type: application/vnd.api+json`头且没有任何其他媒体类型(参数).

服务端**必须**响应一个`415 Unsupported Media Type`状态码如果一个请求头`Content-Type: application/vnd.api+json`附带其他媒体类型(参数).

服务端**必须**响应一个`406 Not Acceptable`状态码如果一个请求接受头包含JSON API媒体类型并但是请求的资源的内容特性无法满足,因而无法生成响应实体.

####Document Structure

This section describes the structure of a JSON API document, which is identified by the media type application/vnd.api+json. JSON API documents are defined in JavaScript Object Notation (JSON) [RFC7159].

Although the same media type is used for both request and response documents, certain aspects are only applicable to one or the other. These differences are called out below.

Unless otherwise noted, objects defined by this specification MUST NOT contain any additional members. Client and server implementations MUST ignore members not recognized by this specification.

    Note: These conditions allow this specification to evolve through additive changes.

####文档结构

这段篇章描述了一个由`application/vnd.api+json`媒体类型声明的JSON API文档的结构,
JSON API文档是由 JavaScript Object Notation(JSON)定义的.

尽管相同的媒体类型被用于请求(request)和响应(reponse)文档,但是一些特定方面仅仅适用于其中请求或者响应.这些不同点将在下面阐述.

除非别特别表述,规范定义的对象**必须不能**包含其他字段.客户端和服务端**必须**忽略除了规范意外的字段.

####Top Level

A JSON object **MUST** be at the root of every JSON API request and response containing data. This object defines a document's "top level".

A document **MUST** contain at least one of the following top-level members:

* `data`: the document's "primary data"
* `errors`: an array of error objects
* `meta`: a meta object that contains non-standard meta-information.

The members `data` and `errors` **MUST NOT** coexist in the same document.

A document **MAY** contain any of these top-level members:

* `jsonapi`: an object describing the server's implementation
* `links`: a links object related to the primary data.
* `included`: an array of resource objects that are related to the primary data and/or each other ("included resources").

If a document does not contain a top-level data key, the included member **MUST NOT** be present either.

The top-level links object **MAY** contain the following members:

* `self`: the link that generated the current response document.
* `related`: a related resource link when the primary data represents a resource relationship.
* `pagination` links for the primary data.

The document's "primary data" is a representation of the resource or collection of resources targeted by a request.

Primary data **MUST** be either:

* a single resource object, a single resource identifier object, or null, for requests that target single resources
* an array of resource objects, an array of resource identifier objects, or an empty array ([]), for requests that target resource collections

For example, the following primary data is a single resource object:

	{
  	  "data": {
     	"type": "articles",
     	"id": "1",
     	"attributes": {
      	// ... this article's attributes
     	},
     	"relationships": {
      	// ... this article's relationships
     	}
      }
	}

The following primary data is a single resource identifier object that references the same resource:

	{
  		"data": {
    	"type": "articles",
   		"id": "1"
  		}
	}

A logical collection of resources **MUST** be represented as an array, even if it only contains one item or is empty.

####Top Level:顶级目录

一个 JSON 对象 **必须** 在每个 JSON API 的请求和响应的数据中的`根`.这个对象被定义为文档的"顶级".

一个文档 **必须** 包含至少下列中得一个"顶级"成员:

* `data`: 文档的`primary data`(主要数据).
* `errors`: 一组数组的`error objects`(错误对象).
* `meta`: 一个`meta object`(元数据)包含非标准的元信息.

成员`data`和`errors` **必须不能** 同时存在于一个相同的文档.

一个文档 **可以** 包含下列中的"顶级"成员:

* `jsonapi`: 一个对象表述服务的设计实现
* `links`: 一个`links object`(链接对象)和`primary data`(主要数据)管理.
* `included`: 一个`resource objects`(资元对象)的数组.和`primary data`(主要数据)相互关联.

如果一个文档没有包含顶级成员`data`键,则`included`成员 **必须不能** 出现.(`data`和`included`总是相辅相成的).

"顶级"成员 `link object`(链接对象) **可以** 包含下列成员:

* `self`: 这个`link`(链接)用于表示当前文档.
* `related`: 一个`related resource link`(关联资源链接)用于表述资源的关联.
* `paginarion`: `primary data`(主要数据)的分页.

文档的"`primary data(主要数据)`"是一个用于请求表现一个目标或者一组目标的资源.

`primary data`(主要数据) **必须** 是下列之一:

* 请求一个单独的资元对象时: 一个单独的`resource object`资元对象,一个单独的`resource identifier object`(资源标示对象)或者 `null(空)`.
* 请求一组资元对象的集合时: 一组资元对象, 一组资源标示对象 或者空的数组 `[]`.

如下示例是表述一个单独的资元对象的主要数据:

		{
  			"data": {
    			"type": "articles",
    			"id": "1",
    			"attributes": {
      			// ... this article's attributes
    			},
    			"relationships": {
      			// ... this article's relationships
    			}
  			}
		}

如下示例是表述一个单独的资源标示对象,和上例表示的是同一个东西:

		{
  			"data": {
    		"type": "articles",
    		"id": "1"
  			}
		}
		
一个逻辑的资源集合 **必须** 表示为一个数组,甚至它仅仅包含一个元素或者它是空得.

####Resource Objects

"Resource objects" appear in a JSON API document to represent resources.

A resource object MUST contain at least the following top-level members:

* `id`
* `type`

Exception: The id member is not required when the resource object originates at the client and represents a new resource to be created on the server.

In addition, a resource object MAY contain any of these top-level members:

* `attributes`: an attributes object representing some of the resource's data.
* `relationships`: a relationships object describing relationships between the resource and other JSON API resources.
* `links`: a links object containing links related to the resource.
* `meta`: a meta object containing non-standard meta-information about a resource that can not be represented as an attribute or relationship.

Here's how an article (i.e. a resource of type "articles") might appear in a document:

		// ...
		{
  		  "type": "articles",
  		  "id": "1",
  		  "attributes": {
    		"title": "Rails is Omakase"
  		  },
  		 "relationships": {
    		"author": {
      		  "links": {
        		"self": "/articles/1/relationships/author",
        		"related": "/articles/1/author"
      		  },
      		  "data": { "type": "people", "id": "9" }
    	    }
  		  }
		}
		// ...

**Identification**

Every resource object **MUST** contain an id member and a type member. The values of the id and type members **MUST** be strings.

Within a given API, each resource object's `type` and `id` pair **MUST** identify a single, unique resource. (The set of URIs controlled by a server, or multiple servers acting as one, constitute an API.)

The `type` member is used to describe resource objects that share common attributes and relationships.

The values of type members **MUST** adhere to the same constraints as member names.

    Note: This spec is agnostic about inflection rules, so the value of type can be either plural or singular. However, the same value should be used consistently throughout an implementation.

**Fields**

A resource object's attributes and its relationships are collectively called its "fields".

Fields for a resource object **MUST** share a common namespace with each other and with type and id. In other words, a resource can not have an attribute and relationship with the same name, nor can it have an attribute or relationship named `type` or `id`.

**Attributes**

The value of the `attributes` key **MUST** be an object (an "attributes object"). Members of the attributes object ("attributes") represent information about the resource object in which it's defined.

Attributes may contain any valid JSON value.

Complex data structures involving JSON objects and arrays are allowed as attribute values. However, any object that constitutes or is contained in an attribute MUST NOT contain a relationships or links member, as those members are reserved by this specification for future use.

Although has-one foreign keys (e.g. author_id) are often stored internally alongside other information to be represented in a resource object, these keys SHOULD NOT appear as attributes.

    Note: See fields and member names for more restrictions on this container.

**Relationships**

The value of the `relationships` key **MUST** be an object (a "relationships object"). Members of the relationships object ("relationships") represent references from the resource object in which it's defined to other resource objects.

Relationships may be to-one or to-many.

A "relationship object" **MUST** contain at least one of the following:

* links: a `links object` containing at least one of the following:
	* self: a link for the relationship itself (a "relationship link"). This link allows the client to directly manipulate the relationship. For example, removing an `author` through an article's relationship URL would disconnect the person from the `article` without deleting the people resource itself. When fetched successfully, this link returns the linkage for the related resources as its primary data. (See Fetching Relationships.)
    * related: a related resource link
* data: `resource linkage`
* meta: a `meta object` that contains non-standard meta-information about the relationship.

A relationship object that represents a to-many relationship **MAY** also contain `pagination` links under the links member, as described below.

    Note: See fields and member names for more restrictions on this container.

**Related Resource Links**

A "related resource link" provides access to [resource objects] linked in a relationship. When fetched, the related resource object(s) are returned as the response's primary data.

For example, an article's comments relationship could specify a link that returns a collection of comment resource objects when retrieved through a GET request.

If present, a related resource link **MUST** reference a valid URL, even if the relationship isn't currently associated with any target resources. Additionally, a related resource link **MUST NOT** change because its relationship's content changes.

**Resource Linkage**

Resource linkage in a compound document allows a client to link together all of the included resource objects without having to GET any URLs via links.

Resource linkage **MUST** be represented as one of the following:

* `null` for empty to-one relationships.
* an empty array (`[]`) for empty to-many relationships.
* a single `resource identifier object` for non-empty to-one relationships.
* an array of `resource identifier objects` for non-empty to-many relationships.

    Note: The spec does not impart meaning to order of resource identifier objects in linkage arrays of to-many relationships, although implementations may do that. Arrays of resource identifier objects may represent ordered or unordered relationships, and both types can be mixed in one response object.

For example, the following article is associated with an `author`:

		// ...
		{
		  "type": "articles",
		  "id": "1",
		  "attributes": {
		    "title": "Rails is Omakase"
		  },
		  "relationships": {
		    "author": {
		      "links": {
		        "self": "http://example.com/articles/1/relationships/author",
		        "related": "http://example.com/articles/1/author"
		      },
		      "data": { "type": "people", "id": "9" }
		    }
		  },
		  "links": {
		    "self": "http://example.com/articles/1"
		  }
		}
		// ...

The `author` relationship includes a link for the relationship itself (which allows the client to change the related author directly), a related resource link to fetch the resource objects, and linkage information.

**Resource Links**

The optional `links` member within each resource object contains links related to the resource.

If present, this links object **MAY** contain a `self` `link` that identifies the resource represented by the resource object.

// ...
{
  "type": "articles",
  "id": "1",
  "attributes": {
    "title": "Rails is Omakase"
  },
  "links": {
    "self": "http://example.com/articles/1"
  }
}
// ...

A server **MUST** respond to a `GET` request to the specified URL with a response that includes the resource as the primary data.
		
####Resource Objects:资元对象

"资元对象"是 JSON API 文档用于展示资源.

一个资元对象 **必须** 包含至少一个下列顶级成员:

* `id`
* `type`

异常: `id` 成员不是必须的,例如当客户端组织一个资源将要发布到服务端(这个时候这个资源还没有id).

还有一个资元对象 **可以** 包含下列顶级成员:

* `attributes`: 一个`attributes object`(属性对象)标示一些资源的数据.
* `relationships`: 一个`relationships object`(关联对象)标示相关资源之间的关联.
* `links`: 一个`links object`(链接对象)包含资源相关的链接.
* `meta`:一个`meta object`(元数据对象)包含关于一个对象既不能用属性对象表述也不能用关联对象表述的非标准的元数据信息.

下例是一片文章(资源类型为文章)的表述:

		// ...
		{
		  "type": "articles",
		  "id": "1",
		  "attributes": {
		    "title": "Rails is Omakase"
		  },
		  "relationships": {
		    "author": {
		      "links": {
		        "self": "/articles/1/relationships/author",
		        "related": "/articles/1/author"
		      },
		      "data": { "type": "people", "id": "9" }
		    }
		  }
		}
		// ...

**Identification:标识**

每个资元对象 **必须** 包含一个`id`成员和一个`type`成员. `id`和`type`的值 **必须** 是字符串.

在一个已有的API, 每个资元对象的`type`和`id` **必须** 标识一个唯一的资源.

`type`成员用于表述一些`resource objects`(资元对象)中共享的一些属性和关联.

`type`成员 **必须** 和`member names`(后续说明)的规则限制一样.

**Fields:字段**

一个资元对象的`attributes`属性和它的`relationships`关联共同被叫做"字段".

字段是用于一个`resource object`(资元对象) **必须** 彼此之间的 `type` 和 `id` 共享一个公共的命名空间.另外,一个资源不能拥有相同名称的`attribute`和`relationship`.`attribute` 和 `relationship` 也不能被命名为 `type` 和 `id`.

**Attributes:属性**

`attributes`的值 **必须** 是一个对象 (`attributes object`属性对象).成员的属性对象(属性`attributes`)表示`resource object`(资元对象)的定义信息.(属性是用于定义资源的)

属性可能包含任意合法的 JSON 值.

复合数据结构的属性数组包含 JSON 对象和数组.但是,任何对象由属性组成或者包含一个属性 **必须不能** 包含`relationships`(关联) 或者 `links`(链接)成员,这些成员(`relationships`和`links`是将来用于特定情况).

尽管`has-one`一些外键(例如`author_id`)经常内部存储起来用于表述一个资元对象(一个外键`author_id`用于表述一个author对象),这些键值 **不应当** 用属性(`attributes`)表示.

**Relationships:关联**

`relationships`的值 **必须** 是一个对象 (一个`relationshops object`关联对象).关联对象的成员("relationships")表示资元对象之间相互定义的引用.

关联可能是 一对一 或 一对多.

一个资元对象 **必须** 包含至少下面一种:

* `links`: 一个`links object`(链接对象)包含至少如下一种:
	* `self`: 一种表示自己的关联链接(`relationship link`).这个链接允许客户直接操作这个关联.如,移除一个`author`(作者)通过一个`article`(文章)的关联链接将会删除用户和文章之间的关联而不是删除用户资源本身.当获取这个链接信息时,这个链接返回相关资源作为其主要数据(`primary data`).
	* `related`:一个相关资源的链接.
* `data`: 资源链接关联.
* `meta`: 一个元对象包含非标准的元信息关于链接.

当一个关联对象(`relationship object`)表示一个对多的关联 **可以** 也包含一个分页(`pagination`)的链接作为一个`links`的成员(下文会描述).

**Related Resource Links:相关资源的链接**

一个相关资源的链接(`related resource link`)提供访问一些[资元对象]的链接.当访问这些链接,相关资元对象会当做主要数据返回.

如,一个文章的评论关联可能被描述为一个用`GET`访问时返回一组评论资元对象链接.

如果表示了,一个资源关联的链接 **必须** 表示为一个合法的 URL,甚至这个关联当前还不能关联人和资源.更多的,一个关联资源的链接 **必须不能** 改变尽管其链接内容产生了变化.(比如一个文章的评论关联链接不会随着评论增多或减少而发生变化).

**Resource Linkage:资源联系**

资源联系在一个复合文档允许一个客户把所有被包含的资源汇总在一起而不用`GET`任何URLs从`links`.

资源联系 **必须** 表现为下面的一种形式:

* `null`表示空的对一关联.
* 空的数组`[]`表示空的对多关联.
* 一个单独的资源标识符表示一个非空的对一关联.
* 一组资源标识符表示非空的对多关联.

		一个规格的没有说明资源标识符的顺序在对多关联的数组内,尽管可以这么做.资源标识符的数组可能会表示排序过的活着没有排序过的关联,2种情况可能会出现在一种表述内.
		
如,下面的文章和作者的关联:

		// ...
		{
		  "type": "articles",
		  "id": "1",
		  "attributes": {
		    "title": "Rails is Omakase"
		  },
		  "relationships": {
		    "author": {
		      "links": {
		        "self": "http://example.com/articles/1/relationships/author",
		        "related": "http://example.com/articles/1/author"
		      },
		      "data": { "type": "people", "id": "9" }
		    }
		  },
		  "links": {
		    "self": "http://example.com/articles/1"
		  }
		}
		// ...
		
作者的关联包含一个关联本身的链接(允许客户改变关联作者,通过操作`http://example.com/articles/1/relationships/author`改变作者与文章之间的关联),一个关联的资源链接获取一组资元对象和关联信息.

**Resource Links:资源链接**

可选的`links`成员是表述每个资元对象包含资源关联的链接.

如果成功表述,`links object` **可以** 包含一个 `self`链接,表示资源自己本身.

		// ...
		{
		  "type": "articles",
		  "id": "1",
		  "attributes": {
		    "title": "Rails is Omakase"
		  },
		  "links": {
		    "self": "http://example.com/articles/1"
		  }
		}
		// ...	

服务端 **必须** 对已这个特定的URL响应一个 `GET` 请求展现其主要数据(primary data).

###Resource Identifier Objects

A "resource identifier object" is an object that identifies an individual resource.

A "resource identifier object" **MUST** contain `type` and `id` members.

A "resource identifier object" **MAY** also include a meta member, whose value is a meta object that contains non-standard meta-information.

###Resource Identifier Objects:资源标示对象

一个资源标示对象是一个表示单独资源的对象.

一个资源表示对象 **必须** 包含 `type` 和 `id` 成员.

一个资源标示对象 **可以** 同样引用一个值是一个包含非标准元信息的元对象作为一个元成员. 

###Compound Documents

To reduce the number of HTTP requests, servers **MAY** allow responses that include related resources along with the requested primary resources. Such responses are called "compound documents".

In a compound document, all included resources **MUST** be represented as an array of resource objects in a top-level included member.

Compound documents require "full linkage", meaning that every included resource **MUST** be identified by at least one resource identifier object in the same document. These resource identifier objects could either be primary data or represent resource linkage contained within primary or included resources.

The only exception to the full linkage requirement is when relationship fields that would otherwise contain linkage data are excluded via sparse fieldsets.

    Note: Full linkage ensures that included resources are related to either the primary data (which could be resource objects or resource identifier objects) or to each other.

A complete example document with multiple included relationships:

		{
		  "data": [{
		    "type": "articles",
		    "id": "1",
		    "attributes": {
		      "title": "JSON API paints my bikeshed!"
		    },
		    "links": {
		      "self": "http://example.com/articles/1"
		    },
		    "relationships": {
		      "author": {
		        "links": {
		          "self": "http://example.com/articles/1/relationships/author",
		          "related": "http://example.com/articles/1/author"
		        },
		        "data": { "type": "people", "id": "9" }
		      },
		      "comments": {
		        "links": {
		          "self": "http://example.com/articles/1/relationships/comments",
		          "related": "http://example.com/articles/1/comments"
		        },
		        "data": [
		          { "type": "comments", "id": "5" },
		          { "type": "comments", "id": "12" }
		        ]
		      }
		    }
		  }],
		  "included": [{
		    "type": "people",
		    "id": "9",
		    "attributes": {
		      "first-name": "Dan",
		      "last-name": "Gebhardt",
		      "twitter": "dgeb"
		    },
		    "links": {
		      "self": "http://example.com/people/9"
		    }
		  }, {
		    "type": "comments",
		    "id": "5",
		    "attributes": {
		      "body": "First!"
		    },
		    "relationships": {
		      "author": {
		        "data": { "type": "people", "id": "2" }
		      }
		    },
		    "links": {
		      "self": "http://example.com/comments/5"
		    }
		  }, {
		    "type": "comments",
		    "id": "12",
		    "attributes": {
		      "body": "I like XML better"
		    },
		    "relationships": {
		      "author": {
		        "data": { "type": "people", "id": "9" }
		      }
		    },
		    "links": {
		      "self": "http://example.com/comments/12"
		    }
		  }]
		}

A compound document **MUST NOT** include more than one resource object for each `type` and `id` pair.

    	Note: In a single document, you can think of the type and id as a composite key that uniquely references resource objects in another part of the document.

    	Note: This approach ensures that a single canonical resource object is returned with each response, even when the same resource is referenced multiple times.

###Compound Documents:复合文档

为了减少HTTP的请求,服务端 **可以** 允许在请求主要资源时包含相关资源的响应.这些响应被称为"复合文档".

在一个复合文,所有被引用的资源 **必须** 表示为一个资元对象的数组在顶级 `included`的成员内.

复合文档需要全关联,意味着每个被包含的资源 **必须** 被标示为至少一个资源标识符对象在相同的文档内.这些资源标识符对象可能是主要数据也可能为包含主要或者引用资源的资源关联.

全关联请求唯一例外的是当`relationship`关联字段可能包含关联数据被`sparse fieldsets`(下文阐述)排除在外.

		注意:全关联必须确定引用的资源是关联 primary data 主要数据(资元对象或者资源标识符对象)或相互之间.
		
一个完全的示例文档和多个包含关联如下:

		{
		  "data": [{
		    "type": "articles",
		    "id": "1",
		    "attributes": {
		      "title": "JSON API paints my bikeshed!"
		    },
		    "links": {
		      "self": "http://example.com/articles/1"
		    },
		    "relationships": {
		      "author": {
		        "links": {
		          "self": "http://example.com/articles/1/relationships/author",
		          "related": "http://example.com/articles/1/author"
		        },
		        "data": { "type": "people", "id": "9" }
		      },
		      "comments": {
		        "links": {
		          "self": "http://example.com/articles/1/relationships/comments",
		          "related": "http://example.com/articles/1/comments"
		        },
		        "data": [
		          { "type": "comments", "id": "5" },
		          { "type": "comments", "id": "12" }
		        ]
		      }
		    }
		  }],
		  "included": [{
		    "type": "people",
		    "id": "9",
		    "attributes": {
		      "first-name": "Dan",
		      "last-name": "Gebhardt",
		      "twitter": "dgeb"
		    },
		    "links": {
		      "self": "http://example.com/people/9"
		    }
		  }, {
		    "type": "comments",
		    "id": "5",
		    "attributes": {
		      "body": "First!"
		    },
		    "relationships": {
		      "author": {
		        "data": { "type": "people", "id": "2" }
		      }
		    },
		    "links": {
		      "self": "http://example.com/comments/5"
		    }
		  }, {
		    "type": "comments",
		    "id": "12",
		    "attributes": {
		      "body": "I like XML better"
		    },
		    "relationships": {
		      "author": {
		        "data": { "type": "people", "id": "9" }
		      }
		    },
		    "links": {
		      "self": "http://example.com/comments/12"
		    }
		  }]
		}

一个复合文档 **必须不能** 包含多于一个资元对象对于一个 `type`和`id`的对(`type`和`id`是复合唯一属性针对一个资元对象).

		注意:一个单一文档,你可以认为 type 和 id 是一个复合键代表一个唯一的资元对象对于其他文档部分.
		
		注意:这个方法是确保一个规范的资源在每个响应中返回,深圳相同的资源被多次关联引用(也就是即使被关联多次,在included内也只呈现一次).
		
###Meta Information

Where specified, a meta member can be used to include non-standard meta-information. The value of each meta member **MUST** be an object (a "meta object").

Any members **MAY** be specified within meta objects.
		
For example:

		{
		  "meta": {
		    "copyright": "Copyright 2015 Example Corp.",
		    "authors": [
		      "Yehuda Katz",
		      "Steve Klabnik",
		      "Dan Gebhardt",
		      "Tyler Kellen"
		    ]
		  },
		  "data": {
		    // ...
		  }
		}
		
###Meta Information:元信息

一个特定的元数据成员可以被用于包含非标准化的元信息.每个元数据的值 **必须** 是一个对象(元对象 `meta object`).

任何成员 **可以** 被指定到元对象内.

示例:

		{
		  "meta": {
		    "copyright": "Copyright 2015 Example Corp.",
		    "authors": [
		      "Yehuda Katz",
		      "Steve Klabnik",
		      "Dan Gebhardt",
		      "Tyler Kellen"
		    ]
		  },
		  "data": {
		    // ...
		  }
		}
		
###Links

Where specified, a links member can be used to represent links. The value of each links member **MUST** be an object (a "links object").

Each member of a links object is a "link". A link **MUST** be represented as either:

* a string containing the link's URL.
* an object ("link object") which can contain the following members:
	* href: a string containing the link's URL.
	* meta: a meta object containing non-standard meta-information about the link.

The following self link is simply a URL:

		"links": {
		  "self": "http://example.com/posts",
		}

The following related link includes a URL as well as meta-information about a related resource collection:

		"links": {
		  "related": {
		    "href": "http://example.com/articles/1/comments",
		    "meta": {
		      "count": 10
		    }
		  }
		}

		Note: Additional members may be specified for links objects and link objects in the future. It is also possible that the allowed values of additional members will be expanded (e.g. a collection link may support an array of values, whereas a self link does not).
		
###Links:链接

一个`links`链接成员可以用于表述链接.每个`links`链接成员的值 **必须** 是一个对象(links object).

每一个连接对象的成员是一个链接`link`. 一个链接 **必须** 表述如下之一:

* 一个字符串包含链接的URL.
* 一个包含如下信息的对象(`link object`链接对象):
	* `href`: 一个字符串包含链接的URL.
	* `meta`: 一个和链接相关,包含非标准元信息的元对象. 
	
如下`self link`是一个简单的URL示例:

		"links": {
		  "self": "http://example.com/posts",
		}
		
如下关联链接(related link)除了包含一个URL也包含一个关于资源集合的元信息:

		"links": {
		  "related": {
		    "href": "http://example.com/articles/1/comments",
		    "meta": {
		      "count": 10
		    }
		  }
		}
		
count代表集合总数.
		
###JSON API Object

A JSON API document **MAY** include information about its implementation under a top level jsonapi member. If present, the value of the jsonapi member **MUST** be an object (a "jsonapi object"). The jsonapi object **MAY** contain a version member whose value is a string indicating the highest JSON API version supported. This object **MAY** also contain a meta member, whose value is a meta object that contains non-standard meta-information.

		{
		  "jsonapi": {
		    "version": "1.0"
		  }
		}

If the version member is not present, clients should assume the server implements at least version 1.0 of the specification.

    	Note: Because JSON API is committed to making additive changes only, the version string primarily indicates which new features a server may support.

###JSON API Object:Json API 对象

一个JSON API文档 **可以** 包含关于自己设计实现的信息在顶级(`top level`)下的 `jsonapi` 成员.如果用于表述,`jsonapi`成员的值 **必须** 是一个对象(jsonapi 对象). `json api`对象 **可以** 包含一个值为字符串的 `version`版本成员.该版本用于表示最高级的 JSON API 版本支持.这个对象 **可以** 也包含一个包含非标准元信息的值为源对象元成员.

		{
		  "jsonapi": {
		    "version": "1.0"
		  }
		}
		
如果这个版本成员没有表述,客户端需要假设服务端至少实现1.0的版本.

		注意: 因为JSON API仅仅是用于提交额外的变更,版本号主要用于表示哪些新的功能服务端可能会支持.

###Member Names

All member names used in a JSON API document **MUST** be treated as case sensitive by clients and servers, and they **MUST** meet all of the following conditions:

* Member names **MUST** contain at least one character.
* Member names **MUST** contain only the allowed characters listed below.
* Member names **MUST** start and end with a "globally allowed character", as defined below.

To enable an easy mapping of member names to URLs, it is **RECOMMENDED** that member names use only non-reserved, URL safe characters specified in RFC 3986.

**Allowed Characters**

The following "globally allowed characters" MAY be used anywhere in a member name:

* U+0061 to U+007A, "a-z"
* U+0041 to U+005A, "A-Z"
* U+0030 to U+0039, "0-9"
* U+0080 and above (non-ASCII Unicode characters; not recommended, not URL safe)

Additionally, the following characters are allowed in member names, except as the first or last character:

* U+002D HYPHEN-MINUS, "-"
* U+005F LOW LINE, "_"
* U+0020 SPACE, " " (not recommended, not URL safe)

**Reserved Characters**

The following characters MUST NOT be used in member names:

* U+002B PLUS SIGN, "+" (used for ordering)
* U+002C COMMA, "," (used as a separator between relationship paths)
* U+002E PERIOD, "." (used as a separator within relationship paths)
* U+005B LEFT SQUARE BRACKET, "[" (used in sparse fieldsets)
* U+005D RIGHT SQUARE BRACKET, "]" (used in sparse fieldsets)
* U+0021 EXCLAMATION MARK, "!"
* U+0022 QUOTATION MARK, '"'
* U+0023 NUMBER SIGN, "#"
* U+0024 DOLLAR SIGN, "$"
* U+0025 PERCENT SIGN, "%"
* U+0026 AMPERSAND, "&"
* U+0027 APOSTROPHE, "'"
* U+0028 LEFT PARENTHESIS, "("
* U+0029 RIGHT PARENTHESIS, ")"
* U+002A ASTERISK, "*"
* U+002F SOLIDUS, "/"
* U+003A COLON, ":"
* U+003B SEMICOLON, ";"
* U+003C LESS-THAN SIGN, "<"
* U+003D EQUALS SIGN, "="
* U+003E GREATER-THAN SIGN, ">"
* U+003F QUESTION MARK, "?"
* U+0040 COMMERCIAL AT, "@"
* U+005C REVERSE SOLIDUS, "\"
* U+005E CIRCUMFLEX ACCENT, "^"
* U+0060 GRAVE ACCENT, "`"
* U+007B LEFT CURLY BRACKET, "{"
* U+007C VERTICAL LINE, "|"
* U+007D RIGHT CURLY BRACKET, "}"
* U+007E TILDE, "~"
* U+007F DELETE
* IU+0000 to U+001F (C0 Controls)

###Member Names:成员名称

所有成员名称用于 JSON API的文档 **必须** 是大小写敏感无论是客户端还是服务端. 并且 **必须** 符合下列情况:

* 成员名称 **必须** 包含至少一个字符.
* 成员名称 **必须** 包含下列被允许的字符.
* 成员名称 **必须** 用全局可用字符(`globally allowed characters`)开始和结尾.

为了更简单的映射成员名称到URL, **建议** 所有成员名称使用未保留的, URL 安全字符. 具体可以参见 `RFC 3986`.

**Allowed Characters:允许的字符**

下列全局可用字符(`globally allowed characters`) **可以** 被用于任何成员名称:

* U+0061 to U+007A, "a-z"
* U+0041 to U+005A, "A-Z"
* U+0030 to U+0039, "0-9"
* U+0080 and above (non-ASCII Unicode characters; 不建议, not URL safe)

额外的,如下字符被允许用于成员名称, 但是起始和结尾字母除外:

* U+002D HYPHEN-MINUS, "-"
* U+005F LOW LINE, "_"
* U+0020 SPACE, " " (not recommended, not URL safe)

**Reserved Characters:保留的字符**

如下字符 **必须不能** 被用于成员名称:

* U+002B PLUS SIGN, "`+`" (用于排序)
* U+002C COMMA, "`,`" (used as a separator between relationship paths)
* U+002E PERIOD, "`.`" (used as a separator within relationship paths)
* U+005B LEFT SQUARE BRACKET, "`[`" (used in sparse fieldsets)
* U+005D RIGHT SQUARE BRACKET, "`]`" (used in sparse fieldsets)
* U+0021 EXCLAMATION MARK, "`!`"
* U+0022 QUOTATION MARK, '`"`'
* U+0023 NUMBER SIGN, "`#`"
* U+0024 DOLLAR SIGN, "`$`"
* U+0025 PERCENT SIGN, "`%`"
* U+0026 AMPERSAND, "`&`"
* U+0027 APOSTROPHE, "`'`"
* U+0028 LEFT PARENTHESIS, "`(`"
* U+0029 RIGHT PARENTHESIS, "`)`"
* U+002A ASTERISK, "`*`"
* U+002F SOLIDUS, "`/`"
* U+003A COLON, "`:`"
* U+003B SEMICOLON, "`;`"
* U+003C LESS-THAN SIGN, "`<`"
* U+003D EQUALS SIGN, "`=`"
* U+003E GREATER-THAN SIGN, "`>`"
* U+003F QUESTION MARK, "`?`"
* U+0040 COMMERCIAL AT, "`@`"
* U+005C REVERSE SOLIDUS, "`\`"
* U+005E CIRCUMFLEX ACCENT, "`^`"
* U+0060 GRAVE ACCENT, "`"
* U+007B LEFT CURLY BRACKET, "`{`"
* U+007C VERTICAL LINE, "`|`"
* U+007D RIGHT CURLY BRACKET, "`}`"
* U+007E TILDE, "`~`"
* U+007F DELETE
* IU+0000 to U+001F (C0 Controls)

###Fetching Data

Data, including resources and relationships, can be fetched by sending a GET request to an endpoint.

Responses can be further refined with the optional features described below.

###Fetching Data:获取数据

数据,包含资源和关联可以被一个 `GET` 请求向接口终端获取.

响应可以被下列一些可选的`features`(特色,功能)重构.

####Fetching Resources

A server **MUST** support fetching resource data for every URL provided as:

* a self link as part of the top-level links object
* a self link as part of a resource-level links object
* a related link as part of a relationship-level links object

For example, the following request fetches a collection of articles:

		GET /articles HTTP/1.1
		Accept: application/vnd.api+json

The following request fetches an article:

		GET /articles/1 HTTP/1.1
		Accept: application/vnd.api+json

And the following request fetches an article's author:

		GET /articles/1/author HTTP/1.1
		Accept: application/vnd.api+json

####Fetching Resources:获取资源

服务端 **必须** 支持获取资源数据对下列URL:

* 一个 `self link(自我链接)` 是一组顶级链接成员的一部分.
* 一个 `self link(自我链接)` 是一组资源级别链接成员的一部分.  
* 一个 `related link(相关链接)` 是一组关联级别链接成员的一部分.

示例如下,下列请求获取一组文章:

		GET /articles HTTP/1.1
		Accept: application/vnd.api+json
		
下列请求获取一篇文章:

		GET /articles/1 HTTP/1.1
		Accept: application/vnd.api+json
		
下列请求获取一篇文章的作者:

		GET /articles/1/author HTTP/1.1
		Accept: application/vnd.api+json
		
#####Responses

**200 OK**

A server **MUST** respond to a successful request to fetch an individual resource or resource collection with a 200 OK response.

A server **MUST** respond to a successful request to fetch a resource collection with an array of resource objects or an empty array (`[]`) as the response document's primary data.

For example, a `GET` request to a collection of articles could return:

		HTTP/1.1 200 OK
		Content-Type: application/vnd.api+json
		
		{
		  "links": {
		    "self": "http://example.com/articles"
		  },
		  "data": [{
		    "type": "articles",
		    "id": "1",
		    "attributes": {
		      "title": "JSON API paints my bikeshed!"
		    }
		  }, {
		    "type": "articles",
		    "id": "2",
		    "attributes": {
		      "title": "Rails is Omakase"
		    }
		  }]
		}

A similar response representing an empty collection would be:

		HTTP/1.1 200 OK
		Content-Type: application/vnd.api+json
		
		{
		  "links": {
		    "self": "http://example.com/articles"
		  },
		  "data": []
		}

A server **MUST** respond to a successful request to fetch an individual resource with a resource object or `null` provided as the response document's primary data.

`null` is only an appropriate response when the requested URL is one that might correspond to a single resource, but doesn't currently.

		Note: Consider, for example, a request to fetch a to-one related resource link. This request would respond with null when the relationship is empty (such that the link is corresponding to no resources) but with the single related resource's resource object otherwise.

For example, a `GET` request to an individual article could return:

		HTTP/1.1 200 OK
		Content-Type: application/vnd.api+json
		
		{
		  "links": {
		    "self": "http://example.com/articles/1"
		  },
		  "data": {
		    "type": "articles",
		    "id": "1",
		    "attributes": {
		      "title": "JSON API paints my bikeshed!"
		    },
		    "relationships": {
		      "author": {
		        "links": {
		          "related": "http://example.com/articles/1/author"
		        }
		      }
		    }
		  }
		}

If the above article's author is missing, then a GET request to that related resource would return:

		HTTP/1.1 200 OK
		Content-Type: application/vnd.api+json
		
		{
		  "links": {
		    "self": "http://example.com/articles/1/author"
		  },
		  "data": null
		}

**404 Not Found**

A server **MUST** respond with 404 Not Found when processing a request to fetch a single resource that does not exist, except when the request warrants a 200 OK response with null as the primary data (as described above).
Other Responses

A server **MAY** respond with other HTTP status codes.

A server **MAY** include error details with error responses.

A server **MUST** prepare responses, and a client **MUST** interpret responses, in accordance with HTTP semantics.

#####响应

**200 ok**

服务端 **必须** 对于一个成功获取单独资源或者一组资源响应 `200 ok`.(响应200状态码)

服务端 **必须** 对于一个成功获取一组资源集合返回一个资源对象数组或者空的数组`[]`在响应文档的主要数据中(`primary data`).(返回空或者一组数据)

如下,一个`GET`请求一组文章可能返回:

		HTTP/1.1 200 OK
		Content-Type: application/vnd.api+json
		
		{
		  "links": {
		    "self": "http://example.com/articles"
		  },
		  "data": [{
		    "type": "articles",
		    "id": "1",
		    "attributes": {
		      "title": "JSON API paints my bikeshed!"
		    }
		  }, {
		    "type": "articles",
		    "id": "2",
		    "attributes": {
		      "title": "Rails is Omakase"
		    }
		  }]
		}
		
一个相似的响应如果返回一个空的集合则会表示如下:

		HTTP/1.1 200 OK
		Content-Type: application/vnd.api+json
		
		{
		  "links": {
		    "self": "http://example.com/articles"
		  },
		  "data": []
		}

一个成功的请求在获取一个单独的资源服务端 **必须** 返回一个资源对象或者 `null`在响应文档的主要数据中.

`null`仅仅是恰当的响应当请求URL获取一个现在不可用的单独资源.

		注意:当一个请求一个关联资源链接.请求可能响应空,当关联链接是空的(不关联到任何资源)但是一个单独的关联资源的资源对象除外.
		
如下示例,一个`GET`请求获取一个单独的文章可能返回:

		HTTP/1.1 200 OK
		Content-Type: application/vnd.api+json
		
		{
		  "links": {
		    "self": "http://example.com/articles/1"
		  },
		  "data": {
		    "type": "articles",
		    "id": "1",
		    "attributes": {
		      "title": "JSON API paints my bikeshed!"
		    },
		    "relationships": {
		      "author": {
		        "links": {
		          "related": "http://example.com/articles/1/author"
		        }
		      }
		    }
		  }
		}
		
如果上面文章的作者不存在,在去 `GET`请求相关资源返回:

		HTTP/1.1 200 OK
		Content-Type: application/vnd.api+json
		
		{
		  "links": {
		    "self": "http://example.com/articles/1/author"
		  },
		  "data": null
		}

**404 Not Found**

服务端 **必须** 响应`404 Not Found`当处理一个获取单独的请求但是不存在,除非请求允许返回一个 `200 ok`响应和一个 `null` 在其主要数据中(上述描述).

**Other Responses:其他响应**

服务端 **可以** 响应返回其他 HTTP 状态码.

服务端 **可以** 响应错误时会包含`error details`错误细节信息.

服务端 **必须** 准备准备着响应,并且客户端 **必须** 解释响应根据 HTTP 语义.

####Fetching Relationships

A server **MUST** support fetching relationship data for every relationship URL provided as a self link as part of a relationship's links object.

For example, the following request fetches data about an article's comments:

		GET /articles/1/relationships/comments HTTP/1.1
		Accept: application/vnd.api+json

And the following request fetches data about an article's author:

		GET /articles/1/relationships/author HTTP/1.1
		Accept: application/vnd.api+json

####Fetching Relationships:获取关联

服务端 **必须** 支持获取关联数据对于每个关联URL以`self link`格式或者一个关联的链接对象的一部分呈现.

示例,下列请求获取关于文章的评论数据:

		GET /articles/1/relationships/comments HTTP/1.1
		Accept: application/vnd.api+json
		
并且下列请求获取关于文章的组着数据:

		GET /articles/1/relationships/author HTTP/1.1
		Accept: application/vnd.api+json				
#####Responses

**200 OK**

A server **MUST** respond to a successful request to fetch a relationship with a 200 OK response.

The primary data in the response document **MUST** match the appropriate value for resource linkage, as described above for relationship objects.

The top-level links object **MAY** contain self and related links, as described above for relationship objects.

For example, a GET request to a URL from a to-one relationship link could return:

		HTTP/1.1 200 OK
		Content-Type: application/vnd.api+json
		
		{
		  "links": {
		    "self": "/articles/1/relationships/author",
		    "related": "/articles/1/author"
		  },
		  "data": {
		    "type": "people",
		    "id": "12"
		  }
		}

If the above relationship is empty, then a GET request to the same URL would return:

		HTTP/1.1 200 OK
		Content-Type: application/vnd.api+json
		
		{
		  "links": {
		    "self": "/articles/1/relationships/author",
		    "related": "/articles/1/author"
		  },
		  "data": null
		}

A `GET` request to a URL from a to-many relationship link could return:

		HTTP/1.1 200 OK
		Content-Type: application/vnd.api+json
		
		{
		  "links": {
		    "self": "/articles/1/relationships/tags",
		    "related": "/articles/1/tags"
		  },
		  "data": [
		    { "type": "tags", "id": "2" },
		    { "type": "tags", "id": "3" }
		  ]
		}

If the above relationship is empty, then a `GET` request to the same URL would return:

		HTTP/1.1 200 OK
		Content-Type: application/vnd.api+json
		
		{
		  "links": {
		    "self": "/articles/1/relationships/tags",
		    "related": "/articles/1/tags"
		  },
		  "data": []
		}

**404 Not Found**

A server **MUST** return `404 Not Found` when processing a request to fetch a relationship link URL that does not exist.

    Note: This can happen when the parent resource of the relationship does not exist. For example, when /articles/1 does not exist, request to /articles/1/relationships/tags returns 404 Not Found.

If a relationship link URL exists but the relationship is empty, then `200 OK` **MUST** be returned, as described above.
Other Responses

A server **MAY** respond with other HTTP status codes.

A server **MAY** include error details with error responses.

A server **MUST** prepare responses, and a client **MUST** interpret responses, in accordance with HTTP semantics.

#####Responses:响应

**200 OK**

服务端 **必须** 对于一个成功获取单独资源或者一组资源响应 `200 ok`.(响应200状态码)

在响应文档的主要数据 **必须** 匹配合适的值对于 `resource linkage`(资源关联),就像上面描述的一组关联对象.

顶级链接对象(`top-level links object`) **可以** 包含 `self` 和 `related` 链接,就上上面描述的关联对象.

示例,一个 `GET` 请求URL从对一关联链接可能返回:

		HTTP/1.1 200 OK
		Content-Type: application/vnd.api+json
		
		{
		  "links": {
		    "self": "/articles/1/relationships/author",
		    "related": "/articles/1/author"
		  },
		  "data": {
		    "type": "people",
		    "id": "12"
		  }
		}
		
如果上面的链接是空的,则 `GET` 请求相同的URL可能返回:

		HTTP/1.1 200 OK
		Content-Type: application/vnd.api+json
		
		{
		  "links": {
		    "self": "/articles/1/relationships/author",
		    "related": "/articles/1/author"
		  },
		  "data": null
		}
		
一个`GET`请求一个URL从对多关联链接可能返回:

		HTTP/1.1 200 OK
		Content-Type: application/vnd.api+json
		
		{
		  "links": {
		    "self": "/articles/1/relationships/tags",
		    "related": "/articles/1/tags"
		  },
		  "data": [
		    { "type": "tags", "id": "2" },
		    { "type": "tags", "id": "3" }
		  ]
		}		

如果上面的链接是空的,则 `GET`请求相同的 URL 可能返回:

		HTTP/1.1 200 OK
		Content-Type: application/vnd.api+json
		
		{
		  "links": {
		    "self": "/articles/1/relationships/tags",
		    "related": "/articles/1/tags"
		  },
		  "data": []
		}

**404 Not Found**

服务端 **必须** 返回 `404 Not Found` 当处理一个请求获取一个关联链接但是不存在.

		注意:这种情况可能发生在父资源关联不存在的时候.比如 当 /articles/1 不存在,请求 /articles/1/relationships/tags 会返回 404 Not Found.
		
如果一个关联链接 URL 存在 但是关联式空的, 则 `200 ok` 必须按照上述情况返回.

**Other Responses:其他相应**

服务端 **可以** 相应返回其他 HTTP 状态码.

服务端 **可以** 响应错误时会包含`error details`错误细节信息.

服务端 **必须** 准备准备着响应,并且客户端 **必须** 解释响应根据 HTTP 语义.

####Inclusion of Related Resources

An endpoint **MAY** return resources related to the primary data by default.

An endpoint **MAY** also support an include request parameter to allow the client to customize which related resources should be returned.

If an endpoint does not support the include parameter, it **MUST** respond with 400 Bad Request to any requests that include it.

If an endpoint supports the include parameter and a client supplies it, the server **MUST NOT** include unrequested resource objects in the included section of the compound document.

The value of the include parameter MUST be a comma-separated (U+002C COMMA, ",") list of relationship paths. A relationship path is a dot-separated (U+002E FULL-STOP, ".") list of relationship names.

If a server is unable to identify a relationship path or does not support inclusion of resources from a path, it **MUST** respond with 400 Bad Request.

    	Note: For example, a relationship path could be comments.author, where comments is a relationship listed under a articles resource object, and author is a relationship listed under a comments resource object.

For instance, comments could be requested with an article:

		GET /articles/1?include=comments HTTP/1.1
		Accept: application/vnd.api+json

In order to request resources related to other resources, a dot-separated path for each relationship name can be specified:

		GET /articles/1?include=comments.author HTTP/1.1
		Accept: application/vnd.api+json

    	Note: Because compound documents require full linkage (except when relationship linkage is excluded by sparse fieldsets), intermediate resources in a multi-part path must be returned along with the leaf nodes. For example, a response to a request for comments.author should include comments as well as the author of each of those comments.

    	Note: A server may choose to expose a deeply nested relationship such as comments.author as a direct relationship with an alias such as comment-authors. This would allow a client to request /articles/1?include=comment-authors instead of /articles/1?include=comments.author. By abstracting the nested relationship with an alias, the server can still provide full linkage in compound documents without including potentially unwanted intermediate resources.

Multiple related resources can be requested in a comma-separated list:

		GET /articles/1?include=author,comments.author HTTP/1.1
		Accept: application/vnd.api+json

Furthermore, related resources can be requested from a relationship endpoint:

		GET /articles/1/relationships/comments?include=comments.author HTTP/1.1
		Accept: application/vnd.api+json

In this case, the primary data would be a collection of resource identifier objects that represent linkage to comments for an article, while the full comments and comment authors would be returned as included data.

    	Note: This section applies to any endpoint that responds with primary data, regardless of the request type. For instance, a server could support the inclusion of related resources along with a POST request to create a resource or relationship.

####Inclusion of Related Resources:相关资源的包含

一个接口终端 **可以** 返回主要数据相关的资源默认.

一个接口终端 **可以** 也支持一个 `include` 请求参数可以让客户自定义哪些相关资源返回.

如果一个接口终端不支持 `include` 参数,它 **必须** 相应 `400 Bad Request`对任何请求包含 `include`.

如果一个接口终端支持 `include` 参数并且客户端也支持它, 这样服务端 **必须不能** 包含非请求资源对象在 复合文档的`include` 区域.

`include` 的值 **必须** 是用逗号(U+002C COMMA, ",")分隔一组关联路径.一个关联路径是用点(U+002E FULL-STOP, ".")分隔列表的名称.

如果一个服务端不能标示一个关联路径的表示或者不支持从一个路径(`xxx.xxx`)包含的资源,它 **必须** 相应 `400 Bad Request`.

		注意:示例,一个关联路径可能是 comments.author, comments 是一个关联罗列在一组文章资源对象下面, author 是一个关联罗列在一组评论资源对象下面.
		
示例,评论可以被请求文章时一起返回:

		GET /articles/1?include=comments HTTP/1.1
		Accept: application/vnd.api+json		
		
为了请求资源和其相关资源,用点作为路径分隔每个关联名称,如下所示:

		GET /articles/1?include=comments.author HTTP/1.1
		Accept: application/vnd.api+json

		注意:因为复合文档需求一个全的链接(除非一个关联链接被用 sparse fieldsets 排除在外),中间资源在一个多部分路径必须和叶节点一起返回.示例,一个相应一个 comments.author 的请求应该包含每个 comments 中的 author.
		
		注意:服务端可能选择暴漏一个深度嵌套关联比如 comments.author 表示为 comment-authors.这个将会允许客户请求 /articles/1?include=comment-authors 代替 /articles/1?include=comments.author.
		
多个关联资源可能被用逗号分隔请求:

		GET /articles/1?include=author,comments.author HTTP/1.1
		Accept: application/vnd.api+json
		
更多的是,关联资源可以被从一个关联接口终端请求:

		GET /articles/1/relationships/comments?include=comments.author HTTP/1.1
		Accept: application/vnd.api+json

在这种情况,主要数据将会被以一个文章的评论关联用资源标识符对象(只有`type`和`id`)展示,同时全部评论和全部评论的作者将会被以`included data`返回.

		注意:这个include可以被应用到任何请求(任何请求类型)主要数据的接口终端.举例,服务端可能支持包含相关资源随着使用 POST 创建一个资源或者关联.

####Sparse Fieldsets

A client **MAY** request that an endpoint return only specific fields in the response on a per-type basis by including a fields[TYPE] parameter.

The value of the fields parameter **MUST** be a comma-separated (U+002C COMMA, ",") list that refers to the name(s) of the fields to be returned.

If a client requests a restricted set of fields for a given resource type, an endpoint **MUST NOT** include additional fields in resource objects of that type in its response.

		GET /articles?include=author&fields[articles]=title,body&fields[people]=name HTTP/1.1
		Accept: application/vnd.api+json

    	Note: The above example URI shows unencoded [ and ] characters simply for readability. In practice, these characters must be percent-encoded, per the requirements in RFC 3986.

    	Note: This section applies to any endpoint that responds with resources as primary or included data, regardless of the request type. For instance, a server could support sparse fieldsets along with a POST request to create a resource.

####Sparse Fieldsets:单独的字段

客户端 **可以** 在使用`fields[TYPE]`请求一个接口终端仅仅返回指定的字段.

一个字段的值 **必须** 是用逗号分隔的(U+002C COMMA, ",")要返回的字段名称.

如果一个客户端对于一个资源类型请求一个有限的字段(使用fields[TYPE]),接口终端 **必须不能** 响应返回包含附加的字段在资源对象中.

		GET /articles?include=author&fields[articles]=title,body&fields[people]=name HTTP/1.1
		Accept: application/vnd.api+json
		
		注意: 上述列子的 URI 显示为编码化的 [ 和 ] 字符只是用于方便阅读.在实际中,这些字符必须编码化根据 RFC 3986 的需求.
		
		注意: 这个章节的功能(Sparse Fieldsets)可以用于任意请求类型的接口终端响应主要数据(primary)部分或者included data(包含数据部分).举例,服务端可以支持 sparse fieldsets 随着 POST 请求一个创建资源.		
####Sorting

A server **MAY** choose to support requests to sort resource collections according to one or more criteria ("sort fields").

    	Note: Although recommended, sort fields do not necessarily need to correspond to resource attribute and association names.

    	Note: It is recommended that dot-separated (U+002E FULL-STOP, ".") sort fields be used to request sorting based upon relationship attributes. For example, a sort field of author.name could be used to request that the primary data be sorted based upon the name attribute of the author relationship.

An endpoint **MAY** support requests to sort the primary data with a sort query parameter. The value for sort **MUST** represent sort fields.

		GET /people?sort=age HTTP/1.1
		Accept: application/vnd.api+json
		
An endpoint **MAY** support multiple sort fields by allowing comma-separated (U+002C COMMA, ",") sort fields. Sort fields **SHOULD** be applied in the order specified.

		GET /people?sort=age,name HTTP/1.1
		Accept: application/vnd.api+json

The sort order for each sort field MUST be ascending unless it is prefixed with a minus (U+002D HYPHEN-MINUS, "-"), in which case it **MUST** be descending.

		GET /articles?sort=-created,title HTTP/1.1
		Accept: application/vnd.api+json

The above example should return the newest articles first. Any articles created on the same date will then be sorted by their title in ascending alphabetical order.

If the server does not support sorting as specified in the query parameter sort, it **MUST** return 400 Bad Request.

If sorting is supported by the server and requested by the client via query parameter sort, the server **MUST** return elements of the top-level data array of the response ordered according to the criteria specified. The server **MAY** apply default sorting rules to top-level data if request parameter sort is not specified.

    	Note: This section applies to any endpoint that responds with a resource collection as primary data, regardless of the request type.

####Sorting:排序

服务端 **可以** 选择支持请求排序主要数据和 `sort` 查询参数. `sort`的值 **必须** 表示为 `sort`字段.

		注意:尽管sort被推荐,但是 sort 字段并不是必须和管理资源的属性名称对应.(这里的意思是比如排序文章的发布时间,发布时间的属性名称为 createTime.则 sort 也最好命名为 createTime).
		
		注意:推荐用点分隔有关联的字段属性名称.比如,排序 author.name 可以被用于请求主要数据按照作者的名称排序.

一个接口终端 **可以** 支持请求排序主要数据和 `sort` 查询参数. `sort`的值 **必须** 表示为排序字段.

		GET /people?sort=age HTTP/1.1
		Accept: application/vnd.api+json
		
一个接口终端 **可以** 支持多种排序字段允许使用逗号(U+002C COMMA, ",")分隔.排序字段 **应该** 被应用指定的排序.

		GET /people?sort=age,name HTTP/1.1
		Accept: application/vnd.api+json
		
对已每一个字段排序 **必须** 是 ascending(升序) 除非添加前缀 `-`减号(U+002D HYPHEN-MINUS, "-")来指定降序.

		GET /articles?sort=-created,title HTTP/1.1
		Accept: application/vnd.api+json

上个列子需要返回最新的文章排在第一个.如果文章的创建时间相同则按照字母的升序排列(这里是针对英文的情况).

如果服务端不支持按照指定 `sort` 参数排序, 这接口 **必须** 返回 `400 Bad Request`.

如果服务端支持排序且客户端传递 `query` 参数, 服务端 **必须** 按照指定的排序返回数据在顶级 `data` 字段中. 服务端 **可以** 应用默认的排序规则在顶级 `data` 字段中. 如果请求参数 `sort` 没有指定.

		注意: 这个章节的功能(sorting)可以应用在响应一组主要数据(primary data)的任何请求类型的接口终端中.

####Pagination

A server **MAY** choose to limit the number of resources returned in a response to a subset ("page") of the whole set available.

A server **MAY** provide links to traverse a paginated data set ("pagination links").

Pagination links **MUST** appear in the links object that corresponds to a collection. To paginate the primary data, supply pagination links in the top-level `links` object. To paginate an included collection returned in a compound document, supply pagination links in the corresponding links object.

The following keys MUST be used for pagination links:

* `first`: the first page of data
* `last`: the last page of data
* `prev`: the previous page of data
* `next`: the next page of data

Keys **MUST** either be omitted or have a `null` value to indicate that a particular link is unavailable.

Concepts of order, as expressed in the naming of pagination links, **MUST** remain consistent with JSON API's sorting rules.

The page query parameter is reserved for pagination. Servers and clients **SHOULD** use this key for pagination operations.

    Note: JSON API is agnostic about the pagination strategy used by a server. Effective pagination strategies include (but are not limited to): page-based, offset-based, and cursor-based. The page query parameter can be used as a basis for any of these strategies. For example, a page-based strategy might use query parameters such as page[number] and page[size], an offset-based strategy might use page[offset] and page[limit], while a cursor-based strategy might use page[cursor].

    Note: The example query parameters above use unencoded [ and ] characters simply for readability. In practice, these characters must be percent-encoded, per the requirements in RFC 3986.

    Note: This section applies to any endpoint that responds with a resource collection as primary data, regardless of the request type.

####Pagination:分页

服务端 **可以** 选择任意数量(这个数量存放在返回的`page`子集中)的资源的在一个响应中返回.

服务端 **可以** 提供链接遍历分页数据集合("`pagination links`分页链接")

分页链接 **必须** 出现在链接对象中用于映射一个集合. 分页主要数据,提供分页链接在顶级 `links` 对象.分页包含的集合返回一个复合文档,在链接对象中提供分页的链接.

下列键 **必须** 用于分页链接:

* `first`: 第一页的数据
* `last`: 最后一页的数据
* `prev`: 上一页数据
* `next`: 下一页数据

键 **必须** 要么省略要么是一个 `null` 值来映射一个不可用的链接.

考虑到排序, **必须** 用分页的链接和 JSON API 的排序规则组合.

`page` 参数是用于预留的分页. 服务端和客户端 **应该** 使用这个键操作分页.

		注意:JSON API 是不知道服务端使用的分页策略.有效的分页策略限于不仅限于: 根据页数,根据偏移量 和 根据当前坐标. page 参数可以被用于任意策略.举例: 一个根据页数的策略可能使用 page[offset] 和 page[limit],然而一个根据坐标的策略可能使用 page[cursor].
		I
		注意:上面的示例使用没有编码过的 [ 和 ] 字符只是用于方便阅读.实际情况,这些字符必须编码根据 RFC 3986.
		
		注意: 这个章节的功能(pagination)可以应用在响应一组主要数据(primary data)的任何请求类型的接口终端中.

####Filtering

The `filter` query parameter is reserved for filtering data. Servers and clients SHOULD use this key for filtering operations.

    Note: JSON API is agnostic about the strategies supported by a server. The filter query parameter can be used as the basis for any number of filtering strategies.

####Filtering:过滤

`filter`过滤参数是用于过滤数据. 服务端和客户端 **应该** 使用该键进行过滤操作.

		注意: JSON API 是不清楚服务端的的过滤策略的. filter 参数可以用被用于任意种过滤策略的.

###Creating, Updating and Deleting Resources

A server **MAY** allow resources of a given type to be created. It MAY also allow existing resources to be modified or deleted.

A request **MUST** completely succeed or fail (in a single "transaction"). No partial updates are allowed.

    	Note: The type member is required in every resource object throughout requests and responses in JSON API. There are some cases, such as when POSTing to an endpoint representing heterogenous data, when the type could not be inferred from the endpoint. However, picking and choosing when it is required would be confusing; it would be hard to remember when it was required and when it was not. Therefore, to improve consistency and minimize confusion, type is always required.
    	
###Creating, Updating and Deleting Resources:创建,更新 和 删除资源    

服务端 **可以** 允许创建一种类型的资源. 服务端 **可以** 更新和删除已存在的资源.

一个请求 **必须** 要么成功要么失败 (一个单独的事物).不允许部分成功(原子性,要么成功要么失败).

		注意: 一个 type 成员是必须的在每个请求或者响应资源对象通过 JSON API.有一些情况. type 是必须存在的.	  	
####Creating Resources

A resource can be created by sending a POST request to a URL that represents a collection of resources. The request **MUST** include a single resource object as primary data. The resource object **MUST** contain at least a type member.

For instance, a new photo might be created with the following request:

		POST /photos HTTP/1.1
		Content-Type: application/vnd.api+json
		Accept: application/vnd.api+json
		
		{
		  "data": {
		    "type": "photos",
		    "attributes": {
		      "title": "Ember Hamster",
		      "src": "http://example.com/images/productivity.png"
		    },
		    "relationships": {
		      "photographer": {
		        "data": { "type": "people", "id": "9" }
		      }
		    }
		  }
		}
		
If a relationship is provided in the `relationships` member of the resource object, its value **MUST** be a relationship object with a data member. The value of this key represents the linkage the new resource is to have.

####Creating Resources:创建资源

一个资源集合可以被通过POST访问一个URL来创建单个资源.一个请求 **必须** 好汉一个单独的资源对象作为主要数据. 资源对象 **必须** 包含至少一个 `type` 成员.

举例,一个新的照片可以被通过下列请求创建:

		POST /photos HTTP/1.1
		Content-Type: application/vnd.api+json
		Accept: application/vnd.api+json
		
		{
		  "data": {
		    "type": "photos",
		    "attributes": {
		      "title": "Ember Hamster",
		      "src": "http://example.com/images/productivity.png"
		    },
		    "relationships": {
		      "photographer": {
		        "data": { "type": "people", "id": "9" }
		      }
		    }
		  }
		}

如果一个关联被通过 `relationships` 中的资源对象提供. 其值 **必须** 是一个关联对象并且附带一个 `data` 成员.这个键的值表示一个新的资源的关联.

#####Client-Generated IDs

A server **MAY** accept a client-generated ID along with a request to create a resource. An ID **MUST** be specified with an id key, the value of which **MUST** be a universally unique identifier. The client **SHOULD** use a properly generated and formatted UUID as described in RFC 4122 [RFC4122].

		NOTE: In some use-cases, such as importing data from another source, it may be possible to use something other than a UUID that is still guaranteed to be globally unique. Do not use anything other than a UUID unless you are 100% confident that the strategy you are using indeed generates globally unique identifiers.
		
For example:

		POST /photos HTTP/1.1
		Content-Type: application/vnd.api+json
		Accept: application/vnd.api+json
		
		{
		  "data": {
		    "type": "photos",
		    "id": "550e8400-e29b-41d4-a716-446655440000",
		    "attributes": {
		      "title": "Ember Hamster",
		      "src": "http://example.com/images/productivity.png"
		    }
		  }
		}

A server **MUST** return `403 Forbidden` in response to an unsupported request to create a resource with a client-generated ID.

####Client-Generated IDs:客户生成id

服务端 **可以** 接收一个客户端生成的ID当创建一个资源的时候.一个ID **必须** 以一个 `id` 键值指定, 值 **必须** 是全局唯一的标示符. 客户端 **可以** 使用合适的格式生 UUID 根据 `RFC 4122`描述.

		注意:在一些用户情况,从一些其他数据源导入数据,可能会使用生成全局唯一的 UUID.不能使用任何除了 UUID 意外的键值,除非你能100%确认所有键值是唯一的.
		
示例:

		POST /photos HTTP/1.1
		Content-Type: application/vnd.api+json
		Accept: application/vnd.api+json
		
		{
		  "data": {
		    "type": "photos",
		    "id": "550e8400-e29b-41d4-a716-446655440000",
		    "attributes": {
		      "title": "Ember Hamster",
		      "src": "http://example.com/images/productivity.png"
		    }
		  }
		}

#####Responses

**201 Created**

If a POST request did not include a Client-Generated ID and the requested resource has been created successfully, the server **MUST** return a `201 Created` status code.

The response **SHOULD** include a Location header identifying the location of the newly created resource.

The response **MUST** also include a document that contains the primary resource created.

If the resource object returned by the response contains a self key in its links member and a Location header is provided, the value of the self member MUST match the value of the Location header.

		HTTP/1.1 201 Created
		Location: http://example.com/photos/550e8400-e29b-41d4-a716-446655440000
		Content-Type: application/vnd.api+json
		
		{
		  "data": {
		    "type": "photos",
		    "id": "550e8400-e29b-41d4-a716-446655440000",
		    "attributes": {
		      "title": "Ember Hamster",
		      "src": "http://example.com/images/productivity.png"
		    },
		    "links": {
		      "self": "http://example.com/photos/550e8400-e29b-41d4-a716-446655440000"
		    }
		  }
		}

**202 Accepted**

If a request to create a resource has been accepted for processing, but the processing has not been completed by the time the server responds, the server MUST return a 202 Accepted status code.

**204 No Content**

If a POST request did include a Client-Generated ID and the requested resource has been created successfully, the server MUST return either a 201 Created status code and response document (as described above) or a 204 No Content status code with no response document.

    Note: If a 204 response is received the client should consider the resource object sent in the request to be accepted by the server, as if the server had returned it back in a 201 response.

**403 Forbidden**

A server **MAY** return `403 Forbidden` in response to an unsupported request to create a resource.

**409 Conflict**

A server **MUST** return `409 Conflict` when processing a POST request to create a resource with a client-generated ID that already exists.

A server **MUST** return `409 Conflict` when processing a POST request in which the resource object's type is not among the type(s) that constitute the collection represented by the endpoint.

A server **SHOULD** include error details and provide enough information to recognize the source of the conflict.

**Other Responses**

A server **MAY** respond with other HTTP status codes.

A server **MAY** include error details with error responses.

A server **MUST** prepare responses, and a client **MUST** interpret responses, in accordance with `HTTP semantics`.

#####Responses:响应

**201 Created**

如果一个 `POST` 请求没有包含一个客户端生成的id并且这个请求已经被创建成功了, 服务端 **必须** 返回`201 Created`状态码.

响应 **应该** 包含一个 `Locaiton` 头标识新创建资源的地址.

响应 **必须** 也包含一个文档包含一个主资源创建的文档.

如果一个资源对象的返回包括一个 `self` 键在 `links`成员并且一个 `Location`头也要被提供, `self`成员的值 **必须** 匹配 `Location`头.

		HTTP/1.1 201 Created
		Location: http://example.com/photos/550e8400-e29b-41d4-a716-446655440000
		Content-Type: application/vnd.api+json
		
		{
		  "data": {
		    "type": "photos",
		    "id": "550e8400-e29b-41d4-a716-446655440000",
		    "attributes": {
		      "title": "Ember Hamster",
		      "src": "http://example.com/images/productivity.png"
		    },
		    "links": {
		      "self": "http://example.com/photos/550e8400-e29b-41d4-a716-446655440000"
		    }
		  }
		}
		
**202 Accepted**

如果一个创建资源的请求已经被接受在处理中,但是服务器响应时间还没有处理完成. 服务端 **必须** 返回一个 `202 Accepted` 状态码.

**204 No Content**

如果一个 `POST` 请求没有包含一个客户端生成的id并且这个请求已经被创建成功了, 服务端 **必须** 返回`201 Created`状态码并且响应文档和上述一致或者`204 No Content` 状态码没有响应文档.

		注意:如果客户端收到一个 204 响应客户端应该考虑资源对象的请求已经被接受,类似服务端返回 201 响应.
		
**403 Forbidden**

服务端 **可以** 返回 `403 Forbidden` 响应一个不支持创建资源请求.

**409 Conflict**

服务端 **必须** 返回 `409 Conflict` 当处理一个 `POST`请求创建资源时,客户生成的id已经存在(重复).

服务端 **必须** 返回 `409 Conflict` 当处理一个 `POST`请求创建资源时,创建资源类型的 `type` 不属于接口终端的 `type(s)`范围。

服务端 **应该** 包含错误细节并且提供足够的信息用于辨识这次冲突.

**Other Responses:其他响应**

服务端 **可以** 相应返回其他 HTTP 状态码.

服务端 **可以** 响应错误时会包含`error details`错误细节信息.

服务端 **必须** 准备准备着响应,并且客户端 **必须** 解释响应根据 HTTP 语义.

####Updating Resources

A resource can be updated by sending a `PATCH` request to the URL that represents the resource.

The URL for a resource can be obtained in the `self` link of the resource object. Alternatively, when a `GET` request returns a single resource object as primary data, the same request URL can be used for updates.

The `PATCH` request **MUST** include a single resource object as primary data. The resource object MUST contain `type` and `id` members.

For example:

		PATCH /articles/1 HTTP/1.1
		Content-Type: application/vnd.api+json
		Accept: application/vnd.api+json
		
		{
		  "data": {
		    "type": "articles",
		    "id": "1",
		    "attributes": {
		      "title": "To TDD or Not"
		    }
		  }
		}

####Updating Resources:更新资源

一个资源可以被通过一个表示资源的URL发起`PATCH`请求.

一个表示资源的url可以通过资源对象中的`self`link中找见.另外,当通过 `GET`请求获取一个单独的资源作为主要数据的时候,相同的URL可以被用于更新.

`PATH`请求 **必须** 包含一个单独的资源对象作为主要数据.资源对象 **必须** 包含 `type`和`id` 成员.

举例:

		PATCH /articles/1 HTTP/1.1
		Content-Type: application/vnd.api+json
		Accept: application/vnd.api+json
		
		{
		  "data": {
		    "type": "articles",
		    "id": "1",
		    "attributes": {
		      "title": "To TDD or Not"
		    }
		  }
		}

#####Updating a Resource's Attributes

Any or all of a resource's `attributes` **MAY** be included in the resource object included in a PATCH request.

If a request does not include all of the attributes for a resource, the server **MUST** interpret the missing attributes as if they were included with their current values. The server **MUST NOT** interpret missing attributes as `null` values.

For example, the following PATCH request is interpreted as a request to update only the `title` and `text` attributes of an article:

		PATCH /articles/1 HTTP/1.1
		Content-Type: application/vnd.api+json
		Accept: application/vnd.api+json
		
		{
		  "data": {
		    "type": "articles",
		    "id": "1",
		    "attributes": {
		      "title": "To TDD or Not",
		      "text": "TLDR; It's complicated... but check your test coverage regardless."
		    }
		  }
		}

#####Updating a Resource's Attributes:更新资源属性

一个或多个资源的属性 **可以** 被包含在 一个`PATCH`请求的资源对象中.

如果一个请求不能包含所有的资源属性,服务端 **必须** 解释遗失的属性如果它们之前被包含在内.服务端 **必须不能** 把遗失的属性作为 `null` 值.

举例,下例 `PATCH` 请求指示更新一片文章的 `title`和`text`.

		PATCH /articles/1 HTTP/1.1
		Content-Type: application/vnd.api+json
		Accept: application/vnd.api+json
		
		{
		  "data": {
		    "type": "articles",
		    "id": "1",
		    "attributes": {
		      "title": "To TDD or Not",
		      "text": "TLDR; It's complicated... but check your test coverage regardless."
		    }
		  }
		}

#####Updating a Resource's Relationships

Any or all of a resource's relationships MAY be included in the resource object included in a PATCH request.

If a request does not include all of the relationships for a resource, the server **MUST** interpret the missing relationships as if they were included with their current values. It **MUST NOT** interpret them as `null` or empty values.

If a relationship is provided in the relationships member of a resource object in a PATCH request, its value **MUST** be a relationship object with a data member. The relationship's value will be replaced with the value specified in this member.

For instance, the following `PATCH` request will update the `author` relationship of an article:

		PATCH /articles/1 HTTP/1.1
		Content-Type: application/vnd.api+json
		Accept: application/vnd.api+json
		
		{
		  "data": {
		    "type": "articles",
		    "id": "1",
		    "relationships": {
		      "author": {
		        "data": { "type": "people", "id": "1" }
		      }
		    }
		  }
		}

Likewise, the following `PATCH` request performs a complete replacement of the `tags` for an article:

		PATCH /articles/1 HTTP/1.1
		Content-Type: application/vnd.api+json
		Accept: application/vnd.api+json
		
		{
		  "data": {
		    "type": "articles",
		    "id": "1",
		    "relationships": {
		      "tags": {
		        "data": [
		          { "type": "tags", "id": "2" },
		          { "type": "tags", "id": "3" }
		        ]
		      }
		    }
		  }
		}

A server **MAY** reject an attempt to do a full replacement of a to-many relationship. In such a case, the server **MUST** reject the entire update, and return a `403 Forbidden` response.

    Note: Since full replacement may be a very dangerous operation, a server may choose to disallow it. For example, a server may reject full replacement if it has not provided the client with the full list of associated objects, and does not want to allow deletion of records the client has not seen.

#####Updating a Resource's Relationships:更新资源的关联

一个或多个资源的关联 **可以** 被包含在 一个`PATCH`请求的资源对象中.

如果一个请求不能包含所有的资源关联,服务端 **必须** 解释遗失的关联如果它们之前被包含在内.服务端 **必须不能** 把遗失的属性作为 `null` 值或者空值.

如果一个`PATCH`请求的`relationships`的资源对象成员包含一个关联,它的值 **必须** 是在一个 `data` 成员的 `relationship object` 关联对象. 这个关联的值将会被指定的值替换.

示例: 下列 `PATCH` 请求将会更新 `author` 一篇文章的作者关联.

		PATCH /articles/1 HTTP/1.1
		Content-Type: application/vnd.api+json
		Accept: application/vnd.api+json
		
		{
		  "data": {
		    "type": "articles",
		    "id": "1",
		    "relationships": {
		      "author": {
		        "data": { "type": "people", "id": "1" }
		      }
		    }
		  }
		}
 
 同样的,下列 `PATCH` 请求将会替换文章的 `tags` 标签(替换关联).
 
		PATCH /articles/1 HTTP/1.1
		Content-Type: application/vnd.api+json
		Accept: application/vnd.api+json
		
		{
		  "data": {
		    "type": "articles",
		    "id": "1",
		    "relationships": {
		      "tags": {
		        "data": [
		          { "type": "tags", "id": "2" },
		          { "type": "tags", "id": "3" }
		        ]
		      }
		    }
		  }
		}

服务端 **可以** 拒绝一个尝试去做一个全部替换一对多关联的请求. 这种情况,服务端 **必须** 拒绝整体更新, 并且返回一个 `403 Forbidden` 响应.

		注意: 因为全部替换请求可能是一个非常危险的操作, 服务端可以选择不允许此操作.举例来说,一个服务端可以拒绝全部替换如果它没有对客户端提供一个全部对象列表,并且也不希望客户端删掉它看不见的记录.

#####Responses

**202 Accepted**

If an update request has been accepted for processing, but the processing has not been completed by the time the server responds, the server MUST return a 202 Accepted status code.

**200 OK**

If a server accepts an update but also changes the resource(s) in ways other than those specified by the request (for example, updating the updated-at attribute or a computed sha), it MUST return a 200 OK response. The response document MUST include a representation of the updated resource(s) as if a GET request was made to the request URL.

A server MUST return a 200 OK status code if an update is successful, the client's current attributes remain up to date, and the server responds only with top-level meta data. In this case the server MUST NOT include a representation of the updated resource(s).

**204 No Content**

If an update is successful and the server doesn't update any attributes besides those provided, the server MUST return either a `200 OK` status code and response document (as described above) or a `204 No Content` status code with no response document.

**403 Forbidden**

A server **MUST** return `403 Forbidden` in response to an unsupported request to update a resource or relationship.

**404 Not Found**

A server **MUST** return `404 Not Found` when processing a request to modify a resource that does not exist.

A server **MUST** return `404 Not Found` when processing a request that references a related resource that does not exist.

**409 Conflict**

A server **MAY** return `409 Conflict` when processing a PATCH request to update a resource if that update would violate other server-enforced constraints (such as a uniqueness constraint on a property other than id).

A server **MUST** return `409 Conflict` when processing a PATCH request in which the resource object's `type` and `id` do not match the server's endpoint.

A server **SHOULD** include error details and provide enough information to recognize the source of the conflict.

#####Response:响应

**202 Accepted**

如果一个更新请求被接受处理,但是处理还没有完成,服务端 **必须** 返回一个 `202 Accepted` 状态码.

**200 OK**

如果服务器接受更新，但是在请求指定内容之外做了资源修改,必须响应200 OK以及更新的资源实例,像是向此URL发出GET请求.

**204 No Content**

如果服务器接受更新,且没有在指定内容之外做了资源更改,服务端 **必须** 要么返回一个 `200 OK`状态吗和一个更新的资源实例,或者一个 `204 No Content`状态码不返回响应文档.

**403 Forbidden**

服务端 **必须** 返回一个 `403 Forbidden` 在响应一个不支持的请求区更新一个资源或关联.

**404 Not Found**

服务端 **必须** 返回 `404 Not Found` 当处理一个请求更新一个不存在的资源.

服务端 **必须** 返回 `404 Not Found` 当处理一个请求更新一个不存在的关联资源.

**409 Conflict**

服务端 **可以** 返回 `409 Conflict` 当处理一个 `PATCH` 请求去更新一个资源,但是该更新会违反服务器的限制(例如唯一限制,尝试更新一个id,但是id是唯一限制,已经存在).

服务端 **必须** 返回 `409 Conflict` 当处理一个 `PATCH` 请求去更新一个资源,该资源的 `type` 和 `id` 与服务器端不匹配.

服务端 **应该** 包含 `error details` 错误细节并且提供足够的信息解释资源的冲突.

**Other Responses**

A server **MAY** respond with other HTTP status codes.

A server **MAY** include error details with error responses.

A server **MUST** prepare responses, and a client MUST interpret responses, in accordance with HTTP semantics.

**Other Responses:其它响应**

服务器使用其它HTTP错误状态码反映错误.客户端必须依据HTTP规范处理这些错误信息.如下所述,错误细节可能会一并返回.

####Updating Relationships

Although relationships can be modified along with resources (as described above), JSON API also supports updating of relationships independently at URLs from relationship links.

    	Note: Relationships are updated without exposing the underlying server semantics, such as foreign keys. Furthermore, relationships can be updated without necessarily affecting the related resources. For example, if an article has many authors, it is possible to remove one of the authors from the article without deleting the person itself. Similarly, if an article has many tags, it is possible to add or remove tags. Under the hood on the server, the first of these examples might be implemented with a foreign key, while the second could be implemented with a join table, but the JSON API protocol would be the same in both cases.

    	Note: A server may choose to delete the underlying resource if a relationship is deleted (as a garbage collection measure).
    	
####Updating Relationships:更新关联

尽管关联可以随着资源一起更新(上述表述), JSON API 仍然支持更新一个关系根据单独的 URL 从 `relationship links`.

#####Updating To-One Relationships

A server **MUST** respond to `PATCH` requests to a URL from a to-one relationship link as described below.

The `PATCH` request **MUST** include a top-level member named `data` containing one of:

* a resource identifier object corresponding to the new related resource.
* `null`, to remove the relationship.

For example, the following request updates the author of an article:

		PATCH /articles/1/relationships/author HTTP/1.1
		Content-Type: application/vnd.api+json
		Accept: application/vnd.api+json
		
		{
		  "data": { "type": "people", "id": "12" }
		}

And the following request clears the author of the same article:

		PATCH /articles/1/relationships/author HTTP/1.1
		Content-Type: application/vnd.api+json
		Accept: application/vnd.api+json
		
		{
		  "data": null
		}

If the relationship is updated successfully then the server **MUST** return a successful response.

#####Updating To-One Relationships:更新单对象关联

服务单 **必须** 响应一个 `PATCH` 请求一个单对象关联的URL.

该 `PATCH` 请求(注意是请求,不是响应) **必须** 包含一个顶级 `data` 成员包含其中一个:

* 一个新关联资源的资源标识符对象
* `null`,移除一个链接

示例,下列请求更新一篇文章的作者:

		PATCH /articles/1/relationships/author HTTP/1.1
		Content-Type: application/vnd.api+json
		Accept: application/vnd.api+json
		
		{
		  "data": { "type": "people", "id": "12" }
		}

下面请求删除该篇文章的作者:

		PATCH /articles/1/relationships/author HTTP/1.1
		Content-Type: application/vnd.api+json
		Accept: application/vnd.api+json
		
		{
		  "data": null
		}

如果一个关联更新成功,服务端 **必须** 返回一个成功的响应.

#####Updating To-Many Relationships

A server **MUST** respond to `PATCH`, `POST`, and `DELETE` requests to a URL from a to-many relationship link as described below.

For all request types, the body **MUST** contain a data member whose value is an empty array or an array of resource identifier objects.

If a client makes a `PATCH` request to a URL from a to-many relationship link, the server MUST either completely replace every member of the relationship, return an appropriate error response if some resources can not be found or accessed, or return a `403 Forbidden` response if complete replacement is not allowed by the server.

For example, the following request replaces every tag for an article:

		PATCH /articles/1/relationships/tags HTTP/1.1
		Content-Type: application/vnd.api+json
		Accept: application/vnd.api+json
		
		{
		  "data": [
		    { "type": "tags", "id": "2" },
		    { "type": "tags", "id": "3" }
		  ]
		}

And the following request clears every tag for an article:

		PATCH /articles/1/relationships/tags HTTP/1.1
		Content-Type: application/vnd.api+json
		Accept: application/vnd.api+json
		
		{
		  "data": []
		}

If a client makes a `POST` request to a URL from a relationship link, the server **MUST** add the specified members to the relationship unless they are already present. If a given `type` and `id` is already in the relationship, the server **MUST NOT** add it again.

    	Note: This matches the semantics of databases that use foreign keys for has-many relationships. Document-based storage should check the has-many relationship before appending to avoid duplicates.

If all of the specified resources can be added to, or are already present in, the relationship then the server **MUST** return a successful response.

  		Note: This approach ensures that a request is successful if the server's state matches the requested state, and helps avoid pointless race conditions caused by multiple clients making the same changes to a relationship.

In the following example, the comment with ID `123` is added to the list of comments for the article with ID `1`:

		POST /articles/1/relationships/comments HTTP/1.1
		Content-Type: application/vnd.api+json
		Accept: application/vnd.api+json
		
		{
		  "data": [
		    { "type": "comments", "id": "123" }
		  ]
		}

If the client makes a `DELETE` request to a URL from a relationship link the server **MUST** delete the specified members from the relationship or return a `403 Forbidden` response. If all of the specified resources are able to be removed from, or are already missing from, the relationship then the server **MUST** return a successful response.

    	Note: As described above for POST requests, this approach helps avoid pointless race conditions between multiple clients making the same changes.

Relationship members are specified in the same way as in the POST request.

In the following example, comments with IDs of `12` and `13` are removed from the list of comments for the article with ID `1`:

		DELETE /articles/1/relationships/comments HTTP/1.1
		Content-Type: application/vnd.api+json
		Accept: application/vnd.api+json
		
		{
		  "data": [
		    { "type": "comments", "id": "12" },
		    { "type": "comments", "id": "13" }
		  ]
		}

    	Note: RFC 7231 specifies that a DELETE request may include a body, but that a server may reject the request. This spec defines the semantics of a server, and we are defining its semantics for JSON API.

#####Updating To-Many Relationships:更新多关联对象 

服务端 **必须** 响应 `PATCH`,`POST` 和 `DELETE` 请求对一个多关联对象的 URL.    	
对于所有请求的类型, `body` **必须** 包含一个 `data` 成员其值是一个空的数组 或者一个 资源标识符的数组.

如果一个客户端发起一个多关联对象的 `PATCH` 请求,服务端 **必须** 要么完全替换所有成员的关系,如果一些资源不能被找见或者访问,则返回一个相关的错误响应.或者返回一个 `403 Forbidden` 响应如果 `complete replacement`(完全替换)服务端不允许.

示例,下列请求替换文章的每个标签:

		PATCH /articles/1/relationships/tags HTTP/1.1
		Content-Type: application/vnd.api+json
		Accept: application/vnd.api+json
		
		{
		  "data": [
		    { "type": "tags", "id": "2" },
		    { "type": "tags", "id": "3" }
		  ]
		}

下列请求删除文章的每个标签:

		PATCH /articles/1/relationships/tags HTTP/1.1
		Content-Type: application/vnd.api+json
		Accept: application/vnd.api+json
		
		{
		  "data": []
		}

如果一个客户端发起一个关于`relationship link`关联链接的 `POST` 请求,服务端 **必须** 添加指定的关联成员除非它们已经存在. 如果 `type` 和 `id` 已经攒在与关联中,服务端 **必须** 不能再次添加.

如果所有指定资源可以被添加,或者已经存在,服务端 **必须** 返回一个成功的响应.

		注意:这样是保证请求总是成功的如果服务端的状态匹配请求状态.并且帮助避免多个客户端发起同样的关联请求引发的无意义的竞争.
		
在下列示例, ID 为 `123`的评论被添加的 ID 为 `1` 的文章.

		POST /articles/1/relationships/comments HTTP/1.1
		Content-Type: application/vnd.api+json
		Accept: application/vnd.api+json
		
		{
		  "data": [
		    { "type": "comments", "id": "123" }
		  ]
		}		

如果客户端向一个 `relationship link` 发起一个 `DELETE` URL 请求,服务端 **必须** 删除指定成员从这个关联中或者返回一个 `403 Forbidden` 响应.如果所有指定资源是可以被删除,或者已经删除了,服务端 **必须** 返回一个成功的响应.

下列示例, ID为 `12` 和 `13` 的评论被从 ID为 `1`的文章评论列表中删除:

		DELETE /articles/1/relationships/comments HTTP/1.1
		Content-Type: application/vnd.api+json
		Accept: application/vnd.api+json
		
		{
		  "data": [
		    { "type": "comments", "id": "12" },
		    { "type": "comments", "id": "13" }
		  ]
		}					
    	
#####Responses

**202 Accepted**

If a relationship update request has been accepted for processing, but the processing has not been completed by the time the server responds, the server **MUST** return a `202 Accepted` status code.

**204 No Content**

A server **MUST** return a `204 No Content` status code if an update is successful and the representation of the resource in the request matches the result.

		Note: This is the appropriate response to a POST request sent to a URL from a to-many relationship link when that relationship already exists. It is also the appropriate response to a DELETE request sent to a URL from a to-many relationship link when that relationship does not exist.
		
**200 OK**

If a server accepts an update but also changes the targeted relationship(s) in other ways than those specified by the request, it **MUST** return a `200 OK` response. The response document MUST include a representation of the updated relationship(s).

A server **MUST** return a `200 OK` status code if an update is successful, the client's current data remain up to date, and the server responds only with top-level meta data. In this case the server **MUST NOT** include a representation of the updated relationship(s).

**403 Forbidden**

A server **MUST** return `403 Forbidden` in response to an unsupported request to update a relationship.

**Other Responses**

A server **MAY** respond with other HTTP status codes.

A server **MAY** include error details with error responses.

A server **MUST** prepare responses, and a client **MUST** interpret responses, in accordance with `HTTP semantics`.

#####Response:响应

**202 Accepted**

如果一个关联更新请求被接受处理,但是处理还没有完成,服务端 **必须** 返回一个 `202 Accepted` 状态码.

**204 No Content**

如果更新成功,且客户端属性保持最新,服务器必须返回`204 No Content`状态码.

		注意: 这个响应发生在,当POST向一个URL发出一个对多关联的请求当这些关联已经存在,当DELETE向一个URL发出对多关联当关联已经不存在.
		
**200 OK**

如果服务器接受更新,但是在请求指定内容之外做了资源修改,，必须响应`200 OK`以及更新的资源实例,像是向此URL发出GET请求.

**403 Forbidden**

服务端 **必须** 返回 `403 Forbidden` 响应对一个不支持的更新关联请求.

**Other Responses:其它响应**

服务器使用其它HTTP错误状态码反映错误.客户端必须依据HTTP规范处理这些错误信息.如下所述,错误细节可能会一并返回.

####Deleting Resources

An individual resource can be deleted by making a `DELETE` request to the resource's URL:

		DELETE /photos/1 HTTP/1.1
		Accept: application/vnd.api+json
		
####Deleting Resources:资源删除

向资源URL发出DELETE请求即可删除单个资源.

		DELETE /photos/1 HTTP/1.1
		Accept: application/vnd.api+json		
		
服务器可以选择性的支持,在一个请求里删除多个资源.
		
		DELETE /photos/1,2,3
		Accept: application/vnd.api+json	
		
#####Responses

**202 Accepted**

If a deletion request has been accepted for processing, but the processing has not been completed by the time the server responds, the server **MUST** return a `202 Accepted` status code.

**204 No Content**

A server **MUST** return a `204 No Content` status code if a deletion request is successful and no content is returned.

**200 OK**

A server **MUST** return a `200 OK` status code if a deletion request is successful and the server responds with only top-level meta data.

**Other Responses**

A server **MAY** respond with other HTTP status codes.

A server **MAY** include error details with error responses.

A server **MUST** prepare responses, and a client **MUST** interpret responses, in accordance with `HTTP semantics`.

#####Responses:响应

**202 Accepted**

如果一个删除请求被接受处理,但是处理还没有完成,服务端 **必须** 返回一个 `202 Accepted` 状态码.

**204 No Content**

没有内容返回,服务端 **必须** 返回一个 `204 No Content` 状态码如果一个删除请求成功.

**200 OK**

如果有内容返回,仅返回`top-level meta data`顶级元数据,服务端 **必须** 返回一个 `200 OK` 状态码如果一个删除请求成功.

**Other Responses:其它响应**

服务器使用其它HTTP错误状态码反映错误.客户端必须依据HTTP规范处理这些错误信息.如下所述,错误细节可能会一并返回.

###Query Parameters

Implementation specific query parameters **MUST** adhere to the same constraints as member names with the additional requirement that they **MUST** contain at least one non a-z character (U+0061 to U+007A). It is **RECOMMENDED** that a U+002D HYPHEN-MINUS, "-", U+005F LOW LINE, "_", or capital letter is used (e.g. camelCasing).

If a server encounters a query parameter that does not follow the naming conventions above, and the server does not know how to process it as a query parameter from this specification, it **MUST** return `400 Bad Request`.

    	Note: This is to preserve the ability of JSON API to make additive additions to standard query parameters without conflicting with existing implementations.

###Query Parameters:查询参数

应用特定查询参数 **必须** 伴随着相同的 `member names`成员名称限制.更多的限制 **必须** 包含至少一个 `a-z`(U+0061 to U+007A)) 字符. 同时 **推荐
** 使用减号 `-`(U+002D HYPHEN-MINUS),下划线 `_` (U+005F LOW LINE) 或者 大写字母(capital letter).

如果服务端遇到一个查询参数不再上述限制之内,请呗服务端不知道怎么处理该参数,服务端 **必须** 返回 `400 Bad Request`.

###Errors

####Processing Errors

A server **MAY** choose to stop processing as soon as a problem is encountered, or it **MAY** continue processing and encounter multiple problems. For instance, a server might process multiple attributes and then return multiple validation problems in a single response.

When a server encounters multiple problems for a single request, the most generally applicable HTTP error code **SHOULD** be used in the response. For instance, `400 Bad Request` might be appropriate for multiple 4xx errors or `500 Internal Server Error` might be appropriate for multiple 5xx errors.

####Error Objects

Error objects provide additional information about problems encountered while performing an operation. Error objects **MUST** be returned as an array keyed by `errors` in the top level of a JSON API document.

An error object **MAY** have the following members:

* `id`: a unique identifier for this particular occurrence of the problem.
* `links`: a links object containing the following members:
	* `about`: a link that leads to further details about this particular occurrence of the problem.
* `status`: the HTTP status code applicable to this problem, expressed as a string value.
* `code`: an application-specific error code, expressed as a string value.
* `title`: a short, human-readable summary of the problem that **SHOULD NOT** change from occurrence to occurrence of the problem, except for purposes of localization.
* `detail`: a human-readable explanation specific to this occurrence of the problem. Like title, this field's value can be localized.
* `source`: an object containing references to the source of the error, optionally including any of the following members:
	* `pointer`: a JSON Pointer [RFC6901] to the associated entity in the request document [e.g. "`/data`" for a primary data object, or "`/data/attributes/title`" for a specific attribute].
	* `parameter`: a string indicating which URI query parameter caused the error.
* `meta`: a meta object containing non-standard meta-information about the error.

###Errors:错误

####Processing Errors:处理错误

服务器 **可以** 选择在第一个问题出现时,立刻终止PATCH 操作,或者继续执行,遇到多个问题.例如,服务器可能多属性更新,然后返回在一个响应里返回多个校验问题.

当服务器单个请求遇到多个问题,响应中**应该**指定最通用可行的HTTP错误码. 举例,`400 Bad Request`适用于多个4xx errors,`500 Internal Server Error`适用于多个5xx errors.

####Error Objects:错误对象

错误对象是特殊化的资源对象,可能在响应中一并返回,用以提供执行操作遭遇问题的额外信息.在在JSON API文档顶层,"errors"对应值即为错误对象集合 **必须** 返回,此时文档不应该包含其它顶层资源.

错误对象可能有以下元素：

* `id`: 特定问题的唯一标示符.
* `links`: 一个连接对象(可以在请求文档中取消应用的关联资源)包含下列成员:
	* `about`: 提供特定问题更多细节的链接.
* `status`: 适用于这个问题的HTTP状态码,使用字符串表示.
* `code`:  应用特定的错误码,以字符串表示.
* `title`: 简短的,可读性高的问题总结.除了国际化本地化处理之外,不同场景下,相同的问题,值是 **不应该** 变动的.
* `detail`: 针对该问题的高可读性解释.
* `source`: 一个包含错误来源的对象,可以包含下列成员:
	* `pointer`: 一个关联整体请求文档的 JSON 节点 [例如 "`/data`" 用于主要数据对象, or "`/data/attributes/title`" 用于一个特定属性].
	* `parameter`: 一个标记引起错误的参数.
* `meta`: 一个和错误相关包含非标准元数据的元对象.