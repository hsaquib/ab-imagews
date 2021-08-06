package service

import (
	"fmt"
	rLog "github.com/hsaquib/rest-log"
	"gopkg.in/gographics/imagick.v3/imagick"
)

type imagickWand struct {
	Log rLog.Logger
}

func NewImagickWand(logger rLog.Logger) *imagickWand {
	imagick.Initialize()
	defer imagick.Terminate()
	return &imagickWand{
		Log: logger,
	}
}

func (i *imagickWand) Resize(fileBytes []byte, size uint) ([]byte, error) {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	err := mw.ReadImageBlob(fileBytes)
	if err != nil {
		return nil, err
	}

	filter := imagick.FILTER_BOX
	err = mw.ResizeImage(size, size, filter)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = mw.SetImageCompressionQuality(95)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return mw.GetImageBlob(), nil
}
