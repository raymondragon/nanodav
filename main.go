package main
import (
        "flag"
        "log"
        "net/http"
        "golang.org/x/net/webdav"
)
var (
        add = flag.String("a","0.0.0.0:2800", "server address:port")
        dir = flag.String("d", "./", "working directory to serve")
        loc = flag.Bool("lock", false, "switching on read-only mode")
        nam = flag.String("n", "admin", "user name for authorization")
        tok = flag.String("t", "adm1n", "user token for authorization")
)
func main() {
        flag.Parse()
        http.HandleFunc("/", DavAuth)
        log.Printf("%v -> starting webdav service for %v, read-only mode: %v", *add, *dir, *loc)
        log.Fatal(http.ListenAndServe(*add, nil))
}
func DavAuth(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
        if name, token, set := r.BasicAuth(); !set || name != *nam || token != *tok {
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
