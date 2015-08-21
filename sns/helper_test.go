package sns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComposeMessageGCM(t *testing.T) {
	assert := assert.New(t)

	opt := make(map[string]interface{})
	msg := composeMessageGCM("test", opt)
	assert.Equal(`{"data":{"message":"test"}}`, msg)

	opt["sound"] = "jazz"
	msg = composeMessageGCM("test", opt)
	assert.Equal(`{"data":{"message":"test","sound":"jazz"}}`, msg)

	delete(opt, "sound")
	opt["badge"] = 5
	msg = composeMessageGCM("test", opt)
	assert.Equal(`{"data":{"badge":5,"message":"test"}}`, msg)

	opt["sound"] = "jazz"
	opt["badge"] = 5
	msg = composeMessageGCM("test", opt)
	assert.Equal(`{"data":{"badge":5,"message":"test","sound":"jazz"}}`, msg)

	opt["x-option"] = "foo"
	msg = composeMessageGCM("test", opt)
	assert.Equal(`{"data":{"badge":5,"message":"test","sound":"jazz","x-option":"foo"}}`, msg)
}

func TestComposeMessageAPNS(t *testing.T) {
	assert := assert.New(t)

	opt := make(map[string]interface{})
	msg := composeMessageAPNS("test", opt)
	assert.Equal(`{"aps":{"alert":"test","sound":"default"}}`, msg)

	opt["sound"] = "jazz"
	msg = composeMessageAPNS("test", opt)
	assert.Equal(`{"aps":{"alert":"test","sound":"jazz"}}`, msg)

	delete(opt, "sound")
	opt["badge"] = 5
	msg = composeMessageAPNS("test", opt)
	assert.Equal(`{"aps":{"alert":"test","badge":5,"sound":"default"}}`, msg)

	opt["sound"] = "jazz"
	opt["badge"] = 5
	msg = composeMessageAPNS("test", opt)
	assert.Equal(`{"aps":{"alert":"test","badge":5,"sound":"jazz"}}`, msg)

	opt["x-option"] = "foo"
	msg = composeMessageAPNS("test", opt)
	assert.Equal(`{"aps":{"alert":"test","badge":5,"sound":"jazz"},"x-option":"foo"}`, msg)
}
