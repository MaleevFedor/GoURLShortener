package redirect_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"GoURLShortener/internal/http-server/handlers/redirect"
	"GoURLShortener/internal/http-server/handlers/redirect/mocks"
	"GoURLShortener/internal/lib/api"
	"GoURLShortener/internal/lib/logger/slogmock"
	"GoURLShortener/internal/storage"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://www.google.com/",
		},
		{
			name:      "NotFound",
			alias:     "bad_alias",
			respError: "not found",
			mockError: storage.ErrURLNotFound,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlGetterMock := mocks.NewURLGetter(t)

			urlGetterMock.
				On("GetUrl", tc.alias).
				Return(tc.url, tc.mockError).
				Once()

			r := chi.NewRouter()
			r.Get("/{alias}", redirect.New(slogmock.NewDiscardLogger(), urlGetterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			if tc.respError == "" {
				redirectedToURL, err := api.GetRedirect(ts.URL + "/" + tc.alias)
				require.NoError(t, err)
				assert.Equal(t, tc.url, redirectedToURL)
				return
			}

			resp, err := http.Get(ts.URL + "/" + tc.alias)
			require.NoError(t, err)
			defer resp.Body.Close()

			bodyBytes, _ := io.ReadAll(resp.Body)
			body := string(bodyBytes)

			assert.Contains(t, body, tc.respError)
		})
	}
}
