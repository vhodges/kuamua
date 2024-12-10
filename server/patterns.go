package server

import (
    "context"
	"encoding/json"
    //"fmt"
	"net/http"
    "log"
	"strconv"

    "github.com/vhodges/kuamua/database"

    "github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/jackc/pgx/v5/pgtype"
)

// PatternCtx middleware is used to load an Pattern object from
// the URL parameters passed through as the request. In case
// the Pattern could not be found, we stop here and return a 404.
func (service *Server) patternCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var Pattern database.Pattern
		var err error

		if PatternID := chi.URLParam(r, "PatternID"); PatternID != "" {
			var id int64
			id, err = strconv.ParseInt(PatternID, 10, 64)
			if err == nil {
				i8 := pgtype.Int8{Int64: id, Valid: true}
				Pattern, err = service.Db.GetPattern(r.Context(), i8)
			}
		} 

		if err != nil || &Pattern == nil {
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "Pattern", &Pattern)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (service *Server) listPatternsHandler(w http.ResponseWriter, r *http.Request) {

	q := database.ListOwnerGroupPatternsParams{
		OwnerID: chi.URLParam(r, "owner"),
		GroupName: chi.URLParam(r, "group"),
		SubGroupName: chi.URLParam(r, "subgroup"),
	}
	
	patterns, err := service.Db.ListOwnerGroupPatterns(context.Background(), q)
    if err != nil || patterns == nil {
		render.Render(w, r, ErrNotFound)
		return
    }

	renderJSON(w, patterns)
}

func (service *Server) getPatternHandler(w http.ResponseWriter, r *http.Request) {

	pattern := r.Context().Value("Pattern").(*database.Pattern)

	renderJSON(w, pattern)
}

func (service *Server) createPatternHandler(w http.ResponseWriter, r *http.Request) {
	newPattern := database.CreatePatternParams{}

    if err := json.NewDecoder(r.Body).Decode(&newPattern); err != nil {
        log.Printf("Decoding body failed: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

	created, err :=  service.Db.CreatePattern(r.Context(), newPattern)

	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
	} else {
		renderJSON(w, created)
	}
}

func (service *Server) updatePatternHandler(w http.ResponseWriter, r *http.Request) {
	pattern := r.Context().Value("Pattern").(*database.Pattern)
	updatePattern := database.UpdatePatternParams{}

	if err := json.NewDecoder(r.Body).Decode(&updatePattern); err != nil {
        log.Printf("Decoding body failed: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

	updatePattern.ID = pattern.ID	
	err :=  service.Db.UpdatePattern(r.Context(), updatePattern)

	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))		
	} else {
		updated, err2 := service.Db.GetPattern(r.Context(), pattern.ID)
		if err2 != nil {
			log.Printf("Error loading updated record: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	
		renderJSON(w, updated)
	}
}

func (service *Server) deletePatternHandler(w http.ResponseWriter, r *http.Request) {

	pattern := r.Context().Value("Pattern").(*database.Pattern)

	err :=  service.Db.DeletePattern(r.Context(), pattern.ID)

	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
	} else {
		render.Render(w, r, &OkResponse{HTTPStatusCode: 200, StatusText: "Deleted"})
	}
}
