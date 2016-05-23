package fabric

import (
	pb "github.com/hyperledger/fabric/protos"
	"github.com/hyperledger/fabric/core/util"
	"github.com/golang/protobuf/proto"
)

type invokeTxImpl struct {
	*txImpl

	query bool
	functionName string
	args []string
}

func (invokeTx *invokeTxImpl) Send() Transaction {
	// Get a transaction handler to be used to submit the execute transaction
	// and bind the chaincode access control logic using the binding
	submittingCertHandler, err := invokeTx.chaincode.member.client.GetTCertificateHandlerNext()
	if err != nil {
		invokeTx.pushError(err)

		return invokeTx
	}
	txHandler, err := submittingCertHandler.GetTransactionHandler()
	if err != nil {
		invokeTx.pushError(err)

		return invokeTx
	}

	clientSDKLog.Debug("Invoke [%s][%s][% x]", invokeTx.query, invokeTx.functionName, invokeTx.args)
	chaincodeInput := &pb.ChaincodeInput{Function: invokeTx.functionName, Args: invokeTx.args}

	// Prepare spec and submit
	spec := &pb.ChaincodeSpec{
		Type:                 1,
		ChaincodeID:          &pb.ChaincodeID{Name: invokeTx.chaincode.getName()},
		CtorMsg:              chaincodeInput,
	}


	// Sec confidentiality
	if invokeTx.confidential {
		spec.ConfidentialityLevel = pb.ConfidentialityLevel_CONFIDENTIAL
	} else {
		spec.ConfidentialityLevel = pb.ConfidentialityLevel_PUBLIC
	}

	// Set metadata
	// Access control. Administrator signs chaincodeInputRaw || binding to confirm his identity
	if invokeTx.withCredential {
		binding, err := txHandler.GetBinding()
		if err != nil {
			invokeTx.pushError(err)

			return invokeTx
		}
		chaincodeInputRaw, err := proto.Marshal(chaincodeInput)
		if err != nil {
			invokeTx.pushError(err)

			return invokeTx
		}
		sigma, err := invokeTx.chaincode.Credential.certificateHandler.Sign(append(chaincodeInputRaw, binding...))
		if err != nil {
			invokeTx.pushError(err)

			return invokeTx
		}
		spec.Metadata = sigma; // Proof of identity
	}

	chaincodeInvocationSpec := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}

	// Now create the Transactions message and send to Peer.
	var fabricTx *pb.Transaction
	if invokeTx.query {
		fabricTx, err = txHandler.NewChaincodeQuery(chaincodeInvocationSpec, util.GenerateUUID())
	} else {
		fabricTx, err = txHandler.NewChaincodeExecute(chaincodeInvocationSpec, util.GenerateUUID())
	}
	if err != nil {
		invokeTx.pushError(err)

		return invokeTx
	}
	invokeTx.fabricTx = fabricTx

	clientSDKLog.Debug("Send..")

	invokeTx.txImpl.Send()

	return invokeTx
}

func (invokeTx *invokeTxImpl) AddArgument(argument string) InvokeTransaction {
	invokeTx.args = append(invokeTx.args, argument)

	return invokeTx
}

func (invokeTx *invokeTxImpl) AddArgumentCredential(credential Credential) InvokeTransaction {
	invokeTx.AddArgument(string(credential.Raw()))

	return invokeTx
}

func (invokeTx *invokeTxImpl) AddArgumentCredentialByAlias(credentialAlias string) InvokeTransaction {
	invokeTx.args = append(invokeTx.args, credentialAlias)

	return invokeTx
}
