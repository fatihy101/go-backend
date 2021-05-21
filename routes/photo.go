package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func addPhoto(w http.ResponseWriter, r *http.Request) {
	// Saves photo to local storage and creates thumbnail.
	isThumbnailStr := r.FormValue("thumbnail")
	isThumbnail, err := strconv.ParseBool(isThumbnailStr)
	if err != nil {
		isThumbnail = false
	}
	imageNames, err := saveImageLocal(r, isThumbnail)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(map[string][]string{"image_names": imageNames})
}

func GetPhotos(r chi.Router, path string, root http.FileSystem) {

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}

	if subStr := strings.Split(path, "/"); subStr[2] == "" || strings.ContainsAny(path, "{}*") {
		r.Get(path, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte("BAD REQUEST"))
			w.WriteHeader(http.StatusBadRequest)
		})
	}
	path += "*"
	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		// photoName := chi.URLParam(r, "*")
		// photoName = strings.Replace(photoName, "thumb-", "", -1) // TODO Search photo name in folder, if it does not exist redirect.
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		w.Header().Set("Content-Type", "image/png")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func replaceFile(tempStr string, oldStr string) {
	tempDir := fmt.Sprintf("%v%v%v", imageFolderDir, PathSeparator, tempStr)
	photoDir := fmt.Sprintf("%v%v%v", imageFolderDir, PathSeparator, oldStr)

	os.Remove(photoDir)
	os.Rename(tempDir, photoDir)
}

func updatePhoto(w http.ResponseWriter, r *http.Request) {
	photoID := r.FormValue("image_id")
	isThumbnailStr := r.FormValue("thumbnail")
	isThumbnail, _ := strconv.ParseBool(isThumbnailStr)

	tempPhotoNames, err := saveImageLocal(r, isThumbnail)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	for _, tempPhotoName := range tempPhotoNames {
		replaceFile(fmt.Sprint(tempPhotoName), photoID)
		if isThumbnail {
			replaceFile(fmt.Sprintf("thumb-%v", tempPhotoName), fmt.Sprintf("thumb-%v", photoID))
		}
	}

}

func deletePhoto(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	isThumbnailStr := chi.URLParam(r, "thumbnail")
	isThumbnail, err := strconv.ParseBool(isThumbnailStr)

	if err != nil {
		fmt.Println("Error on deleting photo:" + err.Error())
		return
	}
	// Delete from local storage.
	imageDir := fmt.Sprintf("%v%v%v", imageFolderDir, PathSeparator, id)
	if err := os.Remove(imageDir); err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusExpectationFailed)
		json.NewEncoder(w).Encode(map[string]bool{"success": false})
		return
	}

	if isThumbnail {
		imageDirThumb := fmt.Sprintf("%v%v%v", imageFolderDir, PathSeparator, fmt.Sprintf("thumb-%v", id))
		if err := os.Remove(imageDirThumb); err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusExpectationFailed)
			json.NewEncoder(w).Encode(map[string]bool{"on deleting thumbnail": false})
			return
		}
	}

	json.NewEncoder(w).Encode(map[string]bool{"success": true})

}
