package channel

import (	
	"fmt"
	"time"
	
)


var debug bool

type Color string

const (
	ColorBlack  Color = "\u001b[30m"
	ColorRed          = "\u001b[31m"
	ColorGreen        = "\u001b[32m"
	ColorYellow       = "\u001b[33m"
	ColorBlue         = "\u001b[34m"
	ColorReset        = "\u001b[0m"
)




func Log(message string) {	 	
	
	start := time.Now()
	fmt.Println(string(ColorGreen), start .Format("2006/01/01  15:04:05"),"\t", string(ColorBlue),  string(ColorYellow), message,string(ColorReset) )
}
