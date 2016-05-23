package fabric

import (
	"fmt"
	"github.com/hyperledger/fabric/core/crypto"
	pb "github.com/hyperledger/fabric/protos"
)

type chainImpl struct {
	*errorHandlerImpl

	Alias    string
	ConfPath string

	nvp      *nvpImpl
	Members  map[string]*memberImpl
	registry *registryImpl
}

func (chain *chainImpl) EnrollMember(alias string, enrollID, enrollPWD string) Member {
	if err := crypto.RegisterClient(alias, nil, enrollID, enrollPWD); err != nil {
		chain.pushError(err)

		return nil
	}

	var err error
	client, err := crypto.InitClient(alias, nil)
	if err != nil {
		chain.pushError(err)

		return nil
	}

	member := &memberImpl{&errorHandlerImpl{}, alias, chain, client, make(map[string]*chaincodeImpl)}
	if err := chain.Members[alias].init(); err != nil {
		chain.pushError(err)

		return nil
	}

	chain.registry.addMember(member)
	chain.Members[alias] = member
	return member
}

func (chain *chainImpl) GetMember(alias string) Member {
	member, ok := chain.Members[alias]
	if !ok {
		chain.pushError(fmt.Errorf("Member not found [%s]", alias))
	}
	return member
}

func (chain *chainImpl) init(alias, confPath string) error {
	// Init registry
	if err := chain.registry.init(chain); err != nil {
		return err
	}

	// Init NVP
	return chain.nvp.init(confPath)
}

func (chain *chainImpl) isDevMode() bool {
	return true
}

func (chain *chainImpl) isConfidential() bool {
	return true
}

func (chain *chainImpl) sendTransaction(tx *pb.Transaction) (*pb.Response, error) {
	return chain.nvp.sendTransaction(tx)
}
