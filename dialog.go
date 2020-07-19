package utils

import (
	"errors"
	"fmt"
)

// GetDialogSubmission is a convenience method for retrieving a specific entry
// in the submission map returned with the dialog submission callback.
func GetDialogSubmission(submissionMap map[string]string, submissionID string) (string, error) {
	if len(submissionMap) == 0 {
		return "", errors.New("empty dialog submission")
	}

	if submissionID == "" {
		return "", errors.New("must specify submission ID")
	}

	if submission, ok := submissionMap[submissionID]; ok {
		return submission, nil
	}

	return "", fmt.Errorf("no submission entry found for provided ID: %s", submissionID)
}
