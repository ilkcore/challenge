package repohttp

import (
	"bufio"
	"log"
	"net/http"
	"time"
)

type Config struct {
	Buffersize  int
	HttpTimeout int
	ImageUrl    string
}

type Repository struct {
	buffersize      int
	imageUrl        string
	getFileFromHttp func() (*http.Response, error)
}

func NewRepoWithUrl(c Config) Repository {
	buffersize := 32 * 1024
	if c.Buffersize > 0 {
		buffersize = c.Buffersize
	}

	httpTimeout := time.Second * 30
	if c.HttpTimeout > 0 {
		httpTimeout = time.Second * time.Duration(c.HttpTimeout)
	}
	imageUrl := "https://upload.wikimedia.org/wikipedia/commons/f/fb/Berliner.Philharmonie.von.Sueden.jpg"
	if c.ImageUrl != "" {
		imageUrl = c.ImageUrl
	}

	return Repository{
		buffersize: buffersize,
		getFileFromHttp: func() (*http.Response, error) {
			client := http.Client{
				Timeout: httpTimeout,
			}
			return client.Get(imageUrl)
		},
	}
}

func NewRepoWithFunction(f func() (*http.Response, error)) Repository {
	return Repository{
		getFileFromHttp: f,
	}
}

func (r Repository) FetchFile() (chan []byte, error) {
	ch := make(chan []byte, 1)

	resp, err := r.getFileFromHttp()
	if err != nil {
		return nil, err
	}

	go func() {
		reader := bufio.NewReader(resp.Body)
		counter := 1
		for {
			b := make([]byte, r.buffersize)
			n, err := reader.Read(b)
			// log.Printf("(%d) read %d", counter, n)

			if err != nil {
				if err.Error() == "EOF" {
					log.Println("reading response body completed")
					ch <- []byte("done")
					err = nil
				}
				if err != nil {
					log.Printf("error while reading local file %v", err)
				}
				break
			}

			ch <- b[:n]
			counter++
		}
		defer func() {
			if resp.Body.Close() != nil {
				log.Println(err)
			}
			log.Println("close body")
		}()
		close(ch)
	}()

	return ch, nil
}
