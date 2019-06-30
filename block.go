package utils

import (
	"time"

	"github.com/nlopes/slack"
)

const (
	datePickTimeFmt = "2006-01-02"

	CancelActionID   = "cancel_action"
	AckActionID      = "acknowledge_action"
	GoActionID       = "go_action"
	ContinueActionID = "continue_action"

	AckBlockID            = "ack_block"
	GoCancelBlockID       = "go_cancel_block"
	ContinueCancelBlockID = "go_continue_block"
)

var (
	cancelBtnTxt   = slack.NewTextBlockObject(slack.PlainTextType, "Cancel", false, false)
	ackBtnTxt      = slack.NewTextBlockObject(slack.PlainTextType, "Got it", false, false)
	goBtnTxt       = slack.NewTextBlockObject(slack.PlainTextType, "Go!", false, false)
	continueBtnTxt = slack.NewTextBlockObject(slack.PlainTextType, "Continue", false, false)
	CancelBtn      = NewButtonWithStyle(CancelActionID, "cancel", cancelBtnTxt, slack.StyleDanger)
	AckBtn         = NewButtonWithStyle(AckActionID, "acknowledge", ackBtnTxt, slack.StylePrimary)
	GoBtn          = NewButtonWithStyle(GoActionID, "go", goBtnTxt, slack.StylePrimary)
	ContinueBtn    = NewButtonWithStyle(ContinueActionID, "continue", continueBtnTxt, slack.StylePrimary)

	DivBlock            = slack.NewDividerBlock()
	AckBlock            = slack.NewActionBlock(AckBlockID, AckBtn)
	GoCancelBlock       = slack.NewActionBlock(GoCancelBlockID, GoBtn, CancelBtn)
	ContinueCancelBlock = slack.NewActionBlock(ContinueCancelBlockID, ContinueBtn, CancelBtn)
)

// NewButtonWithStyle returns a new ButtonBlockElement set to the designated style
func NewButtonWithStyle(actionID, value string, textObj *slack.TextBlockObject, style slack.Style) *slack.ButtonBlockElement {
	btn := slack.NewButtonBlockElement(actionID, value, textObj)
	btn.WithStyle(style)
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
