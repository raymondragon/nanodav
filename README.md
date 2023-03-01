# nanodav
A nano WebDav server writen in golang aims to be lightweight simple and easy to use.

```
~ $ nanodav -h
Usage of nanodav:
  -add string
    	server address:port (default "0.0.0.0:2800")
  -crt string
    	path/to/your/tls.crt (default blank)
  -dir string
    	working directory to serve (default "./")
  -key string
    	path/to/your/tls.key (default blank)
  -lock
    	enable read-only mode (default disabled)
  -name string
    	username for authorization (default blank)
  -pass string
    	password for authorization (default blank)
  -tls
    	enable tls mode (default disabled)
```
