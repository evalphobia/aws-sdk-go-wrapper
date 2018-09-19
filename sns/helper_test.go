package sns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComposeMessageGCM(t *testing.T) {
	assert := assert.New(t)

	opt := make(map[string]interface{})
	msg, err := composeMessageGCM("test", opt)
	assert.NoError(err)
	assert.Equal(`{"android":{"priority":"high"},"data":{"message":"test"}}`, msg)

	opt["sound"] = "jazz"
	msg, err = composeMessageGCM("test", opt)
	assert.NoError(err)
	assert.Equal(`{"android":{"priority":"high"},"data":{"message":"test","sound":"jazz"}}`, msg)

	delete(opt, "sound")
	opt["badge"] = 5
	msg, err = composeMessageGCM("test", opt)
	assert.NoError(err)
	assert.Equal(`{"android":{"priority":"high"},"data":{"badge":5,"message":"test"}}`, msg)

	opt["sound"] = "jazz"
	opt["badge"] = 5
	msg, err = composeMessageGCM("test", opt)
	assert.NoError(err)
	assert.Equal(`{"android":{"priority":"high"},"data":{"badge":5,"message":"test","sound":"jazz"}}`, msg)

	opt["x-option"] = "foo"
	msg, err = composeMessageGCM("test", opt)
	assert.NoError(err)
	assert.Equal(`{"android":{"priority":"high"},"data":{"badge":5,"message":"test","sound":"jazz","x-option":"foo"}}`, msg)
}

func TestComposeMessageAPNS(t *testing.T) {
	assert := assert.New(t)

	opt := make(map[string]interface{})
	msg, err := composeMessageAPNS("test", opt)
	assert.NoError(err)
	assert.Equal(`{"aps":{"alert":"test","sound":"default"}}`, msg)

	opt["sound"] = "jazz"
	msg, err = composeMessageAPNS("test", opt)
	assert.NoError(err)
	assert.Equal(`{"aps":{"alert":"test","sound":"jazz"}}`, msg)

	delete(opt, "sound")
	opt["badge"] = 5
	msg, err = composeMessageAPNS("test", opt)
	assert.NoError(err)
	assert.Equal(`{"aps":{"alert":"test","badge":5,"sound":"default"}}`, msg)

	opt["sound"] = "jazz"
	opt["category"] = "new_message"
	opt["badge"] = 5
	msg, err = composeMessageAPNS("test", opt)
	assert.NoError(err)
	assert.Equal(`{"aps":{"alert":"test","badge":5,"category":"new_message","sound":"jazz"}}`, msg)

	opt["x-option"] = "foo"
	msg, err = composeMessageAPNS("test", opt)
	assert.NoError(err)
	assert.Equal(`{"aps":{"alert":"test","badge":5,"category":"new_message","sound":"jazz"},"x-option":"foo"}`, msg)
}
