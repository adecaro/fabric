package fabric

import "github.com/hyperledger/fabric/core/crypto"

type memberImpl struct {
	*errorHandlerImpl `json:"-"`

	Alias      string
	chain      *chainImpl `json:"-"`
	client     crypto.Client `json:"-"`
	Chaincodes map[string]*chaincodeImpl
}

func (member *memberImpl) GetChaincode(alias string, uri string) Chaincode {
	return member.getChaincodeInternal(alias, uri)
}

func (member *memberImpl) GetChaincodeByAlias(alias string) Chaincode {
	return member.getChaincodeInternalByAlias(alias)
}

func (member *memberImpl) Deploy(alias string, chaincodeURI string) DeployTransaction {
	return member.GetChaincode(alias, chaincodeURI).Deploy()
}

func (member *memberImpl) Invoke(chaincodeURI string, functionName string) InvokeTransaction {
	return member.getChaincodeInternalByAlias(chaincodeURI).Invoke(functionName)
}

func (member *memberImpl) Query(chaincodeURI string, functionName string) QueryTransaction {
	return member.getChaincodeInternalByAlias(chaincodeURI).Query(functionName)
}


func (member *memberImpl) init() error {
	// Load conf

	return nil
}

func (member *memberImpl) getChaincodeInternal(alias string, uri string) Chaincode {
	chaincode, ok := member.Chaincodes[uri]
	if !ok {
		chaincode = &chaincodeImpl{&errorHandlerImpl{}, alias, uri, "", false, member, nil}
		chaincode.SetName(uri)
		member.Chaincodes[alias] =  chaincode
		member.Chaincodes[uri] =  chaincode

		if err := member.chain.registry.addChaincode(chaincode); err != nil {
			member.pushError(err)
		}
	}

	return chaincode
}

func (member *memberImpl) getChaincodeInternalByAlias(alias string) Chaincode {
	return member.getChaincodeInternal(alias, alias)
}
