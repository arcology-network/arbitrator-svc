module github.com/arcology/arbitrator-svc

go 1.13

require (
	github.com/arcology/3rd-party v0.9.2-0.20210626004852-924da2642860
	github.com/arcology/common-lib v0.9.2-0.20210825054709-eadab62563f0
	github.com/arcology/component-lib v0.9.2-0.20210723093653-fdb52426317c
	github.com/arcology/concurrenturl v0.0.0-20210825054146-c09d6c5ad20e
	github.com/arcology/evm v1.10.4-0.20210723080918-610ef3636717
	github.com/arcology/urlarbitrator-engine v0.0.0-20210817183548-1e0e0209734d
	github.com/arcology/vm-adaptor v0.9.2-0.20210825060711-b25f2fd79bf0
	github.com/prometheus/client_golang v1.11.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	go.uber.org/zap v1.15.0
)

// replace github.com/arcology/concurrenturl => ../concurrenturl/

// replace github.com/arcology/concurrentlib => ../concurrentlib/

// replace github.com/arcology/vm-adaptor => ../vm-adaptor/

// replace github.com/arcology/urlarbitrator-engine => ../urlarbitrator-engine/

// replace github.com/arcology/common-lib => ../common-lib/
