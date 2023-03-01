package main
import (
        "flag"
        "log"
        "net/http"
        "golang.org/x/net/webdav"
)
var (
        add = flag.String("add","0.0.0.0:2800", "server address:port")
        crt = flag.String("crt","", "path/to/your/tls.crt")
        dir = flag.String("dir", "./", "working directory to serve")
        key = flag.String("key","", "path/to/your/tls.key")
        loc = flag.Bool("lock", false, "enable read-only mode")
        nam = flag.String("name", "", "username for authorization")
        pas = flag.String("pass", "", "password for authorization")
        tls = flag.Bool("tls", false, "enable tls mode")
)
func main() {
        flag.Parse()
        http.HandleFunc("/", DavAuth)
        if *tls && *crt != "" && *key != "" {
                log.Printf("%v -> webdav service started, read-only mode: %v, tls mode: enabled", *add, *loc)
                log.Fatal(http.ListenAndServeTLS(*add, *crt, *key, nil))
        } else {
                log.Printf("%v -> webdav service started, read-only mode: %v, tls mode: disabled", *add, *loc)
                log.Fatal(http.ListenAndServe(*add, nil))
        }
}
func DavAuth(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
        if name, pass, set := r.BasicAuth(); !set || name != *nam || pass != *pas {
                log.Printf("%v -> unauthorized connection received, authorizing...", r.RemoteAddr)
                w.WriteHeader(401)
        } else {
                log.Printf("%v <- connection authorized, starting transmissions...", r.RemoteAddr)
                w.Header().Set("Timeout", "86399")
                dav := &webdav.Handler {
                FileSystem: webdav.Dir(*dir), 
                LockSystem: webdav.NewMemLS(),
                }
                if *loc {
                        switch r.Method {
                        case "PUT", "DELETE", "PROPPATCH", "MKCOL", "COPY", "MOVE":
                                log.Printf("%v -> unauthorized operations detected, access denied.", r.RemoteAddr)
                                w.WriteHeader(403)
                                return
                                }
                        dav.ServeHTTP(w, r)
                } else {
                        dav.ServeHTTP(w, r)
                }
        }
}
