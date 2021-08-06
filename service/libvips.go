package service

import (
	"github.com/daddye/vips"
	rLog "github.com/hsaquib/rest-log"
)

type VipsThumb struct {
	Log rLog.Logger
}

func NewVipsThumb(logger rLog.Logger) *VipsThumb {
	return &VipsThumb{
		Log: logger,
	}
}

func (v *VipsThumb) Resize(fileBytes []byte, size uint) ([]byte, error) {

	options := vips.Options{
		Width:        int(size),
		Height:       int(size),
		Crop:         false,
		Extend:       vips.EXTEND_WHITE,
		Interpolator: vips.BILINEAR,
		Gravity:      vips.CENTRE,
		Quality:      95,
	}
	buf, err := vips.Resize(fileBytes, options)
	if err != nil {
		v.Log.Error("", "", err.Error())
		return nil, err
	}
	return buf, nil
}
