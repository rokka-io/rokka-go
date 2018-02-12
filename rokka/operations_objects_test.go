package rokka

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestUnmarshalJSON_ContinueOnError(t *testing.T) {
	file, err := ioutil.ReadFile("./fixtures/UnmarshalJSON_ContinueOnError.json")
	if err != nil {
		t.Fatal(err)
	}

	ops := make(Operations, 0)
	err = json.Unmarshal(file, &ops)
	if err != nil {
		t.Fatal(err)
	}

	resizeOp, ok := ops[0].(*ResizeOperation)
	if !ok {
		t.Errorf("Expected operation of type '%T', got '%T'", new(ResizeOperation), ops[0])
	}
	expectedOp := ResizeOperation{
		Height:  IntPtr(0),
		Width:   IntPtr(0),
		Upscale: BoolPtr(false),
	}

	if *resizeOp.Height != *expectedOp.Height {
		t.Errorf("Expected height to be '%d', got '%d'", expectedOp.Height, resizeOp.Height)
	}
	if *resizeOp.Width != *expectedOp.Width {
		t.Errorf("Expected width to be '%d', got '%d'", expectedOp.Width, resizeOp.Width)
	}
	if *resizeOp.Upscale != *expectedOp.Upscale {
		t.Errorf("Expected Upscale to be '%t', got '%t'", *expectedOp.Upscale, *resizeOp.Upscale)
	}
	if resizeOp.Mode != expectedOp.Mode {
		t.Errorf("Expected Upscale to be '%v', got '%v'", expectedOp.Mode, resizeOp.Mode)
	}
}

var newOperationByNameTests = []struct {
	name    string
	isError bool
}{
	{"alpha", false},
	{"autorotate", false},
	{"blur", false},
	{"composition", false},
	{"crop", false},
	{"dropshadow", false},
	{"grayscale", false},
	{"noop", false},
	{"primitive", false},
	{"resize", false},
	{"rotate", false},
	{"sepia", false},
	{"trim", false},
	{"some operation which will never exist", true},
}

func TestNewOperationByName(t *testing.T) {
	for _, v := range newOperationByNameTests {
		t.Run(v.name, func(t *testing.T) {
			_, err := NewOperationByName(v.name)
			hasErr := err != nil
			if hasErr != v.isError {
				if v.isError {
					t.Error("Unexpectd result from NewOperationByName. Expected error, got no error")
				} else {
					t.Errorf("Unexpectd result from NewOperationByName. Expected no error, got error: %s", err)
				}
			}
		})
	}
}

var operationObjectsTests = []struct {
	name    string
	op      Operation
	isValid bool
	urlPath string
}{
	{"AlphaOperation without args", AlphaOperation{}, true, "alpha"},
	{"AlphaOperation with valid args", AlphaOperation{StrPtr("x")}, true, "alpha-mode-x"},
	{"AutorotateOperation without args", AutorotateOperation{}, true, "autorotate"},
	{"AutorotateOperation with single arg", AutorotateOperation{Width: IntPtr(10)}, true, "autorotate-width-10"},
	{"AutorotateOperation with multiple args", AutorotateOperation{Width: IntPtr(10), Height: IntPtr(20), RotationDirection: StrPtr("clockwise")}, true, "autorotate-height-20-rotation_direction-clockwise-width-10"},
	{"BlurOperation with missing args", BlurOperation{}, false, "blur"},
	{"BlurOperation with valid args", BlurOperation{Sigma: Float64Ptr(1.337)}, true, "blur-sigma-1.337"},
	{"CompositionOperation missing args", CompositionOperation{}, false, "composition"},
	{"CompositionOperation with valid args", CompositionOperation{Anchor: StrPtr("top"), Height: IntPtr(10), Width: IntPtr(20), Mode: StrPtr("test")}, true, "composition-anchor-top-height-10-mode-test-width-20"},
	{"CropOperation with missing args", CropOperation{}, false, "crop"},
	{"CropOperation with valid args", CropOperation{Height: IntPtr(100), Width: IntPtr(200)}, true, "crop-height-100-width-200"},
	{"DropshadowOperation without args", DropshadowOperation{}, true, "dropshadow"},
	{"DropshadowOperation with valid args", DropshadowOperation{Color: StrPtr("ffffff"), Vertical: IntPtr(10)}, true, "dropshadow-color-ffffff-vertical-10"},
	{"GrayscaleOperation without args", GrayscaleOperation{}, true, "grayscale"},
	{"NoopOperation without args", NoopOperation{}, true, "noop"},
	{"PrimitiveOperation without args", PrimitiveOperation{}, true, "primitive"},
	{"PrimitiveOperation with args", PrimitiveOperation{Count: IntPtr(5), Mode: IntPtr(2)}, true, "primitive-count-5-mode-2"},
	{"ResizeOperation without args", ResizeOperation{}, false, "resize"},
	{"ResizeOperation with arg (one-of #1)", ResizeOperation{Height: IntPtr(10)}, true, "resize-height-10"},
	{"ResizeOperation with arg (one-of #2)", ResizeOperation{Width: IntPtr(10)}, true, "resize-width-10"},
	{"ResizeOperation with args", ResizeOperation{Height: IntPtr(10), Width: IntPtr(10)}, true, "resize-height-10-width-10"},
	{"ResizeOperation with args and bool arg", ResizeOperation{Height: IntPtr(10), Width: IntPtr(10), Upscale: BoolPtr(true)}, true, "resize-height-10-upscale-true-width-10"},
	{"RotateOperation without args", RotateOperation{}, false, "rotate"},
	{"RotateOperation without required arg", RotateOperation{BackgroundColor: StrPtr("aa9374")}, false, "rotate-background_color-aa9374"},
	{"RotateOperation with args", RotateOperation{Angle: Float64Ptr(45)}, true, "rotate-angle-45"},
	{"SepiaOperation without args", SepiaOperation{}, true, "sepia"},
	{"TrimOperation without args", TrimOperation{}, true, "trim"},
	{"TrimOperation with args", TrimOperation{Fuzzy: Float64Ptr(15)}, true, "trim-fuzzy-15"},
}

func TestOperationsObjects(t *testing.T) {
	for _, v := range operationObjectsTests {
		t.Run(v.name, func(t *testing.T) {
			ok, err := v.op.Validate()
			if ok != v.isValid {
				t.Errorf("Unexpected result from Validate(). Result: %t; Expected: %t; Error: %s", ok, v.isValid, err)
			}

			p := v.op.toURLPath()
			if p != v.urlPath {
				t.Errorf("Unexpected result from toURLPath(). Result: \"%s\"; Expected: \"%s\"", p, v.urlPath)
			}
		})
	}
}
