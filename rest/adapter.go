package rest

import (
	"bufio"
	"challenge/usecase"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type Adapter struct {
	r *mux.Router
}

func NewAdapter() Adapter {
	return Adapter{mux.NewRouter()}
}

func (a Adapter) ListenAndServe() {
	log.Printf("Listening on http://0.0.0.0%s\n", ":8080")
	err := http.ListenAndServe(":8080", a.r)
	if err != nil {
		log.Fatal(err)
	}
}

func (a Adapter) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route {
	return a.r.NewRoute().Path(path).HandlerFunc(f)
}

func (a Adapter) DownloadHandler(fetchfile usecase.File) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("start DownloadHandler")

		ch, err := fetchfile.Run()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Encoding", "chunked")
		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Disposition", "filename=berliner-philharmoni.jpg")

		httpChan := make(chan []byte)
		localChan := make(chan []byte, 1700)

		go func() {
			for incomingBytes := range ch {
				httpChan <- incomingBytes
				localChan <- incomingBytes
			}
			close(httpChan)
			close(localChan)
		}()

		go func() {
			fh, err := os.Create("output.jpg")
			if err != nil {
				log.Println(err)
			}
			defer fh.Close()

			fileWriter := bufio.NewWriter(fh)
			counter := 0
			for data := range localChan {
				time.Sleep(time.Millisecond * 7)
				if _, err := fileWriter.Write(data); err != nil {
					log.Printf("error at DoOnNext write file local %v", err)
				}
				if err = fileWriter.Flush(); err != nil {
					log.Printf("error fileWriter %v", err)
				}
				counter++
			}
			log.Println("finish local delivery")
		}()

		for data := range httpChan {
			_, err := w.Write(data)
			if err != nil {
				log.Printf("error write data to http responsewriter: %v", err)
			}
		}
		log.Println("finish http delivery")
		return
	}
}
