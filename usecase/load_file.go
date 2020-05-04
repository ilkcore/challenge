package usecase

type FilePort interface {
	FetchFile() (chan []byte, error)
}

type File struct {
	repository FilePort
}

func NewCreateFile(repository FilePort) File {
	return File{repository: repository}
}

func (f File) Run() (chan []byte, error) {
	return f.repository.FetchFile()
}
