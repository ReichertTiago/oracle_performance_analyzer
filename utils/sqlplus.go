package utils

import (
	"os/exec"
	"log"
	"io"
	"strings"
)

func RunSqlplus(user, password, sid, query string) string {
	const queryHeader = "set head off \n set feedback off \n set pagesize 999 \n set long 999 \n "

	cmd := exec.Command("sqlplus", "-S", user+"/"+password+"@"+sid)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, queryHeader+query)
	}()



	out, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatal(err)
	}

	return strings.TrimSpace(string(out))
}