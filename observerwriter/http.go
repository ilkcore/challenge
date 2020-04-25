package observerwriter

import (
	"log"
	"net/http"

	"github.com/reactivex/rxgo/v2"
)

func ToHttpWriter(observable rxgo.Observable, w http.ResponseWriter) {
	observable.DoOnNext(func(i interface{}) {
		_, err := w.Write(i.([]byte))
		if err != nil {
			log.Printf("error im DoOnNext write file to response %v", err)
		}
	})
}
