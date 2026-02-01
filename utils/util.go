package utils

import "io/fs"

const AppName = "Okapi Web Framework Example"

var AppVersion = "1.0"

func Must(fsys fs.FS, err error) fs.FS {
	if err != nil {
		panic(err)
	}
	return fsys
}
