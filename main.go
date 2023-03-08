package main
import (
	"flag"
	"log"
	"net/http"
	"golang.org/x/net/webdav"
)
var(
	add = flag.String("add","0.0.0.0:9000", "server address:port")
	crt = flag.String("crt","", "path/to/your/tls.crt (default blank)")
	dir = flag.String("dir", ".", "working directory to serve")
	key = flag.String("key","", "path/to/your/tls.key (default blank)")
	loc = flag.Bool("lock", false, "enable read-only mode (default false)")
	nam = flag.String("name", "", "username for authorization (default blank)")
	pas = flag.String("pass", "", "password for authorization (default blank)")
	pre = flag.String("pre", "/", "webdav prefix: /*")
	noa bool
	tls bool
)
func main() {
	flag.Parse()
	noa = (*nam == "" && *pas == "")
	tls = (*crt != "" && *key != "")
	http.HandleFunc("/", DavCheck)
	if tls {
		log.Printf("-> %v -> [Authorization]=%t [TLS/https]=%t [Read-only]=%v", *add, !noa, tls, *loc)
		log.Fatal(http.ListenAndServeTLS(*add, *crt, *key, nil))
	} else {
		log.Printf("-> %v -> [Authorization]=%t [TLS/https]=%t [Read-only]=%v", *add, !noa, tls, *loc)
		log.Fatal(http.ListenAndServe(*add, nil))
	}
}
func DavCheck(w http.ResponseWriter, r *http.Request) {
	dav := &webdav.Handler {
		Prefix: *pre,
		FileSystem: webdav.Dir(*dir), 
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, e error) {
			if e != nil {
				log.Printf("-> %v -> [%s ERROR]=%v [URL]=%v", r.RemoteAddr, r.Method, e, r.URL)
			} else {
				log.Printf("-> %v -> [%s] [URL]=%v", r.RemoteAddr, r.Method, r.URL)
			}
		},
	}
	switch {
	case noa == false && *loc == true:
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		uname, passwd, _ := r.BasicAuth()
		if uname == *nam && passwd == *pas {
			switch r.Method {
			case "PUT", "DELETE", "PROPPATCH", "MKCOL", "COPY", "MOVE":
				w.WriteHeader(403)
			default:
				dav.ServeHTTP(w, r)
                	}
		} else {
			w.WriteHeader(401)
		}
	case noa == false && *loc == false:
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		uname, passwd, _ := r.BasicAuth()
		if uname == *nam && passwd == *pas {
			dav.ServeHTTP(w, r)
		} else {
			w.WriteHeader(401)
		}
	case noa == true && *loc == true:
		switch r.Method {
		case "PUT", "DELETE", "PROPPATCH", "MKCOL", "COPY", "MOVE":
			w.WriteHeader(403)
		default:
			dav.ServeHTTP(w, r)
                }
	default: dav.ServeHTTP(w, r)
	}
}
