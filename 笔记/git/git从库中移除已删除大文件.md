# git从库中移除已删除大文件

---

我在开发中把`*.xmind`格式文件放到`git`目录下, 因为其本身是二进制文件加上版本化导致`.git`文件夹过于庞大. 所以需要删除.

这里的流程从网上借鉴, 自己记录下.

## 流程

### 运行`count-objects`查看空间使用

`size-pack`是以千字节为单位表示的`packfiles`的大小.

```
git count-objects -v
```

### 运行底层命令`git verify-pack`以识别出大对象, 对输出的第三列信息即文件大小进行排序

假设我们已经知道大文件, 我的场景是知道格式为`*.xmind`文件, 则不需要运行此步骤. 我这里是已经清理过的, 所以呈现的文件只是一个示例.

```
xxx.idx 每个都不同, 选择自己的即可

git verify-pack -v .git/objects/pack/xxxx.idx | sort -k 3 -n | tail -3

e5dc6e3ae9bee288145185de306e32faf882b7d0 blob   32768 3907 880767
68d97738dba6f5e260be2de1665ab78134f452ac blob   70861 3522 286027
f32802056b80e4c31d5a7e2dfb7a4c9912006bd2 blob   130054 15215 307054
```

### 使用`rev-list`命令, 传入`--objects`选项, 它会列出所有`commit SHA`值, `blob SHA`值及相应的文件路径, 这样查看`blob`的文件名

```
git rev-list --objects --all | grep f32802056b80e4c31d5a7e2dfb7a4c9912006bd2
f32802056b80e4c31d5a7e2dfb7a4c9912006bd2 composer.lock
```

### 删除对应的大文件

我这里的示例是`*.xmind`格式的文件.

```
git filter-branch --force --index-filter 'git rm --cached --ignore-unmatch *.xmind' --prune-empty --tag-name-filter cat -- --all
```

### 立刻回收空间

```
rm -rf .git/refs/original/ 
git reflog expire --expire=now --all
git gc --prune=now
git gc --aggressive --prune=now
```

### 最后把改动强制推送到远端

```
git push origin --force --all
```