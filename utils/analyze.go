package utils

import (
	"strings"
	"time"
	"os"
	"strconv"
	"fmt"
)

func QueryOneTime(querysOneTime map[string]string, user string,password string, sid string, relatorio *os.File, osAutentication, header bool) {

	for k, v := range querysOneTime {
		r := RunSqlplus(user,password,sid,v,osAutentication, header)
		if !(strings.Contains(r,"ORA-") || strings.Trim(string(r)," ") == "" ){
			if _, err := strconv.ParseFloat(r,  64); err == nil {
				if strings.Contains(k, "Percent") {
					WriteString(relatorio, k+": \t"+r+"%%")
				} else if strings.Contains(k, "GB") {
					WriteString(relatorio, k+": \t"+r+" GB")
				} else{
					WriteString(relatorio, k+": \t"+r+"")
				}
			} else {
				WriteString(relatorio, k+":\n\t"+strings.Replace(strings.Replace(r, "%", "%%", -1), "\n", "\n\t", -1)+"\n")
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func QueryMoreTimes(queryMoreTimes map[string]string, user string,password string, sid string, relatorio *os.File, osAutentication bool, start time.Time, tempoAnalise float64) {

	resultados := map[string][]string{}
	fmt.Printf("\n")

	for((time.Since(start).Minutes()) < tempoAnalise){

		fmt.Printf("\r...Coleta de informacoes em andamento... [%d%%]", int64(((time.Since(start).Minutes()) / tempoAnalise)*100))

		for k, v := range queryMoreTimes {
			r := RunSqlplus(user,password,sid,v,osAutentication, false)
			if !strings.Contains(r,"ORA-") {
				resultados[k] = Extend(resultados[k], string(r))
			}
			time.Sleep(100 * time.Millisecond)
		}
		time.Sleep(30000 * time.Millisecond)
	}

	fmt.Printf("\n\n")
	SaveAnalyze(relatorio, resultados)
}

func SaveAnalyze(relatorio *os.File, resultados map[string][]string){

	for k, v := range resultados {

		var total, smallest, biggest float64 = 0, StringToFloat(v[0]), StringToFloat(v[0])

		for _, x := range v {
			total += StringToFloat(x)
			if StringToFloat(x) > biggest {
				biggest = StringToFloat(x)
			}
			if StringToFloat(x) < smallest {
				smallest = StringToFloat(x)
			}
		}
		avg := total/float64(len(v))
		if strings.Contains(k, "Percent") {
			WriteString(relatorio, k+": \t\tMenor: "+FloatToString(smallest, 0)+"%%\t\tMaior: "+FloatToString(biggest, 0)+"%%\t\tMedia: "+FloatToString(avg, 0)+"%%")
		}  else if strings.Contains(k, "KB/s") {
			WriteString(relatorio, k+": \t\tMenor: "+FloatToString(smallest, 0)+" KB/s\t\tMaior: "+FloatToString(biggest, 0)+" KB/s\t\tMedia: "+FloatToString(avg, 0)+" KB/s")
		} else if strings.Contains(k, "GB") {
			WriteString(relatorio, k+": \t\tMenor: "+FloatToString(smallest, 0)+" GB\t\tMaior: "+FloatToString(biggest, 0)+" GB\t\tMedia: "+FloatToString(avg, 0)+" GB")
		} else if strings.Contains(k, "IOPS") {
			WriteString(relatorio, k+": \t\tMenor: "+FloatToString(smallest, 0)+" IO/s\t\tMaior: "+FloatToString(biggest, 0)+" IO/s\t\tMedia: "+FloatToString(avg, 0)+" IO/s")
		} else if strings.Contains(k, "Per_Sec") {
			WriteString(relatorio, k+": \t\tMenor: "+FloatToString(smallest, 0)+"/s\t\tMaior: "+FloatToString(biggest, 0)+"/s\t\tMedia: "+FloatToString(avg, 0)+"/s")
		} else if strings.Contains(k, "/s") {
			WriteString(relatorio, k+": \t\tMenor: "+FloatToString(smallest, 0)+"/s\t\tMaior: "+FloatToString(biggest, 0)+"/s\t\tMedia: "+FloatToString(avg, 0)+"/s")
		} else if strings.Contains(k, "MB") {
			WriteString(relatorio, k+": \t\tMenor: "+FloatToString(smallest, 0)+" MB\t\tMaior: "+FloatToString(biggest, 0)+" MB\t\tMedia: "+FloatToString(avg, 0)+" MB")
		} else {
			WriteString(relatorio, k+": \t\tMenor: "+FloatToString(smallest, 2)+"\t\tMaior: "+FloatToString(biggest, 2)+"\t\tMedia: "+FloatToString(avg, 2))
		}
	}
}


func WriteHeader(user, password, sid string, relatorio *os.File, osAutentication bool){

	querysOneTime := map[string]string{
		"Database_Identifier":"select 'DATABASE:  '||db.name||'\nINSTANCE:  '|| ins.instance_name|| '\nCREATED :  '||db.created from v$database db, v$instance ins;",
		"Database_Status":"select 'STARTED :  '|| ins.startup_time ||'\nSTATUS  :  '||ins.status||' - '||db.OPEN_MODE|| '\nMODE    :  '||db.LOG_MODE from v$database db, v$instance ins;",
		"Database_Version":"SELECT 'VERSION :  '||version FROM V$INSTANCE;",
		"Database_Edition":"SELECT 'EDITION :  '||edition FROM V$INSTANCE;",
		"Hostname":"SELECT 'HOSTNAME:  '||host_name FROM V$INSTANCE;",
		"LANGUAGE":"select 'LANGUAGE:  '||VALUE from v$nls_parameters where PARAMETER = 'NLS_LANGUAGE';",
		"COUNTRY":"select 'COUNTRY :  '||VALUE from v$nls_parameters where PARAMETER = 'NLS_TERRITORY';",
		"DATE":"select 'DATE    :  '||VALUE from v$nls_parameters where PARAMETER = 'NLS_DATE_FORMAT';",
		"CHARSET":"select 'CHARSET :  '||VALUE from v$nls_parameters where PARAMETER = 'NLS_CHARACTERSET';",
		"NCHARSET":"select 'NCHARSET:  '||VALUE from v$nls_parameters where PARAMETER = 'NLS_NCHAR_CHARACTERSET';",
	}

	WriteString(relatorio,"------------------------------------------------------------------------------------------------------------------")
	WriteString(relatorio,"-                              LB2 Consultoria - Leading Business 2 the Next Level!                              -")
	WriteString(relatorio,"-                                                                                                                -")
	WriteString(relatorio,"-                            Autor: Tiago M Reichert   Email: tiago.miguel@lb2.com.br                            -")
	WriteString(relatorio,"-                                                                                                                -")
	WriteString(relatorio,"-                                     Oracle Database Performance Analyzer                                       -")
	WriteString(relatorio,"------------------------------------------------------------------------------------------------------------------")
	WriteString(relatorio,"\nHorario inicial da Analise: "+time.Now().Format("02/01/2006 15:04:05"))

	WriteString(relatorio,"\n--------------------------------------[ Informacoes sobre o Banco de Dados ]--------------------------------------\n")

	WriteString(relatorio,"\tHOSTNAME:  " +GetHostname())
	WriteString(relatorio,"\tIP      :  " +GetIP()+"\n")


	for _, v := range querysOneTime {
		r := RunSqlplus(user, password, sid, v, osAutentication, false)

		if !strings.Contains(r,"ORA-") {
			WriteString(relatorio, "\t"+(strings.Replace(r, "\n", "\n\t", -1)))
		}
	}
}