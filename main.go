package main

func main() {
	repository := repolocal.NewRepository("repolocal/feldberg.jpg")
	adapter := rest.NewAdapter()

	file := usecase.NewCreateFile(repository)
	adapter.HandleFunc("/", adapter.DownloadHandler(file))

	adapter.ListenAndServe()
}