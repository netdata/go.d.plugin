package ovpn_status_log

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type openvpnStatus struct {
	totalUsers int
	bytesIn    int
	bytesOut   int
}

type clientInfo struct {
	CommonName    string
	BytesReceived int
	BytesSent     int
}

func (o *VPNStatus) collect() (map[string]int64, error) {
	var err error

	mx := make(map[string]int64)

	conn, err := os.Open(o.StatusPath)
	defer conn.Close()

	if err != nil {
		//TODO: error log
		return nil, err
	}

	reader := bufio.NewReader(conn)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	status, err := parseFile(scanner)
	if err != nil {
		return nil, err
	}

	mx["clients"] = int64(status.totalUsers)
	mx["bytes_in"] = int64(status.bytesIn)
	mx["bytes_out"] = int64(status.bytesOut)
	return mx, nil
}

func parseFile(scanner *bufio.Scanner) (*openvpnStatus, error) {

	if scanner.Scan() {
		line := scanner.Text()
		clients := []clientInfo{}
		if strings.Contains(line, "OpenVPN CLIENT LIST") {
			clients = parseFileV1(scanner)
		} else if strings.Contains(line, "TITLE,Open") {
			clients = parseFilev2(scanner)
		} else {
			return nil,
				fmt.Errorf("The status log file is invalid")
		}
		return &openvpnStatus{
			totalUsers: len(clients),
		}, nil
	}
	return nil,
		fmt.Errorf("The status log file is invalid")
}

func parseFileV1(scanner *bufio.Scanner) []clientInfo {
	checkClientListHeader := func(headers []string) bool {
		var clientListHeaderColumns = [5]string{
			"Common Name",
			"Real Address",
			"Bytes Received",
			"Bytes Sent",
			"Connected Since",
		}
		for i, v := range headers {
			if v != clientListHeaderColumns[i] {
				return false
			}
		}
		return true

	}
	clientListHeader := false

	var clients []clientInfo

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if checkClientListHeader(parts) {
			clientListHeader = true
			continue
		}
		if clientListHeader && len(parts) == 5 {
			in, _ := strconv.Atoi(parts[2])
			out, _ := strconv.Atoi(parts[3])
			c := clientInfo{
				CommonName:    parts[0],
				BytesReceived: in,
				BytesSent:     out,
			}
			clients = append(clients, c)
		} else {
			clientListHeader = false
		}
	}
	return clients
}

func parseFilev2(scanner *bufio.Scanner) []clientInfo {
	var clients []clientInfo

	for scanner.Scan() {
		line := scanner.Text()
		switch parts := strings.Split(line, ","); parts[0] {
		case "HEADER":
		case "END":
			break
		default:
			switch statusType := parts[0]; statusType {
			case "CLIENT_LIST":
				in, _ := strconv.Atoi(parts[5])
				out, _ := strconv.Atoi(parts[6])
				c := clientInfo{
					CommonName:    parts[1],
					BytesReceived: in,
					BytesSent:     out,
				}
				clients = append(clients, c)
			}
		}
	}
	return clients
}

// func (o *VPNStatus) collectLoadStats(mx map[string]int64) error {
// 	stats, err := o.client.LoadStats()
// 	if err != nil {
// 		return err
// 	}

// 	mx["clients"] = stats.NumOfClients
// 	mx["bytes_in"] = stats.BytesIn
// 	mx["bytes_out"] = stats.BytesOut
// 	return nil
// }

// func (o *VPNStatus) collectUsers(mx map[string]int64) error {
// 	users, err := o.client.Users()
// 	if err != nil {
// 		return err
// 	}

// 	now := time.Now().Unix()
// 	var name string

// 	for _, user := range users {
// 		if user.Username == "UNDEF" {
// 			name = user.CommonName
// 		} else {
// 			name = user.Username
// 		}

// 		if !o.perUserMatcher.MatchString(name) {
// 			continue
// 		}
// 		if !o.collectedUsers[name] {
// 			o.collectedUsers[name] = true
// 			if err := o.addUserCharts(name); err != nil {
// 				o.Warning(err)
// 			}
// 		}
// 		mx[name+"_bytes_received"] = user.BytesReceived
// 		mx[name+"_bytes_sent"] = user.BytesSent
// 		mx[name+"_connection_time"] = now - user.ConnectedSince
// 	}
// 	return nil
// }

// func (o *VPNStatus) addUserCharts(userName string) error {
// 	cs := userCharts.Copy()

// 	for _, chart := range *cs {
// 		chart.ID = fmt.Sprintf(chart.ID, userName)
// 		chart.Fam = fmt.Sprintf(chart.Fam, userName)

// 		for _, dim := range chart.Dims {
// 			dim.ID = fmt.Sprintf(dim.ID, userName)
// 		}
// 		chart.MarkNotCreated()
// 	}
// 	return o.charts.Add(*cs...)
// }
