module github.com/muety/stromgedacht-exporter

go 1.21

require github.com/jxsl13/stromgedacht v0.2.0

require (
	github.com/go-resty/resty/v2 v2.7.0 // indirect
	golang.org/x/net v0.20.0 // indirect
)

replace github.com/jxsl13/stromgedacht v0.2.0 => github.com/muety/stromgedacht v0.0.0-20240414201316-bc22bf5eee69
