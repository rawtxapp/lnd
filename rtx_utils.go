package main

// ***NOTE***
// This file isn't actually useful because the
// methods exposed by this file aren't actually used.
// (maybe using getlndversion makes sense...).
// But it's a good example of exposing different
// kinds of functions (%newobject + C.char, etc.).

/*
#include <stdlib.h>

#ifdef SWIG
%newobject GetEnv;
%newobject GetLndVersion;
#endif

*/
import "C"
import (
	"os"
)

//export GetEnv
func GetEnv(v *C.char) *C.char {
	return C.CString(os.Getenv(C.GoString(v)))
}

//export SetEnv
func SetEnv(key *C.char, val *C.char) {
	os.Setenv(C.GoString(key), C.GoString(val))
}
