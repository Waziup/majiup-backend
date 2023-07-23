package main

import (
	"fmt"
	"net/http"

	"github.com/JosephMusya/majiup-backend/api"
	"github.com/julienschmidt/httprouter"
)

func main() {

	router := httprouter.New()

	api.ApiServe(router)

	fmt.Println("Majiup server running at PORT 8081")
	http.ListenAndServe(":8081", router)
}
