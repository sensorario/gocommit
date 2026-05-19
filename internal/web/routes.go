package web

import "net/http"

func RegisterRoutes() {
	http.HandleFunc("/branches", BranchesController)
	http.HandleFunc("/checkout", CheckoutController)
       http.HandleFunc("/remotes", RemotesController)
	http.HandleFunc("/info", InfoController)
	http.HandleFunc("/", MainPageController)
}
