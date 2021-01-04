package services

import (
	"testing"

	"github.com/nocubicles/develytica/src/models"
)

func TestIsSyncInProgressTrue(t *testing.T) {
	sync := models.Sync{
		InProgress: true,
	}

	syncInProgress := isSyncInProgress(sync)

	if syncInProgress == false {
		t.Error("Sync in progress test failed")
	}
}

func TestIsSyncInProgressFalse(t *testing.T) {
	sync := models.Sync{
		InProgress: false,
	}

	syncInProgress := isSyncInProgress(sync)

	if syncInProgress == true {
		t.Error("Sync in progress  test failed")
	}
}
