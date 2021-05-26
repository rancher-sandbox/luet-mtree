package log

import (
	"fmt"
	defaultLog "log"
	"os"
	"time"
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
	log.Println(fmt.Sprintf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), text))
}
