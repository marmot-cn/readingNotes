#Memcached协议

----

[memcached官方协议](http://github.com/memcached/memcached/blob/master/doc/protocol.txt)

##翻译注解


###协议

---

**一些名字注解**

* 文本行  text lines
* 非结构化数据 unstructured data
* 客户端 client
* 服务端 server
* 数据块  data block

Clients of memcached communicate with server through TCP connections.
(A UDP interface is also available; details are below under "UDP
protocol.") A given running memcached server listens on some
(configurable) port; clients connect to that port, send commands to
the server, read responses, and eventually close the connection.

memcached客户端可以通过TCP和UDP连接.memcached服务端监听了一些可以配置的端口,客户端可以连接这些端口,发送命令到服务端,读取反馈,并且最终关闭连接.

There is no need to send any command to end the session. A client may
just close the connection at any moment it no longer needs it. Note,
however, that clients are encouraged to cache their connections rather
than reopen them every time they need to store or retrieve data.  This
is because memcached is especially designed to work very efficiently
with a very large number (many hundreds, more than a thousand if
necessary) of open connections. Caching connections will eliminate the
overhead associated with establishing a TCP connection (the overhead
of preparing for a new connection on the server side is insignificant
compared to this).

不需要发送任何命令去关闭这个会话(session).客户端可以任意关闭连接在任何时刻它不在需要的时候.注意,客户端最好每次存储和读取数据时都缓存他们的连接而不是每次都重新建立连接.这是因为memcached是别设计于非常高效的工作在打开大量的连接场景中.缓存连接可以节省与服务器建立TCP连接的时间开销.(服务器端新建连接的开销可以忽略不计).

There are two kinds of data sent in the memcache protocol: text lines
and unstructured data.  Text lines are used for commands from clients
and responses from servers. Unstructured data is sent when a client
wants to store or retrieve data. The server will transmit back
unstructured data in exactly the same way it received it, as a byte
stream. The server doesn't care about byte order issues in
unstructured data and isn't aware of them. There are no limitations on
characters that may appear in unstructured data; however, the reader
of such data (either a client or a server) will always know, from a
preceding text line, the exact length of the data block being
transmitted.

有2种类型的数据在memcache协议中:文本行和非结构化数据.文本行协议用于`客户端命令`和`服务器端响应`.非结构化数据用于保存和读取数据.服务端会以字符流(byte stream)形式传输回和收到时一样的非结构化数据.服务端既不关心非结构化数据中字节的顺序也不关心字节.非结构化数据的字符没有任何限制,但是,这些数据的用户(客户端和服务端)将会从之前的文本行(命令)知道这些数据传输数据库确切的长度。

Text lines are always terminated by \r\n. Unstructured data is _also_
terminated by \r\n, even though \r, \n or any other 8-bit characters
may also appear inside the data. Therefore, when a client retrieves
data from a server, it must use the length of the data block (which it
will be provided with) to determine where the data block ends, and not
the fact that \r\n follows the end of the data block, even though it
does.

文本行传时总是以 \r\n结尾.非结构化数据也是如此(以 \r\n结尾),尽管\r\n或者其他8位字节可能也会出现在数据内.所以,当客户端从服务端收到数据时,必须使用数据块的长度(会被同时提供)去确定数据块结尾, 并不只是根据结尾处的\r\n, 虽然实际数据流确实用于结尾.
(*主要是根据长度来判断,因为\r\n可能出现在数据中*).


###键(keys)

---

Data stored by memcached is identified with the help of a key. A key
is a text string which should uniquely identify the data for clients
that are interested in storing and retrieving it.  Currently the
length limit of a key is set at 250 characters (of course, normally
clients wouldn't need to use such long keys); the key must not include
control characters or whitespace.

key用于定位存储在memcached中的数据.一个key是字符串,并且是可以唯一定位存储或获取数据的.目前key的长度被限定为250字符(characters).key不能包含`控制字符`(在ASCII码中,第0～31号: LF（换行）、CR（回车）、FF（换页）、DEL（删除）、BS（退格)、BEL（振铃）)和`空格`.

###命令(commands)

---

有三种命令

**存储命令**

Storage commands (there are six: "set", "add", "replace", "append"
"prepend" and "cas") ask the server to store some data identified by a
key. The client sends a command line, and then a data block; after
that the client expects one line of response, which will indicate
success or failure.

`set`,`add`,`replace`,`append`,`prepend`,`cas` 用于服务端通过key存储value.
客户端发送一个命令,和一个数据块.然后客户端期待一行标记成功或者失败的反馈.


**读取命令**

Retrieval commands (there are two: "get" and "gets") ask the server to
retrieve data corresponding to a set of keys (one or more keys in one
request). The client sends a command line, which includes all the
requested keys; after that for each item the server finds it sends to
the client one response line with information about the item, and one
data block with the item's data; this continues until the server
finished with the "END" response line.

`get`和`gets`.用于服务端去获取数据通过一个或一组key.客户端发送一条包含keys的命令.然后服务端根据找到的每一项内容,都会给客户端返回一条关于内容的响应信息,和一个对应数据(value)的数据库.直到服务器以一个"END"作为响应信息结束.

**其他命令**

All other commands don't involve unstructured data. In all of them,
the client sends one command line, and expects (depending on the
command) either one line of response, or several lines of response
ending with "END" on the last line.

其他命令(如`flush_all`)不包含非结构化数据.客户端发送一条命令,并且期望等到一条或多条响应信息(根据不同的命令),或者多响应信息以"END"作为最后一行.

A command line always starts with the name of the command, followed by
parameters (if any) delimited by whitespace. Command names are
lower-case and are case-sensitive.

一条命令行总是以`命令`作为开始,紧跟着以空格分割的参数.命令是`小写`并且`大小写敏感`.

###过期时间(Expiration times)

---

Some commands involve a client sending some kind of expiration time
(relative to an item or to an operation requested by the client) to
the server. In all such cases, the actual value sent may either be
Unix time (number of seconds since January 1, 1970, as a 32-bit
value), or a number of seconds starting from current time. In the
latter case, this number of seconds may not exceed `60*60*24*30` (number
of seconds in 30 days); if the number sent by a client is larger than
that, the server will consider it to be real Unix time value rather
than an offset from current time.

一些命令包含了一个客户发送的各种过期时间(针对一个内容或者客户端的请求操作)到服务端.在所有这些例子中,真正被发送的值是以下2种格式

* Unix时间戳(从 1970.1.1开始起的秒数至失效时间的整型秒数)
* 一个从现在开始计算的秒数.对于后者,秒数不能超过`60*60*24*30`(30天的秒数).

如果客户端传递的值超过这个数字,服务器端会把这个值认为是unix时间戳,而不是当前时间的偏移量.

###错误信息(Error strings)

---

Each command sent by a client may be answered with an error string
from the server. These error strings come in three types:

- "ERROR\r\n"

  means the client sent a nonexistent command name.

- "CLIENT_ERROR <error>\r\n"

  means some sort of client error in the input line, i.e. the input
  doesn't conform to the protocol in some way. <error> is a
  human-readable error string.

- "SERVER_ERROR <error>\r\n"

  means some sort of server error prevents the server from carrying
  out the command. <error> is a human-readable error string. In cases
  of severe server errors, which make it impossible to continue
  serving the client (this shouldn't normally happen), the server will
  close the connection after sending the error line. This is the only
  case in which the server closes a connection to a client.


In the descriptions of individual commands below, these error lines
are not again specifically mentioned, but clients must allow for their possibility.

每个从客户端发送的命令都可能被反馈一个错误信息从服务端.有以下三种类型:

* `"ERROR\r\n`

  表示客户端发送了一个`不存在`的命令
  
* `"CLIENT_ERROR <error>\r\n"`
  
  表示客户端输入的命令行是存在某种错误, 比如: 输入的不符合协议. <error>是一种人类的可读的错误字符串.
  
* `"SERVER_ERROR <error>\r\n"`

  表示着服务器执行命令时候发生了一些错误.<error>是一种人类的可读的错误字符串.一些导致服务器不能正常服务客户端的错误(一般不会发生),服务端将会在发送错误信息后关闭连接.
  
在描述后续各个命令时, 这些错误信息不在特别提出, 但是要清楚这些错误是存在的.

###存储命令(Storage commands)

---

First, the client sends a command line which looks like this:

<command name> <key> <flags> <exptime> <bytes> [noreply]\r\n
cas <key> <flags> <exptime> <bytes> <cas unique> [noreply]\r\n

- <command name> is "set", "add", "replace", "append" or "prepend"

  "set" means "store this data".

  "add" means "store this data, but only if the server *doesn't* already
  hold data for this key".

  "replace" means "store this data, but only if the server *does*
  already hold data for this key".

  "append" means "add this data to an existing key after existing data".

  "prepend" means "add this data to an existing key before existing data".

  The append and prepend commands do not accept flags or exptime.
  They update existing data portions, and ignore new flag and exptime
  settings.

  "cas" is a check and set operation which means "store this data but
  only if no one else has updated since I last fetched it."

- <key> is the key under which the client asks to store the data

- <flags> is an arbitrary 16-bit unsigned integer (written out in
  decimal) that the server stores along with the data and sends back
  when the item is retrieved. Clients may use this as a bit field to
  store data-specific information; this field is opaque to the server.
  Note that in memcached 1.2.1 and higher, flags may be 32-bits, instead
  of 16, but you might want to restrict yourself to 16 bits for
  compatibility with older versions.

- <exptime> is expiration time. If it's 0, the item never expires
  (although it may be deleted from the cache to make place for other
  items). If it's non-zero (either Unix time or offset in seconds from
  current time), it is guaranteed that clients will not be able to
  retrieve this item after the expiration time arrives (measured by
  server time).

- <bytes> is the number of bytes in the data block to follow, *not*
  including the delimiting \r\n. <bytes> may be zero (in which case
  it's followed by an empty data block).

- <cas unique> is a unique 64-bit value of an existing entry.
  Clients should use the value returned from the "gets" command
  when issuing "cas" updates.

- "noreply" optional parameter instructs the server to not send the
  reply.  NOTE: if the request line is malformed, the server can't
  parse "noreply" option reliably.  In this case it may send the error
  to the client, and not reading it on the client side will break
  things.  Client should construct only valid requests.

After this line, the client sends the data block:

<data block>\r\n

- <data block> is a chunk of arbitrary 8-bit data of length <bytes>
  from the previous line.

After sending the command line and the data block the client awaits
the reply, which may be:

- "STORED\r\n", to indicate success.

- "NOT_STORED\r\n" to indicate the data was not stored, but not
because of an error. This normally means that the
condition for an "add" or a "replace" command wasn't met.

- "EXISTS\r\n" to indicate that the item you are trying to store with
a "cas" command has been modified since you last fetched it.

- "NOT_FOUND\r\n" to indicate that the item you are trying to store
with a "cas" command did not exist.



首先,客户端发送一个类似这样的命令:

		<command name> <key> <flags> <exptime> <bytes> [noreply]\r\n
		cas <key> <flags> <exptime> <bytes> <cas unique> [noreply]\r\n


- `<command name>`是`"set"`, `"add"`, `"replace"`, `"append"` or `"prepend"`

  "set" 表示 "存储数据".

  "add" 表示 "存储数据, 但是仅在服务`没有`该键值的时候"(注意是`没有`).

  "replace" 表示 "存储数据, 但是仅在服务`有`该键值的时候"(注意是`有`).

  "append" 表示 "在一个已经存在的数据`后`添加(追加)数据".

  "prepend" 表示 "在一个已经存在的数据`前`添加(追加)数据".

  `append` 和 `prepend` 命令不接受 `flags` 和 `exptime`.他们只是更新数据部分,并且忽略新的`flag`和`exptime`设置.
  
  "cas" 是一个`检查(check)`并`设置(set)`操作意味着只有在最近一次取出后并且没有其他人在更新过该数据时才能`存储(stroe)`数据

- `<key>` 客户端存储数据时用的键值

- `<flags>` is an arbitrary 16-bit unsigned integer (written out in
  decimal) that the server stores along with the data and sends back
  when the item is retrieved. Clients may use this as a bit field to
  store data-specific information; this field is opaque to the server.
  Note that in memcached 1.2.1 and higher, flags may be 32-bits, instead
  of 16, but you might want to restrict yourself to 16 bits for
  compatibility with older versions.
  
###读取命令(Retrieval command)

---

The retrieval commands "get" and "gets" operates like this:

get <key>*\r\n
gets <key>*\r\n

- <key>* means one or more key strings separated by whitespace.

After this command, the client expects zero or more items, each of
which is received as a text line followed by a data block. After all
the items have been transmitted, the server sends the string

"END\r\n"

to indicate the end of response.

Each item sent by the server looks like this:

VALUE <key> <flags> <bytes> [<cas unique>]\r\n
<data block>\r\n

- <key> is the key for the item being sent

- <flags> is the flags value set by the storage command

- <bytes> is the length of the data block to follow, *not* including
  its delimiting \r\n

- <cas unique> is a unique 64-bit integer that uniquely identifies
  this specific item.

- <data block> is the data for this item.

If some of the keys appearing in a retrieval request are not sent back
by the server in the item list this means that the server does not
hold items with such keys (because they were never stored, or stored
but deleted to make space for more items, or expired, or explicitly
deleted by a client).

###删除(Deletion)

---

The command "delete" allows for explicit deletion of items:

delete <key> [noreply]\r\n

- <key> is the key of the item the client wishes the server to delete

- "noreply" optional parameter instructs the server to not send the
  reply.  See the note in Storage commands regarding malformed
  requests.

The response line to this command can be one of:

- "DELETED\r\n" to indicate success

- "NOT_FOUND\r\n" to indicate that the item with this key was not
  found.

See the "flush_all" command below for immediate invalidation
of all existing items.


###增加/减少(Increment/Decrement)

---

Commands "incr" and "decr" are used to change data for some item
in-place, incrementing or decrementing it. The data for the item is
treated as decimal representation of a 64-bit unsigned integer.  If
the current data value does not conform to such a representation, the
incr/decr commands return an error (memcached <= 1.2.6 treated the
bogus value as if it were 0, leading to confusion). Also, the item
must already exist for incr/decr to work; these commands won't pretend
that a non-existent key exists with value 0; instead, they will fail.

The client sends the command line:

incr <key> <value> [noreply]\r\n

or

decr <key> <value> [noreply]\r\n

- <key> is the key of the item the client wishes to change

- <value> is the amount by which the client wants to increase/decrease
the item. It is a decimal representation of a 64-bit unsigned integer.

- "noreply" optional parameter instructs the server to not send the
  reply.  See the note in Storage commands regarding malformed
  requests.

The response will be one of:

- "NOT_FOUND\r\n" to indicate the item with this value was not found

- <value>\r\n , where <value> is the new value of the item's data,
  after the increment/decrement operation was carried out.

Note that underflow in the "decr" command is caught: if a client tries
to decrease the value below 0, the new value will be 0.  Overflow in
the "incr" command will wrap around the 64 bit mark.

Note also that decrementing a number such that it loses length isn't
guaranteed to decrement its returned length.  The number MAY be
space-padded at the end, but this is purely an implementation
optimization, so you also shouldn't rely on that.

###Touch

---

The "touch" command is used to update the expiration time of an existing item
without fetching it.

touch <key> <exptime> [noreply]\r\n

- <key> is the key of the item the client wishes the server to delete

- <exptime> is expiration time. Works the same as with the update commands
  (set/add/etc). This replaces the existing expiration time. If an existing
  item were to expire in 10 seconds, but then was touched with an
  expiration time of "20", the item would then expire in 20 seconds.

- "noreply" optional parameter instructs the server to not send the
  reply.  See the note in Storage commands regarding malformed
  requests.

The response line to this command can be one of:

- "TOUCHED\r\n" to indicate success

- "NOT_FOUND\r\n" to indicate that the item with this key was not
  found.
  
  
###Slabs再分配(Slabs Reassign)

----

NOTE: This command is subject to change as of this writing.

The slabs reassign command is used to redistribute memory once a running
instance has hit its limit. It might be desirable to have memory laid out
differently than was automatically assigned after the server started.

slabs reassign <source class> <dest class>\r\n

- <source class> is an id number for the slab class to steal a page from

A source class id of -1 means "pick from any valid class"

- <dest class> is an id number for the slab class to move a page to

The response line could be one of:

- "OK" to indicate the page has been scheduled to move

- "BUSY [message]" to indicate a page is already being processed, try again
  later.

- "BADCLASS [message]" a bad class id was specified

- "NOSPARE [message]" source class has no spare pages

- "NOTFULL [message]" dest class must be full to move new pages to it

- "UNSAFE [message]" source class cannot move a page right now

- "SAME [message]" must specify different source/dest ids.

###Slabs Automove

---

NOTE: This command is subject to change as of this writing.

The slabs automove command enables a background thread which decides on its
own when to move memory between slab classes. Its implementation and options
will likely be in flux for several versions. See the wiki/mailing list for
more details.

The automover can be enabled or disabled at runtime with this command.

slabs automove <0|1>

- 0|1|2 is the indicator on whether to enable the slabs automover or not.

The response should always be "OK\r\n"

- <0> means to set the thread on standby

- <1> means to run the builtin slow algorithm to choose pages to move

- <2> is a highly aggressive mode which causes pages to be moved every time
  there is an eviction. It is not recommended to run for very long in this
  mode unless your access patterns are very well understood.
  
 
###LRU_Crawler

---

NOTE: This command (and related commands) are subject to change as of this
writing.

The LRU Crawler is an optional background thread which will walk from the tail
toward the head of requested slab classes, actively freeing memory for expired
items. This is useful if you have a mix of items with both long and short
TTL's, but aren't accessed very often. This system is not required for normal
usage, and can add small amounts of latency and increase CPU usage.

lru_crawler <enable|disable>

- Enable or disable the LRU Crawler background thread.

The response line could be one of:

- "OK" to indicate the crawler has been started or stopped.

- "ERROR [message]" something went wrong while enabling or disabling.

lru_crawler sleep <microseconds>

- The number of microseconds to sleep in between each item checked for
  expiration. Smaller numbers will obviously impact the system more.
  A value of "0" disables the sleep, "1000000" (one second) is the max.

The response line could be one of:

- "OK"

- "CLIENT_ERROR [message]" indicating a format or bounds issue.

lru_crawler tocrawl <32u>

- The maximum number of items to inspect in a slab class per run request. This
  allows you to avoid scanning all of very large slabs when it is unlikely to
  find items to expire.

The response line could be one of:

- "OK"

- "CLIENT_ERROR [message]" indicating a format or bound issue.

lru_crawler crawl <classid,classid,classid|all>

- Takes a single, or a list of, numeric classids (ie: 1,3,10). This instructs
  the crawler to start at the tail of each of these classids and run to the
  head. The crawler cannot be stopped or restarted until it completes the
  previous request.

  The special keyword "all" instructs it to crawl all slabs with items in
  them.

The response line could be one of:

- "OK" to indicate successful launch.

- "BUSY [message]" to indicate the crawler is already processing a request.

- "BADCLASS [message]" to indicate an invalid class was specified.

###统计(Statistics)

---

The command "stats" is used to query the server about statistics it
maintains and other internal data. It has two forms. Without
arguments:

stats\r\n

it causes the server to output general-purpose statistics and
settings, documented below.  In the other form it has some arguments:

stats <args>\r\n

Depending on <args>, various internal data is sent by the server. The
kinds of arguments and the data sent are not documented in this version
of the protocol, and are subject to change for the convenience of
memcache developers.


