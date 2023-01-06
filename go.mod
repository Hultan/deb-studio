module github.com/hultan/deb-studio

go 1.19

replace pault.ag/go/debian v0.12.0 => "/home/per/code/go-debian"

require (
	github.com/google/uuid v1.3.0
	github.com/gotk3/gotk3 v0.6.1
	github.com/hultan/dialog v1.0.1
	pault.ag/go/debian v0.12.0
)

require (
	golang.org/x/crypto v0.0.0-20201016220609-9e8e0b390897 // indirect
	pault.ag/go/topsort v0.0.0-20160530003732-f98d2ad46e1a // indirect
)
