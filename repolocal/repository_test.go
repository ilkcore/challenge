package repolocal

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchFile(t *testing.T) {
	observable, err := NewRepository("feldberg.jpg").FetchFile()
	if err != nil {
		t.Fatal(err)
	}

	fh, err := os.Create("compare.jpg")
	if err != nil {
		log.Println(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	fileWriter := bufio.NewWriter(fh)

	observable.DoOnNext(func(i interface{}) {
		if _, err := fileWriter.Write(i.([]byte)); err != nil {
			t.Errorf("error im DoOnNext write file local %v", err)
		}
		if err = fileWriter.Flush(); err != nil {
			t.Errorf("error im DoOnNext write file local %v", err)
		}
	})

	observable.DoOnError(func(err error) {
		if err.Error() == "EOF" {
			wg.Done()
		}
	})

	observable.Connect()
	wg.Wait()
	defer fh.Close()

	compare, err := ioutil.ReadFile("compare.jpg")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "3e01602fc560ffc2b16accba2cec0ad745271725", fmt.Sprintf("%x", sha1.Sum(compare)))
}
