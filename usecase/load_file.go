package usecase

import "github.com/reactivex/rxgo/v2"

type FilePort interface {
	FetchFile() (rxgo.Observable, error)
}

type File struct {
	repository FilePort
}

func NewCreateFile(repository FilePort) File {
	return File{repository: repository}
}

func (f File) Run() (rxgo.Observable, error) {
	return f.repository.FetchFile()
}
