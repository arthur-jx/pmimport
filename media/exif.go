package media

import (
	"fmt"
	"os"
	"pmimport/global"

	"github.com/rwcarlsen/goexif/exif"
	"go.uber.org/zap"
)

func GetExif(path string) *exif.Exif {
	f, err := os.Open(path)
	if err != nil {
		global.LOG.Error("open file", zap.Any("error", err))
		return nil
	}

	defer f.Close()

	// Optionally register camera makenote data parsing - currently Nikon and
	// Canon are supported.
	// exif.RegisterParsers(mknote.All...)

	x, err := exif.Decode(f)
	if err != nil {
		global.LOG.Error("decode exif", zap.Any("error", err))
		return nil
	}

	return x
}

func ShowExit(info *exif.Exif) {
	if info == nil {
		return
	}
	t, _ := info.DateTime()
	fmt.Println("Create Date:", t)
	lat, long, _ := info.LatLong()
	fmt.Println("LatLong    :", lat, long)
	fmt.Println("String     :", info.String())
	json, err := info.MarshalJSON()
	if err == nil {
		fmt.Println("JSON       :", string(json))
	}
}
