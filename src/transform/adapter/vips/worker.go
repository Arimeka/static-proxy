package vips

import (
	"gopkg.in/h2non/bimg.v1"

	"os"
)

type Worker struct {
	OriginFilename string
	ResultFilename string

	Height  int
	Width   int
	Quality int

	Crop bool

	Type string
}

func (w Worker) Process() (*os.File, error) {
	buffer, err := bimg.Read(w.OriginFilename)
	if err != nil {
		return nil, err
	}

	opts := bimg.Options{
		Height:  w.Height,
		Width:   w.Width,
		Quality: w.Quality,
		Gravity: bimg.GravityNorth,
		Enlarge: true,
		Crop:    w.Crop,
	}

	result, err := bimg.NewImage(buffer).Process(opts)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(w.ResultFilename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = file.Write(result)

	return file, err
}
