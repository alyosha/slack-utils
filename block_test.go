package utils

import (
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/nlopes/slack"
)

const (
	mockActionID             = "fAkE123"
	mockValue                = "click_me_btn"
	fakeStyle    slack.Style = "fake"
)

var mockTextObj = &slack.TextBlockObject{
	Type:     slack.PlainTextType,
	Text:     "Click me",
	Emoji:    false,
	Verbatim: false,
}

func TestNewButtonWithStyle(t *testing.T) {
	testCases := []struct {
		description string
		style       slack.Style
		wantButton  *slack.ButtonBlockElement
	}{
		{
			description: "returns default style button, no error",
			style:       slack.StyleDefault,
			wantButton: &slack.ButtonBlockElement{
				Type:     slack.METButton,
				ActionID: mockActionID,
				Text:     mockTextObj,
				Value:    mockValue,
				Style:    slack.StyleDefault,
			},
		},
		{
			description: "returns danger style button, no error",
			style:       slack.StyleDanger,
			wantButton: &slack.ButtonBlockElement{
				Type:     slack.METButton,
				ActionID: mockActionID,
				Text:     mockTextObj,
				Value:    mockValue,
				Style:    slack.StyleDanger,
			},
		},
		{
			description: "returns primary style button, no error",
			style:       slack.StylePrimary,
			wantButton: &slack.ButtonBlockElement{
				Type:     slack.METButton,
				ActionID: mockActionID,
				Text:     mockTextObj,
				Value:    mockValue,
				Style:    slack.StylePrimary,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			btn := NewButtonWithStyle(mockActionID, mockValue, mockTextObj, tc.style)
			if diff := pretty.Compare(btn, tc.wantButton); diff != "" {
				t.Fatalf("expected to receive button: %v, got: %v", tc.wantButton, btn)
			}
		})
	}
}
