package main 

import (
	"fmt"
	"time"
)

var deposits = make(chan int) // send amount to deposit 
var balances = make(chan int) // receive balance

func Deposit(amount int) { 
	time.Sleep(1*time.Second)
	deposits <- amount 
} 

func Balance() int { return <-balances } 

func teller() { 
 	var balance int // balance is confined to teller goroutine 
 	for {
 		select { 
 		case amount := <-deposits: 
 			balance += amount 
 		case balances <- balance: 
 		} 
 	} 
 }

 func main() { 
 	go teller()

 	fmt.Println("amount", Balance())
 	go Deposit(100)
 	go Deposit(100)
 	go Deposit(100)
 	go Deposit(100)
 	time.Sleep(2*time.Second)
 	fmt.Println("amount", Balance())


 }