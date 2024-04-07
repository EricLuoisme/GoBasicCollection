package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	// 使用Echo框架
	e := echo.New()
	e.GET("/user", handleGetUserEcho)
	// 如果是apiErr可做额外逻辑将其parse为json返回等
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if apiErr, ok := err.(APIError); ok {
			c.JSON(apiErr.Status, map[string]any{"error": apiErr.Msg})
		}
		c.JSON(http.StatusInternalServerError, map[string]any{"error": "Internal Server Error"})
	}

	// 传统使用net/http包的执行方式
	http.HandleFunc("/user", handlerGetUser)
	fmt.Println("finished")
}

type APIError struct {
	Status int
	Msg    string
}

// Error 与error interface相同, 相当于继承同一个class而来的
func (e APIError) Error() string {
	return e.Msg
}

// NotFoundError 抽离出来的错误码
func NotFoundError() APIError {
	return APIError{
		Status: http.StatusNotFound,
		Msg:    "User Not found",
	}
}

// handleGetUserEcho 这个是Echo框架下的controller, 它最大的不同是返回error, 并且入参只有context
func handleGetUserEcho(c echo.Context) error {
	// logic
	user, err := getUser()
	if err != nil {
		// 可以再加上自己的Error
		return NotFoundError()
	}
	// echo提供的简单直接的返回, 不需要手动执行write等操作
	return c.JSON(http.StatusOK, user)
}

// handlerGetUser 原始直接使用net/http包实现的controller
func handlerGetUser(w http.ResponseWriter, r *http.Request) {
	// logic
	_, err := getUser()
	// error
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{"err": err.Error()})
		return
	}
	// success
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{"ok": "ok"})
}

type User struct {
}

func getUser() (*User, error) {
	return nil, nil
}
