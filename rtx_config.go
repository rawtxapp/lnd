package main

/*
#ifdef SWIG
%newobject InitLnd;
#endif
*/
import "C"
import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lightningnetwork/lnd/channeldb"
)

var (
	channelDB              *channeldb.DB
	shutdownSuccessChannel = make(chan bool)
	fout                   *os.File
	ferr                   *os.File
)

type Shutdown struct{}

//InitLnd initializes lnd, lndHomeDir is coming from host app.
// lndHomeDir could be for example in android /data/user/0/com.rtxwallet/files.
//export InitLnd
func InitLnd(lndHomeDir *C.char) *C.char {
	lndHomeDirString := C.GoString(lndHomeDir)
	// setStdout(lndHomeDirString)
	err := initLnd(lndHomeDirString)
	if err != nil {
		// shutdownStdout()
		return C.CString(err.Error())
	}
	return C.CString("")
}

//export StopLnd
func StopLnd() bool {
	shutdownRequestChannel <- struct{}{}
	success := <-shutdownSuccessChannel
	// shutdownStdout()
	return success
}

//export TestPanic
func TestPanic() {
	panic("Testing panic!")
}

//export StartLnd
func StartLnd() *C.char {
	defer func() {
		if x := recover(); x != nil {
			ltndLog.Errorf("run time panic: %v", x)
		}
	}()
	err := lndMain()
	if err != nil {
		return C.CString(err.Error())
	}
	return C.CString("")
}

func setStdout(lndHomeDir string) {
	fileout := filepath.Join(lndHomeDir, "stdout")
	fout, _ = os.Create(fileout)
	os.Stdout = fout

	fileerr := filepath.Join(lndHomeDir, "stdout")
	ferr, _ = os.Create(fileerr)
	os.Stderr = ferr
}

func shutdownStdout() {
	fout.Close()
	ferr.Close()
}

func initLnd(lndHomeDir string) error {
	setDefaultVars(lndHomeDir)

	lndCfg, err := loadConfig()
	if err != nil {
		fmt.Println(err)
		return err
	}
	cfg = lndCfg
	return nil
}

func setDefaultVars(lndHomeDir string) {
	if lndHomeDir == "" {
		// If lndHomeDir is null, just leave the defaults as is.
		return
	}
	defaultLndDir = lndHomeDir
	defaultConfigFile = filepath.Join(defaultLndDir, defaultConfigFilename)
	defaultDataDir = filepath.Join(defaultLndDir, defaultDataDirname)
	defaultTLSCertPath = filepath.Join(defaultLndDir, defaultTLSCertFilename)
	defaultTLSKeyPath = filepath.Join(defaultLndDir, defaultTLSKeyFilename)
	defaultAdminMacPath = filepath.Join(defaultLndDir, defaultAdminMacFilename)
	defaultReadMacPath = filepath.Join(defaultLndDir, defaultReadMacFilename)
	defaultLogDir = filepath.Join(defaultLndDir, defaultLogDirname)

	defaultBtcdDir = filepath.Join(lndHomeDir, "btcd", "default")
	defaultBtcdRPCCertFile = filepath.Join(defaultBtcdDir, "rpc.cert")

	defaultLtcdDir = filepath.Join(lndHomeDir, "ltcd", "default")
	defaultLtcdRPCCertFile = filepath.Join(defaultLtcdDir, "rpc.cert")

	defaultBitcoindDir = filepath.Join(lndHomeDir, "bitcoin", "default")
	defaultLitecoindDir = filepath.Join(lndHomeDir, "litecoin", "default")
}
