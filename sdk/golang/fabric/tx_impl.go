package fabric

import (
	"fmt"
	pb "github.com/hyperledger/fabric/protos"
)

type txImpl struct {
	*errorHandlerImpl

	chaincode    *chaincodeImpl
	confidential bool
	withCredential bool
	fabricTx *pb.Transaction
	response *pb.Response
}

func (tx *txImpl) Confidential() Transaction {
	tx.confidential = true

	return tx
}

func (tx *txImpl) WithCredential() Transaction {
	tx.withCredential = true
	tx.chaincode.BindCredential()

	return tx
}

func (tx *txImpl) Send() Transaction {

	if tx.fabricTx == nil {
		tx.pushError(fmt.Errorf("Transaction type not specified"))

		return tx
	}

	clientSDKLog.Debug("Sending [%s]", tx.fabricTx.String())

	resp, err := tx.chaincode.member.chain.sendTransaction(tx.fabricTx);
	if err != nil {
		clientSDKLog.Error("Error sending fabric transaction: [%s]", err)
		tx.pushError(fmt.Errorf("Error sending fabric transaction: [%s]", err))

		return tx
	}
	tx.response = resp

	clientSDKLog.Debug("Checking response status [%d]", resp.Status)

	if resp.Status == pb.Response_FAILURE || resp.Status == pb.Response_UNDEFINED {
		clientSDKLog.Debug("Failure. [%d][%s]", resp.Status, string(resp.Msg))
		tx.pushError(fmt.Errorf("Failure. [%d][%s]", resp.Status, string(resp.Msg)))
	} else if resp.Status == pb.Response_SUCCESS {
		clientSDKLog.Debug("Success. [%d][%s]", resp.Status, string(resp.Msg))
	} else {
		clientSDKLog.Debug("Response code unknown. [%d][%s]", resp.Status, string(resp.Msg))
		tx.pushError(fmt.Errorf("Response code unknown. [%d][%s]", resp.Status, string(resp.Msg)))
	}

	clientSDKLog.Debug("Sending done!")

	return tx
}

func (tx *txImpl) GetResponse() []byte {

	if tx.response == nil {
		tx.pushError(fmt.Errorf("Response not set yet!"))

		return nil
	}

	return tx.response.Msg
}