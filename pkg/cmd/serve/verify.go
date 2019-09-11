/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package serve

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/vchain-us/vcn/pkg/api"
	"github.com/vchain-us/vcn/pkg/cmd/internal/types"
)

func verify(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	var keys []string
	org := r.URL.Query().Get("org")
	if org != "" {
		bo, err := api.GetBlockChainOrganisation(org)
		if err != nil {
			writeErrorResponse(w, "organization error", err, http.StatusBadRequest)
			return
		}
		keys = bo.MembersIDs()
	} else {
		ks := r.URL.Query().Get("signers")
		if ks != "" {
			keys = strings.Split(ks, ",")
			// add 0x if missing, lower case
			for i, k := range keys {
				if !strings.HasPrefix(k, "0x") {
					keys[i] = "0x" + k
				}
				keys[i] = strings.ToLower(keys[i])
			}
		}
	}

	var err error
	var verification *api.BlockchainVerification
	user, err := currentUser()
	if err != nil {
		writeErrorResponse(w, "user error", err, http.StatusBadRequest)
		return
	}

	// if keys have been passed, check for a verification matching them
	if len(keys) > 0 {
		verification, err = api.VerifyMatchingSignerIDs(hash, keys)
	} else {
		// if we have an user, check for verification matching user's key first
		userKey := ""
		if hasAuth, _ := user.IsAuthenticated(); hasAuth {
			userKey = user.Config().PublicAddress()
		}
		if userKey != "" {
			verification, err = api.VerifyMatchingSignerIDWithFallback(hash, userKey)
		} else {
			// if no passed keys nor user,
			// just get the last with highest level available verification
			verification, err = api.Verify(hash)
		}
	}

	if err != nil {
		writeErrorResponse(w, "verification error", err, http.StatusBadRequest)
		return
	}

	var artifact *api.ArtifactResponse
	if !verification.Unknown() {
		artifact, _ = api.LoadArtifact(user, hash, verification.MetaHash())
	}

	writeResponse(w, types.NewResult(nil, artifact, verification))
}
