/*
 Author:chloroplast
 */
#WEB API 设计方法论

---

官方资料可以翻阅 RESTful Web APIs中文版

该文章是采摘自[七步法设计API][id]
[id]:http://www.infoq.com/cn/articles/web-api-design-methodology/


###前序

---

通读完整本书,其实restful是让人关注于一个资源状态移交的过程(representational state transfer).终归其实是想做出一个可让人像理解网站一样(你访问淘宝,不需要看任何说明书就能访问)去理解一个接口.在作出结论,一切的规范性都是为了更好的约束代码和降低沟通成本.

###1.罗列语意描述符

第一步是列出客户端程序可能要从我们的服务中获取的，或要放到我们的服务中的所有数据片段。我们将这些称为语义描述符。语义是指它们处理数据在应用程序中的含义，描述符是指它们描述了在应用程序自身中发生了什么。

**这里的视点是客户端，不是服务器端。将API设计成客户端使用的东西很重要。** 

*rest很多时候强调从客户端观察,亦或是从服务器端观察.*

比如在一个简单的待办事项列表应用中，你可能会找到下面这些语义描述符：

> id : 系统中每条记录的唯一标识符   
> title : 每个待办事项的标题  
> dateDue : 待办事项应该完成的日期  
> complete : 一个是/否标记，表明待办事项是否已经完成了。  

###2.画状态图

下一步是根据建议的API绘制出状态图。图中的每个框都表示一种可能的表示--一个包含在步骤1中确定的一或多个语义描述符的文档。你可以用箭头表示从一个框到下一个的转变--从一个状态到下一个状态。这些转变是由协议请求触发的。

**在每次变化中还不用急着指明用哪个协议方法。只要标明变化是`安全的`（比如HTTP GET），还是`不安全/非幂等`的（比如HTTP.POST），或者`不安全/幂等`的（PUT）。**

幂等动作是指重复执行时不会有无法预料的副作用。比如HTTP PUT ，因为规范说服务器应该用客户端传来的状态值替换目标资源的已有值，所以说它是幂等的。而 HTTP POST 是非幂等的，因为规范指出提交的值应该是追加到已有资源集合上的，而不是替换。

![toDoList](./img/restful/restful-1.png "toDoList")

这个状态图中展示的这些动作也是语义描述符-- 它们描述了这个服务的语义动作。

> read-list  
> filter-list  
> read-item  
> create-item  
> mark-complete

在你做这个状态图的过程中，你可能会发现自己漏掉了客户端真正想要或需要的动作或数据项。这是退回到步骤1的机会，添加一些新的描述符，并/或者在步骤2中改进状态图。


在你重新迭代过这两步之后，你应该对客户端跟服务交互所需的所有数据点和动作有了好的认识和想法

###3.调整命名

下一步是调和服务接口中的所有“魔法字符串”。“魔法字符串” 全是描述符的名称--它们没有内在的含义，只是表示客户端跟你的服务通讯时将要访问的动作或数据元素。调和这些描述符名称的意思是指采用源自下面这些地方的，知名度更高的公共名称：

> Schema.org  
> microformats.org  
> Dublin Core  
> IANA Link Relation Values  

这些全是明确定义的、共享的名称库。当你服务接口使用来自这些源头的名称时，开发人员很可能之前见过并知道它们是什么意思。这可以提高API的可用性。


说明：尽管在服务接口上使用共享名称是个好主意，但在内部实现里可以不用（比如数据库里的数据域名称）。服务自身可以毫不困难地将公共接口名称映射为内部存储名称。


以待办事项服务为例，除了一个语义描述符- create-item，我能找到所有可接受的已有名称。为此我根据Web Linking RFC5988中的规则创建了一个具有唯一性的URI。在给接口描述符选择知名的名称时需要折中。它们极少能跟你的内部数据存储元素完美匹配，不过那没关系。

> id -> [来自Dublin Core的identifier](http://purl.org/dc/elements/1.1/identifier)  
> title -> [来自Schema.org的name](https://schema.org/name)  
> dueDate -> [来自Schema.org的scheduledTime](https://schema.org/scheduledtime)  
> complete -> [来自Schema.org的status](https://schema.org/status,)  
> read-list -> [来自IANA Link Relation Values的collection](http://www.iana.org/assignments/link-relations/link-relations.xhtml)  
> create-item ->用RFC5988的http://mamund.com/rels/create-item  
> mark-complete -> [来自IANA Link Relation Values的edit](http://www.iana.org/assignments/link-relations/link-relations.xhtml)  



![经过调和的图](./img/restful/restful-2.png "经过调和的图")

###4.选一个媒体类型

API设计过程的下一步是选一个媒体类型，用来在客户端和服务器端之间传递消息。Web的特点之一是数据是通过统一的接口作为标准化文档传输的。选择同时支持数据描述符（比如"identifier"、"status"等）和动作描述符（比如"search"、"edit"等）的媒体类型很重要。有相当多可用的格式。

> [Collection+JSON](/guide/restful-CollectionJson)  

**个人建议选择此媒体类型**


###5.编写profile

语义档案是一个文档，其中列出了设计中的所有描述符，包括对开发人员构建客户端和服务器端实现有帮助的所有细节。这个档案是一个实现指南，不是实现描述。这个差别很重要。

> [JSON-LD](/guide/restful-JsonLD)  


**个人建议选择此媒体类型 jsonLD配合Hydra一起**

###6.实现
代码实现

###7.发布
接口发布


###设计建议

---
1. 设计流程关注的是状态转换和语意描述符.
2. 像设计网站一样将同样的思想应用于API设计中. **不要**发布基于数据库模式的api,这样改变数据库模式就基本上变成不可能的事情了.
3. **URL设计并不重要**.一个URL仅仅是摸个资源的地址而已,客户端可以使用它得一个表述.URL并不能在技术上说明任何关于资源或者表述的信息.<<万维网的架构,第一卷>>.用更可靠的方式来描述这些内容:媒体类型定义以及机器可读的profile.
4. 标准名称**优于**自定义名称.
5. vnd.前缀(用于商业项目). prs.前缀(用于个人或者试验工作)

