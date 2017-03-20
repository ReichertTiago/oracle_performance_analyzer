package utils

import (
	"strings"
	"time"
	"os"
	"strconv"
)

func QueryOneTime(querysOneTime map[string]string, user string,password string, sid string, relatorio *os.File) {

	// Query's executadas somente uma veze durante a analise (valores estáticos)
	for k, v := range querysOneTime {
		r := RunSqlplus(user,password,sid,v)
		if _, err := strconv.ParseInt(r, 10, 64); err == nil {
			if strings.Contains(k,"Percent"){
				WriteString(relatorio,k+": \t"+r+"%%\n")
			}else{
				WriteString(relatorio,k+": \t"+r+"\n")
			}
		}else{
			WriteString(relatorio,k+":\n\t"+strings.Replace(strings.Replace(r, "%", "%%", -1),"\n", "\n\t", -1)+"\n")
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func QueryMoreTimes(queryMoreTimes map[string]string, user string,password string, sid string, relatorio *os.File, start time.Time, tempoAnalise float64) {

	resultados := map[string][]string{}

	for((time.Since(start).Minutes()) < tempoAnalise){

		for k, v := range queryMoreTimes {
			r := RunSqlplus(user,password,sid,v)
			resultados[k] = Extend(resultados[k], string(r))
			time.Sleep(100 * time.Millisecond)
		}
		time.Sleep(10000 * time.Millisecond)
	}
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
		WriteString(relatorio,k+": \t\tMenor: "+FloatToString(smallest,2)+"\tMaior: "+FloatToString(biggest,2)+"\tMedia: "+FloatToString(avg,2))
	}
}