package weblog

//
//// TODO: add cases
//func Test_patterns(t *testing.T) {
//	var cases = []struct {
//		name    string
//		line    string
//		pattern csvPattern
//	}{
//		{
//			name:    "nginx netdata format",
//			line:    `10.254.254.3 - - [09/Nov/2018:00:36:19 +0900] "GET /cacti HTTP/1.1" 305 44 484 0.120 0.555 "10.254.254.1" "-" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36"`,
//			pattern: logFormatNetdata,
//		},
//		{
//			name:    "nginx default format",
//			line:    `10.254.254.3 - - [09/Nov/2018:00:36:19 +0900] "GET /cacti HTTP/1.1" 305 44 484 "10.254.254.1" "-" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36"`,
//			pattern: logFormatDefault,
//		},
//		{
//			name:    "apache netdata format",
//			line:    `127.0.0.1 - - [20/Dec/2018:00:50:01 +0900] "GET /viewvc.cgi/*docroot*/images/favicon.ico HTTP/1.0" 200 562 506 93496 "http://viewvc.ttk-chita.lan/" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36"`,
//			pattern: logFormatNetdata,
//		},
//		{
//			name:    "apache default format",
//			line:    `127.0.0.1 - - [20/Dec/2018:00:50:01 +0900] "GET /viewvc.cgi/*docroot*/images/favicon.ico HTTP/1.0" 200 562 "http://viewvc.ttk-chita.lan/" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36"`,
//			pattern: logFormatDefault,
//		},
//		{
//			name:    "apache vhost netdata format",
//			line:    `viewvc.ttk-chita.lan:80 127.0.0.1 - - [18/Dec/2018:17:47:28 +0900] "GET /viewvc.cgi/*docroot*/images/log.png HTTP/1.0" 200 719 506 93496 "http://viewvc.ttk-chita.lan/viewvc.cgi/rancid/ttk-current/configs/cta-core-cta101-c7606?revision=1.225&view=markup" "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:64.0) Gecko/20100101 Firefox/64.0"`,
//			pattern: logFormatNetdataVhost,
//		},
//		{
//			name:    "apache vhost default format",
//			line:    `viewvc.ttk-chita.lan:80 127.0.0.1 - - [18/Dec/2018:17:47:28 +0900] "GET /viewvc.cgi/*docroot*/images/log.png HTTP/1.0" 200 719 "http://viewvc.ttk-chita.lan/viewvc.cgi/rancid/ttk-current/configs/cta-core-cta101-c7606?revision=1.225&view=markup" "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:64.0) Gecko/20100101 Firefox/64.0"`,
//			pattern: logFormatDefaultVhost,
//		},
//	}
//
//	for _, c := range cases {
//		_, err := newParser(c.line, c.pattern)
//		assert.NoErrorf(t, err, c.name)
//	}
//}
