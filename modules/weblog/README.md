# web_log

## Supported Log Format

### Apache
```apache
LogFormat    "%h %l %u %t \"%r\" %>s %b" common
LogFormat    "%h %l %u %t ¥"%r¥" %>s %b \"%{Referer}i\" \"%{User-Agent}i\"" combined
LogFormat    "%h %l %u %t \"%r\" %>s %b \"%{Referer}i\" \"%{User-Agent}i\" %{cookie}n %D" custom1
LogFormat    "%h %l %u %t \"%r\" %>s %O %I %D" costom2
LogFormat "%v %h %l %u %t \"%r\" %>s %b" vhost_common
LogFormat "%v %h %l %u %t ¥"%r¥" %>s %b \"%{Referer}i\" \"%{User-Agent}i\"" vhost_combined
LogFormat "%v %h %l %u %t \"%r\" %>s %b \"%{Referer}i\" \"%{User-Agent}i\" %{cookie}n %D" vhost_custom1
LogFormat "%v %h %l %u %t \"%r\" %>s %O %I %D" costom2
```

### Nginx
```nginx
log_format combined '$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"';
log_format custom1  '$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent $request_length $request_time "$upstream_response_time" "$http_referer" "$http_user_agent"';
log_format custom2  '$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" $uid_got $request_time';
```

### Examples

|format  |remote_addr |logname|remote_user|date                 |TZ    |request                 |status|bytes_sent|referer |User-Agent|
|--------|------------|:-----:|:---------:|---------------------|------|------------------------|------|----------|--------|----------|
|Index   |0           |1      |2          |3                    |4     |5                       |6     |7         |8       |9         |
|common  |64.242.88.10|-      |-          |[07/Mar/2004:16:47:12|-0800]|GET /robots.txt HTTP/1.1|200   |68        |        |          |
|combined|64.242.88.10|-      |-          |[07/Mar/2004:16:47:12|-0800]|GET /robots.txt HTTP/1.1|200   |68        |<refer> |<UA>      |
