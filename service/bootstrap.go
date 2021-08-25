package service

import (
	"net/http"

	"github.com/HPISTechnologies/arbitrator-svc/service/workers"
	"github.com/HPISTechnologies/component-lib/actor"
	"github.com/HPISTechnologies/component-lib/streamer"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type KafkaDownloaderCreator func(concurrency int, groupid string, topics, messageTypes []string, mqaddr string) actor.IWorker
type RPCServerCreator func(serviceAddr, basepath string, zkAddrs []string, rcvrs, fns []interface{})

type Config struct {
	concurrency        int
	groupid            string
	topicMsgExch       string
	topicInclusiveTxs  string
	topicAccessRecords string
	kafka1             string
	kafka2             string
	localIP            string
	zkURL              string
	kdc                KafkaDownloaderCreator
	rsc                RPCServerCreator
	downloader1        actor.IWorker
	downloader2        actor.IWorker
	openPrometheus     bool
}

//return a Subscriber struct
func NewConfig(
	concurrency int,
	topicMsgExch string,
	topicInclusiveTxs string,
	topicAccessRecords string,
	kafka1 string,
	kafka2 string,
	localIP string,
	zkURL string,
	kdc KafkaDownloaderCreator,
	rsc RPCServerCreator,
) *Config {
	return &Config{
		concurrency:        concurrency,
		groupid:            "arbitrator",
		topicMsgExch:       topicMsgExch,
		topicInclusiveTxs:  topicInclusiveTxs,
		topicAccessRecords: topicAccessRecords,
		kafka1:             kafka1,
		kafka2:             kafka2,
		localIP:            localIP,
		zkURL:              zkURL,
		kdc:                kdc,
		rsc:                rsc,
		openPrometheus:     true,
	}
}

func (cfg *Config) Start() {

	// logrus.SetLevel(logrus.DebugLevel)
	// logrus.SetFormatter(&logrus.JSONFormatter{})

	if cfg.openPrometheus {
		http.Handle("/streamer", promhttp.Handler())
		go http.ListenAndServe(":29001", nil)
	}

	broker := streamer.NewStatefulStreamer()
	//00 initializer
	initializer := actor.NewActor(
		"initializer",
		broker,
		[]string{actor.MsgStarting},
		[]string{
			actor.MsgStartSub,
		},
		[]int{1},
		workers.NewInitializer(cfg.concurrency, cfg.groupid),
	)
	initializer.Connect(streamer.NewDisjunctions(initializer, 1))

	receiveMseeages := []string{
		actor.MsgInclusive,
		actor.MsgAppHash,
	}

	receiveTopics := []string{
		cfg.topicMsgExch,
		cfg.topicInclusiveTxs,
	}
	//01 kafkaDownloader
	cfg.downloader1 = cfg.kdc(cfg.concurrency, cfg.groupid, receiveTopics, receiveMseeages, cfg.kafka1)
	kafkaDownloader := actor.NewActor(
		"kafkaDownloader",
		broker,
		[]string{actor.MsgStartSub},
		receiveMseeages,
		[]int{1, 1},
		cfg.downloader1,
	)
	kafkaDownloader.Connect(streamer.NewDisjunctions(kafkaDownloader, 2))

	//01-01 kafkaDownloader
	cfg.downloader2 = cfg.kdc(cfg.concurrency, cfg.groupid, []string{cfg.topicAccessRecords}, []string{actor.MsgTxAccessRecords}, cfg.kafka2)
	kafkaDownloaderEu := actor.NewActor(
		"kafkaDownloaderEu",
		broker,
		[]string{actor.MsgStartSub},
		[]string{actor.MsgTxAccessRecords},
		[]int{100},
		cfg.downloader2,
	)
	kafkaDownloaderEu.Connect(streamer.NewDisjunctions(kafkaDownloaderEu, 100))

	// EuResult Pre-processor.
	euResultPreProcessor := actor.NewActor(
		"euResultPreProcessor",
		broker,
		[]string{actor.MsgTxAccessRecords},
		[]string{actor.MsgPreProcessedEuResults},
		[]int{100},
		workers.NewEuResultPreProcessor(cfg.concurrency, cfg.groupid),
	)
	euResultPreProcessor.Connect(streamer.NewDisjunctions(euResultPreProcessor, 100))

	//02 aggre
	euresultAggreSelector := actor.NewActor(
		"euresultAggreSelector",
		broker,
		[]string{
			actor.MsgPreProcessedEuResults,
			actor.MsgReapinglist,
			actor.MsgAppHash,
			actor.MsgInclusive,
		},
		[]string{actor.MsgEuResultSelected},
		[]int{1},
		workers.NewEuResultsAggreSelector(cfg.concurrency, cfg.groupid),
	)
	euresultAggreSelector.Connect(streamer.NewDisjunctions(euresultAggreSelector, 4))

	//03 rpcService
	rpcsvc := workers.NewRpcService(cfg.concurrency, cfg.groupid)
	rpcService := actor.NewActor(
		"rpcService",
		broker,
		[]string{
			actor.MsgStartSub,
			actor.MsgEuResultSelected,
		},
		[]string{actor.MsgReapinglist},
		[]int{1},
		rpcsvc,
	)
	rpcService.Connect(streamer.NewDisjunctions(rpcService, 1))

	//starter
	selfStarter := streamer.NewDefaultProducer("selfStarter", []string{actor.MsgStarting}, []int{1})
	broker.RegisterProducer(selfStarter)
	broker.Serve()

	cfg.rsc(cfg.localIP+":8972", "arbitrator", []string{cfg.zkURL}, []interface{}{rpcsvc}, nil)

	//start signel
	streamerStarting := actor.Message{
		Name:   actor.MsgStarting,
		Height: 0,
		Round:  0,
		Data:   "start",
	}
	broker.Send(actor.MsgStarting, &streamerStarting)
}

func (cfg *Config) Stop() {

}
