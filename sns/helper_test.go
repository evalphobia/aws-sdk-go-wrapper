package sns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComposeMessageGCM(t *testing.T) {
	assert := assert.New(t)

	opt := make(map[string]interface{})
	msg, err := composeMessageGCM("test", opt, false)
	assert.NoError(err)
	assert.Equal(`{"data":{"message":"test"}, "android": {"priority":"normal"}}`, msg)

	opt["sound"] = "jazz"
	msg, err = composeMessageGCM("test", opt, false)
	assert.NoError(err)
	assert.Equal(`{"data":{"message":"test","sound":"jazz"}, "android": {"priority":"normal"}}`, msg)

	delete(opt, "sound")
	opt["badge"] = 5
	msg, err = composeMessageGCM("test", opt, false)
	assert.NoError(err)
	assert.Equal(`{"data":{"badge":5,"message":"test"}, "android": {"priority":"normal"}}`, msg)

	opt["sound"] = "jazz"
	opt["badge"] = 5
	msg, err = composeMessageGCM("test", opt, false)
	assert.NoError(err)
	assert.Equal(`{"data":{"badge":5,"message":"test","sound":"jazz"}, "android": {"priority":"normal"}}`, msg)

	opt["x-option"] = "foo"
	msg, err = composeMessageGCM("test", opt, false)
	assert.NoError(err)
	assert.Equal(`{"data":{"badge":5,"message":"test","sound":"jazz","x-option":"foo"}, "android": {"priority":"normal"}}`, msg)

	opt = make(map[string]interface{})
	msg, err = composeMessageGCM("test", opt, true)
	assert.NoError(err)
	assert.Equal(`{"data":{"message":"test"}, "android": {"priority":"high"}}`, msg)

	opt["sound"] = "jazz"
	msg, err = composeMessageGCM("test", opt, true)
	assert.NoError(err)
	assert.Equal(`{"data":{"message":"test","sound":"jazz"}, "android": {"priority":"high"}}`, msg)

	delete(opt, "sound")
	opt["badge"] = 5
	msg, err = composeMessageGCM("test", opt, true)
	assert.NoError(err)
	assert.Equal(`{"data":{"badge":5,"message":"test"}, "android": {"priority":"high"}}`, msg)

	opt["sound"] = "jazz"
	opt["badge"] = 5
	msg, err = composeMessageGCM("test", opt, true)
	assert.NoError(err)
	assert.Equal(`{"data":{"badge":5,"message":"test","sound":"jazz"}, "android": {"priority":"high"}}`, msg)

	opt["x-option"] = "foo"
	msg, err = composeMessageGCM("test", opt, true)
	assert.NoError(err)
	assert.Equal(`{"data":{"badge":5,"message":"test","sound":"jazz","x-option":"foo"}, "android": {"priority":"high"}}`, msg)
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
