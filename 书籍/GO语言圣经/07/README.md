# README

## 笔记

### 7.1 接口约定

### 7.2 接口类型

### 7.3 实现接口的条件

一个类型如果拥有一个接口需要的所有方法，那么这个类型就实现了这个接口。

### 7.4 flag.Value 接口

### 7.5 接口值

### 7.6 sort.Interface 接口

### 7.7 http.Handler 接口

### 7.8 error接口

### 7.9 表达式求值

### 7.10 类型断言

如果是对象

```
# 类型判断
if f, ok := w.(*os.File); ok { 
	// ...use f... 
}
```

### 7.11 基于类型断言区别错误类型

```
import (
	"errors" 
	"syscall" 
)

var ErrNotExist = errors.New("file does not exist") 

// IsNotExist returns a boolean indicating whether the error is known to 
// report that a file or directory does not exist. It is satisfied by
 // ErrNotExist as well as some syscall errors. 
 func IsNotExist(err error) bool { 
 	if pe, ok := err.(*PathError); ok { 
 		err = pe.Err 
 	}
 	
 	return err == syscall.ENOENT || err == ErrNotExist 
 }
```

### 7.12. 通过类型断言询问行为

```
package fmt 

func formatOneValue(x interface{}) string { 
	if err, ok := x.(error); ok { 
		return err.Error() 
	}
	if str, ok := x.(Stringer); ok { 
	return str.String() 
	}
	// ...all other types... 
}
```

### 7.13 类型开关

`x.(type)`只能用于`siwtch`

```
switch x.(type) { 
	case nil: // ...
	case int, uint: // ... 
	case bool: // ... 
	case string: // ... 
	default: // ... 
}
```