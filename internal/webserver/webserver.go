package webserver

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/fsgh42/meshmap.net/internal/meshtastic"
)

var (
	HTMLtemplatePath = "template.html"
)

type WebServer struct {
	addr  string
	mux   *http.ServeMux
	Nodes meshtastic.NodeDB
}

func (ws *WebServer) Run() {
	log.Printf("[srv] start listening on %s", ws.addr)
	http.ListenAndServe(ws.addr, ws.mux)
	log.Printf("[srv] end listening on %s", ws.addr)
}

func handleMap(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write(meshmapHTML)
}

func handleNodes(w http.ResponseWriter, r *http.Request, nodes meshtastic.NodeDB) {
	data := []byte("{}")
	if len(nodes) > 0 {
		_data, err := json.Marshal(nodes)
		if err != nil {
			log.Printf("[err] error while marshalling database: %s", err)
			http.Error(w, "error serializing database", http.StatusInternalServerError)
			return
		}
		data = _data
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/map", http.StatusPermanentRedirect)
}

func logClients(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[srv] %s %s %s\n", r.RemoteAddr, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func NewWebServer(addr string) *WebServer {
	ws := &WebServer{addr: addr, mux: http.NewServeMux()}

	ws.mux.HandleFunc("/", redirectHandler)
	ws.mux.Handle("/map", logClients(http.HandlerFunc(handleMap)))
	ws.mux.Handle("/nodes", logClients(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				handleNodes(w, r, ws.Nodes)
			},
		),
	))
	ws.mux.HandleFunc("/site.webmanifest", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(webmanifest)
	})
	ws.mux.HandleFunc("/android-chrome-192x192.png", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(iconChrome192)
	})
	ws.mux.HandleFunc("/android-chrome-512x512.png", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(iconChrome512)
	})
	ws.mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(iconFav)
	})

	return ws
}
