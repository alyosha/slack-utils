package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/slack-go/slack"
)

const (
	DeleteActionID = "delete_action"
)

const (
	datePickTimeFmt      = "2006-01-02"
	valueStringMaxLength = 2000
)

type concreteTyper interface {
	ConcreteTypePtr() interface{}
	ConcreteTypeVal() interface{}
}

func (msg Msg) ConcreteTypePtr() interface{} {
	return &msg
}

func (msg Msg) ConcreteTypeVal() interface{} {
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

// NewAttributeEmbeddedValue takes a map of attributes to embed and generates
// a value string for use with block actions. When the block is interacted
// with, the value will be sent back into the callback endpoint, where its
// embedded values can be extracted and used for response actions. This
// functionality is intended for situations where the time-to-interact with
// the block is hard to predict (surveys, etc.) and a time-limited memory cache
// might not be a good option.
func NewAttributeEmbeddedValue(attributesToEmbed map[string]interface{}) (string, error) {
	if err := validateAttributes(attributesToEmbed); err != nil {
		return "", fmt.Errorf("validateAttributes > %w", err)
	}

	marshalledAttributes, err := json.Marshal(attributesToEmbed)
	if err != nil {
		return "", fmt.Errorf("json.Marshal > %w", err)
	}

	attributeEmbeddedVal := string(marshalledAttributes)

	if len(attributeEmbeddedVal) > valueStringMaxLength {
		return "", fmt.Errorf("embedded value string cannot be longer than %d characters", valueStringMaxLength)
	}

	return attributeEmbeddedVal, nil
}

// ExtractEmbeddedAttributes is used with value strings generated via
// GetAttributeEmbeddedValue to restore the embedded attributes back to their
// original types. The destination map should contain the same keys used when
// embedding the attributes mapped to the zero value for the type.
func ExtractEmbeddedAttributes(embeddedValueString string, dest map[string]interface{}) error {
	embedded := make(map[string]interface{})
	if err := json.Unmarshal([]byte(embeddedValueString), &embedded); err != nil {
		return fmt.Errorf("json.Unmarshal > %w", err)
	}

	typeMap := map[string]reflect.Type{}

	for k, v := range embedded {
		destVal, ok := dest[k]
		if !ok {
			return fmt.Errorf("emedded attribute of key: %s not found in destination map", k)
		}

		destMapType := reflect.TypeOf(destVal)
		typeMap[k] = destMapType

		switch val := destVal.(type) {
		case string, bool, float64, []interface{}, map[string]interface{}:
			dest[k] = v
		case concreteTyper:
			data, err := json.Marshal(v)
			if err != nil {
				return fmt.Errorf("json.Marshal > %w", err)
			}
			concreteTypePtr := val.ConcreteTypePtr()
			err = json.Unmarshal(data, concreteTypePtr)
			if err != nil {
				return fmt.Errorf("json.Unmarshal > %w", err)
			}
			if reflect.ValueOf(val).Kind() == reflect.Ptr {
				dest[k] = concreteTypePtr
				continue
			}
			dest[k] = concreteTypePtr.(concreteTyper).ConcreteTypeVal()
		default:
			switch reflect.ValueOf(destVal).Kind() {
			case reflect.Struct:
				return errors.New("unsupported embedded attribute type: only structs implementing concreteTyper interface can be embedded")
			case reflect.Map:
				convertedMap, err := convertJSONInterfaceMap(v, val, destMapType)
				if err != nil {
					return fmt.Errorf("convertJSONInterfaceMap > %w", err)
				}
				dest[k] = convertedMap
			case reflect.Slice, reflect.Array:
				convertedSliceOrArray, err := convertJSONInterfaceSliceOrArray(v, val, destMapType)
				if err != nil {
					return fmt.Errorf("convertJSONInterfaceSliceOrArray > %w", err)
				}
				dest[k] = convertedSliceOrArray
			default:
				if !reflect.TypeOf(v).ConvertibleTo(destMapType) {
					return fmt.Errorf("cannot convert type %T to type %T", v, val)
				}
				dest[k] = reflect.ValueOf(v).Convert(destMapType).Interface()
			}
		}
	}

	for k, v := range dest {
		if finalType := reflect.TypeOf(v); finalType != typeMap[k] {
			return fmt.Errorf("entry at key: %s should be type: %v, but is: %v", k, typeMap[k], finalType)
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

func convertJSONInterfaceMap(embeddedVal, destVal interface{}, destMapType reflect.Type) (interface{}, error) {
	stringInterfaceMap, ok := embeddedVal.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("dest map and embedded attribute vals at same key are of different types: %T => %T", destVal, embeddedVal)
	}

	destMapValue := reflect.ValueOf(destVal)

	for k, v := range stringInterfaceMap {
		keyVal, elemVal := reflect.ValueOf(k), reflect.ValueOf(v)
		elemType := reflect.TypeOf(v)
		mapType := reflect.MapOf(reflect.TypeOf(k), elemType)

		if destMapType.Elem() != elemType {
			if !elemType.ConvertibleTo(destMapType.Elem()) {
				return nil, fmt.Errorf("attempted to append elem of type: %v to map of type: %v", elemType, destMapType.Elem())
			}
			destMapValue.SetMapIndex(keyVal, elemVal.Convert(destMapType.Elem()))
			continue
		}

		if mapType != destMapType {
			return nil, fmt.Errorf("dest map and embedded attribute map types do not match: %v => %v", mapType, destMapType)
		}

		destMapValue.SetMapIndex(keyVal, elemVal)
	}

	return destMapValue.Interface(), nil
}

func convertJSONInterfaceSliceOrArray(embeddedVal, destVal interface{}, destMapType reflect.Type) (interface{}, error) {
	interfaceSlice, ok := embeddedVal.([]interface{})
	if !ok {
		return nil, fmt.Errorf("dest map and embedded attribute vals at same key are of different types: %T => %T", destVal, embeddedVal)
	}

	newSlice := reflect.MakeSlice(reflect.SliceOf(destMapType.Elem()), 0, cap(interfaceSlice))

	for _, elem := range interfaceSlice {
		elemType := reflect.TypeOf(elem)
		if destMapType.Elem() != elemType {
			if !elemType.ConvertibleTo(destMapType.Elem()) {
				return nil, fmt.Errorf("attempted to append elem of type: %T to slice of type: %T", elem, destVal)
			}
			newSlice = reflect.Append(newSlice, reflect.ValueOf(elem).Convert(destMapType.Elem()))
			continue
		}
		newSlice = reflect.Append(newSlice, reflect.ValueOf(elem))
	}

	if destMapType.Kind() == reflect.Slice {
		return newSlice.Interface(), nil
	}

	newArray := reflect.New(reflect.ArrayOf(newSlice.Len(), destMapType.Elem())).Elem()
	reflect.Copy(newArray, newSlice)

	return newArray.Interface(), nil
}
