package main

import (
	"bytes"
	"log"
	"net/http"

	pngquant "github.com/manhtai/gopngquant"
)

func homePageHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		buff, err := GetImage(r)
		if err != nil {
			log.Fatalf("GetImage error: %s", err.Error())
		}

		err = pngquant.CompressPng(bytes.NewReader(buff), w, 3)
		if err != nil {
			log.Fatalf("ComporessPng error: %s", err.Error())
		}
	}
	Templ.ExecuteTemplate(w, "index.html", nil)
}
