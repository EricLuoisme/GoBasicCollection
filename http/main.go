package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

// http handler 应该返回一个error
func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /user/{id}", HandleGetUserByID)
	router.HandleFunc("GET /user", makeHandler(handleListUsers))

	_ = http.ListenAndServe(":9099", router)
}

type APIError struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

func (e APIError) Error() string {
	slog.Error("[REQUEST FAILED]", e.Error())
	return e.Msg
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

// makeHandler 是一个特殊的方法, 通过包裹传入的apiFunc, 使得apiFunc包含error, 也不会印象将makeHandler包裹到原来的golang中net/http
// 并且, 针对不同类型error的处理, 也可以将其汇总到这个部分
func makeHandler(h apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			// 这里 err.(APIError) 是类型转换, 确认err是否符合APIError这个类,
			// e和ok分别代表转换成功的error, 以及是否转换成功
			if e, ok := err.(APIError); ok {
				slog.Error("http handler error", "err", err, "status", e.Status, "msg", e.Msg)
				writeJson(w, e.Status, e)
			}
		}
	}
}

// HandleGetUserByID 例子是不返回error, 那么可见的遇到不同的error情况都需要手动打印日志并且指定返回什么错误码
func HandleGetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		// 规范化日志使用slog
		// 这个地方可以想到, 如果根据不同的错误, 我们都需要进行日志打印定位错误, 则不同的err我们都需要进行定制化的报错日志处理
		// 但如果你的handler会返回错误, 看下面的例子
		slog.Error("http handler error", "handler", "/getUserByID", "err", err)
		writeJson(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	writeJson(w, http.StatusOK, User{ID: id})
}

// handleListUsers 例子返回error, 但修改方法参数会导致无法放入原来http/net, 需要额外的type进行转换
func handleListUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := getUsers()
	if err != nil {
		return APIError{
			Status: http.StatusBadRequest,
			Msg:    "Bad request my friend",
		}
	}
	return writeJson(w, http.StatusOK, users)
}

func writeJson(w http.ResponseWriter, code int, v any) error {
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func getUsers() ([]User, error) {
	return []User{}, nil
}

type User struct {
	ID uuid.UUID
}
