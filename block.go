package utils

import (
	"time"

	"github.com/slack-go/slack"
)

const datePickTimeFmt = "2006-01-02"

const (
	CancelActionID = "cancel_action"
)

var (
	DoneBtn   = NewButton(CancelActionID, "done", "Done", slack.StylePrimary)
	CancelBtn = NewButton(CancelActionID, "cancel", "Cancel", slack.StyleDanger)
)

var DivBlock = slack.NewDividerBlock()

// NewTextBlock returns a section block of common configuration.
func NewTextBlock(body string, accessory *slack.Accessory) *slack.SectionBlock {
	text := slack.NewTextBlockObject(slack.MarkdownType, body, false, false)
	return slack.NewSectionBlock(text, nil, accessory)
}

// NewButton returns a new ButtonBlockElement set to the designated style
func NewButton(actionID, value string, text string, style slack.Style) *slack.ButtonBlockElement {
	btnText := slack.NewTextBlockObject(slack.PlainTextType, text, false, false)
	btn := slack.NewButtonBlockElement(actionID, value, btnText)
	if style != slack.StyleDefault {
		btn.WithStyle(style)
	}
	return btn
}

// NewDatePickerWithOpts returns a new DatePickerBlockElement initialized with
// its date/placeholder text set as specified by one of the two parameters
func NewDatePickerWithOpts(actionID string, placeholder *slack.TextBlockObject, initialDate time.Time) *slack.DatePickerBlockElement {
	picker := slack.NewDatePickerBlockElement(actionID)
	if placeholder != nil {
		picker.Placeholder = placeholder
		return picker
	}
	dateStr := initialDate.Format(datePickTimeFmt)
	picker.InitialDate = dateStr
	return picker
}

// DateOptToTime parses the selected date opt back to time.Time, but owing to
// its format will always initialize to zero time for everything more
// granular than a day. Add time.Duration as needed in downstream packages
func DateOptToTime(opt string) (time.Time, error) {
	return time.Parse(datePickTimeFmt, opt)
}
