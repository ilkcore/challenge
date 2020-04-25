package rest

import (
	"log"
	"net/http"
	"springmedia/observerwriter"
	"springmedia/usecase"
	"sync"

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
	_ = http.ListenAndServe(":8080", a.r)
}

func (a Adapter) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route {
	return a.r.NewRoute().Path(path).HandlerFunc(f)
}

func (a Adapter) DownloadHandler(fetchfile usecase.File) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("start DownloadHandler")

		var wg sync.WaitGroup
		wg.Add(1)
		observable, err := fetchfile.Run()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Encoding", "chunked")
		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Disposition", "filename=berliner-philharmoni.jpg")

		observerwriter.ToHttpWriter(observable, w)
		closeFiles := observerwriter.WriteToLocalFile(observable)
		defer closeFiles()

		observable.DoOnError(func(err error) {
			if err.Error() == "EOF" {
				wg.Done()
			}
		})

		observable.Connect()
		wg.Wait()
	}
}
