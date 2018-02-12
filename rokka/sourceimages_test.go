package rokka

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/rokka-io/rokka-go/test"
)

func TestListSourceImages(t *testing.T) {
	org := "test"
	ts := test.NewMockAPI(test.Routes{"GET /sourceimages/" + org: test.Response{http.StatusOK, "./fixtures/ListSourceImages.json", nil}})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.ListSourceImages(org, ListSourceImagesOptions{})
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestListSourceImagesWithQueryParams(t *testing.T) {
	org := "test"
	ts := test.NewMockAPI(test.Routes{"GET /sourceimages/" + org: test.Response{http.StatusOK, "./fixtures/ListSourceImagesWithLimitAndOffset.json", nil}})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.ListSourceImages(org, ListSourceImagesOptions{
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
	org := "test"
	hash := "hash"
	ts := test.NewMockAPI(test.Routes{"GET /sourceimages/" + org + "/" + hash: test.Response{http.StatusOK, "./fixtures/GetSourceImage.json", nil}})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.GetSourceImage(org, hash)
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestCreateSourceImage(t *testing.T) {
	org := "test"
	ts := test.NewMockAPI(test.Routes{"POST /sourceimages/" + org: test.Response{http.StatusOK, "./fixtures/CreateSourceImage.json", nil}})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	file, err := os.Open("./fixtures/image.png")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	res, err := c.CreateSourceImage(org, "image", file)
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestCreateSourceImageWithMetadata(t *testing.T) {
	org := "test"
	ts := test.NewMockAPI(test.Routes{"POST /sourceimages/" + org: test.Response{http.StatusOK, "./fixtures/CreateSourceImageWithMetadata.json", nil}})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	file, err := os.Open("./fixtures/image.png")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	userMetadata := map[string]interface{}{"key1": "value1"}
	dynamicMetadata := map[string]interface{}{"subject_area": map[string]int{"x": 50, "y": 50, "width": 10, "height": 10}}
	res, err := c.CreateSourceImageWithMetadata(org, "image", file, userMetadata, dynamicMetadata)
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestAddDynamicMetadata(t *testing.T) {
	loc := "https://api.example.org/test/1234-2"
	org := "test"
	hash := "1234"
	metaName := "test-name"
	path := fmt.Sprintf("/sourceimages/%s/%s/meta/dynamic/%s", org, hash, metaName)
	headers := map[string]string{"Location": loc}
	ts := test.NewMockAPI(test.Routes{"PUT " + path: test.Response{http.StatusCreated, "", headers}})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.AddDynamicMetadata(org, hash, metaName, bytes.NewBufferString("{\"test\": \"testing\""), DynamicMetadataOptions{})
	if err != nil {
		t.Error(err)
	}

	if res.Location != loc {
		t.Errorf("Expected location to be parsed in response, want: '%s', got: '%s'", loc, res.Location)
	}
}

func TestDeleteDynamicMetadata(t *testing.T) {
	loc := "https://api.example.org/test/1234-2"
	org := "test"
	hash := "1234"
	metaName := "test-name"
	path := fmt.Sprintf("/sourceimages/%s/%s/meta/dynamic/%s", org, hash, metaName)
	headers := map[string]string{"Location": loc}
	ts := test.NewMockAPI(test.Routes{"DELETE " + path: test.Response{http.StatusNoContent, "", headers}})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.DeleteDynamicMetadata(org, hash, metaName, DynamicMetadataOptions{})
	if err != nil {
		t.Error(err)
	}

	if res.Location != loc {
		t.Errorf("Expected location to be parsed in response, want: '%s', got: '%s'", loc, res.Location)
	}
}
