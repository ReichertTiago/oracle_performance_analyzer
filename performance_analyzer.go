package main

import (
	"fmt"
	"time"
	"./utils"
	"strings"
)

func buildQuerys() (map[string]string, (map[string]string), map[string]string) {



	// estas query's serão executadas somente uma vez
	querysOneTimeMoreLines := map[string]string{
		"Tablespaces_Size ":"col SPACE_USED format a10 JUS R WRA \n col SPACE_ALLOCATED format a15 \n col TOTAL_SIZE format a10 \n col FREE_PERCENT format a12 \n col TABLESPACE_NAME format a25 \n select  a.tablespace_name as TABLESPACE_NAME, round((tbs_size-a.free_space) ,2)||' GB' SPACE_USED, round(tbs_size,2)||' GB' SPACE_ALLOCATED, round(tbs_max_size)||' GB' as TOTAL_SIZE, round(((tbs_max_size-(tbs_size-a.free_space))/tbs_max_size)*100)||' %' as FREE_PERCENT from  (select tablespace_name, round(sum(bytes)/1024/1024/1024 ,2) as free_space from dba_free_space group by tablespace_name) a, (select tablespace_name, sum(bytes)/1024/1024/1024 as tbs_size, sum(max_size)/1024/1024/1024 as tbs_max_size from (select tablespace_name, bytes, case when maxbytes > bytes then maxbytes else bytes end  max_size from dba_data_files) group by tablespace_name) b where a.tablespace_name(+)=b.tablespace_name;",
		"Biggest_System_Tables ":"col SEGMENT_NAME format a50 \n col TABLESPACE_NAME format a25 \n col TYPE format a15 \n col SEGMENT_SIZE format 999,999,999 \n col SEGMENT_SIZE format a13 \n  select tablespace_name TABLESPACE_NAME, segment_type TYPE, owner||'.'||segment_name SEGMENT_NAME, mb ||' MB' SEGMENT_SIZE from( select tablespace_name, owner, segment_name, segment_type, bytes/1024/1024 MB from dba_segments where tablespace_name in ('SYSTEM','SYSAUX') order by bytes desc) where rownum < 11;",
	}

	querysOneTimeOneLine := map[string]string{
		"Buffer_Cache_Hit_Percent  ":"SELECT ROUND((1 - (phy.value / (cur.value + con.value))) * 100) FROM v$sysstat cur, v$sysstat con, v$sysstat phy WHERE cur.name = 'db block gets' AND con.name = 'consistent gets' AND phy.name = 'physical reads';",
		"Library_Cache_Hit_Percent ":"SELECT ROUND((sum(pinhits) / sum(pins))*100) FROM v$librarycache WHERE namespace in ('SQL AREA', 'TABLE/PROCEDURE', 'BODY', 'TRIGGER');",
		"Tablespace_Used_GB        ":"select round(sum((tbs_size-a.free_space)),3) from (select tablespace_name, sum(bytes)/1024/1024/1024 as free_space from dba_free_space group by tablespace_name) a, (select tablespace_name, sum(bytes)/1024/1024/1024 as tbs_size from dba_data_files group by tablespace_name) b where a.tablespace_name(+)=b.tablespace_name;",
		"Tablespace_Allocated_GB   ":"select round(sum((tbs_size)),3) from (select tablespace_name, sum(bytes)/1024/1024/1024 as free_space from dba_free_space group by tablespace_name) a, (select tablespace_name, sum(bytes)/1024/1024/1024 as tbs_size from dba_data_files group by tablespace_name) b where a.tablespace_name(+)=b.tablespace_name;",
		"SGA_Target_GB             ":"select round(value/1024/1024/1024,2) from v$parameter where name like 'sga_target';",
		"SGA_Max_Size_GB           ":"select round(value/1024/1024/1024,2) from v$parameter where name like 'sga_max_size';",
		"PGA_Max_Size_GB           ":"select round(value/1024/1024/1024,2) from v$parameter where name like 'pga_aggregate_limit';",
		"PGA_Target_GB             ":"select round(value/1024/1024/1024,2) from v$parameter where name like 'pga_aggregate_target';",
		"MEMORY_Max_Size_GB        ":"select round(value/1024/1024/1024,2) from v$parameter where name like 'memory_max_target';",
		"MEMORY_Target_GB          ":"select round(value/1024/1024/1024,2) from v$parameter where name like 'memory_target';",

	}

	// estas query's serão executadas a cada 30 segundos durante a analise
	querysMoreTimes := map[string]string{
		"How_Much_Current_Locks    ":"SELECT count(*) FROM gv$lock WHERE block = 1 AND request > 0;",
		"Physical_Read_KB/s        ":"select  sum(case metric_name when 'Physical Read Total Bytes Per Sec' then value end)/1024 as kbytes from V$SYSMETRIC;",
		"Physical_Read_IOPS        ":"select  sum(case metric_name when 'Physical Read Total IO Requests Per Sec' then value end) as kbytes from V$SYSMETRIC;",
		"Physical_Write_KB/s       ":"select  sum(case metric_name when 'Physical Write Total Bytes Per Sec' then value end)/1024 as kbytes  from V$SYSMETRIC;",
		"Physical_Write_IOPS       ":"select  sum(case metric_name when 'Physical Write Total IO Requests Per Sec' then value end) as kbytesfrom V$SYSMETRIC;",
		"Redo_Generated_KB/s       ":"select  sum(case metric_name when 'Redo Generated Per Sec' then value end)/1024 as kbytes from V$SYSMETRIC;",
		"Redo_Generated_IOPS       ":"select  sum(case metric_name when 'Redo Writes Per Sec' then value end) as kbytes from V$SYSMETRIC;",
		"OS_Load                   ":"select  sum(case metric_name when 'Current OS Load' then value end) as kbytes from V$SYSMETRIC;",
		"DB_CPU_Usage_Per_Second   ":"select  sum(case metric_name when 'CPU Usage Per Sec' then value end) as kbytes from V$SYSMETRIC;",
		"OS_CPU_Utilization_Percent":"select  sum(case metric_name when 'Host CPU Utilization (%)' then value end) as kbytes from V$SYSMETRIC;",
		"Network_Usage_KB/s        ":"select  sum(case metric_name when 'Network Traffic Volume Per Sec' then value end)/1024 as kbytes from V$SYSMETRIC;",
		"DB_Wait_Time_Ratio        ":"select sum(case metric_name when 'Database Wait Time Ratio' then value end) as kbytes from V$SYSMETRIC;",
		"DB_CPU_Time_Ratio         ":"select sum(case metric_name when 'Database CPU Time Ratio' then value end) as kbytes from V$SYSMETRIC;",
		"Temp_Space_Used_MB        ":"select sum(case metric_name when 'Temp Space Used' then value end)/1024/1024 as kbytes from V$SYSMETRIC;",
		"Shared_Pool_Free_Percent  ":"select sum(case metric_name when 'Shared Pool Free %' then value end) as kbytes from V$SYSMETRIC;",
		"User_Commits_Per_Sec      ":"select sum(case metric_name when 'User Commits Per Sec' then value end) as kbytes from V$SYSMETRIC;",
		"User_Rollbacks_Per_Sec    ":"select sum(case metric_name when 'User Rollbacks Per Sec' then value end) as kbytes from V$SYSMETRIC;",
		"Open_Cursors_Count        ":"select sum(case metric_name when 'Current Open Cursors Count' then value end) as kbytes from V$SYSMETRIC;",
		"Enqueue_Deadlocks_Per_Sec ":"select sum(case metric_name when 'Enqueue Deadlocks Per Sec' then value end) as kbytes from V$SYSMETRIC;",
		"Cache_Blocks_Corrupted    ":"select sum(case metric_name when 'Global Cache Blocks Corrupted' then value end) as kbytes from V$SYSMETRIC;",
		"Total_Parse_Count_Per_Sec ":"select sum(case metric_name when 'Total Parse Count Per Sec' then value end) as kbytes from V$SYSMETRIC;",
		"DBWR_Checkpoints_Per_Sec  ":"select sum(case metric_name when 'DBWR Checkpoints Per Sec' then value end) as kbytes from V$SYSMETRIC;",
		"DB_Block_Changes_Per_Sec  ":"select sum(case metric_name when 'DB Block Changes Per Sec' then value end) as kbytes from V$SYSMETRIC;",
		"Temp_Used_MB              ":"SELECT SUM (A.used_blocks * D.block_size)/1024/1024 as mb_used FROM v$sort_segment A, ( SELECT B.name, C.block_size, SUM (C.bytes) / 1024 / 1024 as mb_total FROM v$tablespace B, v$tempfile C WHERE B.ts#= C.ts# GROUP BY B.name, C.block_size ) D WHERE A.tablespace_name = D.name",
		"Temp_Used_Percent         ":"SELECT ROUND((AVG(SUM(A.used_blocks * D.block_size)/1024/1024/D.mb_total)*100)) as percent_used FROM v$sort_segment A, ( SELECT B.name, C.block_size, SUM (C.bytes) / 1024 / 1024 mb_total FROM v$tablespace B, v$tempfile C WHERE B.ts#= C.ts# GROUP BY B.name, C.block_size ) D WHERE A.tablespace_name = D.name GROUP by A.tablespace_name, D.mb_total;",
	}
	return querysOneTimeMoreLines, querysOneTimeOneLine, querysMoreTimes
}


func main() {
	user := ""
	password := ""
	sid := ""
	var tempoAnalise float64
	var osAutentication bool

	// Lendo variaveis a serem informadas pelo usuario
	fmt.Print("Utilizar autenticacao de SO: [n] ")
	var autenticationType string
	fmt.Scanln(&autenticationType)
	if strings.Contains(autenticationType,"y") || strings.Contains(autenticationType,"Y") || strings.Contains(autenticationType,"S") || strings.Contains(autenticationType,"s") {
		osAutentication = true
	}else{
		osAutentication = false
	}

	if !osAutentication {
		fmt.Print("Favor informar o USUARIO do banco de dados: [system] ")
		fmt.Scanln(&user)
		if len(user) == 0 {
			user = "system"
		}
		fmt.Print("Favor informar a SENHA do banco de dados: [oracle] ")
		fmt.Scanln(&password)
		if len(password) == 0 {
			password = "oracle"
		}

		fmt.Print("Favor informar o SID do banco de dados: [orcl] ")
		fmt.Scanln(&sid)
		if len(sid) == 0{
			sid = "orcl"
		}
	}

	fmt.Print("Favor informar o tempo da analise em minutos: ")
	fmt.Scanln(&tempoAnalise)


	// Coletando horario de inicio e criando arquivo de log
	start := time.Now()
	relatorio := utils.GetFileWriter(start, sid)
	defer relatorio.Close()

	// Write header with database info's
	utils.WriteHeader(user,password,sid,relatorio,osAutentication)

	// maps e vetores com query's e resultados
	querysOneTimeMoreLines, querysOneTimeOneLine, querysMoreTimes := buildQuerys()


	utils.WriteString(relatorio,"\n---------------------------------------------[ Relatorio da Analise ]---------------------------------------------\n")

	// Querys executadas somente uma vez
	utils.QueryOneTime(querysOneTimeMoreLines,user,password,sid,relatorio,osAutentication, true)
	utils.QueryOneTime(querysOneTimeOneLine,user,password,sid,relatorio,osAutentication, false)

	// Query's executadas varias vezes durante a analise (para calcular min, max e média)
	utils.QueryMoreTimes(querysMoreTimes,user,password,sid,relatorio,osAutentication,start,tempoAnalise)


	utils.WriteString(relatorio,"\nHorario final da Analise: "+time.Now().Format("02/01/2006 15:04:05"))

}
