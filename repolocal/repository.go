package repolocal

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/reactivex/rxgo/v2"
)

type Repository struct {
	Filename string
}

func NewRepository(filename string) Repository {
	return Repository{Filename: filename}
}

func (r Repository) FetchFile() (rxgo.Observable, error) {
	ch := make(chan rxgo.Item, 1)

	fh, err := os.Open(r.Filename)
	if err != nil {
		return nil, err
	}

	fmt.Println("start FetchFile ...")

	go func() {
		reader := bufio.NewReader(fh)
		for {
			b := make([]byte, 100000)
			n, err := reader.Read(b)
			time.Sleep(time.Millisecond * 20)

			if err != nil {
				if err.Error() == "EOF" {
					ch <- rxgo.Item{
						V: nil,
						E: err,
					}
					err = nil
				}
				if err != nil {
					log.Printf("error while reading local file %v", err)
				}
				break
			}
			ch <- rxgo.Item{
				V: b[:n],
				E: err,
			}
		}
		defer func() {
			if err = fh.Close(); err != nil {
				log.Fatal(err)
			}
		}()
		close(ch)
	}()

	return rxgo.FromChannel(ch, rxgo.WithPublishStrategy()), nil
}
