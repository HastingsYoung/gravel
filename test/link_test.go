package test

import (
	. "github.com/gravel/models"
	"testing"
)

func TestLink(t *testing.T) {
	var (
		link_root          = NewLink("/")
		link_level_one     = NewLink("/level1")
		link_level_two     = NewLink("/level1/level2")
		link_level_three   = NewLink("/level1/level2/level3")
		link_not_compliant = NewLink("\\.level1")
	)

	if ok := link_root.Parsed(); !ok {
		t.Error(
			"Failed to pass test => [",
			link_root.Raw(),
			"]",
			" Expected:",
			true,
			" But got:",
			ok,
		)
	}

	if ok := link_level_one.Parsed(); !ok {
		t.Error(
			"Failed to pass test => [",
			link_level_one.Raw(),
			"]",
			" Expected:",
			true,
			" But got:",
			ok,
		)
	}

	if ok := link_level_two.Parsed(); !ok {
		t.Error(
			"Failed to pass test => [",
			link_level_two.Raw(),
			"]",
			" Expected:",
			true,
			" But got:",
			ok,
		)
	}

	if ok := link_level_three.Parsed(); !ok {
		t.Error(
			"Failed to pass test => [",
			link_level_three.Raw(),
			"]",
			" Expected:",
			true,
			" But got:",
			ok,
		)
	}

	if ok := link_not_compliant.Parsed(); ok {
		t.Error(
			"Failed to pass test => [",
			link_not_compliant.Raw(),
			"]",
			" Expected:",
			false,
			" But got:",
			ok,
		)
	}
}
