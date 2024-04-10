package main

import "fmt"

type User struct {
	email string
	age   int
}

func main() {

	// 正常声明
	user := User{
		email: "a@gg.com",
		age:   101,
	}
	fmt.Printf("%+v\n", user)

	// 空声明, golang会有默认值
	var userEmpty User
	fmt.Printf("%+v\n", userEmpty)

	// pointer声明, 如果没有赋值, 可以发现这个是nil, 很可能引发空指针问题
	var userP *User
	fmt.Printf("%+v\n", userP)

	// 这样就空指针了, runtime error
	// fmt.Println(userP.age)
	// 可以这样理解, 这个地方相当于告知compiler需要de-reference取出里面的值,
	// 但pointer是没有默认值的, 直接为nil, 那么compiler获取的时候就会报runtime error

	// 引申出 -> 不要返回pointer, 应该更多的使用'谁'允许调用这个方法, 并作为指针调用(保存update数据)
}
