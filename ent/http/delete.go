// Code generated by entc, DO NOT EDIT.

package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"viecco.dev/awesome/ent"
)

// Delete removes a ent.User from the database.
func (h UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	l := h.log.With(zap.String("method", "Delete"))
	// ID is URL parameter.
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		l.Error("error getting id from url parameter", zap.String("id", chi.URLParam(r, "id")), zap.Error(err))
		BadRequest(w, "id must be an integer")
		return
	}
	err = h.client.User.DeleteOneID(id).Exec(r.Context())
	if err != nil {
		switch {
		case ent.IsNotFound(err):
			msg := stripEntError(err)
			l.Info(msg, zap.Error(err), zap.Int("id", id))
			NotFound(w, msg)
		default:
			l.Error("could-not-delete-user", zap.Error(err), zap.Int("id", id))
			InternalServerError(w, nil)
		}
		return
	}
	l.Info("user deleted", zap.Int("id", id))
	w.WriteHeader(http.StatusNoContent)
}
