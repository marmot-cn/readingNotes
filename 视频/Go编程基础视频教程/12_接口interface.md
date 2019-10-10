#12 接口interface

[地址](http://edu.51cto.com/lesson/id-32319.html "地址")

###笔记

---

不用显示的声明实现了哪些接口,而是只要实现这个接口的方法,即代表你实现了该接口.

**接口interface**

* 接口是一个或多个方法签名的集合
* 只要某个类型拥有该接口的所有方法签名, 即算实现该接口, 无需显示生命实现了哪个接口, 称为 `Structual Typing`
* 接口只有方法声明, 没有实现, 没有数据字段
* 接口可以`匿名嵌入`其他`接口`, 或嵌入到`结构`中
	
		type USB interface {
			Name() string
			Connecter
		}
		
		type Connecter interface {
			Connect()
		}	

* 将对象赋值给接口时, 会发生拷贝, 而接口内部存储的是指向这个复制品的指针, 即无法修改复制品的状态, 也无法获取指针
	
		...
		pc := PhoneConnecter{"PhoneConnecter"}
		var a Connecter
		//
		a = Connecter(pc)//a通过强制类型转换把 PhoneConnecter 转化为 Connecter类型
		a.Connect()
		
		pc.name = "pc"//修改pc的名字不会影响后续a.Connect()的输出
		a.Connect()

* 只有当接口存储的类型和对象都为`nil`时, 接口才等于`nil`

		 var a interace{}// = nil
		 
		 var p *int = nil
		 a = p //存储的类型是指针,不为nil a != nil
		 
* 接口调用不会做`receiver`的自动转换(涉及方法集问题)
* 接口同样支持`匿名字段`方法
* 接口也可实现类似OOP中的多态
* 空接口可以作为任何类型数据的容器
		
		//go语言中所有的类型都实现了空接口
		type empty interface {
		}


**示例:实现接口**

		//定义接口
		type USB interface {
			Name() string
			Connect()
		}
		
		type PhoneConnecter struct {
			name string
		}
		
		func (pc PhoneConnecter) Name() string {
			return pc.name
		}
		
		func (pc PhoneConnecter) Connect() {
			fmt.Println("Connected:", pc.name)
		}
		
		
		func main(){
			
			var a USB
			
			//不能这样定义,因为USB没有name字段
			//a = PhoneConnecter{}
			//a.name = "PhoneConnecter"
			
			a = PhoneConnecter{"PhoneConnecter"}
			a.Connect()
			Disconnect(a)
			
		}
		
		
		func Disconnect(usb USB){
			if pc,ok := usb.(PhoneConnecter);ok{
				fmt.Println("Disconnected:", pc.name)
				return 
			}
			
			fmt.Println("Unknown Device")
		}
		
		修改 Disconnect 方法,可以接受任何类型传参(空接口):
		func Disconnect(usb interface{}){
			
			//使用switch 判断类型,接受空接口时用的比较多
			switch v:=usb.(type)
				case PhoneConnecter:
					fmt.Println("Disconnected:", v.name)
				default:
					fmt.Println("Unknown Device")
				
		}
		

###整理知识点

---
