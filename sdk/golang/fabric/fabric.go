/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fabric

type ErrorHandler interface {

	Flush() error

}

type Chain interface {
	ErrorHandler

	EnrollMember(alias string, enrollID, enrollPWD string) Member

	GetMember(alias string) Member

}

type Credential interface {

	Raw() []byte

}

type Chaincode interface {
	ErrorHandler

	BindCredential() Credential

	Deploy() DeployTransaction

	Invoke(functionName string) InvokeTransaction

	Query(functionName string) QueryTransaction
}

type Transaction interface {
	ErrorHandler

	Confidential() Transaction

	WithCredential() Transaction

	Send() Transaction

	GetResponse() []byte

}

type DeployTransaction interface {
	Transaction

}

type InvokeTransaction interface {
	Transaction

	AddArgument(argument string) InvokeTransaction

	AddArgumentCredential(credential Credential) InvokeTransaction

	AddArgumentCredentialByAlias(credentialAlias string) InvokeTransaction
}

type QueryTransaction interface {
	InvokeTransaction

}

type Member interface {
	ErrorHandler

	GetChaincode(alias string, uri string) Chaincode

	GetChaincodeByAlias(alias string) Chaincode

	Deploy(alias string, chaincodeURI string) DeployTransaction

	Invoke(chaincodeURI string, functionName string) InvokeTransaction

	Query(chaincodeURI string, functionName string) QueryTransaction
}