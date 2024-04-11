package main

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"strings"
)

// PaymentStrategy 策略模式的核心interface
type PaymentStrategy interface {
	Pay(amount float64) error
}

type CreditCard struct {
}

// Pay 每一个struct, 都需要有自己实现interface的方法
func (c *CreditCard) Pay(amount float64) error {
	fmt.Printf("Paid $%.2f using Credit card\n", amount)
	return nil
}

type Paypal struct {
}

// Pay 每一个struct, 都需要有自己实现interface的方法
func (c *Paypal) Pay(amount float64) error {
	fmt.Printf("Paid $%.2f using Paypal\n", amount)
	return nil
}

// PaymentHandler 是处理web请求的Fiber框架的handler
func PaymentHandler(c fiber.Ctx) error {

	// 抽离
	var strategy PaymentStrategy
	payMethod := strings.ToLower(c.Params("payMethod"))

	// 判断执行方法
	switch payMethod {
	case "paypal":
		strategy = &Paypal{}
	default:
		strategy = &CreditCard{}
	}

	// 进行执行并返回内容
	_ = strategy.Pay(100)
	return c.JSON(fiber.Map{"message": fmt.Sprintf("Paid using %s", payMethod)})
}

func main() {
	app := fiber.New()
	app.Get("/pay/:payMethod", PaymentHandler)
	_ = app.Listen(":8989")
}
