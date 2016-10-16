package pointers

import (
	"github.com/Sirupsen/logrus"
	"github.com/evalphobia/go-log-wrapper/log"
	"github.com/evalphobia/go-log-wrapper/log/sentry"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	sentry.AddLevel(logrus.WarnLevel)
	sentry.Set("https://9ea8861ff68d4288839f37c645e2aa62:d219f4fe7f5148798a4a126287d6655a@sentry.io/59738")

	list1 := []int64{1, 2, 3, 4}
	list2 := []int64{1001, 1002, 1003, 1004}

	log.Packet{
		Title: "hoge2",
		Data:  struct{ UserID, ExistedIDs, CurrentIDs interface{} }{999, list1, list2},
	}.Error()

	time.Sleep(1 * time.Second)
}
