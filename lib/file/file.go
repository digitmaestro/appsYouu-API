package file

import (
	"image"
	"io"
	"log"
	"mime"
	"os"
	"path/filepath"

	"github.com/99designs/gqlgen/graphql"
	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
	"github.com/robojones/iid"
)

// BasePath is the base http URL path
const BasePath = "/api/assets/"

// BaseDir is the os path where images are saved
var BaseDir = "/data/images/"

// Path returns the absolute OS path of a given image ID
func Path(id string) string {
	return filepath.Join(BaseDir, id)
}

const MaxImageWidth = 1400
const MaxImageHeight = 1400

func Upload(upload graphql.Upload, scaleToFit bool) (*string, error) {
	extensions, err := mime.ExtensionsByType(upload.ContentType)
	if err != nil {
		return nil, errors.Wrap(err, "could not create file")
	}

	imageID := iid.New().String()

	imageData, err := imaging.Decode(upload.File, imaging.AutoOrientation(true))

	if err == nil {
		imageID += ".jpg"
		log.Printf("saving file to %s", Path(imageID))

		if scaleToFit {
			imageData = ScaleToFit(imageData)
		}

		imaging.Save(imageData, Path(imageID), imaging.JPEGQuality(80))
	} else {
		imageID += extensions[0]
		log.Printf("saving file to %s", Path(imageID))

		file, err := os.Create(Path(imageID))
		if err != nil {
			return nil, errors.Wrap(err, "could not create file")
		}

		defer file.Close()

		if _, err = io.Copy(file, upload.File); err != nil {
			return nil, errors.Wrap(err, "could not copy to file")
		}
	}

	return &imageID, nil
}

func MaybeUpload(upload *graphql.Upload, scaleToFit bool) (*string, error) {
	if upload == nil {
		return nil, nil
	}

	return Upload(*upload, scaleToFit)
}

func ScaleToFit(imageData image.Image) image.Image {
	bounds := imageData.Bounds()
	if bounds.Dx() > MaxImageWidth || bounds.Dy() > MaxImageHeight {
		resizedImage := imaging.Fit(imageData, int(MaxImageWidth), int(MaxImageHeight), imaging.Lanczos)
		return resizedImage
	} else {
		return imageData
	}
}
