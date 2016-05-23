package fabric

import (
	"encoding/json"
	"io/ioutil"
)

type registryImpl struct {
	chain *chainImpl
}

func (r *registryImpl) init(chain *chainImpl) error {
	r.chain = chain

	// load from file

	return nil
}

func (r *registryImpl) addMember(m *memberImpl) error {
	if err := r.store(); err != nil {
		return err
	}
	return nil
}

func (r *registryImpl) addChaincode(c *chaincodeImpl) error {
	if err := r.store(); err != nil {
		return err
	}
	return nil
}

func (r *registryImpl) addTransaction(t *txImpl) error {
	return nil
}

func (r *registryImpl) load() (err error) {
	s, err := ioutil.ReadFile(r.chain.Alias)
	if err == nil {
		err = json.Unmarshal(s, r.chain)
		if err == nil {
			clientSDKLog.Info("Loading done!")
		}
	}

	return
}

func (r *registryImpl) store() error {
	s, err := json.MarshalIndent(r.chain, "", "    ")

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(r.chain.Alias, s, 0700)
	if err != nil {
		return err
	}


	return nil
}