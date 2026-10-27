package downloaders

import (
	"errors"
	"testing"

	"github.com/rclone/rclone/fs/fserrors"
)

// zdrive-desktop patch: a fatal/no-retry error should jump straight past
// maxErrorCount instead of incrementing by one like any other error, so a
// caller (vfs cache: too many errors in kickWaiters) gives up on the very
// first attempt rather than waiting out the full retry budget for
// something that can never succeed.
func TestCountErrorsNoRetryErrorSkipsStraightToMaxErrorCount(t *testing.T) {
	dls := &Downloaders{}

	dls._countErrors(0, fserrors.NoRetryError(errors.New("permanently blocked")))

	if dls.errorCount <= maxErrorCount {
		t.Errorf("errorCount = %d, want > %d (maxErrorCount) after a single NoRetryError", dls.errorCount, maxErrorCount)
	}
}

func TestCountErrorsFatalErrorSkipsStraightToMaxErrorCount(t *testing.T) {
	dls := &Downloaders{}

	dls._countErrors(0, fserrors.FatalError(errors.New("permanently blocked")))

	if dls.errorCount <= maxErrorCount {
		t.Errorf("errorCount = %d, want > %d (maxErrorCount) after a single FatalError", dls.errorCount, maxErrorCount)
	}
}

// An ordinary (non-fatal, non-no-retry) error must still behave exactly as
// before the patch - increment by one, not jump to maxErrorCount+1.
func TestCountErrorsOrdinaryErrorStillIncrementsByOne(t *testing.T) {
	dls := &Downloaders{}

	dls._countErrors(0, errors.New("transient network blip"))

	if dls.errorCount != 1 {
		t.Errorf("errorCount = %d, want 1 after a single ordinary error", dls.errorCount)
	}
}
