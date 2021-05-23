package routes

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"enstrurent.com/server/db"
	"enstrurent.com/server/flags"
	"github.com/dgrijalva/jwt-go"
	"github.com/disintegration/imaging"
	"golang.org/x/crypto/bcrypt"
)

const PathSeparator = string(os.PathSeparator)

var imageFolderDir = fmt.Sprintf("assets%vimages", PathSeparator)

func HashPassword(password string) string {
	val, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return ""
	}
	return string(val)
}

func CompareHashAndPassword(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getDB(r *http.Request) *db.DBHandle {
	return r.Context().Value(DBContext).(*db.DBHandle)
}

func generateToken(email string, role string, expires time.Duration) (token string, err error) {
	claims := jwt.MapClaims{}
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(expires).Unix()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(flags.GetConfig().JWT_KEY))
}

func checkToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(flags.GetConfig().JWT_KEY), nil
	})
}

func saveImageLocal(r *http.Request, isThumbnail bool) ([]string, error) {
	// Get Image from form.
	err := r.ParseMultipartForm(200000)
	if err != nil {
		return nil, err
	}
	formdata := r.MultipartForm
	files := formdata.File["images"] // grab the filenames
	var imageNames []string
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()
		// Get file extension and validate the file is image.
		if fileHeader.Size == 0 {
			return nil, errors.New("there's no photo or corrupted file")
		}
		image, err := ioutil.TempFile(imageFolderDir, fmt.Sprintf("p-*.%v", extractExtension(fileHeader.Filename)))

		if err != nil {
			return nil, err
		}
		defer image.Close()
		// read all of the contents of our uploaded file into a byte array
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		// write this byte array to the file
		image.Write(fileBytes)

		if isThumbnail {
			if err = createThumbnail(image.Name()); err != nil {
				return nil, err
			}
		}
		imageNames = append(imageNames, extractImageName(image.Name())) // collect Image names
	}

	return imageNames, nil
}

func createThumbnail(dir string) error {
	filename := extractImageName(dir)
	img, err := imaging.Open(dir)

	if err != nil {
		fmt.Printf("Error on createThumbnail:  %v \n", err.Error())
		return err
	}

	resizeImg := imaging.Resize(img, 300, 0, imaging.Lanczos)
	saveLocation := fmt.Sprintf("%v\\%v", imageFolderDir, fmt.Sprintf("thumb-%v", filename))
	imaging.Save(resizeImg, saveLocation)

	return nil
}

func extractImageName(val string) string {
	return strings.Split(val, PathSeparator)[2]
}

func extractExtension(val string) string {
	splitted := strings.Split(val, ".")
	return splitted[len(splitted)-1]
}
