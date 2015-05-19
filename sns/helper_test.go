package sns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComposeMessageGCM(t *testing.T) {
	msg := composeMessageGCM("test")
	assert.Equal(t, `{"data": {"message": "test"}}`, msg)
}

func TestComposeMessageAPNS(t *testing.T) {
	opt := make(map[string]interface{})
	msg := composeMessageAPNS("test", opt)
	assert.Equal(t, `{"aps":{"alert": "test", "sound": "default"}}`, msg)

	opt["sound"] = "jazz"
	msg = composeMessageAPNS("test", opt)
	assert.Equal(t, `{"aps":{"alert": "test", "sound": "jazz"}}`, msg)

	delete(opt, "sound")
	opt["badge"] = 5
	msg = composeMessageAPNS("test", opt)
	assert.Equal(t, `{"aps":{"alert": "test", "sound": "default", "badge": 5}}`, msg)

	opt["sound"] = "jazz"
	opt["badge"] = 5
	msg = composeMessageAPNS("test", opt)
	assert.Equal(t, `{"aps":{"alert": "test", "sound": "jazz", "badge": 5}}`, msg)

}
