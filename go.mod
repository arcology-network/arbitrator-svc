module github.com/arcology-network/arbitrator-svc

go 1.13

require (
	github.com/arcology-network/3rd-party v0.9.2-0.20210626004852-924da2642860
	github.com/arcology-network/common-lib v0.9.2-0.20210910023057-e170e0ae1807
	github.com/arcology-network/component-lib v0.9.3-0.20210914004816-85bd308a8467
	github.com/arcology-network/concurrenturl v0.0.0-20210913021258-ef03ce074986
	github.com/arcology-network/evm v1.10.4-0.20210723080918-610ef3636717
	github.com/arcology-network/urlarbitrator-engine v1.1.1-0.20210915075702-800ae97d1071
	github.com/arcology-network/vm-adaptor v0.9.2-0.20210913050346-015f98047606
	github.com/prometheus/client_golang v1.11.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	go.uber.org/zap v1.15.0
)

// 	github.com/arcology-network/common-lib => ../common-lib/
// replace	github.com/arcology-network/component-lib => ../component-lib/
// replace github.com/arcology-network/concurrentlib => ../concurrentlib/
// replace	github.com/arcology-network/concurrenturl => ../concurrenturl/
// replace github.com/arcology-network/urlarbitrator-engine => ../urlarbitrator-engine/
// replace github.com/arcology-network/vm-adaptor => ../vm-adaptor/
