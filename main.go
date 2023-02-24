package main
import (
        "flag"
        "log"
        "net/http"
        "golang.org/x/net/webdav"
)
var (
        add = flag.String("a",":2800", "address:port to listen")
        dir = flag.String("d", ".", "directory to serve")
        nam = flag.String("n", "admin", "user name")
        tok = flag.String("t", "adm1n", "user token")
)
func main() {
        flag.Parse()
        http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
                w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
                if name, token, _ := r.BasicAuth(); name != *nam || token != *tok {
                        log.Printf("%v -> unauthorized connection received, authorizing...", r.RemoteAddr)
                        w.WriteHeader(401)
                } else {
                        log.Printf("%v <- connection authorized, starting transmissions...", r.RemoteAddr)
                        w.Header().Set("Timeout", "86399")
                        svr := &webdav.Handler {
                                FileSystem: webdav.Dir(*dir),
                                LockSystem: webdav.NewMemLS(),
                        }
                        svr.ServeHTTP(w, r)
                }
        })
        if err := http.ListenAndServe(*add, nil); err != nil {
                log.Fatal(err)
        }
}
