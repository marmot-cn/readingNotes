# git撤销提交

---

## 撤销,同时将代码恢复到前一`commit_id`对应的版本

1. `git log` 找到你想撤销的`commit_id`.
2. `git reset --hard commit_id`  完成撤销,同时将代码恢复到前一`commit_id`对应的版本.

## 撤销,但是不对代码修改进行撤销

1. `git log` 找到你想撤销的`commit_id`.
2. 完成Commit命令的撤销, 但是不对代码修改进行撤销, 可以直接通过`git commit`重新提交对本地代码的修改.