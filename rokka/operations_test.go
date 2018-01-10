package rokka

import (
	"net/http"
	"testing"

	"github.com/rokka-io/rokka-go/test"
)

func TestGetOperations(t *testing.T) {
	ts := test.NewMockAPI("./fixtures/GetOperations.json", http.StatusOK)
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.GetStackOptions()
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func sref(v string) *string { return &v }
func iref(v int) *int { return &v }
func fref(v float64) *float64 { return &v }

var operationObjectsTests = []struct{
	name string
	op Operation
	isValid bool
	urlPath string
}{
	{"AlphaOperation without args", AlphaOperation{}, true, "alpha"},
	{"AlphaOperation with valid args", AlphaOperation{sref("x")}, true, "alpha-mode-x"},
	{"AutorotateOperation without args", AutorotateOperation{}, true, "autorotate"},
	{"AutorotateOperation with single arg", AutorotateOperation{Width: iref(10)}, true, "autorotate-width-10"},
	{"AutorotateOperation with multiple args", AutorotateOperation{Width: iref(10), Height: iref(20), RotationDirection: sref("clockwise")}, true, "autorotate-height-20-rotation_direction-clockwise-width-10"},
	{"BlurOperation with missing args", BlurOperation{}, false, "blur"},
	{"BlurOperation with valid args", BlurOperation{Sigma: fref(1.337)}, true, "blur-sigma-1.337"},
	{"CompositionOperation missing args", CompositionOperation{}, false, "composition"},
	{"CompositionOperation with valid args", CompositionOperation{Anchor: sref("top"), Height: iref(10), Width: iref(20), Mode: sref("test")}, true, "composition-anchor-top-height-10-mode-test-width-20"},
	{"CropOperation with missing args", CropOperation{}, false, "crop"},
	{"CropOperation with valid args", CropOperation{Height: iref(100), Width: iref(200)}, true, "crop-height-100-width-200"},
	{"DropshadowOperation without args", DropshadowOperation{}, true, "dropshadow"},
	{"DropshadowOperation with valid args", DropshadowOperation{Color: sref("ffffff"), Vertical: iref(10)}, true, "dropshadow-color-ffffff-vertical-10"},
	{"GrayscaleOperation without args", GrayscaleOperation{}, true, "grayscale"},
	{"NoopOperation without args", NoopOperation{}, true, "noop"},
	{"PrimitiveOperation without args", PrimitiveOperation{}, true, "primitive"},
	{"PrimitiveOperation with args", PrimitiveOperation{Count: iref(5), Mode: iref(2)}, true, "primitive-count-5-mode-2"},
	{"ResizeOperation without args", ResizeOperation{}, false, "resize"},
	{"ResizeOperation with arg (one-of #1)", ResizeOperation{Height: iref(10)}, true, "resize-height-10"},
	{"ResizeOperation with arg (one-of #2)", ResizeOperation{Width: iref(10)}, true, "resize-width-10"},
	{"ResizeOperation with args", ResizeOperation{Height: iref(10), Width: iref(10)}, true, "resize-height-10-width-10"},
	{"RotateOperation without args", RotateOperation{}, false, "rotate"},
	{"RotateOperation without required arg", RotateOperation{BackgroundColor: sref("aa9374")}, false, "rotate-background_color-aa9374"},
	{"RotateOperation with args", RotateOperation{Angle: fref(45)}, true, "rotate-angle-45"},
	{"SepiaOperation without args", SepiaOperation{}, true, "sepia"},
	{"TrimOperation without args", TrimOperation{}, true, "trim"},
	{"TrimOperation with args", TrimOperation{Fuzzy: fref(15)}, true, "trim-fuzzy-15"},
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