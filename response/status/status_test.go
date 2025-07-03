package status

import (
	"net/http"
	"testing"
)

func TestStatus(t *testing.T) {
	var actual int

	actual = new(Ok).Status()

	if actual != http.StatusOK {
		t.Logf("status.Ok unexpected result %v", actual)
	}

	actual = new(Created).Status()

	if actual != http.StatusCreated {
		t.Logf("status.Created unexpected result %v", actual)
	}

	actual = new(Accepted).Status()

	if actual != http.StatusAccepted {
		t.Logf("status.Accepted unexpected result %v", actual)
	}

	actual = new(NoContent).Status()

	if actual != http.StatusNoContent {
		t.Logf("status.NoContent unexpected result %v", actual)
	}

	actual = new(ResetContent).Status()

	if actual != http.StatusResetContent {
		t.Logf("status.StatusResetContent unexpected result %v", actual)
	}

	actual = new(ServiceUnavailable).Status()

	if actual != http.StatusServiceUnavailable {
		t.Logf("status.ServiceUnavailable unexpected result %v", actual)
	}
}
