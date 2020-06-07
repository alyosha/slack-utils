package utils

import (
	"time"

	"github.com/slack-go/slack"
)

const (
	DeleteActionID = "delete_action"
)

const datePickTimeFmt = "2006-01-02"

var DivBlock = slack.NewDividerBlock()

// NewTextBlock returns a section block of common configuration
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

// NewDoneBtn returns a new primary button which will callback to
// DeleteActionID on select. Configure the button title via param
func NewDoneButton(text string) *slack.ButtonBlockElement {
	return NewButton(DeleteActionID, "done", text, slack.StylePrimary)
}

// NewCancelBtn returns a new danger button which will callback to
// DeleteActionID on select. Configure the button title via param
func NewCancelButton(text string) *slack.ButtonBlockElement {
	return NewButton(DeleteActionID, "cancel", text, slack.StyleDanger)
}

// NewDatePickerWithOpts returns a new DatePickerBlockElement initialized with
// its date/placeholder text set to either placeholder text or an initial date
func NewDatePickerWithPlaceholder(actionID string, placeholder interface{}) *slack.DatePickerBlockElement {
	picker := slack.NewDatePickerBlockElement(actionID)

	switch placeholder.(type) {
	case string:
		picker.Placeholder = slack.NewTextBlockObject(slack.PlainTextType, placeholder.(string), false, false)
	case time.Time:
		picker.InitialDate = placeholder.(time.Time).Format(datePickTimeFmt)
	}

	return picker
}

// DateOptToTime parses the selected date opt back to time.Time, but owing to
// its format always resolve anything more granular than a day to zero
func DateOptToTime(opt string) (time.Time, error) {
	return time.Parse(datePickTimeFmt, opt)
}
