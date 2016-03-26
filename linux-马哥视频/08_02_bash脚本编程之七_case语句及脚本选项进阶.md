#08_02 bash脚本编程之七 case语句及脚本选项进阶

###笔记

---

**面向过程**

`控制结构`:

* 顺序结构
* 选择结构
* 循环结构

`选择结构`:

	if: 单分支,双分支,多分支		if CONDITION; then			statement			…		fi				if CONDITION; then			statement			…		else			statement			…		fi		if CONDITION1; then			statement			…		elif CONDITION2; then			statement			…		else			statement			…		fi`case`语句:选择结构:
		case SWITCH in		value1)			statement			…			;;		value2)			statement			…			;;		*)			statement			…			;;		esac
###整理知识点

---