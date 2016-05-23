package fabric

type chaincodeImpl struct {
	*errorHandlerImpl `json:"-"`

	Alias         string
	ChaincodeURI  string
	Name          string
	Deployed      bool
	member        *memberImpl `json:"-"`
	Credential    *credentialImpl
}

func (chaincode *chaincodeImpl) Deploy() DeployTransaction {
	clientSDKLog.Debug("Deploy")
	return &deployTxImpl{&txImpl{&errorHandlerImpl{}, chaincode, chaincode.member.chain.isConfidential(), false, nil, nil}}
}

func (chaincode *chaincodeImpl) Invoke(functionName string) InvokeTransaction {
	clientSDKLog.Debug("Invoke [%s]", functionName)
	return &invokeTxImpl{&txImpl{&errorHandlerImpl{}, chaincode, chaincode.member.chain.isConfidential(), false, nil, nil}, false, functionName, make([]string, 0)}
}

func (chaincode *chaincodeImpl) Query(functionName string) QueryTransaction {
	clientSDKLog.Debug("Query [%s]", functionName)
	return &queryTxImpl{&invokeTxImpl{&txImpl{&errorHandlerImpl{}, chaincode, chaincode.member.chain.isConfidential(), false, nil, nil}, true, functionName, make([]string, 0)}}
}

func (chaincode *chaincodeImpl) BindCredential() Credential {
	if chaincode.Credential == nil {
		cert, err := chaincode.member.client.GetTCertificateHandlerNext()
		if err != nil {
			chaincode.pushError(err)

			return nil
		}
		chaincode.Credential = &credentialImpl{&errorHandlerImpl{}, cert, cert.GetCertificate()}

		// update the registry
		if err = chaincode.member.chain.registry.addChaincode(chaincode); err != nil {
			chaincode.pushError(err)
		}
	}
	return chaincode.Credential
}


func (chaincode *chaincodeImpl) getName() string {
	return chaincode.Name
}

func (chaincode *chaincodeImpl) SetName(name string) {
	chaincode.Name = name
}

func (chaincode *chaincodeImpl) SetDeployed() {
	chaincode.Deployed = true
}

