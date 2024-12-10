package server

import (
	"io"
    "log"
	"net/http"
)

func (service *Server) matchHandler(w http.ResponseWriter, req *http.Request) {

    query := req.URL.Query()

    q, err := service.store.GetQuamina(query.Get("owner"), 
	  query.Get("group"), 
	  query.Get("subgroup"))

    if err != nil {
        log.Printf("Error from getQuamina: '%v'\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
    }
    var doc []byte

	doc, err = io.ReadAll(req.Body)
    if err != nil {
        log.Printf("Error in ReadBody: '%v'\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
    }

    matches, qerr := q.MatchesForEvent(doc)
    if qerr != nil {
        log.Printf("Error in Matches for Event: '%v'\n", qerr)
        http.Error(w, qerr.Error(), http.StatusInternalServerError)
        return
    }
    renderJSON(w, matches)
}

