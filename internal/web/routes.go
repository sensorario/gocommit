package web

import "net/http"

func RegisterRoutes() {
	http.HandleFunc("/branches", BranchesController)
	http.HandleFunc("/checkout", CheckoutController)
	http.HandleFunc("/remotes", RemotesController)
	http.HandleFunc("/pull-requests", PullRequestsController)
	http.HandleFunc("/info", InfoController)
	http.HandleFunc("/version", VersionController)
	http.HandleFunc("/shutdown", ShutdownController)
	http.HandleFunc("/", MainPageController)
}
