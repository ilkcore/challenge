package main

import (
	"springmedia/repolocal"
	"springmedia/rest"
	"springmedia/usecase"
)

func main() {
	repository := repolocal.NewRepository("repolocal/feldberg.jpg")
	adapter := rest.NewAdapter()

	file := usecase.NewCreateFile(repository)
	adapter.HandleFunc("/", adapter.DownloadHandler(file))

	adapter.ListenAndServe()
}
