# sed

----

### 参数

* `-e<script>`或`--expression=<script>`: 以选项中的指定的script来处理输入的文本文件
* `-f<script文件>`或`--file=<script文件>`: 以选项中指定的script文件来处理输入的文本文件

### `s`和`y`

#### `s`

		[address1[ ,address2]] s/pattern/replacement/[flag]
		
flag : 主要用它来控制一些替换情况:

* 当 flag 为 g 时, 代表替换所有符合(match)的字串 。
* 当 flag 数字时, 代表替换行内第 m 个符合的字串。
* 当 flag 为 p 时, 代表替换第一个符合 pattern 的字串後 , 将该行再在原始行下面再输出一次。
* 当 flag 为 w wfile 时, 代表替换第一个符合 pattern 的字串後 , 输出到 wfile 档内(如果 wfile 不存在 , 则会 重新开启名为 wfile 的档案)

#### `y`

`y`表示转换资料中的字元.处理单个字符.

	   echo "This 1 is a test 2 try." | sed 'y/123/456/'
	   This 4 is a test 5 try.
	   
### `p`和`P`

#### `p`

`p`打印当前模式空间内容,追加到默认输出之后.

		cat NUM
		1
		2
		3
		4
		
		sed 'N;p' NUM
		1
		2
		1
		2
		3
		4
		3
		4

p 会把模式控件中的`1\n`

#### `P`

P打印当前模式空间`开端至\n`的内容,并`追加`到`默认输出之前`.

		cat NUM
		1
		2
		3
		4
		
		sed 'N;P' NUM
		1
		1
		2
		3
		3
		4
		
首先sed默认的读取1,模式空间为1,执行N,模式空间变成1\n2\n,然后执行P,也就是打印1\n.当前行的处理，打印模式空间(pattern space)也就是1\n2\n.这时sed再从文件中读取下一行,也就到了3\n,执行N;模式空间变成了3\n4\n;
执行P;打印3\n;继续执行当前行的处理,打印模式空间3\n4\n;sed再从文件中读取下一行,发现没有了,结束处理流程.

### `d`和`D`

#### `d`

d命令是删除当前模式空间内容(不在传至标准输出),并放弃之后的命令,并对新读取的内容,重头执行sed.

		cat aaa   
		This is 1   
		This is 2   
		This is 3   
		This is 4   
		This is 5   
		                                                           
		sed 'n;d' aaa           
		This is 1   
		This is 3   
		This is 5

		注释：读取1，执行n，得出2，执行d，删除2，得空，以此类推，读取3，执行n，得出4，执行d，删除4，得空，但是读取5时，因为n无法执行，所以d不执行。因无-n参数，故输出1\n3\n5
		
#### `D`

D命令是删除当前模式空间开端至\n的内容(不在传至标准输出),放弃之后的命令,但是对剩余模式空间重新执行sed.

		cat aaa   
		This is 1   
		This is 2   
		This is 3   
		This is 4   
		This is 5   
		                                                
		sed 'N;D' aaa           
		This is 5
		
		读取1,执行N,得出1\n2,执行D,得出2,执行N,得出2\n3,执行D,得出3,依此类推,得出5,执行N,条件失败退出,因无-n参数,故输出5.

