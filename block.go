package utils

import (
	"time"

	"github.com/nlopes/slack"
)

var (
	cancelActionID = "cancel_action"
	cancelBtnTxt   = slack.NewTextBlockObject(slack.PlainTextType, "Cancel", false, false)
	ackActionID    = "acknowledge_action"
	ackBtnTxt      = slack.NewTextBlockObject(slack.PlainTextType, "Got it", false, false)

	CancelBtn = NewButtonWithStyle(cancelActionID, "cancel", cancelBtnTxt, slack.StyleDanger)
	AckBtn    = NewButtonWithStyle(cancelActionID, "acknowledge", ackBtnTxt, slack.StylePrimary)
)

func NewButtonWithStyle(actionID, value string, textObj *slack.TextBlockObject, style slack.Style) *slack.ButtonBlockElement {
	btn := slack.NewButtonBlockElement(actionID, value, textObj)
	btn.WithStyle(style)
	return btn
}

func NewDatePickerWithOpts(actionID string, placeholder *slack.TextBlockObject, initialDate time.Time) *slack.DatePickerBlockElement {
	picker := slack.NewDatePickerBlockElement(actionID)
	if placeholder != nil {
		picker.Placeholder = placeholder
		return picker
	}
	dateStr := initialDate.Format("2006-01-02")
	picker.InitialDate = dateStr
	return picker
}
