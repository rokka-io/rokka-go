package rokka

// Code generated by go generate; DO NOT EDIT.
// This file was generated at 2018-03-06 17:30:00.325584678 +0100 CET m=+0.274580001

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Operation is an interface all operation structs implement.
type Operation interface {
	// Name returns the operation's name known by the API.
	Name() string
	// Validate checks if required properties are set.
	// Otherwise it returns false with an error indicating the missing property.
	Validate() (bool, error)
	// toURLPath generates a part of the URL used for dynamic rendering of a stack.
	toURLPath() string
}

type rawStack struct {
	Name    string          `json:"name"`
	Options json.RawMessage `json:"options"`
}

// Operations is a slice of Operation implementing json.Unmarshaler and json.Marshaler in order to create
// the correct operation types for JSON.
type Operations []Operation

// UnmarshalJSON implements json.Unmarshaler.
func (o *Operations) UnmarshalJSON(data []byte) error {
	ops := make([]rawStack, 0)
	if err := json.Unmarshal(data, &ops); err != nil {
		return err
	}
	for _, v := range ops {
		op, err := NewOperationByName(v.Name)
		if err != nil {
			return err
		}
		*o = append(*o, op.(Operation))
		if err := json.Unmarshal(v.Options, op); err != nil {
			// BUG(mweibel): We continue here when such an error is reached because rokka sometimes (legacy reasons)
			//               has options on an operation which are not of the correct type. Should we write something to stdout? also not nice though..
			continue
		}
	}
	return nil
}

// MarshalJSON implements json.Marshaler
func (o Operations) MarshalJSON() ([]byte, error) {
	ops := make([]map[string]interface{}, len(o))
	for i, v := range o {
		ops[i] = make(map[string]interface{})
		ops[i]["name"] = v.Name()
		ops[i]["options"] = v
	}

	return json.Marshal(ops)
}

var errOperationNotImplemented = errors.New("Operation not implemented")

// NewOperationByName creates a struct of the respective type based on the name given.
func NewOperationByName(name string) (Operation, error) {
	switch name {
	case "alpha":
		return new(AlphaOperation), nil

	case "autorotate":
		return new(AutorotateOperation), nil

	case "blur":
		return new(BlurOperation), nil

	case "composition":
		return new(CompositionOperation), nil

	case "crop":
		return new(CropOperation), nil

	case "dropshadow":
		return new(DropshadowOperation), nil

	case "grayscale":
		return new(GrayscaleOperation), nil

	case "noop":
		return new(NoopOperation), nil

	case "primitive":
		return new(PrimitiveOperation), nil

	case "resize":
		return new(ResizeOperation), nil

	case "rotate":
		return new(RotateOperation), nil

	case "sepia":
		return new(SepiaOperation), nil

	case "trim":
		return new(TrimOperation), nil

	}
	return nil, errOperationNotImplemented
}

// AlphaOperation is an auto-generated Operation as specified by the rokka API.
//
// See: https://rokka.io/documentation/references/operations.html
type AlphaOperation struct {
	Mode *string `json:"mode,omitempty"`
}

// Name implements rokka.Operation.Name
func (o AlphaOperation) Name() string { return "alpha" }

// Validate implements rokka.Operation.Validate.
func (o AlphaOperation) Validate() (bool, error) {
	return true, nil
}

// toURLPath implements rokka.Operation.toURLPath.
func (o AlphaOperation) toURLPath() string {
	options := make([]string, 0)
	if o.Mode != nil {
		options = append(options, fmt.Sprintf("%s-%v", "mode", *o.Mode))
	}

	if len(options) == 0 {
		return o.Name()
	}
	return fmt.Sprintf("%s-%s", o.Name(), strings.Join(options, "-"))
}

// AutorotateOperation is an auto-generated Operation as specified by the rokka API.
//
// See: https://rokka.io/documentation/references/operations.html
type AutorotateOperation struct {
	Height            *int    `json:"height,omitempty"`
	RotationDirection *string `json:"rotation_direction,omitempty"`
	Width             *int    `json:"width,omitempty"`
}

// Name implements rokka.Operation.Name
func (o AutorotateOperation) Name() string { return "autorotate" }

// Validate implements rokka.Operation.Validate.
func (o AutorotateOperation) Validate() (bool, error) {
	return true, nil
}

// toURLPath implements rokka.Operation.toURLPath.
func (o AutorotateOperation) toURLPath() string {
	options := make([]string, 0)
	if o.Height != nil {
		options = append(options, fmt.Sprintf("%s-%v", "height", *o.Height))
	}
	if o.RotationDirection != nil {
		options = append(options, fmt.Sprintf("%s-%v", "rotation_direction", *o.RotationDirection))
	}
	if o.Width != nil {
		options = append(options, fmt.Sprintf("%s-%v", "width", *o.Width))
	}

	if len(options) == 0 {
		return o.Name()
	}
	return fmt.Sprintf("%s-%s", o.Name(), strings.Join(options, "-"))
}

// BlurOperation is an auto-generated Operation as specified by the rokka API.
// Calling .Validate() will return false if required properties are missing.
//
// See: https://rokka.io/documentation/references/operations.html
type BlurOperation struct {
	Sigma *float64 `json:"sigma,omitempty"`
}

// Name implements rokka.Operation.Name
func (o BlurOperation) Name() string { return "blur" }

// Validate implements rokka.Operation.Validate.
func (o BlurOperation) Validate() (bool, error) {
	if o.Sigma == nil {
		return false, errors.New("option \"Sigma\" is required")
	}
	return true, nil
}

// toURLPath implements rokka.Operation.toURLPath.
func (o BlurOperation) toURLPath() string {
	options := make([]string, 0)
	if o.Sigma != nil {
		options = append(options, fmt.Sprintf("%s-%v", "sigma", *o.Sigma))
	}

	if len(options) == 0 {
		return o.Name()
	}
	return fmt.Sprintf("%s-%s", o.Name(), strings.Join(options, "-"))
}

// CompositionOperation is an auto-generated Operation as specified by the rokka API.
//
// See: https://rokka.io/documentation/references/operations.html
type CompositionOperation struct {
	Anchor           *string `json:"anchor,omitempty"`
	Height           *int    `json:"height,omitempty"`
	Mode             *string `json:"mode,omitempty"`
	SecondaryColor   *string `json:"secondary_color,omitempty"`
	SecondaryOpacity *int    `json:"secondary_opacity,omitempty"`
	Width            *int    `json:"width,omitempty"`
}

// Name implements rokka.Operation.Name
func (o CompositionOperation) Name() string { return "composition" }

// Validate implements rokka.Operation.Validate.
func (o CompositionOperation) Validate() (bool, error) {
	return true, nil
}

// toURLPath implements rokka.Operation.toURLPath.
func (o CompositionOperation) toURLPath() string {
	options := make([]string, 0)
	if o.Anchor != nil {
		options = append(options, fmt.Sprintf("%s-%v", "anchor", *o.Anchor))
	}
	if o.Height != nil {
		options = append(options, fmt.Sprintf("%s-%v", "height", *o.Height))
	}
	if o.Mode != nil {
		options = append(options, fmt.Sprintf("%s-%v", "mode", *o.Mode))
	}
	if o.SecondaryColor != nil {
		options = append(options, fmt.Sprintf("%s-%v", "secondary_color", *o.SecondaryColor))
	}
	if o.SecondaryOpacity != nil {
		options = append(options, fmt.Sprintf("%s-%v", "secondary_opacity", *o.SecondaryOpacity))
	}
	if o.Width != nil {
		options = append(options, fmt.Sprintf("%s-%v", "width", *o.Width))
	}

	if len(options) == 0 {
		return o.Name()
	}
	return fmt.Sprintf("%s-%s", o.Name(), strings.Join(options, "-"))
}

// CropOperation is an auto-generated Operation as specified by the rokka API.
// Calling .Validate() will return false if required properties are missing.
//
// See: https://rokka.io/documentation/references/operations.html
type CropOperation struct {
	Anchor *string  `json:"anchor,omitempty"`
	Height *int     `json:"height,omitempty"`
	Mode   *string  `json:"mode,omitempty"`
	Scale  *float64 `json:"scale,omitempty"`
	Width  *int     `json:"width,omitempty"`
}

// Name implements rokka.Operation.Name
func (o CropOperation) Name() string { return "crop" }

// Validate implements rokka.Operation.Validate.
func (o CropOperation) Validate() (bool, error) {
	if o.Width == nil {
		return false, errors.New("option \"Width\" is required")
	}
	if o.Height == nil {
		return false, errors.New("option \"Height\" is required")
	}
	return true, nil
}

// toURLPath implements rokka.Operation.toURLPath.
func (o CropOperation) toURLPath() string {
	options := make([]string, 0)
	if o.Anchor != nil {
		options = append(options, fmt.Sprintf("%s-%v", "anchor", *o.Anchor))
	}
	if o.Height != nil {
		options = append(options, fmt.Sprintf("%s-%v", "height", *o.Height))
	}
	if o.Mode != nil {
		options = append(options, fmt.Sprintf("%s-%v", "mode", *o.Mode))
	}
	if o.Scale != nil {
		options = append(options, fmt.Sprintf("%s-%v", "scale", *o.Scale))
	}
	if o.Width != nil {
		options = append(options, fmt.Sprintf("%s-%v", "width", *o.Width))
	}

	if len(options) == 0 {
		return o.Name()
	}
	return fmt.Sprintf("%s-%s", o.Name(), strings.Join(options, "-"))
}

// DropshadowOperation is an auto-generated Operation as specified by the rokka API.
//
// See: https://rokka.io/documentation/references/operations.html
type DropshadowOperation struct {
	BlurRadius *float64 `json:"blur_radius,omitempty"`
	Color      *string  `json:"color,omitempty"`
	Horizontal *int     `json:"horizontal,omitempty"`
	Opacity    *int     `json:"opacity,omitempty"`
	Sigma      *float64 `json:"sigma,omitempty"`
	Vertical   *int     `json:"vertical,omitempty"`
}

// Name implements rokka.Operation.Name
func (o DropshadowOperation) Name() string { return "dropshadow" }

// Validate implements rokka.Operation.Validate.
func (o DropshadowOperation) Validate() (bool, error) {
	return true, nil
}

// toURLPath implements rokka.Operation.toURLPath.
func (o DropshadowOperation) toURLPath() string {
	options := make([]string, 0)
	if o.BlurRadius != nil {
		options = append(options, fmt.Sprintf("%s-%v", "blur_radius", *o.BlurRadius))
	}
	if o.Color != nil {
		options = append(options, fmt.Sprintf("%s-%v", "color", *o.Color))
	}
	if o.Horizontal != nil {
		options = append(options, fmt.Sprintf("%s-%v", "horizontal", *o.Horizontal))
	}
	if o.Opacity != nil {
		options = append(options, fmt.Sprintf("%s-%v", "opacity", *o.Opacity))
	}
	if o.Sigma != nil {
		options = append(options, fmt.Sprintf("%s-%v", "sigma", *o.Sigma))
	}
	if o.Vertical != nil {
		options = append(options, fmt.Sprintf("%s-%v", "vertical", *o.Vertical))
	}

	if len(options) == 0 {
		return o.Name()
	}
	return fmt.Sprintf("%s-%s", o.Name(), strings.Join(options, "-"))
}

// GrayscaleOperation is an auto-generated Operation as specified by the rokka API.
//
// See: https://rokka.io/documentation/references/operations.html
type GrayscaleOperation struct {
}

// Name implements rokka.Operation.Name
func (o GrayscaleOperation) Name() string { return "grayscale" }

// Validate implements rokka.Operation.Validate.
func (o GrayscaleOperation) Validate() (bool, error) {
	return true, nil
}

// toURLPath implements rokka.Operation.toURLPath.
func (o GrayscaleOperation) toURLPath() string {
	options := make([]string, 0)

	if len(options) == 0 {
		return o.Name()
	}
	return fmt.Sprintf("%s-%s", o.Name(), strings.Join(options, "-"))
}

// NoopOperation is an auto-generated Operation as specified by the rokka API.
//
// See: https://rokka.io/documentation/references/operations.html
type NoopOperation struct {
}

// Name implements rokka.Operation.Name
func (o NoopOperation) Name() string { return "noop" }

// Validate implements rokka.Operation.Validate.
func (o NoopOperation) Validate() (bool, error) {
	return true, nil
}

// toURLPath implements rokka.Operation.toURLPath.
func (o NoopOperation) toURLPath() string {
	options := make([]string, 0)

	if len(options) == 0 {
		return o.Name()
	}
	return fmt.Sprintf("%s-%s", o.Name(), strings.Join(options, "-"))
}

// PrimitiveOperation is an auto-generated Operation as specified by the rokka API.
//
// See: https://rokka.io/documentation/references/operations.html
type PrimitiveOperation struct {
	Count *int `json:"count,omitempty"`
	Mode  *int `json:"mode,omitempty"`
}

// Name implements rokka.Operation.Name
func (o PrimitiveOperation) Name() string { return "primitive" }

// Validate implements rokka.Operation.Validate.
func (o PrimitiveOperation) Validate() (bool, error) {
	return true, nil
}

// toURLPath implements rokka.Operation.toURLPath.
func (o PrimitiveOperation) toURLPath() string {
	options := make([]string, 0)
	if o.Count != nil {
		options = append(options, fmt.Sprintf("%s-%v", "count", *o.Count))
	}
	if o.Mode != nil {
		options = append(options, fmt.Sprintf("%s-%v", "mode", *o.Mode))
	}

	if len(options) == 0 {
		return o.Name()
	}
	return fmt.Sprintf("%s-%s", o.Name(), strings.Join(options, "-"))
}

// ResizeOperation is an auto-generated Operation as specified by the rokka API.
// Calling .Validate() will return false if required properties are missing.
//
// See: https://rokka.io/documentation/references/operations.html
type ResizeOperation struct {
	Height     *int    `json:"height,omitempty"`
	Mode       *string `json:"mode,omitempty"`
	Upscale    *bool   `json:"upscale,omitempty"`
	UpscaleDpr *bool   `json:"upscale_dpr,omitempty"`
	Width      *int    `json:"width,omitempty"`
}

// Name implements rokka.Operation.Name
func (o ResizeOperation) Name() string { return "resize" }

// Validate implements rokka.Operation.Validate.
func (o ResizeOperation) Validate() (bool, error) {
	valid := false
	if o.Width != nil {
		valid = true
	}
	if o.Height != nil {
		valid = true
	}
	if !valid {
		return false, errors.New("one of \"[width height]\" is required")
	}
	return true, nil
}

// toURLPath implements rokka.Operation.toURLPath.
func (o ResizeOperation) toURLPath() string {
	options := make([]string, 0)
	if o.Height != nil {
		options = append(options, fmt.Sprintf("%s-%v", "height", *o.Height))
	}
	if o.Mode != nil {
		options = append(options, fmt.Sprintf("%s-%v", "mode", *o.Mode))
	}
	if o.Upscale != nil {
		options = append(options, fmt.Sprintf("%s-%v", "upscale", *o.Upscale))
	}
	if o.UpscaleDpr != nil {
		options = append(options, fmt.Sprintf("%s-%v", "upscale_dpr", *o.UpscaleDpr))
	}
	if o.Width != nil {
		options = append(options, fmt.Sprintf("%s-%v", "width", *o.Width))
	}

	if len(options) == 0 {
		return o.Name()
	}
	return fmt.Sprintf("%s-%s", o.Name(), strings.Join(options, "-"))
}

// RotateOperation is an auto-generated Operation as specified by the rokka API.
// Calling .Validate() will return false if required properties are missing.
//
// See: https://rokka.io/documentation/references/operations.html
type RotateOperation struct {
	Angle             *float64 `json:"angle,omitempty"`
	BackgroundColor   *string  `json:"background_color,omitempty"`
	BackgroundOpacity *float64 `json:"background_opacity,omitempty"`
}

// Name implements rokka.Operation.Name
func (o RotateOperation) Name() string { return "rotate" }

// Validate implements rokka.Operation.Validate.
func (o RotateOperation) Validate() (bool, error) {
	if o.Angle == nil {
		return false, errors.New("option \"Angle\" is required")
	}
	return true, nil
}

// toURLPath implements rokka.Operation.toURLPath.
func (o RotateOperation) toURLPath() string {
	options := make([]string, 0)
	if o.Angle != nil {
		options = append(options, fmt.Sprintf("%s-%v", "angle", *o.Angle))
	}
	if o.BackgroundColor != nil {
		options = append(options, fmt.Sprintf("%s-%v", "background_color", *o.BackgroundColor))
	}
	if o.BackgroundOpacity != nil {
		options = append(options, fmt.Sprintf("%s-%v", "background_opacity", *o.BackgroundOpacity))
	}

	if len(options) == 0 {
		return o.Name()
	}
	return fmt.Sprintf("%s-%s", o.Name(), strings.Join(options, "-"))
}

// SepiaOperation is an auto-generated Operation as specified by the rokka API.
//
// See: https://rokka.io/documentation/references/operations.html
type SepiaOperation struct {
}

// Name implements rokka.Operation.Name
func (o SepiaOperation) Name() string { return "sepia" }

// Validate implements rokka.Operation.Validate.
func (o SepiaOperation) Validate() (bool, error) {
	return true, nil
}

// toURLPath implements rokka.Operation.toURLPath.
func (o SepiaOperation) toURLPath() string {
	options := make([]string, 0)

	if len(options) == 0 {
		return o.Name()
	}
	return fmt.Sprintf("%s-%s", o.Name(), strings.Join(options, "-"))
}

// TrimOperation is an auto-generated Operation as specified by the rokka API.
//
// See: https://rokka.io/documentation/references/operations.html
type TrimOperation struct {
	Fuzzy *float64 `json:"fuzzy,omitempty"`
}

// Name implements rokka.Operation.Name
func (o TrimOperation) Name() string { return "trim" }

// Validate implements rokka.Operation.Validate.
func (o TrimOperation) Validate() (bool, error) {
	return true, nil
}

// toURLPath implements rokka.Operation.toURLPath.
func (o TrimOperation) toURLPath() string {
	options := make([]string, 0)
	if o.Fuzzy != nil {
		options = append(options, fmt.Sprintf("%s-%v", "fuzzy", *o.Fuzzy))
	}

	if len(options) == 0 {
		return o.Name()
	}
	return fmt.Sprintf("%s-%s", o.Name(), strings.Join(options, "-"))
}
