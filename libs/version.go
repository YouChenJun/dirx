package libs

import "fmt"

const (
	VERSION = "0.0.2"
	DESC    = "目录扫描器"
	AUTHOR  = "@Chen_Dark"
	BINARY  = "dirx"
)

var LOGDIR = fmt.Sprintf("./%s-log/", BINARY)

var Banner string = fmt.Sprintf(
	`
      _   _               
     | | (_)              
   __| |  _   _ __  __  __
  / _' | | | | '__| \ \/ /
 |(_|  | | | | |     >  <
 \_,___| |_| |_|    /_/\_\
         version:%v by:%v
`, VERSION, AUTHOR)
