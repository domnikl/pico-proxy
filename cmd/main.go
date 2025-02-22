package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var paths map[string]string

func extractPath(url *url.URL) string {
	parts := strings.Split(url.Path, "/")
	parts = parts[2:]

	p := strings.Join(parts, "/")
	if url.RawQuery != "" {
		p += "?" + url.RawQuery
	}

	return p
}

func proxy(url string, w http.ResponseWriter, r *http.Request) {
	url = url + "/" + extractPath(r.URL)
	fmt.Printf("Proxying %s to %s\n", r.URL, url)

	request, err := http.NewRequest("GET", url, r.Body)
	if err != nil {
		fmt.Println("Error creating request: " + err.Error())

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Copy headers
	for name, value := range r.Header {
		request.Header.Set(name, value[0])
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Printf("Error requesting %s: %s\n", url, err.Error())

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	// Copy response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response from %s: %s\n", url, err.Error())

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)

	// Copy response headers
	for name, value := range resp.Header {
		w.Header().Set(name, value[0])
	}

	w.Write(body)
}

func handler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path[1:])

	if path == "" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Pico Proxy"))
		return
	}

	for k, v := range paths {
		if strings.HasPrefix(path, k) {
			proxy(v, w, r)
			return
		}
	}

	http.Error(w, "Not found", http.StatusNotFound)
}

func main() {
	if strings.TrimSpace(os.Getenv("INSECURE_SKIP_VERIFY")) == "yes" {
		fmt.Println("InsecureSkipVerify enabled")
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Starting pico-proxy on port %s\n", port)

	setupPaths := strings.Split(os.Getenv("PATHS"), ",")
	paths = make(map[string]string)

	for _, path := range setupPaths {
		v := strings.SplitN(path, ":", 2)
		paths[v[0]] = v[1]

		fmt.Printf("Path: %s, URL: %s\n", v[0], v[1])
	}

	http.HandleFunc("GET /", handler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
