package utils

import "testing"

func TestGetDialogSubmissionText(t *testing.T) {
	testCases := []struct {
		description   string
		submissionMap map[string]string
		submissionID  string
		want          string
		wantErrString string
	}{
		{
			description: "basic success case",
			submissionMap: map[string]string{
				"input_comment_submission_id": "This is my comment input from dialog",
			},
			submissionID: "input_comment_submission_id",
			want:         "This is my comment input from dialog",
		},
		{
			description:   "empty submission map",
			submissionMap: map[string]string{},
			submissionID:  "input_comment_submission_id",
			wantErrString: "empty dialog submission",
		},
		{
			description: "empty submission ID",
			submissionMap: map[string]string{
				"input_comment_submission_id": "This is my comment input from dialog",
			},
			wantErrString: "must specify submission ID",
		},
		{
			description: "submission entry not found",
			submissionMap: map[string]string{
				"input_comment_submission_id": "This is my comment input from dialog",
			},
			submissionID:  "free_form_submission_id",
			wantErrString: "no submission entry found for provided ID: free_form_submission_id",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			got, err := GetDialogSubmission(tc.submissionMap, tc.submissionID)
			if err != nil {
				if tc.wantErrString == "" {
					t.Fatal("unexpected error", err)
				}
				if err.Error() != tc.wantErrString {
					t.Fatalf("expected error: %s, got: %s", tc.wantErrString, err)
				}
			}
			if got != tc.want {
				t.Fatalf("expected submission: %s, got: %s", tc.want, got)
			}
		})
	}
}
