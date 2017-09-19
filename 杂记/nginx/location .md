# nginx location 

---

## location 语法

在nginx的location和配置中location的顺序没有太大关系. 和location表达式的类型有关. 相同类型的表达式, 字符串长的会优先匹配.

`location [=|~|~*|^~] /uri/ { … }`

* `=`: 严格匹配. 如果这个查询匹配, 那么将停止搜索并立即处理此请求.
* `~`: 区分大小写**匹配**(可用正则表达式).
* `~*`: 不区分大小写**匹配**(可用正则表达式).
* `^~`: 表示普通字符匹配. 使用前缀匹配. 如果匹配成功, 则不再匹配其他`location`.
* `@`: 它定义一个命名的 location，使用在内部定向时，例如 error_page, try_files.
* `/`: 通用匹配, 如果没有其它匹配,任何请求都会匹配到

```shell
(location =) > (location 完整路径) > (location ^~ 路径) > (location ~,~* 正则顺序) > (location 部分起始路径) > (/)
```

### 匹配优先级

* `=`精确匹配会第一个被处理. 如果发现精确匹配, nginx停止搜索其他匹配.
* `^~`类型表达式. 一旦匹配成功，则不再查找其他匹配项.
* 正则表达式类型(~ ~*)的优先级次之. 如果有多个location的正则能匹配的话, 则使用正则表达式最长的那个.
* 常规字符串匹配类型. 按前缀匹配.

```shell
location  = / {
  # 精确匹配 / ，主机名后面不能带任何字符串
  [ configuration A ] 
}
location  / {
  # 因为所有的地址都以 / 开头，所以这条规则将匹配到所有请求
  # 但是正则和最长字符串会优先匹配
  [ configuration B ] 
}
location /documents/ {
  # 匹配任何以 /documents/ 开头的地址，匹配符合以后，还要继续往下搜索
  # 只有后面的正则表达式没有匹配到时，这一条才会采用这一条
  [ configuration C ] 
}
location ~ /documents/Abc {
  # 匹配任何以 /documents/ 开头的地址，匹配符合以后，还要继续往下搜索
  # 只有后面的正则表达式没有匹配到时，这一条才会采用这一条
  [ configuration CC ] 
}
location ^~ /images/ {
  # 匹配任何以 /images/ 开头的地址，匹配符合以后，停止往下搜索正则，采用这一条。
  [ configuration D ] 
}
location ~* \.(gif|jpg|jpeg)$ {
  # 匹配所有以 gif,jpg或jpeg 结尾的请求
  # 然而，所有请求 /images/ 下的图片会被 config D 处理，因为 ^~ 到达不了这一条正则
  [ configuration E ] 
}
location /images/ {
  # 字符匹配到 /images/，继续往下，会发现 ^~ 存在
  [ configuration F ] 
}
location /images/abc {
  # 最长字符匹配到 /images/abc，继续往下，会发现 ^~ 存在
  # F与G的放置顺序是没有关系的
  [ configuration G ] 
}
location ~ /images/abc/ {
  # 只有去掉 config D 才有效：先最长匹配 config G 开头的地址，继续往下搜索，匹配到这一条正则，采用
    [ configuration H ] 
}
location ~* /js/.*/\.js
```


* / -> config A
精确完全匹配，即使/index.html也匹配不了
* /downloads/download.html -> config B
匹配B以后，往下没有任何匹配，采用B
* /images/1.gif -> configuration D
匹配到F，往下匹配到D，停止往下
* /images/abc/def -> config D
最长匹配到G，往下匹配D，停止往下
你可以看到 任何以/images/开头的都会匹配到D并停止，FG写在这里是没有任何意义的，H是永远轮不到的，这里只是为了说明匹配顺序
* /documents/document.html -> config C
匹配到C，往下没有任何匹配，采用C
* /documents/1.jpg -> configuration E
匹配到C，往下正则匹配到E
* /documents/Abc.jpg -> config CC
最长匹配到C，往下正则顺序匹配到CC，不会往下到E