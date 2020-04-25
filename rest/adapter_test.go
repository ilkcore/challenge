package rest

import (
	"challenge/usecase"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/reactivex/rxgo/v2"
	"github.com/stretchr/testify/assert"
)

type TestRepository struct{}

func NewTestRepository() TestRepository {
	return TestRepository{}
}

func (r TestRepository) FetchFile() (rxgo.Observable, error) {

	ch := make(chan rxgo.Item, 1)
	go func() {
		ch <- rxgo.Of([]byte("1"))
		ch <- rxgo.Of([]byte("2"))
		ch <- rxgo.Of([]byte("3"))
		ch <- rxgo.Item{
			V: nil,
			E: errors.New("EOF"),
		}
		close(ch)
	}()

	return rxgo.FromChannel(ch, rxgo.WithPublishStrategy()), nil
}

func TestMakeDownloadHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", strings.NewReader(""))
	if err != nil {
		t.Error(err)
	}
	u := usecase.NewCreateFile(NewTestRepository())
	resp := httptest.NewRecorder()
	handler := NewAdapter().DownloadHandler(u)
	handler.ServeHTTP(resp, req)

	expected := []byte("123")
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, expected, resp.Body.Bytes())
}
