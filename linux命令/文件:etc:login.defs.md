# 文件/etc/login.defs

---

创建用户时的一些规划.

```shell
[ansible@iZ94ebqp9jtZ ~]$ cat /etc/login.defs
#
# Please note that the parameters in this configuration file control the
# behavior of the tools from the shadow-utils component. None of these
# tools uses the PAM mechanism, and the utilities that use PAM (such as the
# passwd command) should therefore be configured elsewhere. Refer to
# /etc/pam.d/system-auth for more information.
#

# *REQUIRED*
#   Directory where mailboxes reside, _or_ name of file, relative to the
#   home directory.  If you _do_ define both, MAIL_DIR takes precedence.
#   QMAIL_DIR is for Qmail
#
#QMAIL_DIR	Maildir
MAIL_DIR	/var/spool/mail
#MAIL_FILE	.mail

# Password aging controls:
#
#	PASS_MAX_DAYS	Maximum number of days a password may be used.
#	PASS_MIN_DAYS	Minimum number of days allowed between password changes.
#	PASS_MIN_LEN	Minimum acceptable password length.
#	PASS_WARN_AGE	Number of days warning given before a password expires.
#
PASS_MAX_DAYS	99999
PASS_MIN_DAYS	0
PASS_MIN_LEN	5
PASS_WARN_AGE	7

#
# Min/max values for automatic uid selection in useradd
#
UID_MIN                  1000
UID_MAX                 60000
# System accounts
SYS_UID_MIN               201
SYS_UID_MAX               999

#
# Min/max values for automatic gid selection in groupadd
#
GID_MIN                  1000
GID_MAX                 60000
# System accounts
SYS_GID_MIN               201
SYS_GID_MAX               999

#
# If defined, this command is run when removing a user.
# It should remove any at/cron/print jobs etc. owned by
# the user to be removed (passed as the first argument).
#
#USERDEL_CMD	/usr/sbin/userdel_local

#
# If useradd should create home directories for users by default
# On RH systems, we do. This option is overridden with the -m flag on
# useradd command line.
#
CREATE_HOME	yes

# The permission mask is initialized to this value. If not specified,
# the permission mask will be initialized to 022.
UMASK           077

# This enables userdel to remove user groups if no members exist.
#
USERGROUPS_ENAB yes

# Use SHA512 to encrypt password.
ENCRYPT_METHOD SHA512
```

* `MAIL_DIR /var/spool/mail` 创建用户时创建相应的mail文件.
* `PASS_MAX_DAYS 99999`: 用户的密码不过期最多的天数.
* `PASS_MIN_DAYS 0`: 密码修改之间最小的天数.
* `PASS_MIN_LEN 5`: 密码最小长度.
* `PASS_WARN_AGE 7`: 密码过期之前7天开始提示.
* `UID_MIN 500`: 最小UID为500.
* `UID_MAX 60000`: 最大UID为60000.
* `GID_MIN 500`: GID 是从500开始.
* `GID_MAX 60000`: 最大GID为60000.
* `CREATE_HOME yes`: 是否创用户家目录,默认创建.
* `UMASK`: 通过umask决定新建用户HOME目录的权限(注意这里是家目录的权限).

## `UMASK`

```shell
默认 077 我们创建后用户家目录是 700 权限
[ansible@iZ94ebqp9jtZ ~]$ ll -d /home/ansible
drwx------ 4 ansible ansible 4096 Sep  5 11:40 /home/ansible

设定/etc/login.defs文件中UMASK为000

权限变为 777
[ansible@iZ94ebqp9jtZ ~]$ ll -d /home/wangcai3/
drwxrwxrwx 2 wangcai3 wangcai3 4096 Sep  5 12:01 /home/wangcai3/
```

