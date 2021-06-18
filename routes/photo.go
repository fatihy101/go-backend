package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"cloud.google.com/go/storage"
	"enstrurent.com/server/utils"
	"github.com/go-chi/chi/v5"
)

func addPhoto(w http.ResponseWriter, r *http.Request) {
	// Saves photo to local storage and creates thumbnail.
	isThumbnailStr := r.FormValue("thumbnail")
	isThumbnail, err := strconv.ParseBool(isThumbnailStr)
	if err != nil {
		isThumbnail = false
	}
	fmt.Print(isThumbnail) // FIXME create thumbnail
	// imageNames, err := saveImageLocal(r, isThumbnail)
	imageNames, err := utils.UploadToCloud(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(map[string][]string{"image_names": imageNames})
}

func getPhoto(w http.ResponseWriter, r *http.Request) {
	photoName := chi.URLParam(r, "photoName")
	bytes, err := utils.DownloadFromCloud(photoName)
	if errors.Is(err, storage.ErrObjectNotExist) {
		http.Error(w, "photo doesn't exist", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(bytes)
}

func deletePhoto(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	isThumbnailStr := chi.URLParam(r, "thumbnail")
	isThumbnail, err := strconv.ParseBool(isThumbnailStr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = utils.DeleteFromCloud(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if isThumbnail {
		if err = utils.DeleteFromCloud("thumb-" + id); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	json.NewEncoder(w).Encode(map[string]bool{"success": true})

}
