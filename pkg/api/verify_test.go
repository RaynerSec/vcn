/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package api

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vchain-us/vcn/pkg/meta"
	"gopkg.in/yaml.v2"
)

const (
	zeroMetahash        = "04d05179816469d9d1d3493b695df3396b10617035c9621bbb8757cfa0efb0f4"
	zeroUnknownMetahash = "83e7f28a74ada954937e609a2b7966edc39497d53203094b9097d3684f9a27ec"
	emptyJSON           = "null"
	emptyYAML           = "null\n"
)

func TestBlockchainVerification(t *testing.T) {
	var bv *BlockchainVerification

	// Nil
	assert.False(t, bv.Trusted())
	assert.True(t, bv.Unknown())
	assert.Empty(t, bv.MetaHash())
	assert.Empty(t, bv.Key())
	assert.Empty(t, bv.LevelName())
	assert.Empty(t, bv.Date())
	j, err := json.Marshal(bv)
	assert.NoError(t, err)
	assert.Equal(t, emptyJSON, string(j))
	y, err := yaml.Marshal(bv)
	assert.NoError(t, err)
	assert.Equal(t, emptyYAML, string(y))

	bv = &BlockchainVerification{}
	// Zero value
	assert.True(t, bv.Trusted())
	assert.False(t, bv.Unknown())
	assert.Equal(t, zeroMetahash, bv.MetaHash())
	assert.Empty(t, bv.Key())
	assert.NotEmpty(t, bv.LevelName())
	assert.Empty(t, bv.Date())
	j, err = json.Marshal(bv)
	assert.NoError(t, err)
	assert.NotEqual(t, emptyJSON, string(j))
	y, err = yaml.Marshal(bv)
	assert.NoError(t, err)
	assert.NotEqual(t, emptyYAML, string(y))

	bv = &BlockchainVerification{
		Status: meta.StatusUnknown,
	}
	// Default status
	assert.False(t, bv.Trusted())
	assert.True(t, bv.Unknown())
	assert.Equal(t, zeroUnknownMetahash, bv.MetaHash())
	assert.Empty(t, bv.Key())
	assert.NotEmpty(t, bv.LevelName())
	assert.Empty(t, bv.Date())
	j, err = json.Marshal(bv)
	assert.NoError(t, err)
	assert.NotEqual(t, emptyJSON, string(j))
	y, err = yaml.Marshal(bv)
	assert.NoError(t, err)
	assert.NotEqual(t, emptyYAML, string(y))
}
