package main

import ("fmt")

func main(){
   //可以同时定义多个变量
   var user1,user2="Bob","Alice"
   //定义变量时可以指定类型，也可以不指定类型
   var age1,age2 int=1,2
   //%s 字符串占位符     %d 数字占位符
   //其它参考https://studygolang.com/articles/2880
   fmt.Printf("Hello World %s %d %s %d",user1,age1,user2,age2); 
}