package utils

import (
	"path/filepath"
	"os"
	"log"
	"time"
	"strings"
	"fmt"
	"runtime"
	"net"
)



func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func GetFileWriter(start time.Time, sid string) *os.File {

	f, err := os.Create(GetCurrentDirectory()+"/LB2_Oracle_Analyzer_"+sid+"_"+start.Format("20060102-15H04M.log"))
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func WriteString(f *os.File, str string){

	f.WriteString(strings.Replace(str, "%%", "%", -1) +"\n")
	fmt.Printf(str+"\n")
}

func GetOsVersion() string {
	return runtime.GOOS
}


func GetIP() string {

	ips := ""
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	for _, i := range ifaces {
		if !strings.Contains(i.Name, "vm") {
			addrs, err := i.Addrs()
			if err != nil {
				log.Fatal(err)
			}
			for _, addr := range addrs {
				if !strings.Contains(addr.String(), ":") {
					if !strings.Contains(addr.String(), "127.0.0.1") {
						ips += addr.String()+" "
					}

				}
			}
		}
	}

	return ips
}

func GetHostname() string{
	name, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	return name
}

