package main

import (
	"fmt"
	"time"
	"./utils"
)

func buildQuerys() (map[string]string, map[string]string) {

	// estas query's serão executadas somente uma vez
	querysOneTime := map[string]string{
		"Database_Identifier":"select 'Banco de dados: '||db.name||', instancia: '|| ins.instance_name|| ', criado dia: '||db.created from v$database db, v$instance ins;",
		"Database_Status":"select 'Status: '||ins.status||' '||db.OPEN_MODE||', iniciado desde: '|| ins.startup_time|| ', modo atual: '||db.LOG_MODE from v$database db, v$instance ins;",
		"Tablespaces_Size":"set lines 300 pages 200 \n col SPACE_USED format a10 JUS R WRA \n col SPACE_ALLOCATED format a15 \n col TOTAL_SIZE format a10 \n col FREE_PERCENT format a12 \n col TABLESPACE_NAME format a25 \n select  RPAD('Tablespace: '||a.tablespace_name,40)||RPAD(' Utilizando: '|| round((tbs_size-a.free_space),2)||' GB',23)||RPAD(' Alocado: '||round(tbs_size,2)||' GB',20)||RPAD('Total: '||round(tbs_max_size)||' GB',17)||round(((tbs_max_size-(tbs_size-a.free_space))/tbs_max_size)*100)||'% Livre' from  (select tablespace_name, round(sum(bytes)/1024/1024/1024 ,2) as free_space from dba_free_space group by tablespace_name) a, (select tablespace_name, sum(bytes)/1024/1024/1024 as tbs_size, sum(maxbytes)/1024/1024/1024 as tbs_max_size from dba_data_files group by tablespace_name) b where a.tablespace_name(+)=b.tablespace_name;",
		"Biggest_System_Tables":"col owner format a15 \n col segment_name format a30 \n col segment_type format a15 \n col mb format 999,999,999 \n select 'Tablespace: '||rpad(tablespace_name,7)||'  Tipo: '||rpad(segment_type,12)||rpad(owner||'.'||segment_name,30)|| mb ||' MB' from( select tablespace_name, owner, segment_name, segment_type, bytes/1024/1024 MB from dba_segments where tablespace_name in ('SYSTEM','SYS') order by bytes desc) where rownum < 11;",
		"Buffer_Cache_Hit_Ratio_Percent" : "SELECT ROUND((1 - (phy.value / (cur.value + con.value))) * 100) FROM v$sysstat cur, v$sysstat con, v$sysstat phy WHERE cur.name = 'db block gets' AND con.name = 'consistent gets' AND phy.name = 'physical reads';",
		"Library_Cache_Hit_Ration_Percent" : "SELECT ROUND((sum(pinhits) / sum(pins))*100) FROM v$librarycache WHERE namespace in ('SQL AREA', 'TABLE/PROCEDURE', 'BODY', 'TRIGGER');",

	}

	// estas query's serão executadas a cada 30 segundos durante a analise
	querysMoreTimes := map[string]string{
		"How_Much_Locks" : "SELECT count(*) FROM gv$lock WHERE block = 1 AND request > 0;",
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
	fmt.Print("Favor informar o tempo da analise em minutos: ")
	fmt.Scanln(&tempoAnalise)


	// Coletando horario de inicio e criando arquivo de log
	start := time.Now()
	relatorio := utils.GetFileWriter(start)
	defer relatorio.Close()
	utils.WriteHeader(relatorio)


	utils.WriteString(relatorio,"Horario inicial da Analise: "+start.Format("02/01/2006 15:04:05\n"))

	// maps e vetores com query's e resultados
	querysOneTime, querysMoreTimes := buildQuerys()

	// Querys executadas somente uma vez
	utils.WriteString(relatorio,"\n------------------------[ Informações sobre o Banco de Dados ]------------------------\n")
	utils.QueryOneTime(querysOneTime,user,password,sid,relatorio)


	// Query's executadas varias vezes durante a analise (para calcular min, max e média)
	utils.WriteString(relatorio,"\n-------------------------------[ Relatório da Analise ]-------------------------------\n")
	utils.QueryMoreTimes(querysMoreTimes,user,password,sid,relatorio,start,tempoAnalise)

	utils.WriteString(relatorio,"\nHorario final da Analise: "+time.Now().Format("02/01/2006 15:04:05"))

}
