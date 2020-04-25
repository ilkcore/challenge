package observerwriter

import (
	"bufio"
	"log"
	"os"

	"github.com/reactivex/rxgo/v2"
)

func WriteToLocalFile(observable rxgo.Observable) func() {
	fh, err := os.Create("output.jpg")
	if err != nil {
		log.Println(err)
	}
	fileWriter := bufio.NewWriter(fh)
	observable.DoOnNext(func(i interface{}) {
		if _, err := fileWriter.Write(i.([]byte)); err != nil {
			log.Printf("error im DoOnNext write file local %v", err)
		}
		if err = fileWriter.Flush(); err != nil {
			log.Printf("error fileWriter %v", err)
		}
	})
	return func() {
		if err := fh.Close(); err != nil {
			panic(err)
		}
	}
}
