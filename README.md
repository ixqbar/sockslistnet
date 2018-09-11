
[TOC]

### desc
* 定时抓取 https://sockslist.net/list/proxy-socks-5-list/ 代理并rpush到目标redis服务器对应List缓存key=proxy下

### usage
```
./proxy --config=config.xml
```

### version
```
v0.0.1
```

### config
```
<?xml version="1.0" encoding="UTF-8" ?>
<config>
	<task>
		<startup>true</startup>
		<schedule>0 0 * * * *</schedule>
	</task>
	<redis_server></redis_server>
</config>
```
* schedule 语法可以参考 https://godoc.org/github.com/robfig/cron

### target website
* https://sockslist.net/list/proxy-socks-5-list/

### faq
* 更多疑问请+qq群 233415606
