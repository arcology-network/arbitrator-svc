package service

import (
	tmCommon "github.com/arcology-network/3rd-party/tm/common"
	mainConfig "github.com/arcology-network/component-lib/config"
	"github.com/arcology-network/component-lib/kafka"
	"github.com/arcology-network/component-lib/log"
	"github.com/arcology-network/component-lib/rpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start arbitrator service Daemon",
	RunE:  startCmd,
}

func init() {
	flags := StartCmd.Flags()

	flags.String("mqaddr", "localhost:9092", "host:port of kafka ")
	flags.String("mqaddr2", "localhost:9092", "host:port of kafka ")
	flags.String("mqaddr3", "localhost:9092", "host:port of kafka ")

	//common
	flags.Int("concurrency", 4, "num of threads")
	flags.String("logcfg", "./log.toml", "log conf path")
	//flags.String("arbitrateAddr", "localhost:8972", "arbitrator server address")
	flags.String("msgexch", "msgexch", "topic for receive msg exchange")
	flags.String("accessrecords", "access-records", "topic for received accessrecords")
	flags.String("inclusive-txs", "inclusive-txs", "topic of received txlist")

	flags.Int("nidx", 0, "node index in cluster")
	flags.String("nname", "node1", "node name in cluster")

	flags.String("zkUrl", "127.0.0.1:2181", "url of zookeeper")
	flags.String("localIp", "127.0.0.1", "local ip of server")
	flags.String("maincfg", "./monaco.toml", "main conf path")
}

func startCmd(cmd *cobra.Command, args []string) error {
	mainConfig.InitCfg(viper.GetString("maincfg"))
	log.InitLog("arbitrator.log", viper.GetString("logcfg"), "arbitrator", viper.GetString("nname"), viper.GetInt("nidx"))

	en := NewConfig(
		viper.GetInt("concurrency"),
		viper.GetString("msgexch"),
		viper.GetString("inclusive-txs"),
		viper.GetString("accessrecords"),
		viper.GetString("mqaddr"),
		viper.GetString("mqaddr3"),
		viper.GetString("localIp"),
		viper.GetString("zkUrl"),
		kafka.NewKafkaDownloader,
		rpc.InitZookeeperRpcServer,
	)
	en.Start()

	// Wait forever
	tmCommon.TrapSignal(func() {
		// Cleanup
		en.Stop()
	})

	return nil
}
