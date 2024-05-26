module github.com/oosawy/magichost

go 1.22.2

replace github.com/hashicorp/mdns v1.0.5 => ./_vendor/mdns

require github.com/hashicorp/mdns v1.0.5

require (
	github.com/miekg/dns v1.1.55 // indirect
	golang.org/x/mod v0.8.0 // indirect
	golang.org/x/net v0.14.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/tools v0.6.0 // indirect
)
