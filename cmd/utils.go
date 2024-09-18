package cmd

import "os"


func handleExit(sig chan os.Signal, f func(interface{}), param interface{}) {
	<-sig
 	f(param)
    os.Exit(0)
}
