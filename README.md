# http-server
A very basic HTTP server for quickly serving a directory similar to python's SimpleHTTPServer.  Serves a default
page of index.html if it exists, otherwise a directory listing if the path is a directory or the appropriate file
with root directory safety-checks.

# Example

    $ http-server
    Serving /home/chris/html on :8080

    $ http-server :8081
    Serving /home/chris/html on :8081
