# supervisor

---

`supervisor`就是用`Python`开发的一套通用的进程管理程序, 能将一个普通的命令行进程变为后台`daemon`, 并监控进程状态, 异常退出时能自动重启.

可以是我们一些需要跑得守护进程使用`supervisor`来管理:

这是原先的使用方法, 新方法可以使用`supervisor`

```shell
nohup your-program &>nohup.log &
```

## 安装

### 安装Python包管理工具

```shell
yum install python-setuptools -y
```

### 安装`supervisor`

```shell
easy_install supervisor
```

### 初始化配置文件

```shell
mkdir /etc/supervisor
echo_supervisord_conf > /etc/supervisor/supervisord.conf

[root@localhost ~]# cat /etc/supervisor/supervisord.conf
; Sample supervisor config file.
;
; For more information on the config file, please see:
; http://supervisord.org/configuration.html
;
; Notes:
;  - Shell expansion ("~" or "$HOME") is not supported.  Environment
;    variables can be expanded using this syntax: "%(ENV_HOME)s".
;  - Quotes around values are not supported, except in the case of
;    the environment= options as shown below.
;  - Comments must have a leading space: "a=b ;comment" not "a=b;comment".
;  - Command will be truncated if it looks like a config file comment, e.g.
;    "command=bash -c 'foo ; bar'" will truncate to "command=bash -c 'foo ".

[unix_http_server]
file=/tmp/supervisor.sock   ; the path to the socket file
;chmod=0700                 ; socket file mode (default 0700)
;chown=nobody:nogroup       ; socket file uid:gid owner
;username=user              ; default is no username (open server)
;password=123               ; default is no password (open server)

;[inet_http_server]         ; inet (TCP) server disabled by default
;port=127.0.0.1:9001        ; ip_address:port specifier, *:port for all iface
;username=user              ; default is no username (open server)
;password=123               ; default is no password (open server)

[supervisord]
logfile=/tmp/supervisord.log ; main log file; default $CWD/supervisord.log
logfile_maxbytes=50MB        ; max main logfile bytes b4 rotation; default 50MB
logfile_backups=10           ; # of main logfile backups; 0 means none, default 10
loglevel=info                ; log level; default info; others: debug,warn,trace
pidfile=/tmp/supervisord.pid ; supervisord pidfile; default supervisord.pid
nodaemon=false               ; start in foreground if true; default false
minfds=1024                  ; min. avail startup file descriptors; default 1024
minprocs=200                 ; min. avail process descriptors;default 200
;umask=022                   ; process file creation umask; default 022
;user=chrism                 ; default is current user, required if root
;identifier=supervisor       ; supervisord identifier, default is 'supervisor'
;directory=/tmp              ; default is not to cd during start
;nocleanup=true              ; don't clean up tempfiles at start; default false
;childlogdir=/tmp            ; 'AUTO' child log dir, default $TEMP
;environment=KEY="value"     ; key value pairs to add to environment
;strip_ansi=false            ; strip ansi escape codes in logs; def. false

; The rpcinterface:supervisor section must remain in the config file for
; RPC (supervisorctl/web interface) to work.  Additional interfaces may be
; added by defining them in separate [rpcinterface:x] sections.

[rpcinterface:supervisor]
supervisor.rpcinterface_factory = supervisor.rpcinterface:make_main_rpcinterface

; The supervisorctl section configures how supervisorctl will connect to
; supervisord.  configure it match the settings in either the unix_http_server
; or inet_http_server section.

[supervisorctl]
serverurl=unix:///tmp/supervisor.sock ; use a unix:// URL  for a unix socket
;serverurl=http://127.0.0.1:9001 ; use an http:// url to specify an inet socket
;username=chris              ; should be same as in [*_http_server] if set
;password=123                ; should be same as in [*_http_server] if set
;prompt=mysupervisor         ; cmd line prompt (default "supervisor")
;history_file=~/.sc_history  ; use readline history if available

; The sample program section below shows all possible program subsection values.
; Create one or more 'real' program: sections to be able to control them under
; supervisor.

;[program:theprogramname]
;command=/bin/cat              ; the program (relative uses PATH, can take args)
;process_name=%(program_name)s ; process_name expr (default %(program_name)s)
;numprocs=1                    ; number of processes copies to start (def 1)
;directory=/tmp                ; directory to cwd to before exec (def no cwd)
;umask=022                     ; umask for process (default None)
;priority=999                  ; the relative start priority (default 999)
;autostart=true                ; start at supervisord start (default: true)
;startsecs=1                   ; # of secs prog must stay up to be running (def. 1)
;startretries=3                ; max # of serial start failures when starting (default 3)
;autorestart=unexpected        ; when to restart if exited after running (def: unexpected)
;exitcodes=0,2                 ; 'expected' exit codes used with autorestart (default 0,2)
;stopsignal=QUIT               ; signal used to kill process (default TERM)
;stopwaitsecs=10               ; max num secs to wait b4 SIGKILL (default 10)
;stopasgroup=false             ; send stop signal to the UNIX process group (default false)
;killasgroup=false             ; SIGKILL the UNIX process group (def false)
;user=chrism                   ; setuid to this UNIX account to run the program
;redirect_stderr=true          ; redirect proc stderr to stdout (default false)
;stdout_logfile=/a/path        ; stdout log path, NONE for none; default AUTO
;stdout_logfile_maxbytes=1MB   ; max # logfile bytes b4 rotation (default 50MB)
;stdout_logfile_backups=10     ; # of stdout logfile backups (0 means none, default 10)
;stdout_capture_maxbytes=1MB   ; number of bytes in 'capturemode' (default 0)
;stdout_events_enabled=false   ; emit events on stdout writes (default false)
;stderr_logfile=/a/path        ; stderr log path, NONE for none; default AUTO
;stderr_logfile_maxbytes=1MB   ; max # logfile bytes b4 rotation (default 50MB)
;stderr_logfile_backups=10     ; # of stderr logfile backups (0 means none, default 10)
;stderr_capture_maxbytes=1MB   ; number of bytes in 'capturemode' (default 0)
;stderr_events_enabled=false   ; emit events on stderr writes (default false)
;environment=A="1",B="2"       ; process environment additions (def no adds)
;serverurl=AUTO                ; override serverurl computation (childutils)

; The sample eventlistener section below shows all possible eventlistener
; subsection values.  Create one or more 'real' eventlistener: sections to be
; able to handle event notifications sent by supervisord.

;[eventlistener:theeventlistenername]
;command=/bin/eventlistener    ; the program (relative uses PATH, can take args)
;process_name=%(program_name)s ; process_name expr (default %(program_name)s)
;numprocs=1                    ; number of processes copies to start (def 1)
;events=EVENT                  ; event notif. types to subscribe to (req'd)
;buffer_size=10                ; event buffer queue size (default 10)
;directory=/tmp                ; directory to cwd to before exec (def no cwd)
;umask=022                     ; umask for process (default None)
;priority=-1                   ; the relative start priority (default -1)
;autostart=true                ; start at supervisord start (default: true)
;startsecs=1                   ; # of secs prog must stay up to be running (def. 1)
;startretries=3                ; max # of serial start failures when starting (default 3)
;autorestart=unexpected        ; autorestart if exited after running (def: unexpected)
;exitcodes=0,2                 ; 'expected' exit codes used with autorestart (default 0,2)
;stopsignal=QUIT               ; signal used to kill process (default TERM)
;stopwaitsecs=10               ; max num secs to wait b4 SIGKILL (default 10)
;stopasgroup=false             ; send stop signal to the UNIX process group (default false)
;killasgroup=false             ; SIGKILL the UNIX process group (def false)
;user=chrism                   ; setuid to this UNIX account to run the program
;redirect_stderr=false         ; redirect_stderr=true is not allowed for eventlisteners
;stdout_logfile=/a/path        ; stdout log path, NONE for none; default AUTO
;stdout_logfile_maxbytes=1MB   ; max # logfile bytes b4 rotation (default 50MB)
;stdout_logfile_backups=10     ; # of stdout logfile backups (0 means none, default 10)
;stdout_events_enabled=false   ; emit events on stdout writes (default false)
;stderr_logfile=/a/path        ; stderr log path, NONE for none; default AUTO
;stderr_logfile_maxbytes=1MB   ; max # logfile bytes b4 rotation (default 50MB)
;stderr_logfile_backups=10     ; # of stderr logfile backups (0 means none, default 10)
;stderr_events_enabled=false   ; emit events on stderr writes (default false)
;environment=A="1",B="2"       ; process environment additions
;serverurl=AUTO                ; override serverurl computation (childutils)

; The sample group section below shows all possible group values.  Create one
; or more 'real' group: sections to create "heterogeneous" process groups.

;[group:thegroupname]
;programs=progname1,progname2  ; each refers to 'x' in [program:x] definitions
;priority=999                  ; the relative start priority (default 999)

; The [include] section can just contain the "files" setting.  This
; setting can list multiple files (separated by whitespace or
; newlines).  It can also contain wildcards.  The filenames are
; interpreted as relative to this file.  Included files *cannot*
; include files themselves.

;[include]
;files = relative/directory/*.ini
```

我们修改配置, 让每个需要管理的程序使用独立的配置文件.修改最后一行`[include]`即可. 改为:

```shell
[include]
files = /etc/supervisor/conf.d/*.ini
```

## 启动

```shell
[root@localhost ~]# supervisord -c /etc/supervisor/supervisord.conf
```

关闭:

```shell
supervisorctl shutdown
```

## 测试

编写`helloworld.sh`脚本

```shell
[root@localhost helloworld]# pwd
/root/helloworld
[root@localhost helloworld]# cat helloworld.sh
#!/bin/bash
while true
do
    echo `date`
    sleep 1
done
```

编写`supervisor`配置文件

* `stopasgroup`: 如果设置为 true，当进程收到 stop 信号时，会自动将该信号发给该进程的子进程
* `stopsignal`: 用来杀死进程的信号
* `redirect_stderr=true`: 重定向`stderr`到`stdout`
* `autostart=true`: 随着`supervisord`的启动而启动
* `autorestart=true`: 自动重启
* `startretries=10`: 启动失败时的最多重试次数
* `numprocs`: 启动几个进程, 如果`>=2`时`process_name`必须包含`%(process_num)`

```shell
/etc/supervisor/conf.d/helloworld.ini

[program:helloworld]
directory=/root/helloworld
command= bash helloworld.sh
process_name=%(program_name)s_%(process_num)02d
numprocs=4
autostart=true
autorestart=true
startretries=10
redirect_stderr=true
stopsignal=TERM
stopasgroup=true
```

```shell
如果更新了配置文件需要运行 update
[root@localhost ~]# supervisorctl update
helloworld: stopped
helloworld: updated process group

[root@localhost ~]# supervisorctl status
helloworld:helloworld_00         RUNNING   pid 23781, uptime 0:00:04
helloworld:helloworld_01         RUNNING   pid 23782, uptime 0:00:04
helloworld:helloworld_02         RUNNING   pid 23779, uptime 0:00:04
helloworld:helloworld_03         RUNNING   pid 23780, uptime 0:00:04

可见我们4个进程都在运行

我们测试关闭其中一个, 会自动生成一个
[root@localhost ~]# kill 23781
[root@localhost ~]# supervisorctl status
helloworld:helloworld_00         RUNNING   pid 24132, uptime 0:00:01
helloworld:helloworld_01         RUNNING   pid 23782, uptime 0:01:28
helloworld:helloworld_02         RUNNING   pid 23779, uptime 0:01:28
helloworld:helloworld_03         RUNNING   pid 23780, uptime 0:01:28
```

## 常见命令 

* `supervisorctl stop programxxx`: 停止某一个进程(programxxx)，programxxx为[program:chatdemon]里配置的值，这个示例就是chatdemon.
* `supervisorctl start programxxx`: 启动某个进程.
* `supervisorctl restart programxxx`: 重启某个进程
* `supervisorctl stop groupworker`: 重启所有属于名为groupworker这个分组的进程(start,restart同理)
* `supervisorctl stop all`: 停止全部进程，注：start、restart、stop都不会载入最新的配置文件
* `supervisorctl reload`: 载入最新的配置文件，停止原有进程并按新的配置启动、管理所有进程
* `supervisorctl update`: 根据最新的配置文件，启动新配置或有改动的进程，配置没有改动的进程不会受影响而重启。注意：显示用stop停止掉的进程，用reload或者update都不会自动重启

### 查看日志

```shell
[root@localhost ~]# supervisorctl status
helloworld:helloworld_00         RUNNING   pid 25011, uptime 0:00:05
helloworld:helloworld_01         RUNNING   pid 25012, uptime 0:00:05
helloworld:helloworld_02         RUNNING   pid 25009, uptime 0:00:05
helloworld:helloworld_03         RUNNING   pid 25010, uptime 0:00:05
[root@localhost ~]# supervisorctl tail helloworld:helloworld_00
Thu Oct 5 16:21:46 CST 2017
Thu Oct 5 16:21:47 CST 2017
Thu Oct 5 16:21:48 CST 2017
Thu Oct 5 16:21:49 CST 2017
Thu Oct 5 16:21:50 CST 2017
Thu Oct 5 16:21:51 CST 2017
Thu Oct 5 16:21:52 CST 2017
Thu Oct 5 16:21:53 CST 2017
Thu Oct 5 16:21:54 CST 2017
Thu Oct 5 16:21:55 CST 2017
Thu Oct 5 16:21:56 CST 2017
```
