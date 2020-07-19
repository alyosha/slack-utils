package utils

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kylelemons/godebug/pretty"
	"github.com/slack-go/slack"
)

const (
	mockActionID             = "fAkE123"
	mockValue                = "click_me_btn"
	mockText                 = "Do something"
	fakeStyle    slack.Style = "fake"
)

var mockTextObj = &slack.TextBlockObject{
	Type:     slack.PlainTextType,
	Text:     "Do something",
	Emoji:    false,
	Verbatim: false,
}

type fakeConcreteTyper struct{}

func (f fakeConcreteTyper) concreteTypePtr() interface{} {
	return &Msg{}
}

func (f fakeConcreteTyper) concreteTypeVal() interface{} {
	return Msg{}
}

func TestNewButton(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description string
		style       slack.Style
		wantButton  *slack.ButtonBlockElement
	}{
		{
			description: "ignores default style",
			style:       slack.StyleDefault,
			wantButton: &slack.ButtonBlockElement{
				Type:     slack.METButton,
				ActionID: mockActionID,
				Text:     mockTextObj,
				Value:    mockValue,
			},
		},
		{
			description: "returns danger style button",
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
			description: "returns primary style button",
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
			btn := NewButton(mockActionID, mockValue, mockText, tc.style)
			if diff := pretty.Compare(btn, tc.wantButton); diff != "" {
				t.Fatalf("-got +want %s\n", diff)
			}
		})
	}
}

func TestNewDatePickerAtTime(t *testing.T) {
	t.Parallel()

	now := time.Now()
	expectedStr := now.Format("2006-01-02")

	testCases := []struct {
		description    string
		placeholder    interface{}
		wantDatePicker *slack.DatePickerBlockElement
	}{
		{
			description: "returns new date picker with placeholder text",
			placeholder: "Do something",
			wantDatePicker: &slack.DatePickerBlockElement{
				Type:        slack.METDatepicker,
				ActionID:    mockActionID,
				Placeholder: mockTextObj,
			},
		},
		{
			description: "returns new date picker with proper initial date",
			placeholder: now,
			wantDatePicker: &slack.DatePickerBlockElement{
				Type:        slack.METDatepicker,
				ActionID:    mockActionID,
				InitialDate: expectedStr,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			picker := NewDatePickerWithPlaceholder(mockActionID, tc.placeholder)
			if diff := pretty.Compare(picker, tc.wantDatePicker); diff != "" {
				t.Fatalf("-got +want %s\n", diff)
			}
		})
	}

}

func TestNewDateOptToTime(t *testing.T) {
	t.Parallel()

	now := time.Now()
	dateOptStr := now.Format(datePickTimeFmt)
	unrecognizedLayoutOptStr := now.Format(time.ANSIC)
	expectedTime, err := time.Parse(datePickTimeFmt, dateOptStr)
	if err != nil {
		t.Fatalf("received unexpected error: %s", err)
	}

	testCases := []struct {
		description string
		opt         string
		wantErr     bool
	}{
		{
			description: "returns parsed time value with no error",
			opt:         dateOptStr,
		},
		{
			description: "unrecognized layout string causes error",
			opt:         unrecognizedLayoutOptStr,
			wantErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			parsedTime, err := DateOptToTime(dateOptStr)
			if err != nil && !tc.wantErr {
				t.Fatalf("received unexpected error: %s", err)
			}
			if diff := pretty.Compare(parsedTime, expectedTime); diff != "" {
				t.Fatalf("-got +want %s\n", diff)
			}
		})
	}
}

func TestTextBlock(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description string
		text        string
		accessory   *slack.Accessory
		wantBlock   *slack.SectionBlock
	}{
		{
			description: "returns expected text-only section block",
			text:        "Basic message",
			wantBlock: &slack.SectionBlock{
				Type: slack.MBTSection,
				Text: &slack.TextBlockObject{
					Type:     slack.MarkdownType,
					Text:     "Basic message",
					Emoji:    false,
					Verbatim: false,
				},
			},
		},
		{
			description: "returns text block with accessory",
			text:        "Basic message",
			wantBlock: &slack.SectionBlock{
				Type: slack.MBTSection,
				Text: &slack.TextBlockObject{
					Type:     slack.MarkdownType,
					Text:     "Basic message",
					Emoji:    false,
					Verbatim: false,
				},
				Accessory: &slack.Accessory{
					ButtonElement: &slack.ButtonBlockElement{
						Type:     slack.METButton,
						ActionID: mockActionID,
						Text:     mockTextObj,
						Value:    mockValue,
					},
				},
			},
			accessory: slack.NewAccessory(NewButton(mockActionID, mockValue, mockText, slack.StyleDefault)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			block := NewTextBlock(tc.text, tc.accessory)
			if diff := pretty.Compare(block, tc.wantBlock); diff != "" {
				t.Fatalf("-got +want %s\n", diff)
			}
		})
	}
}

func TestEmbedAndExtractAttribute(t *testing.T) {
	t.Parallel()

	msgTimestampKey := "message_timestamp"
	timestamp := 1592117784

	requestingUserKey := "requesting_user"
	requestingUser := "U123456"

	requestedChannelsKey := "requested_channels"
	requestedChannels := [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	requestedNumsKey := "requested_nums"
	requestedNums := []int32{1, 2, 3}

	processedUsersKey := "processed_users"
	processedUsers := []string{"U123456", "U234567"}

	responseMsgKey := "response_msg"
	responseMsg := Msg{
		Body: "yo",
	}

	isEnabledKey := "is_enabled"
	isEnabled := true

	userMapKey := "user_map"
	userMap := map[string]string{
		"U12345": "steve",
		"U23456": "hal",
	}

	channelMapKey := "channel_map"
	channelMap := map[string]uint32{
		"C12345": 1,
		"C23456": 1,
	}

	reverseMapKey := "reverse_map"
	reverseMap := map[int]string{
		1: "why would you do this",
	}

	msgPtrKey := "msg_ptr"
	msgPtr := &Msg{Body: "testing attention please"}

	userChanKey := "user_chan"
	userChan := make(chan string)

	userFuncKey := "user_func"
	userFunc := func(user string) {
		fmt.Println(user)
	}

	readerKey := "reader_key"
	reader := bytes.Buffer{}
	reader.WriteString("yo")

	blockID := "actionblock"

	actionID1 := "action1"
	actionID2 := "action2"

	testCases := []struct {
		description      string
		embedAttributes1 map[string]interface{}
		embedAttributes2 map[string]interface{}
		dest1            map[string]interface{}
		dest2            map[string]interface{}
		wantErrString    string
	}{
		{
			description: "attributes embedded and extracted successfully",
			embedAttributes1: map[string]interface{}{
				msgPtrKey:            msgPtr,
				userMapKey:           userMap,
				processedUsersKey:    processedUsers,
				requestedChannelsKey: requestedChannels,
				isEnabledKey:         isEnabled,
				msgTimestampKey:      timestamp,
				requestingUserKey:    requestingUser,
				responseMsgKey:       responseMsg,
			},
			dest1: map[string]interface{}{
				msgPtrKey:            &Msg{},
				userMapKey:           map[string]string{},
				requestedChannelsKey: [10]int{},
				processedUsersKey:    []string{},
				isEnabledKey:         false,
				msgTimestampKey:      0,
				requestingUserKey:    "",
				responseMsgKey:       Msg{},
			},
			embedAttributes2: map[string]interface{}{
				msgPtrKey:         msgPtr,
				userMapKey:        userMap,
				processedUsersKey: processedUsers,
				isEnabledKey:      isEnabled,
				msgTimestampKey:   timestamp,
				requestingUserKey: requestingUser,
				responseMsgKey:    responseMsg,
			},
			dest2: map[string]interface{}{
				msgPtrKey:         &Msg{},
				userMapKey:        map[string]string{},
				processedUsersKey: []string{},
				isEnabledKey:      false,
				msgTimestampKey:   0,
				requestingUserKey: "",
				responseMsgKey:    Msg{},
			},
		},
		{
			description: "slice of ambiguous JSON type properly converted",
			embedAttributes1: map[string]interface{}{
				requestedNumsKey: requestedNums,
			},
			dest1: map[string]interface{}{
				requestedNumsKey: []int{},
			},
			embedAttributes2: map[string]interface{}{
				requestedNumsKey: requestedNums,
			},
			dest2: map[string]interface{}{
				requestedNumsKey: []int{},
			},
		},
		{
			description: "map of ambiguous JSON type properly converted",
			embedAttributes1: map[string]interface{}{
				channelMapKey: channelMap,
			},
			dest1: map[string]interface{}{
				channelMapKey: map[string]int{},
			},
			embedAttributes2: map[string]interface{}{
				channelMapKey: channelMap,
			},
			dest2: map[string]interface{}{
				channelMapKey: map[string]int{},
			},
		},
		{
			description: "unsupported type chan",
			embedAttributes1: map[string]interface{}{
				userChanKey: userChan,
			},
			embedAttributes2: map[string]interface{}{
				userChanKey: userChan,
			},
			dest1: map[string]interface{}{
				userChanKey: make(chan string),
			},
			dest2: map[string]interface{}{
				userChanKey: make(chan string),
			},
			wantErrString: "json.Marshal() > json: unsupported type: chan string",
		},
		{
			description: "unsupported type func",
			embedAttributes1: map[string]interface{}{
				userFuncKey: userFunc,
			},
			embedAttributes2: map[string]interface{}{
				userFuncKey: userFunc,
			},
			dest1: map[string]interface{}{
				userFuncKey: func(string) {},
			},
			dest2: map[string]interface{}{
				userFuncKey: func(string) {},
			},
			wantErrString: "json.Marshal() > json: unsupported type: func(string)",
		},
		{
			description: "unsupported non-concrete typer struct",
			embedAttributes1: map[string]interface{}{
				readerKey: reader,
			},
			embedAttributes2: map[string]interface{}{
				readerKey: reader,
			},
			dest1: map[string]interface{}{
				readerKey: bytes.Buffer{},
			},
			dest2: map[string]interface{}{
				readerKey: bytes.Buffer{},
			},
			wantErrString: "validateAttributes() > due to JSON marshalling restrictions, all structs must implement concreteTyper interface to be embedded",
		},
		{
			description: "dest map and embedded attributes have mismatched types - default case - no panic",
			embedAttributes1: map[string]interface{}{
				responseMsgKey: responseMsg,
			},
			dest1: map[string]interface{}{
				responseMsgKey: 0,
			},
			embedAttributes2: map[string]interface{}{
				responseMsgKey: responseMsg,
			},
			dest2: map[string]interface{}{
				responseMsgKey: Msg{},
			},
			wantErrString: "cannot convert type map[string]interface {} to type int",
		},
		{
			description: "dest map and embedded attributes have mismatched types - slice - no panic",
			embedAttributes1: map[string]interface{}{
				processedUsersKey: processedUsers,
			},
			dest1: map[string]interface{}{
				processedUsersKey: []int{},
			},
			embedAttributes2: map[string]interface{}{
				processedUsersKey: processedUsers,
			},
			dest2: map[string]interface{}{
				processedUsersKey: []string{},
			},
			wantErrString: "convertJSONInterfaceSliceOrArray() > attempted to append elem of type: string to slice of type: []int",
		},
		{
			description: "fails for unsupported map key types - no panic",
			embedAttributes1: map[string]interface{}{
				reverseMapKey: reverseMap,
			},
			dest1: map[string]interface{}{
				reverseMapKey: map[int]string{},
			},
			embedAttributes2: map[string]interface{}{
				reverseMapKey: reverseMap,
			},
			dest2: map[string]interface{}{
				reverseMapKey: map[int]string{},
			},
			wantErrString: "validateAttributes() > cannot have map with key of type: int",
		},
		{
			description: "fails when provided/return dest maps have different types - no panic",
			embedAttributes1: map[string]interface{}{
				"concrete_typer_key": fakeConcreteTyper{},
			},
			dest1: map[string]interface{}{
				"concrete_typer_key": fakeConcreteTyper{},
			},
			wantErrString: "entry at key: concrete_typer_key should be type: utils.fakeConcreteTyper, but is: utils.Msg",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			val1, err1 := GetAttributeEmbeddedValue(tc.embedAttributes1)
			val2, err2 := GetAttributeEmbeddedValue(tc.embedAttributes2)

			if err1 != nil || err2 != nil {
				if err1.Error() != tc.wantErrString || err2.Error() != tc.wantErrString {
					t.Fatal("unexpected err", err1, err2)
				}
				return
			}

			mux := http.NewServeMux()

			mux.HandleFunc("/chat.postMessage", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(mockSuccessResp))
			})

			testServSlack := httptest.NewServer(mux)
			defer testServSlack.Close()

			client := &Client{
				SlackAPI: slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServSlack.URL))),
			}

			btn1 := NewButton(actionID1, val1, "Click me", slack.StylePrimary)
			btn2 := NewButton(actionID2, val2, "Dont click me", slack.StyleDanger)

			_, err := client.PostMsg(Msg{Blocks: []slack.Block{slack.NewActionBlock(blockID, btn1, btn2)}}, "C1H9RESGL")

			if err != nil {
				t.Fatal("unexpected err", err)
			}

			err = ExtractEmbeddedAttributes(val1, tc.dest1)

			if err != nil {
				if err.Error() != tc.wantErrString {
					t.Fatal("unexpected err", err)
				}
				return
			}

			for k, v := range tc.embedAttributes1 {
				attr, ok := tc.dest1[k]
				if !ok {
					t.Fatalf("missing expected attribute: %v", v)
				}
				if diff := pretty.Compare(attr, v); diff != "" {
					t.Fatalf("-got +want %s\n", diff)
				}
			}

			err = ExtractEmbeddedAttributes(val2, tc.dest2)

			if err != nil {
				if err.Error() != tc.wantErrString {
					t.Fatal("unexpected err", err)
				}
				return
			}

			for k, v := range tc.embedAttributes2 {
				attr, ok := tc.dest2[k]
				if !ok {
					t.Fatalf("missing expected attribute: %v", v)
				}
				if diff := pretty.Compare(attr, v); diff != "" {
					t.Fatalf("-got +want %s\n", diff)
				}
			}
		})
	}
}
