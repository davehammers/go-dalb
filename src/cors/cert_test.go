package cors

// Copyright (c) 2018 by Extreme Networks Inc.

import "testing"

func TestCertKeys(t *testing.T) {
	var err error

	if _, _, err = CertKeys(); err != nil {
		t.Error("Cannot find public/private certificates for HTTPS")
		t.FailNow()
	}

	t.Log("Found public/private certificates for HTTPS")
}
