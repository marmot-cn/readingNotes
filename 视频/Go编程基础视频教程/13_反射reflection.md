#13 反射reflection

[地址](http://edu.51cto.com/lesson/id-32320.html "地址")

###笔记

---

**反射**

* 配合`interface{}`
* 使用`TypeOf`和`ValueOf`函数从接口中获取目标对象信息
* 反射会将匿名字段作为独立字段(匿名字段本质)
* 通过反射可以`动态`调用方法

**示例**

		引入 "reflect" 包

		type User struct {
			Id int
			Name string
			Age int
		}
		
		func (u User) Hello() {
			fmt.Println("Hello world.")
		}
		
		func main(){
			u := User{1,"OK",12}
			info(u)
		}
		
		func Info(o Interface){
			t := reflect.Typeof(o)
			fmt.Prinln("type:", t.Name())
			
			v := reflect.Valueof(o)
			fmt.Prinln("Fields:")
			for i:=0; i < t.NumField(); i++ {
				f := t.Field(i)
				//使用 Interface 方法还原接口值
				val := v.Field(i).Interface()
				
				fmt.Prinf("%6s : %v = %v", f.Name, f.Type, val)
			}
			//获取方法信息
			for i:=0; i< t.NumMethod(); i++ {
				m := t.Method(i)
				fmt.Printf("%6s : %v = %v", m.Name, m.Type)
			}
		}
		
		输出:
		Type: User
		Fields:
			Id: int = 1
			Name: string = ok
			Age : int = 12
		Hello: func(main.User)//receiver 就是方法接收的第一个参数
		
		
		
		匿名字段:
		
		type Manager struct {
			User //User 是匿名字段,等同于 User User
			title string
		}
		
**示例:通过反射修改内容**

		x := 123
		v := reflect.ValueOf(&x)
		v.Elem().SetInt(999)
		
		fmr.Prnintln(x)//输出 999


**reflect.ValueOf**

将真实的数据类型用 `reflect.ValueOf` 转换成反射的数据类型.

###整理知识点

---