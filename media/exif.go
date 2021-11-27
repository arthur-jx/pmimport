package media

import (
	"fmt"
	"os"
	"pmimport/global"

	exiftool "github.com/barasher/go-exiftool"
	"github.com/rwcarlsen/goexif/exif"
	"go.uber.org/zap"
)

func goexif_getexif(path string) *exif.Exif {
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

func barasher_exif_getexif(path string) (info map[string]interface{}, err error) {
	et, e := exiftool.NewExiftool()
	if e != nil {
		fmt.Printf("Error when intializing: %v\n", err)
		err = e
		return
	}

	defer et.Close()

	fileInfos := et.ExtractMetadata(path)

	for _, fileInfo := range fileInfos {
		if fileInfo.Err != nil {
			fmt.Printf("Error concerning %v: %v\n", fileInfo.File, fileInfo.Err)
			continue
		}

		// for k, v := range fileInfo.Fields {
		// fmt.Printf("[%v] %v\n", k, v)
		// }
		info = fileInfo.Fields
	}
	return
}

func GetExif(path string) (info map[string]interface{}, err error) {
	return barasher_exif_getexif(path)
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

func GetExifInfoString(info map[string]interface{}, key string) (string, bool) {
	v, ok := info[key]
	if ok {
		switch v.(type) {
		case string:
			return v.(string), true
			break
		}
	}

	return "", false
}
