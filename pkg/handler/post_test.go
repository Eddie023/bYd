package handler_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/eddie023/byd/pkg/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateNewPost(t *testing.T) {
	tests := []struct {
		name           string
		httpMethod     string
		queryParam     []string
		wantStatusCode int
		request        string
		wantBody       string
	}{
		{
			name:           "valid request",
			httpMethod:     http.MethodPost,
			wantStatusCode: http.StatusOK,
			request:        `{"title": "test", "description": "test description", "isAnon": false, "type": "1"}`,
			wantBody:       `{"description":"test description", "id": 1, "userId": "user_001", "isAnon":false,"title":"test","type":1, "createdAt": "0001-01-01T00:00:00Z"}`,
		},
		{
			name:           "invalid request",
			httpMethod:     http.MethodPost,
			wantStatusCode: http.StatusBadRequest,
			request:        `{"title": "title", "description": "description"}`,
			wantBody:       `{"code":400,"message":"request body has an error: doesn't match schema: Error at \"/isAnon\": property \"isAnon\" is missing"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := newHandlerFixture(tc.httpMethod, tc.queryParam, tc.request)
			f.path = "/v1/posts"

			// convert request JSON body to fakeDBProvider response
			var dbCreatePostReq store.CreatePost
			err := json.Unmarshal([]byte(tc.request), &dbCreatePostReq)
			require.NoError(t, err)
			f.fakeDBCreateNewPost = &dbCreatePostReq

			f.execute(t)

			require.Equal(t, tc.wantStatusCode, f.gotStatusCode)
			assert.JSONEq(t, tc.wantBody, f.gotBody)
		})
	}
}
