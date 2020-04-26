package main

import (
	"challenge/env"
	"challenge/repohttp"
	"challenge/rest"
	"challenge/usecase"
)

func main() {
	cfg := env.Get()
	// repository := repolocal.NewRepository("repolocal/feldberg.jpg")
	repository := repohttp.NewRepoWithUrl(
		repohttp.Config{
			Buffersize:  cfg.Buffersize,
			HttpTimeout: cfg.Timeout,
			ImageUrl:    cfg.ImageUrl,
		},
	)
	adapter := rest.NewAdapter()

	file := usecase.NewCreateFile(repository)
	adapter.HandleFunc("/", adapter.DownloadHandler(file))

	adapter.ListenAndServe()
}
