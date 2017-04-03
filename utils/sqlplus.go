package utils

import (
	"os/exec"
	"log"
	"io"
	"strings"
	"fmt"
)

func RunSqlplus(user, password, sid, query string, osAutentication, header bool) string {

	var queryHeader string
	if header {
		queryHeader =  "set feedback off \n set pagesize 999 \n set lines 900 \n set long 999999999 \n column kbytes format 9999999999999999999 \n "
	}else {
		queryHeader = "set head off \n set feedback off \n set pagesize 999 \n set lines 900 \n set long 999999999 \n column kbytes format 9999999999999999999 \n "
	}
	var connStr string

	if !osAutentication{
		connStr = user+"/"+password+"@"+sid
	}else{
		connStr = "/ as sysdba"
	}


	cmd := exec.Command("sqlplus", "-S", connStr)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
		fmt.Printf("ERRO "+err.Error())
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, queryHeader+query)
	}()



	out, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatal(err)
		fmt.Printf("ERRO "+err.Error())
	}

	return strings.TrimSpace(string(out))
}
