package fabric

import "github.com/hyperledger/fabric/core/crypto"

type credentialImpl struct {
	*errorHandlerImpl `json:"-"`

	certificateHandler crypto.CertificateHandler `json:"-"`
	RawCertificate                []byte
}

func (credential *credentialImpl) Raw() []byte {
	return credential.RawCertificate
}