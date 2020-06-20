package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/slack-go/slack"
)

const (
	DeleteActionID = "delete_action"
)

const datePickTimeFmt = "2006-01-02"

type concreteTyper interface {
	concreteTypePtr() interface{}
	concreteTypeVal() interface{}
}

func (msg Msg) concreteTypePtr() interface{} {
	return &msg
}

func (msg Msg) concreteTypeVal() interface{} {
	return msg
}

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

// GetAttributeEmbeddedValue takes a map of attributes to embed and unique prefix
// to generate a value string for use with block actions. When the block is
// interacted with, the value will be sent back into the callback endpoint,
// where its embedded value can be extracted and used for response actions.
// This functionality is intended for situations where the time-to-interact
// with the block is hard to predict (surveys, etc.) and where a memory cache
// might not be a good option.
func GetAttributeEmbeddedValue(prefix string, attributesToEmbed map[string]interface{}) (string, error) {
	if err := validateAttributes(attributesToEmbed); err != nil {
		return "", fmt.Errorf("validateAttributes() > %w", err)
	}

	marshalledAttributes, err := json.Marshal(attributesToEmbed)
	if err != nil {
		return "", fmt.Errorf("json.Marshal() > %w", err)
	}

	return fmt.Sprintf("%s_%s", prefix, string(marshalledAttributes)), nil
}

// ExtractEmbeddedAttributes is used with value strings generated via
// GetAttributeEmbeddedValue to restore the embedded attributes back to their
// original types. The destination map should contain the same keys used when
// embedding the attributes mapped to the zero value for the type.
func ExtractEmbeddedAttributes(val, prefix string, dest map[string]interface{}) error {
	trimmedVal := strings.Replace(val, prefix+"_", "", -1)

	tmpDest := make(map[string]interface{})
	if err := json.Unmarshal([]byte(trimmedVal), &tmpDest); err != nil {
		return fmt.Errorf("json.Unmarshal() > %w", err)
	}

	for k, v := range tmpDest {
		destVal, ok := dest[k]
		if !ok {
			return fmt.Errorf("emedded attribute of key: %s not found in destination map", k)
		}
		switch val := destVal.(type) {
		case string, bool, float64, []interface{}, map[string]interface{}:
			dest[k] = v
		case concreteTyper:
			data, err := json.Marshal(v)
			if err != nil {
				return fmt.Errorf("json.Marshal() > %w", err)
			}
			concreteTypePtr := val.concreteTypePtr()
			err = json.Unmarshal(data, concreteTypePtr)
			if err != nil {
				return fmt.Errorf("json.Unmarshal() > %w", err)
			}
			dest[k] = concreteTypePtr.(concreteTyper).concreteTypeVal()
		default:
			origType, jsonType := reflect.TypeOf(val), reflect.TypeOf(v)
			origVal, jsonVal := reflect.ValueOf(val), reflect.ValueOf(v)
			switch valKind := origVal.Kind(); valKind {
			case reflect.Struct:
				return errors.New("unsupported embedded attribute type: only structs implementing concreteTyper interface can be embedded")
			case reflect.Map:
				convertedMap, err := convertJSONInterfaceMap(v, val, origVal, origType)
				if err != nil {
					return fmt.Errorf("convertJSONInterfaceMap() > %w", err)
				}
				dest[k] = convertedMap
			case reflect.Slice, reflect.Array:
				convertedSliceOrArray, err := convertJSONInterfaceSliceOrArray(v, val, origType, valKind)
				if err != nil {
					return fmt.Errorf("convertJSONInterfaceSliceOrArray() > %w", err)
				}
				dest[k] = convertedSliceOrArray
			default:
				if !jsonType.ConvertibleTo(origType) {
					return fmt.Errorf("cannot convert type %T to type %T", v, val)
				}
				dest[k] = jsonVal.Convert(origType).Interface()
			}
		}
	}
	return nil
}

func validateAttributes(attributesToEmbed map[string]interface{}) error {
	for _, v := range attributesToEmbed {
		switch reflect.ValueOf(v).Kind() {
		case reflect.Struct:
			if _, ok := v.(concreteTyper); !ok {
				return errors.New("due to JSON marshalling restrictions, all structs must implement concreteTyper interface to be embedded")
			}
		case reflect.Map:
			mapKeys := reflect.ValueOf(v).MapKeys()
			for _, mapKey := range mapKeys {
				if keyKind := mapKey.Kind(); keyKind != reflect.String {
					return fmt.Errorf("cannot have map with key of type: %v", keyKind)
				}
			}
		}
	}
	return nil
}

func convertJSONInterfaceMap(interfaceVal, destVal interface{}, origVal reflect.Value, origType reflect.Type) (interface{}, error) {
	stringInterfaceMap, ok := interfaceVal.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("dest map and embedded attribute vals at same key are of different types: %T => %T", destVal, interfaceVal)
	}

	for k, v := range stringInterfaceMap {
		keyVal, elemVal := reflect.ValueOf(k), reflect.ValueOf(v)
		elemType := reflect.TypeOf(v)
		mapType := reflect.MapOf(reflect.TypeOf(k), elemType)

		if origType.Elem() != elemType {
			if !elemType.ConvertibleTo(origType.Elem()) {
				return nil, fmt.Errorf("attempted to append elem of type: %v to slice of type: %v", elemType, origType.Elem())
			}
			origVal.SetMapIndex(keyVal, elemVal.Convert(origType.Elem()))
			continue
		}

		if mapType != origType {
			return nil, fmt.Errorf("dest map and embedded attribute map types do not match: %v => %v", mapType, origType)
		}

		origVal.SetMapIndex(keyVal, elemVal)
	}

	return origVal.Interface(), nil
}

func convertJSONInterfaceSliceOrArray(interfaceVal, destVal interface{}, origType reflect.Type, valKind reflect.Kind) (interface{}, error) {
	interfaceSlice, ok := interfaceVal.([]interface{})
	if !ok {
		return nil, fmt.Errorf("dest map and embedded attribute vals at same key are of different types: %T => %T", destVal, interfaceVal)
	}

	newSlice := reflect.MakeSlice(reflect.SliceOf(origType.Elem()), 0, cap(interfaceSlice))

	for _, elem := range interfaceSlice {
		elemType := reflect.TypeOf(elem)
		if origType.Elem() != elemType {
			if !elemType.ConvertibleTo(origType.Elem()) {
				return nil, fmt.Errorf("attempted to append elem of type: %T to slice of type: %T", elem, destVal)
			}
			newSlice = reflect.Append(newSlice, reflect.ValueOf(elem).Convert(origType.Elem()))
			continue
		}
		newSlice = reflect.Append(newSlice, reflect.ValueOf(elem))
	}

	if valKind == reflect.Slice {
		return newSlice.Interface(), nil
	}

	newArray := reflect.New(reflect.ArrayOf(newSlice.Len(), origType.Elem())).Elem()
	reflect.Copy(newArray, newSlice)
	return newArray.Interface(), nil
}
