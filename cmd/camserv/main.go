package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/dhowden/raspicam"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("Received request from: %v\n", req.Host)

		s := raspicam.NewStill()
		s.Width = 1640
		s.Height = 1233
		s.Quality = 100

		q := req.URL.Query()
		if val := q.Get("w"); val != "" {
			if v, err := strconv.Atoi(val); err == nil {
				s.Width = v
			}
		}
		if val := q.Get("h"); val != "" {
			if v, err := strconv.Atoi(val); err == nil {
				s.Height = v
			}
		}
		if val := q.Get("r"); val != "" {
			if v, err := strconv.Atoi(val); err == nil {
				s.Camera.Rotation = v
			}
		}
		if val := q.Get("hflip"); val != "" {
			if v, err := strconv.ParseBool(val); err == nil {
				s.Camera.HFlip = v
			}
		}
		if val := q.Get("vflip"); val != "" {
			if v, err := strconv.ParseBool(val); err == nil {
				s.Camera.VFlip = v
			}
		}

		errCh := make(chan error)
		go func() {
			for x := range errCh {
				log.Printf("%v", x)
			}
		}()

		log.Printf("Capturing image with params: %s", strings.Join(s.Params(), ", "))

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "image/jpeg")

		raspicam.Capture(s, w, errCh)
		log.Println("Done")
	})

	err := http.ListenAndServe(":6666", nil)
	if err != nil {
		log.Fatalf("listen and serve: %s", err)
	}
}
