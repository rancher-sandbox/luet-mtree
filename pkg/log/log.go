package log

import (
	defaultLog "log"
	"os"
)

var log defaultLog.Logger

func init(){
	file, err := os.OpenFile("/tmp/luet-mtree.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
}

func Log(text string){
	log.Println(text)
}
