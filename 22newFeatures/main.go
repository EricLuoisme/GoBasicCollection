package main

import "net/http"

func main() {
	router := http.NewServeMux()
	router.HandleFunc("GET /user/{id}", handlerGetUserByID)
}

// handlerGetUserByID 新版本go1.22可以直接获取path, 但针对HandleFunc并没有返回error
func handlerGetUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	_ = id
}
