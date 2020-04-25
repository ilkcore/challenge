package repohttp

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchFile(t *testing.T) {
	testfile := "../testdata/testresult-http.jpg"
	fakeHttp := func() (*http.Response, error) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}
			_, err = w.Write(data)
			if err != nil {
				t.Fatal(err)
			}
		}

		fh, err := os.Open("../testdata/feldberg.jpg")
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest("GET", "/", fh)
		resp := httptest.NewRecorder()
		handler(resp, req)

		return resp.Result(), nil
	}

	observable, err := NewRepoWithFunction(fakeHttp).FetchFile()
	if err != nil {
		t.Fatal(err)
	}

	fh, err := os.Create(testfile)
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

	compare, err := ioutil.ReadFile(testfile)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "3e01602fc560ffc2b16accba2cec0ad745271725", fmt.Sprintf("%x", sha1.Sum(compare)))
}
