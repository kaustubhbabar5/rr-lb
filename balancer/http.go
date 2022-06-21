package balancer

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
	cerrors "github.com/kaustubhbabar5/rr-lb/pkg/errors"
	chttp "github.com/kaustubhbabar5/rr-lb/pkg/http"
)

type Handler struct {
	s IService
}

func NewHandler(s IService) *Handler {
	return &Handler{s}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte{})
}
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		chttp.JSON(w, http.StatusBadRequest, map[string]any{"error": []string{err.Error()}})
		return
	}
	request := RegisterRequest{}
	err = json.Unmarshal(body, &request)
	if err != nil {
		chttp.JSON(w, http.StatusBadRequest, map[string]any{"error": []string{err.Error()}})
		return
	}

	err = h.s.AddServer(request.Endpoint)
	if err != nil {
		//TODO: handle case differently when server is already registered
		chttp.JSON(w, http.StatusInternalServerError, map[string]any{"error": []string{err.Error()}})
		return
	}
}

func (h *Handler) Proxy(w http.ResponseWriter, r *http.Request) {

	sUrl, err := h.s.GetServer()
	if err != nil {
		var nfError *cerrors.NotFound

		ok := errors.As(err, &nfError)
		if ok {
			chttp.JSON(w, http.StatusServiceUnavailable, nil)
			return
		}

		chttp.JSON(w, http.StatusInternalServerError, map[string]any{"error": []string{err.Error()}})
		return
	}

	serverUrl, err := url.Parse(sUrl)
	if err != nil {
		chttp.JSON(w, http.StatusInternalServerError, map[string]any{"error": []string{err.Error()}})
		return
	}
	reverseProxy := httputil.NewSingleHostReverseProxy(serverUrl)

	r.URL.Host = serverUrl.Host
	r.URL.Scheme = serverUrl.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = serverUrl.Host

	r.URL.Path = mux.Vars(r)["rest"]
	reverseProxy.ServeHTTP(w, r)

}
