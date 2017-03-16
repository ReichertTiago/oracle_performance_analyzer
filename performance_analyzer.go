package main

import (
"os/exec"
"io"
"log"
"fmt"
"strings"
"time"
"os"
"path/filepath"
"strconv"
)

func buildQuerys() (map[string]string, map[string]string) {

	querysOneTime := map[string]string{
		"Buffer_Cache_Hit_Ratio_Percent" : "SELECT ROUND((1 - (phy.value / (cur.value + con.value))) * 100) FROM v$sysstat cur, v$sysstat con, v$sysstat phy WHERE cur.name = 'db block gets' AND con.name = 'consistent gets' AND phy.name = 'physical reads';",
		"Library_Cache_Hit_Ration_Percent" : "SELECT ROUND((sum(pinhits) / sum(pins))*100) FROM v$librarycache WHERE namespace in ('SQL AREA', 'TABLE/PROCEDURE', 'BODY', 'TRIGGER');",
	}

	querysMoreTimes := map[string]string{
		"How_Much_Current_Locks" : "SELECT count(*) FROM gv$lock WHERE block = 1 AND request > 0;",

	}

	return querysOneTime, querysMoreTimes
}


func main() {
	var user string
	var password string
	var sid string
	var tempoAnalise float64


	// Lendo variaveis a serem informadas pelo usuario
	fmt.Print("Favor informar o USUARIO do banco de dados: [system] ")
	fmt.Scanln(&user)
	if len(user) == 0{
		user = "system"
	}
	fmt.Print("Favor informar a SENHA do banco de dados: [oracle] ")
	fmt.Scanln(&password)
	if len(password) == 0{
		password = "oracle"
	}

	fmt.Print("Favor informar o SID do banco de dados: [orcl] ")
	fmt.Scanln(&sid)
	if len(sid) == 0{
		sid = "orcl"
	}
	fmt.Print("Favor informar o tempo deseja para a analise em minutos: ")
	fmt.Scanln(&tempoAnalise)

	//user = "system"
	//password = "oracle"
	//sid = "reichert"
	//tempoAnalise = 1


	// Coletando horario de inicio e criando arquivo de log
	start := time.Now()
	relatorio := getFileWriter(start)
	defer relatorio.Close()
	writeHeader(relatorio)
	writeString(relatorio,"Horario inicial da Analise: "+start.Format("02/01/2006 15:04:05\n"))


	// maps e vetores com query's e resultados
	querysOneTime, querysMoreTimes := buildQuerys()
	resultados := map[string][]string{}


	// Query's executadas somente uma veze durante a analise (valores estáticos) TODO
	for k, v := range querysOneTime {
		r := run_sqlplus(user,password,sid,v)
		if _, err := strconv.ParseInt(r, 10, 64); err == nil {
			if strings.Contains(k,"Percent"){
				writeString(relatorio,k+": \t"+r+"%%")
			}else{
				writeString(relatorio,k+": \t"+r)
			}
		}else{
			writeString(relatorio,k+"\n"+string(r))
		}

		time.Sleep(100 * time.Millisecond)
	}


	// Query's executadas varias vezes durante a analise (para calcular min, max e média)
	for((time.Since(start).Minutes()) < tempoAnalise){

		for k, v := range querysMoreTimes {
			r := run_sqlplus(user,password,sid,v)
			resultados[k] = Extend(resultados[k], string(r))
			time.Sleep(100 * time.Millisecond)
		}
		time.Sleep(10000 * time.Millisecond)
	}

	//fmt.Println(resultados)

	writeString(relatorio,"\nHorario final da Analise: "+time.Now().Format("02/01/2006 15:04:05"))

}

func run_sqlplus(user, password, sid, query string) string {
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

func Extend(slice []string, element string) []string {
	n := len(slice)
	if n == cap(slice) {
		// Slice is full; must grow.
		// We double its size and add 1, so if the size is zero we still grow.
		newSlice := make([]string, len(slice), 2*len(slice)+1)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0 : n+1]
	slice[n] = element
	return slice
}


func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func getFileWriter(start time.Time) *os.File {
	f, err := os.Create(getCurrentDirectory()+"/LB2_Oracle_Analyzer_"+start.Format("20060102-15H04M"))
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func writeHeader(relatorio *os.File){
	writeString(relatorio,"------------------------------------------------------------------------------------")
	writeString(relatorio,"-               LB2 Consultoria - Leading Business 2 the Next Level!               -")
	writeString(relatorio,"-                                                                                  -")
	writeString(relatorio,"-             Autor: Tiago M Reichert   Email: tiago.miguel@lb2.com.br             -")
	writeString(relatorio,"-                                                                                  -")
	writeString(relatorio,"-                      Oracle Database Performance Analyzer                        -")
	writeString(relatorio,"------------------------------------------------------------------------------------")
}

func writeString(f *os.File, str string){

	f.WriteString(strings.Replace(str, "%%", "%", -1) +"\n")
	fmt.Printf(str+"\n")
}