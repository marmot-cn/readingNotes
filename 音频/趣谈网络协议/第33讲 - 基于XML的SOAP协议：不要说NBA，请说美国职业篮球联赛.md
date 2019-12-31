# 第33讲 | 基于XML的SOAP协议：不要说NBA，请说美国职业篮球联赛

## 笔记

### ONC RPC 存在哪些问题?

客户端要发送的参数, 以及服务端要发送的回复, 都压缩为一个二进制串.

问题:

1. 需要双方的压缩格式完全一致, 可以用传输层的可靠性以及加入校验值等方式, 来减少传输过程中的差错.
2. 协议修改不灵活, 服务端的业务逻辑升级, 增加或删除了字段, 双方没有及时通知, 就造成解压缩不成功.

当业务发生改变, 需要多传输一些参数或者少传输一些参数的时候, 都需要及时通知对方, 并且根据约定好的协议文件重新生成双方的`Stub`程序.

版本问题: 一个服务端对应多个客户端, 因为一个客户端有一个需要添加一个字段. 如果服务端改了, 所有客户端都需要适配.

**ONC RPC 的设计是面向函数的, 而非是面向对象的**

RPC框架适合客户端和服务端的开发人员要密切沟通的场景.

### XML 与 SOAP

使用`NBA`(压缩后的语言)可能会听不懂, 但是说**美国职业篮球赛**(文本语言)所有人都会听懂.

使用**文本类**的方式进行传输.

```xml

<?xml version="1.0" encoding="UTF-8"?>
<geek:purchaseOrder xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:geek="http://www.example.com/geek">
    <order>
        <date>2018-07-01</date>
        <className>趣谈网络协议</className>
        <Author>刘超</Author>
        <price>68</price>
    </order>
</geek:purchaseOrder>
```

* 格式没必要完全一致, `author`和`price`字段就算换了位置, 也不影响解析.
* 客户端想增加一个字段, 只需要在上面的文件中加一行, 对于不需要这个字段的客户端, 只要不解析这一行就是.
* 这种表述方式是描述一个订单对象的, 是一种面向对象的, 更加接近用户场景的表示方式.

### 传输协议问题

基于`XML`的最著名的就是通信协议就是`SOAP`, **简单对象访问协议(Simple Object Access Protocol)**. 它使用`XML`编写简单的请求和回复消息, 并用`HTTP`协议进行传输.

`SOAP`将请求和回复放在一个信封里面, 就像传递一个邮件一样. 信封里面的信分**抬头**和**正文**.

```
POST /purchaseOrder HTTP/1.1
Host: www.geektime.com
Content-Type: application/soap+xml; charset=utf-8
Content-Length: nnn
```

```
<?xml version="1.0"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2001/12/soap-envelope"
soap:encodingStyle="http://www.w3.org/2001/12/soap-encoding">
    <soap:Header>
        <m:Trans xmlns:m="http://www.w3schools.com/transaction/"
          soap:mustUnderstand="1">1234
        </m:Trans>
    </soap:Header>
    <soap:Body xmlns:m="http://www.geektime.com/perchaseOrder">
        <m:purchaseOrder">
            <order>
                <date>2018-07-01</date>
                <className>趣谈网络协议</className>
                <Author>刘超</Author>
                <price>68</price>
            </order>
        </m:purchaseOrder>
    </soap:Body>
</soap:Envelope>
```

这个请求使用`POST`方法, 发送一个格式为`application/soap+xml`的`XML`正文给`www.geektime.com`, 从而下一个订单, 这个订单封装在`SOAP`的信封里面, 并且表明这是一笔交易(`transaction`), 而且订单的详情都已经写明了.

### 协议约定问题

双方的协议约定是什么样的? 如果写文档也不一定及时更新. 需要一种相对比较严谨的**Web 服务描述语言, WSDL(Seb Service Description Languages)**. 它也是一个`XML`文件.

要定义个类型`order`, 与上面的`XML`对应起来.

```
 <wsdl:types>
  <xsd:schema targetNamespace="http://www.example.org/geektime">
   <xsd:complexType name="order">
    <xsd:element name="date" type="xsd:string"></xsd:element>
<xsd:element name="className" type="xsd:string"></xsd:element>
<xsd:element name="Author" type="xsd:string"></xsd:element>
    <xsd:element name="price" type="xsd:int"></xsd:element>
   </xsd:complexType>
  </xsd:schema>
 </wsdl:types>
```

定一个`message`的结构.

```
 <wsdl:message name="purchase">
  <wsdl:part name="purchaseOrder" element="tns:order"></wsdl:part>
 </wsdl:message>
```

暴露一个端口.

```
 <wsdl:portType name="PurchaseOrderService">
  <wsdl:operation name="purchase">
   <wsdl:input message="tns:purchase"></wsdl:input>
   <wsdl:output message="......"></wsdl:output>
  </wsdl:operation>
 </wsdl:portType>
```

编写一个`binding`, 将上面定义的信息绑定到`SOAP`请求的`body`里面.

```
 <wsdl:binding name="purchaseOrderServiceSOAP" type="tns:PurchaseOrderService">
  <soap:binding style="rpc"
   transport="http://schemas.xmlsoap.org/soap/http" />
  <wsdl:operation name="purchase">
   <wsdl:input>
    <soap:body use="literal" />
   </wsdl:input>
   <wsdl:output>
    <soap:body use="literal" />
   </wsdl:output>
  </wsdl:operation>
 </wsdl:binding>
```

编写`service`

```
 <wsdl:binding name="purchaseOrderServiceSOAP" type="tns:PurchaseOrderService">
  <soap:binding style="rpc"
   transport="http://schemas.xmlsoap.org/soap/http" />
  <wsdl:operation name="purchase">
   <wsdl:input>
    <soap:body use="literal" />
   </wsdl:input>
   <wsdl:output>
    <soap:body use="literal" />
   </wsdl:output>
  </wsdl:operation>
 </wsdl:binding>
```

有的工具可以根据`WSDL`生成客户端的`Stub`, 让客户端通过`Stub`进行远程调用, 就跟调用本地的方法一样.

### 服务发现问题

`UDDI(Universal Description, Descovery, and Integration)`, 统一描述, 发现和集成协议. 是一个注册中心, 服务提供方可以将上面的`WSDL`描述文件, 发布到这个注册中心, 注册完毕后, 服务使用方可以查找到服务的描述, 封装为本地的客户端调用.

### 总结

`SOAP`的三大要素:

* 协议约定用`WSDL`
* 传输协议用`HTTP`
* 服务发现用`UDDL`

## 扩展