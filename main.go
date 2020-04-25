package main

import (
	"challenge/repohttp"
	"challenge/rest"
	"challenge/usecase"
)

func main() {
	// repository := repolocal.NewRepository("repolocal/feldberg.jpg")
	repository := repohttp.NewRepoWithUrl("https://upload.wikimedia.org/wikipedia/commons/f/fb/Berliner.Philharmonie.von.Sueden.jpg")
	adapter := rest.NewAdapter()

	file := usecase.NewCreateFile(repository)
	adapter.HandleFunc("/", adapter.DownloadHandler(file))

	adapter.ListenAndServe()
}
