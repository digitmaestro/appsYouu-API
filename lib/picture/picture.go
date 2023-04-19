package picture

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/lib/file"
)

// FromID returns a new image object for graphql usage
func FromID(id *string) *gqlgen.Image {
	if id == nil {
		return nil
	}

	// temporary workaround because prisma can't set fields to null, so "" is in the DB instead
	if *id == "" {
		return nil
	}

	return &gqlgen.Image{
		ID:  *id,
		URL: path.Join(file.BasePath, *id),
	}
}

func IDFromFileName(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func LoadImage(path string) (image.Image, error) {
	imageData, err := imaging.Open(path, imaging.AutoOrientation(true))

	if err != nil {
		return nil, err
	}

	return imageData, nil
}

var IconSizes = []uint{16, 32, 64, 152, 192, 512}

func CreatePwaIconSizes(iconID *string) error {
	if iconID == nil {
		return nil
	}

	iconImagePath := file.Path(*iconID)
	iconImage, err := LoadImage(iconImagePath)

	if err != nil {
		return err
	}

	for _, size := range IconSizes {
		resizedIcon := imaging.Fill(iconImage, int(size), int(size), imaging.Center, imaging.Lanczos)

		resizedIconFile, err := os.Create(file.Path(fmt.Sprintf("%s_%v.png", IDFromFileName(*iconID), size)))
		if err != nil {
			return err
		}
		defer resizedIconFile.Close()

		err = png.Encode(resizedIconFile, resizedIcon)
		if err != nil {
			return err
		}
	}

	return nil
}
