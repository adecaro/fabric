package fabric

import (
	"github.com/spf13/viper"
	"github.com/hyperledger/fabric/core/peer"
	"fmt"
	"flag"
	"strings"
	"runtime"
	"google.golang.org/grpc"
	pb "github.com/hyperledger/fabric/protos"

	"github.com/op/go-logging"
	"golang.org/x/net/context"
	"github.com/hyperledger/fabric/core/crypto"
)

type nvpImpl struct {
	confPath string

	peerClientConn *grpc.ClientConn
	peerClient     pb.PeerClient
}

func (nvp *nvpImpl) init(confPath string) (err error) {
	nvp.confPath = confPath

	// Init configuration and logging
	if err := crypto.Init(); err != nil {
		fmt.Printf("Failed initializing crypto layer: [%s]", err)
		return err
	}
	nvp.initConfiguration()
	nvp.initLogging()

	// Init Peer Client
	nvp.peerClientConn, err = peer.NewPeerClientConnection()
	if err != nil {
		fmt.Printf("Error connection to server at host:port [%s]: [%s]", viper.GetString("peer.address"), err)
		return
	}
	nvp.peerClient = pb.NewPeerClient(nvp.peerClientConn)

	return
}

func (nvp *nvpImpl) initLogging() {
	var formatter = logging.MustStringFormatter(
		`%{color}[%{module}] %{shortfunc} [%{shortfile}] -> %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	logging.SetFormatter(formatter)

	level, err := logging.LogLevel(viper.GetString("logging.peer"))
	if err == nil {
		// No error, use the setting
		logging.SetLevel(level, "main")
		logging.SetLevel(level, "server")
		logging.SetLevel(level, "peer")
	} else {
		clientSDKLog.Warning("Log level not recognized '%s', defaulting to %s: %s", viper.GetString("logging.peer"), logging.ERROR, err)
		logging.SetLevel(logging.ERROR, "main")
		logging.SetLevel(logging.ERROR, "server")
		logging.SetLevel(logging.ERROR, "peer")
	}
}

func (nvp *nvpImpl) initConfiguration() error {
	flag.Parse()

	// Now set the configuration file
	viper.SetEnvPrefix("HYPERLEDGER")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetConfigName("core")         // name of config file (without extension)
	viper.AddConfigPath(nvp.confPath) // path to look for the config file in
	err := viper.ReadInConfig()              // Find and read the config file
	if err != nil {                          // Handle errors reading the config file
		return fmt.Errorf("Fatal error config file: [%s]", err)
	}

	viper.Set("ledger.blockchain.deploy-system-chaincode", "false")
	viper.Set("peer.validator.validity-period.verification", "false")

	// Set the number of maxprocs
	var numProcsDesired = viper.GetInt("peer.gomaxprocs")
	clientSDKLog.Debug("setting Number of procs to %d, was %d\n", numProcsDesired, runtime.GOMAXPROCS(2))

	return nil
}

func (nvp *nvpImpl) sendTransaction(tx *pb.Transaction) (*pb.Response, error) {
	return nvp.peerClient.ProcessTransaction(context.Background(), tx)
}