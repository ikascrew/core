module github.com/ikascrew/core

go 1.21

require (
	github.com/ikascrew/plugin v0.0.0-20200715234203-87c9c5b19416
	gocv.io/x/gocv v0.38.0
	golang.org/x/net v0.23.0
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
)

require golang.org/x/sys v0.18.0 // indirect

replace github.com/ikascrew/plugin => ../plugin
