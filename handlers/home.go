package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func HomeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte("Это мейн страница"))
}
