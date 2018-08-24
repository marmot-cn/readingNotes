# crontab禁发邮件

### 在第一行加上`MAILTO=""`

```
MAILTO=""
*/1 * * * * /shell/shell.sh >/dev/null 2>&1
```

### 在尾部追加`&> /dev/null`

把**标准输出**和**错误输出**都指向到`/dev/null`

```
*/1 * * * * /shell/shell.sh &> /dev/null1
```
