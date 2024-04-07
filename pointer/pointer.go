package main

import "fmt"

type User struct {
	email    string
	username string
	age      int
}

// When to use pointer:
// 1) when we need to update state (only pointer could remain the updates)
// 2) when input is too large (pointer is only fixed size 8 bytes), and be called a lot

func (u User) Email() string {
	return u.email
}

// updateEmail BIG PROBLEM 因为这里并不是Pointer 无法进行内容的更新
// 这里的User u是正常的receiver, golang为它生成了一个副本, 用来接收, 所以这里的更新并没有对内存中的实例进行变更
func (u User) updateEmail(email string) {
	u.email = email
}

// updateEmailPointer 由于使用pointer, 更新内容将会被保存
// 这里的User
func (u *User) updateEmailPointer(email string) {
	u.email = email
}

func updateEmailFromInputPointer(u *User, email string) {
	u.email = email
}

// 这个地方没有传pointer, 同样返回后无法影响到调用方的传入, 处理的时候都是拷贝了一个副本进行处理的
func updateEmailFromInput(u User, email string) {
	u.email = email
}

// Email 当我们传入pointer, 固定的大小是8bytes
func Email(user *User) string {
	return user.email
}

func main() {
	user := User{email: "abc@foo.com"}

	// 没有进行内容的更新
	user.updateEmail("ccc@qq.com")
	fmt.Printf("%+v\n", user)

	// 发现使用pointer的部分, 是能够成功更新内容的
	user.updateEmailPointer("1g1@protonmail.com")
	fmt.Printf("%+v\n", user)

	// 发现传入指针进行操作, 也会使得影响持久化
	updateEmailFromInputPointer(&user, "25gq@yes.com")
	fmt.Printf("%+v\n", user)

	// 可以看到, 如果传入的不是指针, 同样会被复制一个副本过去, 修改是不能影响原来的调用方实例
	updateEmailFromInput(user, "345@outlook.com")
	fmt.Printf("%+v\n", user)

}
