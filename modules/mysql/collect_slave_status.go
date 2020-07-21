package mysql

import (
	"strconv"
	"strings"
)

const querySlaveStatus = "SHOW SLAVE STATUS"

/*
mysql> SHOW SLAVE STATUS\G
*************************** 1. row ***************************
               Slave_IO_State: Waiting for master to send event
                  Master_Host: localhost
                  Master_User: repl
                  Master_Port: 13000
                Connect_Retry: 60
              Master_Log_File: source-bin.000002
          Read_Master_Log_Pos: 1307
               Relay_Log_File: replica-relay-bin.000003
                Relay_Log_Pos: 1508
        Relay_Master_Log_File: source-bin.000002
             Slave_IO_Running: Yes
            Slave_SQL_Running: Yes
              Replicate_Do_DB:
          Replicate_Ignore_DB:
           Replicate_Do_Table:
       Replicate_Ignore_Table:
      Replicate_Wild_Do_Table:
  Replicate_Wild_Ignore_Table:
                   Last_Errno: 0
                   Last_Error:
                 Skip_Counter: 0
          Exec_Master_Log_Pos: 1307
              Relay_Log_Space: 1858
              Until_Condition: None
               Until_Log_File:
                Until_Log_Pos: 0
           Master_SSL_Allowed: No
           Master_SSL_CA_File:
           Master_SSL_CA_Path:
              Master_SSL_Cert:
            Master_SSL_Cipher:
               Master_SSL_Key:
        Seconds_Behind_Master: 0
Master_SSL_Verify_Server_Cert: No
                Last_IO_Errno: 0
                Last_IO_Error:
               Last_SQL_Errno: 0
               Last_SQL_Error:
  Replicate_Ignore_Server_Ids:
             Master_Server_Id: 1
                  Master_UUID: 3e11fa47-71ca-11e1-9e33-c80aa9429562
             Master_Info_File:
                    SQL_Delay: 0
          SQL_Remaining_Delay: NULL
      Slave_SQL_Running_State: Reading event from the relay log
           Master_Retry_Count: 10
                  Master_Bind:
      Last_IO_Error_Timestamp:
     Last_SQL_Error_Timestamp:
               Master_SSL_Crl:
           Master_SSL_Crlpath:
           Retrieved_Gtid_Set: 3e11fa47-71ca-11e1-9e33-c80aa9429562:1-5
            Executed_Gtid_Set: 3e11fa47-71ca-11e1-9e33-c80aa9429562:1-5
                Auto_Position: 1
         Replicate_Rewrite_DB:
                 Channel_name:
           Master_TLS_Version: TLSv1.2
       Master_public_key_path: public_key.pem
        Get_master_public_key: 0
*/

var slaveStatusMetrics = []string{
	"Seconds_Behind_Master",
	"Slave_SQL_Running",
	"Slave_IO_Running",
}

func (m *MySQL) collectSlaveStatus(collected map[string]int64) error {
	// https://dev.mysql.com/doc/refman/18.0/en/show-slave-status.html
	// https://dev.mysql.com/doc/refman/8.0/en/replication-channels.html
	// https://github.com/gdaws/mysql-slave-status
	rows, err := m.db.Query(querySlaveStatus)
	if err != nil {
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	values := nullStringsFromColumns(columns)

	for rows.Next() {
		if err := rows.Scan(values...); err != nil {
			return err
		}

		set := rowAsMap(columns, values)
		channel := set["Channel_Name"]
		prefix := slaveMetricPrefix(channel)

		if !m.collectedChannels[channel] {
			m.collectedChannels[channel] = true
			m.addSlaveReplicationChannelCharts(channel)
		}

		for _, name := range slaveStatusMetrics {
			strValue, ok := set[name]
			if !ok {
				continue
			}
			value, err := parseSlaveStatusValue(name, strValue)
			if err != nil {
				continue
			}
			collected[prefix+strings.ToLower(name)] = value
		}
	}
	return rows.Err()
}

func parseSlaveStatusValue(name, value string) (int64, error) {
	switch name {
	case "Slave_SQL_Running", "Slave_IO_Running":
		value = convertSlaveRunning(value)
	}
	return strconv.ParseInt(value, 10, 64)
}

func convertSlaveRunning(val string) string {
	if val == "Yes" {
		return "1"
	}
	return "0"
}

func slaveMetricPrefix(channel string) string {
	if channel == "" {
		return ""
	}
	return channel + "_"
}
