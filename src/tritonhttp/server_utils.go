package tritonhttp
import (
	"bufio"
	"strings"
	"os"
	"fmt"
	"time"
	"path/filepath"
	"strconv"
)

/** 
	Load and parse the mime.types file 
**/
func ParseMIME(MIMEPath string) (MIMEMap map[string]string, err error) {
	//panic("todo - ParseMIME")

	MIMEMap = make(map[string]string) 
	file, err := os.Open(MIMEPath)
	if err != nil {
		fmt.Println(err)
	} else {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			elements := strings.Split(line," ")
			MIMEMap[elements[0]] = elements[1]
		}
	}
	file.Close()
	return MIMEMap,err
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}

//get last modified time of the file in <day-name>,<day><month><year><hour>:<minute>:<second><time-offset> format
func getLastModifiedTime(filename string) (string) {
	file, err := os.Stat(filename)
	if err != nil {
		return ""
	}
	mtime := file.ModTime()
	return mtime.Format(time.RFC1123Z)
}

func getContentLength(filename string) (string){
	file, err := os.Stat(filename)
	if err != nil {
		return ""
	}
	return strconv.FormatInt(file.Size(), 10) 
}

func getContentType(filename string, MIMEMap map[string]string) string {
	ext := filepath.Ext(filename)
	if val,ok := MIMEMap[ext];ok{
		return val
	}
	return "application/octet-stream"
}

func getErrorPage(msg string)(error string){
	error = "<html><title>OMS Error</title><head><h1>Oops!You have sent a "+msg+"</h1></head></html>"
	return error
}

