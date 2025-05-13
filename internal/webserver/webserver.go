package webserver

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/fsgh42/meshmap.net/internal/meshtastic"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

var (
	HTMLtemplatePath = "template.html"
)

type WebServer struct {
	mux         *http.ServeMux
	Nodes       meshtastic.NodeDB
	certManager *autocert.Manager
}

func (ws *WebServer) Run() {
	log.Printf("[srv] start HTTP serve on ports 8080(http) and 8443(https)")

	go func() {
		// port is mapped to external 80 via docker-compose
		http.ListenAndServe(":8080", ws.certManager.HTTPHandler(ws.mux))
	}()

	getCertificate := func(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
		cert, err := ws.certManager.GetCertificate(chi)
		if err != nil {
			log.Printf("[tls]: error: %v", err)
		}
		return cert, err
	}

	httpsSrv := http.Server{
		// port is mapped to external 443 via docker-compose
		Addr:    ":8443",
		Handler: ws.mux,
		TLSConfig: &tls.Config{
			GetCertificate: getCertificate,
		},
	}

	if err := httpsSrv.ListenAndServeTLS("", ""); err != nil {
		log.Fatal(err)
	}

	log.Printf("[srv] end HTTP serve")
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

func NewAcmeManager(certPath, domain string) *autocert.Manager {
	m := &autocert.Manager{
		Cache:      autocert.DirCache(certPath),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domain),
	}

	if acmeUrl := os.Getenv("ACME_URL"); acmeUrl != "" {
		log.Printf("[acme]: using URL: \"%s\"", acmeUrl)
		m.Client = &acme.Client{
			DirectoryURL: acmeUrl,
		}
	}

	return m
}

func NewWebServer(certManager *autocert.Manager) *WebServer {

	ws := &WebServer{
		mux:         http.NewServeMux(),
		Nodes:       meshtastic.NodeDB{},
		certManager: certManager,
	}

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
