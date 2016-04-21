/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package crypto

import (
	"fmt"
	"sync"
)

type tCertPoolSingleThreadImpl struct {
	client *clientImpl

	len    int
	tCerts []TransactionCertificate
	m      sync.Mutex
}

func (tCertPool *tCertPoolSingleThreadImpl) Start() (err error) {
	tCertPool.m.Lock()
	defer tCertPool.m.Unlock()

	tCertPool.client.debug("Starting TCert Pool...")

	// Load unused TCerts if any
	tCertDERs, err := tCertPool.client.ks.loadUnusedTCerts()
	if err != nil {
		tCertPool.client.error("Failed loading TCerts from cache: [%s]", err)

		return
	}
	if len(tCertDERs) == 0 {
		tCertPool.client.debug("No more TCerts in cache! Load new from TCA.")

		tCertPool.client.getTCertsFromTCA(tCertPool.client.conf.getTCertBatchSize())
	} else {
		tCertPool.client.debug("TCerts in cache found! Loading them...")

		for _, tCertDER := range tCertDERs {
			tCert, err := tCertPool.client.getTCertFromDER(tCertDER, true)
			if err != nil {
				tCertPool.client.error("Failed paring TCert [% x]: [%s]", tCertDER, err)

				continue
			}
			tCertPool.AddTCert(tCert)
		}
	}

	return
}

func (tCertPool *tCertPoolSingleThreadImpl) Stop() (err error) {
	tCertPool.m.Lock()
	defer tCertPool.m.Unlock()

	tCertPool.client.debug("Found %d unused TCerts...", tCertPool.len)

	tCertPool.client.ks.storeUnusedTCerts(tCertPool.tCerts[:tCertPool.len])

	tCertPool.client.debug("Store unused TCerts...done!")

	return
}

func (tCertPool *tCertPoolSingleThreadImpl) GetNextTCert() (tCert TransactionCertificate, err error) {
	tCertPool.m.Lock()
	defer tCertPool.m.Unlock()

	if tCertPool.len <= 0 {
		// Reload
		if err := tCertPool.client.getTCertsFromTCA(tCertPool.client.conf.getTCertBatchSize()); err != nil {

			return nil, fmt.Errorf("Failed loading TCerts from TCA")
		}
	}

	tCert = tCertPool.tCerts[tCertPool.len-1]
	tCertPool.len--

	return
}

func (tCertPool *tCertPoolSingleThreadImpl) AddTCert(tCert TransactionCertificate) (err error) {
	tCertPool.client.debug("Adding new Cert [% x].", tCert.GetRaw())

	tCertPool.len++
	tCertPool.tCerts[tCertPool.len-1] = tCert

	return nil
}

func (tCertPool *tCertPoolSingleThreadImpl) init(client *clientImpl) (err error) {
	tCertPool.client = client

	tCertPool.client.debug("Init TCert Pool...")

	tCertPool.tCerts = make([]TransactionCertificate, tCertPool.client.conf.getTCertBatchSize())
	tCertPool.len = 0

	return
}
