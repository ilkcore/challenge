package repolocal

import (
	"bufio"
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
	ch := make(chan rxgo.Item, 2)

	fh, err := os.Open(r.Filename)
	if err != nil {
		return nil, err
	}

	log.Println("start FetchFile ...")

	go func() {
		reader := bufio.NewReader(fh)
		counter := 1
		for {
			b := make([]byte, 512*1024)
			n, err := reader.Read(b)
			log.Printf("(%d) read %d", counter, n)
			time.Sleep(time.Millisecond * 20)

			if err != nil {
				if err.Error() == "EOF" {
					log.Printf("reach end of %s", r.Filename)
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
			counter++
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
