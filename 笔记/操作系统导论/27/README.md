# README

测试`myret_t *rvals = malloc(sizeof(myret_t));`与`myarg_t args = { 10, 20 };`声明的区别

替换代码

### 原

```c
myret_t *rvals = malloc(sizeof(myret_t));
rvals->x = 1;
rvals->y = 2;
return (void *) rvals;
```

### 新

```c
myret_t *rvals = malloc(sizeof(myret_t));
rvals->x = 1;
rvals->y = 2;
return (void *) rvals;
```