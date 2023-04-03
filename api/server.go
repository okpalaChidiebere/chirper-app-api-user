package api

import (
	"net/http"
	"strings"

	userv1connect "github.com/okpalaChidiebere/chirper-app-gen-protos/user/v1/userv1connect"
)

type ServerConnectHandlers map[string]func() (path string, handler http.Handler)

type Servers struct {
	UserServer  userv1connect.UserServiceHandler
}

type APIServer struct {
	httpMux *http.ServeMux
}

func (a Servers) NewAPIServer(httpMux *http.ServeMux) *APIServer{
	server := &APIServer{ httpMux: httpMux, }
	return server
}

func (s Servers) GetAllServiceHandlers () ServerConnectHandlers {
	return map[string]func() (path string, handlers http.Handler){
		userv1connect.UserServiceName:  func() (path string, handler http.Handler) {
			return userv1connect.NewUserServiceHandler(s.UserServer)
		},
	}
}

func (a *APIServer) RegisterEndpoints (sh ServerConnectHandlers) {
	for _, f := range sh {
		p, h := f()
    	a.httpMux.Handle(p, allowCORS(h))
	}
}

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Add("Access-Control-Allow-Credentials", "true")
			headers := []string{"Content-Type", "Accept", "X-Requested-With", "Origin"}
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
			methods := []string{"post"}
			w.Header().Set("Access-Control-Allow-Methods", strings.ToUpper(strings.Join(methods, ",")))

			if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}

		h.ServeHTTP(w, r)
	})
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Max-Age", "1728000")
	w.Header().Add("Content-Type", "text/plain; charset=UTF-8")
    w.Header().Add("Content-Length", "0")
	w.WriteHeader(http.StatusNoContent)
}