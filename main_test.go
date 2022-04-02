package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alekstet/social_graph/conf"

	"github.com/stretchr/testify/assert"
)

func TestSocial(t *testing.T) {
	params, err := conf.Cnf()
	if err != nil {
		t.Errorf("%s", err)
	}
	s := New(params)
	err = s.DB()
	if err != nil {
		t.Errorf("%s", err)
	}
	_, err = s.db.Query("DELETE FROM public.social")
	if err != nil {
		t.Errorf("%s", err)
	}

	matrix := [][]int{{0, 2, 1}, {2, 0, 0}, {1, 0, 0}}
	info := Info{2, 1, 1.5}
	r := Resp{matrix, info}
	expected, err := json.Marshal(r)
	if err != nil {
		t.Errorf("%s", err)
	}
	handler := http.HandlerFunc(s.Social)

	t.Run("PUT empty data", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("PUT", "/social", nil)
		handler.ServeHTTP(w, r)
		assert.Equal(t, 400, w.Result().StatusCode)
	})

	t.Run("PUT negative data", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("PUT", "/social?from=-1&to=-2", nil)
		handler.ServeHTTP(w, r)
		assert.Equal(t, 400, w.Result().StatusCode)
	})

	t.Run("PUT equal data", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("PUT", "/social?from=1&to=1", nil)
		handler.ServeHTTP(w, r)
		assert.Equal(t, 400, w.Result().StatusCode)
	})

	t.Run("PUT valid data 1", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("PUT", "/social?from=1&to=2", nil)
		handler.ServeHTTP(w, r)
		assert.Equal(t, []byte(nil), w.Body.Bytes())
	})

	t.Run("PUT valid data 2", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("PUT", "/social?from=1&to=2", nil)
		handler.ServeHTTP(w, r)
		assert.Equal(t, []byte(nil), w.Body.Bytes())
	})

	t.Run("PUT valid data 3", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("PUT", "/social?from=1&to=3", nil)
		handler.ServeHTTP(w, r)
		assert.Equal(t, []byte(nil), w.Body.Bytes())
	})

	t.Run("GET", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/social", nil)
		handler.ServeHTTP(w, r)
		assert.Equal(t, expected, w.Body.Bytes())
	})
}
