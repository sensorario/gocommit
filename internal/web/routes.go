package web

import "net/http"

func RegisterRoutes() {
	http.HandleFunc("/branches", BranchesController)
	http.HandleFunc("/current-branch", CurrentBranchController)
	http.HandleFunc("/checkout", CheckoutController)
	http.HandleFunc("/merge-into-next", MergeIntoNextController)
	http.HandleFunc("/remotes", RemotesController)
	http.HandleFunc("/pull-requests", PullRequestsController)
	http.HandleFunc("/info", InfoController)
	http.HandleFunc("/modified-files", ModifiedFilesController)
	http.HandleFunc("/diff", DiffController)
	http.HandleFunc("/version", VersionController)
	http.HandleFunc("/shutdown", ShutdownController)
	http.HandleFunc("/repositories", RepositoriesController)
	http.HandleFunc("/repo-versions", RepoVersionsController)
	http.HandleFunc("/current-repo", CurrentRepoController)
	http.HandleFunc("/switch-repo", SwitchRepoController)
	http.HandleFunc("/", MainPageController)
}
