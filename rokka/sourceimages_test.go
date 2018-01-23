package rokka

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/rokka-io/rokka-go/test"
)

func TestListSourceImages(t *testing.T) {
	ts := test.NewMockAPI("./fixtures/ListSourceImages.json", http.StatusOK)
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.ListSourceImages("test", ListSourceImagesOptions{})
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestListSourceImagesWithQueryParams(t *testing.T) {
	ts := test.NewMockAPI("./fixtures/ListSourceImagesWithLimitAndOffset.json", http.StatusOK)
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.ListSourceImages("test", ListSourceImagesOptions{
		Limit:      10,
		Offset:     20,
		Hash:       "73ecc577d1c51941647378f3460675b6ad7c4fff",
		BinaryHash: "b9914b12d668dfb6e35fe85fd4a52be1df4aa9ff",
		Size:       "39189",
		Format:     "png",
		Width:      "1920",
		Height:     "960",
		Created:    "2017-11-14T10:10:40+00:00",
	})
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestGetSourceImage(t *testing.T) {
	ts := test.NewMockAPI("./fixtures/GetSourceImage.json", http.StatusOK)
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.GetSourceImage("test", "hash")
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestCreateSourceImage(t *testing.T) {
	ts := test.NewMockAPI("./fixtures/CreateSourceImage.json", http.StatusOK)
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	file, err := os.Open("./fixtures/image.png")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	res, err := c.CreateSourceImage("test", "image", file)
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestCreateSourceImageWithMetadata(t *testing.T) {
	ts := test.NewMockAPI("./fixtures/CreateSourceImageWithMetadata.json", http.StatusOK)
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	file, err := os.Open("./fixtures/image.png")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	userMetadata := map[string]interface{}{"key1": "value1"}
	dynamicMetadata := map[string]interface{}{"subject_area": map[string]int{"x": 50, "y": 50, "width": 10, "height": 10}}
	res, err := c.CreateSourceImageWithMetadata("test", "image", file, userMetadata, dynamicMetadata)
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestAddDynamicMetadata(t *testing.T) {
	loc := "https://api.example.org/test/1234-2"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", loc)
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte{})
	}))
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.AddDynamicMetadata("test", "1234", "test-name", bytes.NewBufferString("{\"test\": \"testing\""), AddDynamicMetadataOptions{})
	if err != nil {
		t.Error(err)
	}

	if res.Location != loc {
		t.Errorf("Expected location to be parsed in response, want: '%s', got: '%s'", loc, res.Location)
	}
}
