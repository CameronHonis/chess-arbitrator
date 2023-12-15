module github.com/CameronHonis/chess-arbitrator

go 1.18

require (
	github.com/CameronHonis/chess v0.0.0-20231104040721-1fa63f099091
	github.com/CameronHonis/log v0.0.0-20231110230333-7c1ee849db4a
	github.com/CameronHonis/set v0.0.0-20231110043107-dace21619137
	github.com/google/uuid v1.4.0
	github.com/gorilla/websocket v1.5.0
	github.com/onsi/ginkgo/v2 v2.13.0
	github.com/onsi/gomega v1.30.0
)

replace github.com/CameronHonis/chess => ../chess

replace github.com/CameronHonis/log => ../log

replace github.com/CameronHonis/set => ../set

replace github.com/CameronHonis/service => ../service

require (
	github.com/CameronHonis/service v0.0.0-20231215050504-e639f9883805 // indirect
	github.com/go-logr/logr v1.3.0 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/pprof v0.0.0-20231101202521-4ca4178f5c7a // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	golang.org/x/tools v0.14.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
