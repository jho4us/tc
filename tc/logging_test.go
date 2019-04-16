package tc

import (
	"github.com/go-kit/kit/log"

	"os"
	"testing"
)

func TestLogging(t *testing.T) {
	rep, err := testCreateTestRepo(t)
	if err != nil {
		t.Fatal(err)
	}
	var logger log.Logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	tcs := NewService(rep)
	tcs = NewLoggingService(log.With(logger, "component", "tc"), tcs)

	_, err = tcs.CreateTest("")
	if err == nil {
		t.Fatal(err)
	}
	_, err = tcs.LoadTest("")
	if err == nil {
		t.Fatal(err)
	}
	_, err = tcs.LoadTest("unexistent")
	if err == nil {
		t.Fatal(err)
	}
	err = tcs.DeleteTest("")
	if err == nil {
		t.Fatal(err)
	}
	ta := tcs.Tests()
	if ta == nil {
		t.Errorf("cant query all tests")
	}
	testInsertTestAndGetTest(tcs, t, 0)
	testUpdateAndGetTest(tcs, t)
}
