package handlers

import (
	"github.com/gin-gonic/gin"
)

// Learn is the HTTP handler to process incoming Learn requests
// It persists the agreed upon value to its local KV Store
func Learn(c *gin.Context) {
	// // Obtain the key & value from URL params
	// key := mux.Vars(r)["key"]
	// value := mux.Vars(r)["value"]

	// err := Store.Set(key, value)
	// if err != nil {
	// 	log.WithFields(log.Fields{"error": err}).Error("failed to set value")
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// log.WithFields(log.Fields{
	// 	"key":   key,
	// 	"value": value,
	// }).Debug("successful learn")

	// w.WriteHeader(http.StatusOK)
}
