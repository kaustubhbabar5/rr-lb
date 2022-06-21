package balancer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
	chttp "github.com/kaustubhbabar5/rr-lb/pkg/http"
	"github.com/kaustubhbabar5/rr-lb/stratergy/robin"
)

type Handler struct {
	s              IService
	balancingStrat robin.IService
}

func NewHandler(s IService, robinService robin.IService) *Handler {
	return &Handler{s, robinService}
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

	err = h.s.AddNode(request.Url)
	if err != nil {
		chttp.JSON(w, http.StatusInternalServerError, map[string]any{"error": []string{err.Error()}})
		return
	}
}

func (h *Handler) Proxy(w http.ResponseWriter, r *http.Request) {

	Url, err := h.balancingStrat.GetServer()
	if err != nil {
		chttp.JSON(w, http.StatusInternalServerError, map[string]any{"error": []string{err.Error()}})
		return
	}

	serverUrl, err := url.Parse(Url)
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
	fmt.Println(r.URL.Path)
	reverseProxy.ServeHTTP(w, r)

}
