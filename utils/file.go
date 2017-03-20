package utils

import (
	"path/filepath"
	"os"
	"log"
	"time"
	"strings"
	"fmt"
	"runtime"
)

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func GetFileWriter(start time.Time) *os.File {
	f, err := os.Create(GetCurrentDirectory()+"/LB2_Oracle_Analyzer_"+start.Format("20060102-15H04M"))
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func WriteHeader(relatorio *os.File){
	WriteString(relatorio,"------------------------------------------------------------------------------------")
	WriteString(relatorio,"-               LB2 Consultoria - Leading Business 2 the Next Level!               -")
	WriteString(relatorio,"-                                                                                  -")
	WriteString(relatorio,"-             Autor: Tiago M Reichert   Email: tiago.miguel@lb2.com.br             -")
	WriteString(relatorio,"-                                                                                  -")
	WriteString(relatorio,"-                      Oracle Database Performance Analyzer                        -")
	WriteString(relatorio,"------------------------------------------------------------------------------------")
}

func WriteString(f *os.File, str string){

	f.WriteString(strings.Replace(str, "%%", "%", -1) +"\n")
	fmt.Printf(str+"\n")
}

func CheckOsVersion() string {
	return runtime.GOOS
}