package repohttp

import (
	"bufio"
	"log"
	"net/http"

	"github.com/reactivex/rxgo/v2"
)

type Repository struct {
	GetFileFromHttp func() (*http.Response, error)
}

func NewRepoWithUrl(url string) Repository {
	return Repository{
		GetFileFromHttp: func() (*http.Response, error) {
			return http.Get(url)
		},
	}
}

func NewRepoWithFunction(f func() (*http.Response, error)) Repository {
	return Repository{
		GetFileFromHttp: f,
	}
}

func (r Repository) FetchFile() (rxgo.Observable, error) {
	ch := make(chan rxgo.Item, 1)

	resp, err := r.GetFileFromHttp()
	if err != nil {
		return nil, err
	}

	go func() {
		reader := bufio.NewReader(resp.Body)
		for {
			b := make([]byte, 512*1024)
			n, err := reader.Read(b)

			if err != nil {
				log.Println("Error reading file:", err)
				if err.Error() == "EOF" {
					log.Println("reading response body completed")
					ch <- rxgo.Item{
						V: nil,
						E: err,
					}
					err = nil
				}
				break
			}

			ch <- rxgo.Item{
				V: b[:n],
				E: err,
			}
		}
		defer func() {
			if resp.Body.Close() != nil {
				log.Println(err)
			}
			log.Println("Close body")
		}()
		close(ch)
	}()

	return rxgo.FromChannel(ch, rxgo.WithPublishStrategy()), nil
}
