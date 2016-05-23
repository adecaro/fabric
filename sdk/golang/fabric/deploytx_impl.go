package fabric

import (
	"fmt"
	pb "github.com/hyperledger/fabric/protos"
	"github.com/spf13/viper"
	"github.com/hyperledger/fabric/core/chaincode"
	"github.com/hyperledger/fabric/core/container"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/platforms"
)

type deployTxImpl struct {
	*txImpl
}

func (deployTx *deployTxImpl) Send() Transaction {
	// Prepare fabric transaction

	// Prepare chaincode spec
	spec := &pb.ChaincodeSpec{
		Type:    1,
		CtorMsg: &pb.ChaincodeInput{Function: "init", Args: []string{}},
	}

	// Set chaincode ID
	if deployTx.chaincode.member.chain.isDevMode() {
		spec.ChaincodeID = &pb.ChaincodeID{Name: deployTx.chaincode.ChaincodeURI}
	} else {
		spec.ChaincodeID = &pb.ChaincodeID{Path: deployTx.chaincode.ChaincodeURI}
	}

	// Sec confidentiality
	if deployTx.confidential {
		spec.ConfidentialityLevel = pb.ConfidentialityLevel_CONFIDENTIAL
	} else {
		spec.ConfidentialityLevel = pb.ConfidentialityLevel_PUBLIC
	}

	// Set metadata
	if deployTx.withCredential {
		spec.Metadata = []byte(deployTx.chaincode.BindCredential().Raw())
	}

	// First build the deployment spec
	var cds *pb.ChaincodeDeploymentSpec
	var err error

	if deployTx.chaincode.member.chain.isDevMode() {
		cds = &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec, CodePackage: nil}
	} else {
		if cds, err = deployTx.getChaincodeBytes(spec); err != nil {
			deployTx.pushError(fmt.Errorf("Error getting deployment spec: [%s]", err))

			return deployTx
		}
	}

	// Now create the Transactions message and send to Peer.
	fabricTx, err := deployTx.chaincode.member.client.NewChaincodeDeployTransaction(
		cds,
		cds.ChaincodeSpec.ChaincodeID.Name,
	)
	if err != nil {
		deployTx.pushError(fmt.Errorf("Error creating fabric transaction: [%s]", err))

		return deployTx
	}
	deployTx.fabricTx = fabricTx

	deployTx.txImpl.Send()

	deployTx.chaincode.SetName(string(deployTx.GetResponse()))
	deployTx.chaincode.SetDeployed()
	deployTx.chaincode.member.chain.registry.addChaincode(deployTx.chaincode)

	return deployTx
}

func (deployTx *deployTxImpl) getChaincodeBytes(spec *pb.ChaincodeSpec) (*pb.ChaincodeDeploymentSpec, error) {
	mode := viper.GetString("chaincode.mode")
	var codePackageBytes []byte
	if mode != chaincode.DevModeUserRunsChaincode {
		clientSDKLog.Debug("Received build request for chaincode spec: %v", spec)
		var err error
		if err = deployTx.checkSpec(spec); err != nil {
			return nil, err
		}

		codePackageBytes, err = container.GetChaincodePackageBytes(spec)
		if err != nil {
			err = fmt.Errorf("Error getting chaincode package bytes: %s", err)
			clientSDKLog.Error(fmt.Sprintf("%s", err))
			return nil, err
		}
	}
	chaincodeDeploymentSpec := &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec, CodePackage: codePackageBytes}
	return chaincodeDeploymentSpec, nil
}

func (deployTx *deployTxImpl) checkSpec(spec *pb.ChaincodeSpec) error {
	// Don't allow nil value
	if spec == nil {
		return errors.New("Expected chaincode specification, nil received")
	}

	platform, err := platforms.Find(spec.Type)
	if err != nil {
		return fmt.Errorf("Failed to determine platform type: %s", err)
	}

	return platform.ValidateSpec(spec)
}