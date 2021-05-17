package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func addPhoto(w http.ResponseWriter, r *http.Request) {
	// Saves photo to local storage and creates thumbnail.
	isThumbnailStr := r.FormValue("thumbnail")
	isThumbnail, err := strconv.ParseBool(isThumbnailStr)
	if err != nil {
		isThumbnail = false
	}
	fmt.Println(isThumbnail)
	imageNames, err := saveImageLocal(r, isThumbnail)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(map[string][]string{"image_names": imageNames})
}
