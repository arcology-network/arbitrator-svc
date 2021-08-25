module github.com/HPISTechnologies/arbitrator-svc

go 1.13

require (
	github.com/HPISTechnologies/3rd-party v0.9.2-0.20210626004852-924da2642860
	github.com/HPISTechnologies/common-lib v0.9.2-0.20210825054709-eadab62563f0
	github.com/HPISTechnologies/component-lib v0.9.2-0.20210723093653-fdb52426317c
	github.com/HPISTechnologies/concurrenturl v0.0.0-20210825054146-c09d6c5ad20e
	github.com/HPISTechnologies/evm v1.10.4-0.20210723080918-610ef3636717
	github.com/HPISTechnologies/urlarbitrator-engine v0.0.0-20210817183548-1e0e0209734d
	github.com/HPISTechnologies/vm-adaptor v0.9.2-0.20210825060711-b25f2fd79bf0
	github.com/prometheus/client_golang v1.11.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	go.uber.org/zap v1.15.0
)

// replace github.com/HPISTechnologies/concurrenturl => ../concurrenturl/

// replace github.com/HPISTechnologies/concurrentlib => ../concurrentlib/

// replace github.com/HPISTechnologies/vm-adaptor => ../vm-adaptor/

// replace github.com/HPISTechnologies/urlarbitrator-engine => ../urlarbitrator-engine/

// replace github.com/HPISTechnologies/common-lib => ../common-lib/
