package ethtool

import (
	"testing"
)

func TestParseNegotiationsSuccess(t *testing.T) {
	ecmd := &EthtoolLinkSettings{
		Cmd:                    ETHTOOL_GLINKSETTINGS,
		Link_mode_masks_nwords: 4,
		Link_mode_masks:        [381]uint32{0x62ef, 0x8000, 0x0, 0x0, 0x62ef, 0x8000, 0x0, 0x0, 0x607f, 0x8000, 0x0, 0x0, 0x0},
	}

	negotiations, err := ecmd.ParseNegotiations()
	if err != nil {
		t.Fatal("unexpected error during parsing")
	}

	if negotiations.Supported != uint64(uint32(0x8000))<<32+uint64(0x62ef) {
		t.Fatal("wrong parsing supported negotiation")
	}

	if negotiations.Advertising != uint64(uint32(0x8000))<<32+uint64(0x62ef) {
		t.Fatal("wrong parsing advertising negotiation")
	}

	if negotiations.LpAdvertising != uint64(uint32(0x8000))<<32+uint64(0x607f) {
		t.Fatal("wrong parsing lp advertising negotiation")
	}
}

func TestParseNegotiationsWrongCmd(t *testing.T) {
	ecmd := &EthtoolLinkSettings{}

	_, err := ecmd.ParseNegotiations()
	if err == nil {
		t.Fatal("expected error without cmd")
	}

	expectedError := "must CmdGet before parsing"
	if err.Error() != expectedError {
		t.Fatalf("unexpected error want: %s, got: %s", expectedError, err.Error())
	}
}

func TestParseNegotiationsWrongNword(t *testing.T) {
	ecmd := &EthtoolLinkSettings{
		Cmd: ETHTOOL_GLINKSETTINGS,
	}

	_, err := ecmd.ParseNegotiations()
	if err == nil {
		t.Fatal("expected error without cmd")
	}

	expectedError := "nwords not set"
	if err.Error() != expectedError {
		t.Fatalf("unexpected error want: %s, got: %s", expectedError, err.Error())
	}
}
