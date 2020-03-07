package main

import "os"

func FileInfoHook(path string, fileInfo os.FileInfo) () {
	if fileInfo.Mode().IsRegular() {
		Log.Debug("[R] ", path)
	} else if fileInfo.Mode().IsDir() {
		Log.Debug("[D] ", path)
	} else if fileInfo.Mode() & os.ModeSymlink != 0 {
		Log.Debug("[L] ", path)
	} else {
		Log.Warn("[U] ", path)
	}

}
