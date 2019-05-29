package utils

import (
	"time"

	"github.com/nlopes/slack"
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
