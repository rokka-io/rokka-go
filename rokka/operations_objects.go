// Code generated by go generate; DO NOT EDIT.
// This file was generated at 2018-01-10 12:00:18.554761 +0100 CET m=+0.259618653
package rokka

import (
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

// AlphaOperation is an auto-generated Operation as specified by the rokka API.
//
// See: https://rokka.io/documentation/references/operations.html
type AlphaOperation struct {
	Mode *string
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
		options = append(options, fmt.Sprintf("%s", *o.Mode))
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
	Width             *int
	Height            *int
	RotationDirection *string
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
	if o.Width != nil {
		options = append(options, fmt.Sprintf("%s", *o.Width))
	}
	if o.Height != nil {
		options = append(options, fmt.Sprintf("%s", *o.Height))
	}
	if o.RotationDirection != nil {
		options = append(options, fmt.Sprintf("%s", *o.RotationDirection))
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
	Sigma *float64
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
		options = append(options, fmt.Sprintf("%s", *o.Sigma))
	}

	if len(options) == 0 {
		return o.Name()
	}
	return fmt.Sprintf("%s-%s", o.Name(), strings.Join(options, "-"))
}

// CompositionOperation is an auto-generated Operation as specified by the rokka API.
// Calling .Validate() will return false if required properties are missing.
//
// See: https://rokka.io/documentation/references/operations.html
type CompositionOperation struct {
	Mode             *string
	Width            *int
	Height           *int
	Anchor           *string
	SecondaryColor   *string
	SecondaryOpacity *int
}

// Name implements rokka.Operation.Name
func (o CompositionOperation) Name() string { return "composition" }

// Validate implements rokka.Operation.Validate.
func (o CompositionOperation) Validate() (bool, error) {
	if o.Mode == nil {
		return false, errors.New("option \"Mode\" is required")
	}
	if o.Width == nil {
		return false, errors.New("option \"Width\" is required")
	}
	if o.Height == nil {
		return false, errors.New("option \"Height\" is required")
	}
	return true, nil
}

// toURLPath implements rokka.Operation.toURLPath.
func (o CompositionOperation) toURLPath() string {
	options := make([]string, 0)
	if o.Mode != nil {
		options = append(options, fmt.Sprintf("%s", *o.Mode))
	}
	if o.Width != nil {
		options = append(options, fmt.Sprintf("%s", *o.Width))
	}
	if o.Height != nil {
		options = append(options, fmt.Sprintf("%s", *o.Height))
	}
	if o.Anchor != nil {
		options = append(options, fmt.Sprintf("%s", *o.Anchor))
	}
	if o.SecondaryColor != nil {
		options = append(options, fmt.Sprintf("%s", *o.SecondaryColor))
	}
	if o.SecondaryOpacity != nil {
		options = append(options, fmt.Sprintf("%s", *o.SecondaryOpacity))
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
	Scale  *float64
	Width  *int
	Height *int
	Anchor *string
	Mode   *string
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
	if o.Scale != nil {
		options = append(options, fmt.Sprintf("%s", *o.Scale))
	}
	if o.Width != nil {
		options = append(options, fmt.Sprintf("%s", *o.Width))
	}
	if o.Height != nil {
		options = append(options, fmt.Sprintf("%s", *o.Height))
	}
	if o.Anchor != nil {
		options = append(options, fmt.Sprintf("%s", *o.Anchor))
	}
	if o.Mode != nil {
		options = append(options, fmt.Sprintf("%s", *o.Mode))
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
	Horizontal *int
	Vertical   *int
	Opacity    *int
	Sigma      *float64
	BlurRadius *float64
	Color      *string
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
	if o.Horizontal != nil {
		options = append(options, fmt.Sprintf("%s", *o.Horizontal))
	}
	if o.Vertical != nil {
		options = append(options, fmt.Sprintf("%s", *o.Vertical))
	}
	if o.Opacity != nil {
		options = append(options, fmt.Sprintf("%s", *o.Opacity))
	}
	if o.Sigma != nil {
		options = append(options, fmt.Sprintf("%s", *o.Sigma))
	}
	if o.BlurRadius != nil {
		options = append(options, fmt.Sprintf("%s", *o.BlurRadius))
	}
	if o.Color != nil {
		options = append(options, fmt.Sprintf("%s", *o.Color))
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
	Count *int
	Mode  *int
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
		options = append(options, fmt.Sprintf("%s", *o.Count))
	}
	if o.Mode != nil {
		options = append(options, fmt.Sprintf("%s", *o.Mode))
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
	Width      *int
	Height     *int
	Mode       *string
	Upscale    *bool
	UpscaleDpr *bool
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
	if o.Width != nil {
		options = append(options, fmt.Sprintf("%s", *o.Width))
	}
	if o.Height != nil {
		options = append(options, fmt.Sprintf("%s", *o.Height))
	}
	if o.Mode != nil {
		options = append(options, fmt.Sprintf("%s", *o.Mode))
	}
	if o.Upscale != nil {
		options = append(options, fmt.Sprintf("%s", *o.Upscale))
	}
	if o.UpscaleDpr != nil {
		options = append(options, fmt.Sprintf("%s", *o.UpscaleDpr))
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
	Angle             *float64
	BackgroundColor   *string
	BackgroundOpacity *float64
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
		options = append(options, fmt.Sprintf("%s", *o.Angle))
	}
	if o.BackgroundColor != nil {
		options = append(options, fmt.Sprintf("%s", *o.BackgroundColor))
	}
	if o.BackgroundOpacity != nil {
		options = append(options, fmt.Sprintf("%s", *o.BackgroundOpacity))
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
	Fuzzy *float64
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
		options = append(options, fmt.Sprintf("%s", *o.Fuzzy))
	}

	if len(options) == 0 {
		return o.Name()
	}
	return fmt.Sprintf("%s-%s", o.Name(), strings.Join(options, "-"))
}
