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
log_format custom3  '$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent $request_length $request_time "$upstream_response_time" "$http_referer" "$http_user_agent"';
log_format custom2  '$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" $uid_got $request_time';
```

### Examples


|Index   |0              |1      |2          |3   |4        |5          |6         |7             |8                 |9            |10                    |11           |12        |
|:------:|:-------------:|:-----:|:---------:|:--:|:-------:|:---------:|:--------:|:------------:|:----------------:|:-----------:|:--------------------:|:-----------:|:--------:|
|common  |**remote_addr**|logname|remote_user|time|time_zone|**request**|**status**|**bytes_sent**|                  |             |                      |             |          |
|combined|**remote_addr**|logname|remote_user|time|time_zone|**request**|**status**|**bytes_sent**|refer             |User-Agent   |                      |             |          |
|custom1 |**remote_addr**|logname|remote_user|time|time_zone|**request**|**status**|**bytes_sent**|refer             |User-Agent   |Cookie                |**resp_time**|          |
|custom2 |**remote_addr**|logname|remote_user|time|time_zone|**request**|**status**|**bytes_sent**|**request_length**|**resp_time**|                      |             |          |
|custom3 |**remote_addr**|logname|remote_user|time|time_zone|**request**|**status**|**bytes_sent**|**request_length**|**resp_time**|**upstream_resp_time**|refer        |User-Agent|

* remote_addr: `64.242.88.10`
* logname: `-`
* remote_user: `-`
* time: `[07/Mar/2004:16:47:12`
* time_zone: `+09:00]`
* request: `GET /robots.txt HTTP/1.1`
* status: `200`
* bytes_sent: `56`
* request_length: `32`
* refer: `http://www.example.com`
* User-Agent: `Mozilla/5.0`
* Cookie: `uid=xxxxxx`
* resp_time: `0.05`
* upstream_resp_time: `0.05, 0.03`