package handler_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/eddie023/byd/core/logger"
	"github.com/eddie023/byd/pkg/auth"
	"github.com/eddie023/byd/pkg/handler"
	"github.com/eddie023/byd/pkg/store"
	"github.com/stretchr/testify/require"
)

type handlerFixture struct {
	path   string
	method string
	body   string

	query []string

	gotStatusCode int
	gotBody       string

	fakeDBCreateNewPost *store.CreatePost
	fakeDBErr           error
}

func newHandlerFixture(method string, query []string, body string) *handlerFixture {
	return &handlerFixture{
		method: method,
		query:  query,
		body:   body,
	}
}

func (f *handlerFixture) execute(t *testing.T) {
	log := logger.New(os.Stderr, "")

	db := &FakeDBProvider{
		err: f.fakeDBErr,
	}
	if f.fakeDBCreateNewPost != nil {
		db.createNewPost = f.fakeDBCreateNewPost
	}

	localAuth := auth.NewLocalAuth(TestUser.ID, TestUser.Email)

	h, err := handler.NewAPIHandler(db, log, localAuth)
	require.NoError(t, err)

	u, err := url.Parse(f.path)
	require.NoError(t, err)
	// make sure query param has key,value pair
	require.Equal(t, 0, len(f.query)%2)

	q := url.Values{}
	for i := 0; i < len(f.query); i += 2 {
		q.Add(f.query[i], f.query[i+1])
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(context.Background(), f.method, u.String(), strings.NewReader(f.body))
	require.NoError(t, err)

	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	f.gotStatusCode = res.StatusCode
	bodyBuf, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	f.gotBody = string(bodyBuf)
}
