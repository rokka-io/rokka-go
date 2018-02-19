package rokka

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/rokka-io/rokka-go/test"
)

func TestListSourceImages(t *testing.T) {
	org := "test"
	r := test.NewResponse(http.StatusOK, "./fixtures/ListSourceImages.json")
	ts := test.NewMockAPI(t, test.Routes{"GET /sourceimages/" + org: r})
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
	r := test.NewResponse(http.StatusOK, "./fixtures/ListSourceImagesWithLimitAndOffset.json")
	ts := test.NewMockAPI(t, test.Routes{"GET /sourceimages/" + org: r})
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
	r := test.NewResponse(http.StatusOK, "./fixtures/GetSourceImage.json")
	ts := test.NewMockAPI(t, test.Routes{"GET /sourceimages/" + org + "/" + hash: r})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.GetSourceImage(org, hash)
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestDeleteSourceImage(t *testing.T) {
	org := "test"
	hash := "hash"
	r := test.NewResponse(http.StatusNoContent, "")
	ts := test.NewMockAPI(t, test.Routes{"DELETE /sourceimages/" + org + "/" + hash: r})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	err := c.DeleteSourceImage(org, hash)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteSourceImageByBinaryHash(t *testing.T) {
	org := "test"
	hash := "hash"
	r := test.NewResponse(http.StatusNoContent, "")
	ts := test.NewMockAPI(t, test.Routes{"DELETE /sourceimages/" + org + "?binaryHash=" + hash: r})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	err := c.DeleteSourceImageByBinaryHash(org, hash)
	if err != nil {
		t.Error(err)
	}
}

func TestDownloadSourceImage(t *testing.T) {
	org := "test"
	hash := "hash"
	r := test.NewResponse(http.StatusOK, "./fixtures/image.png")
	r.Headers["Content-Disposition"] = `attachment; filename="image.png"`
	ts := test.NewMockAPI(t, test.Routes{"GET /sourceimages/" + org + "/" + hash + "/download": r})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.DownloadSourceImage(org, hash)
	if err != nil {
		t.Error(err)
	}
	if res.FileName != "image.png" {
		t.Errorf("Expected FileName to be '%s', got '%s'", "image.png", res.FileName)
	}

	d, err := ioutil.ReadAll(res.Data)
	if err != nil {
		t.Error(err)
	}
	defer res.Data.Close()

	if len(d) != 289 {
		t.Errorf("Expected length of Data to be '%d', got '%d'", 289, len(d))
	}
}

func TestCreateSourceImage(t *testing.T) {
	org := "test"
	r := test.NewResponse(http.StatusOK, "./fixtures/CreateSourceImage.json")
	ts := test.NewMockAPI(t, test.Routes{"POST /sourceimages/" + org: r})
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
	r := test.NewResponse(http.StatusOK, "./fixtures/CreateSourceImageWithMetadata.json")
	ts := test.NewMockAPI(t, test.Routes{"POST /sourceimages/" + org: r})
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
	r := test.NewResponse(http.StatusCreated, "")
	r.Headers["Location"] = loc
	ts := test.NewMockAPI(t, test.Routes{"PUT " + path: r})
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
	r := test.NewResponse(http.StatusNoContent, "")
	r.Headers["Location"] = loc
	ts := test.NewMockAPI(t, test.Routes{"DELETE " + path: r})
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

func TestUpdateUserMetadata(t *testing.T) {
	org := "test"
	hash := "1234"
	path := fmt.Sprintf("/sourceimages/%s/%s/meta/user", org, hash)
	r := test.NewResponse(http.StatusNoContent, "")
	ts := test.NewMockAPI(t, test.Routes{"PATCH " + path: r})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	err := c.UpdateUserMetadata(org, hash, bytes.NewBufferString("{\"test\": \"testing\""))
	if err != nil {
		t.Error(err)
	}
}

func TestSetUserMetadata(t *testing.T) {
	org := "test"
	hash := "1234"
	path := fmt.Sprintf("/sourceimages/%s/%s/meta/user", org, hash)
	r := test.NewResponse(http.StatusNoContent, "")
	ts := test.NewMockAPI(t, test.Routes{"PUT " + path: r})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	err := c.SetUserMetadata(org, hash, bytes.NewBufferString("{\"test\": \"testing\""))
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteUserMetadata(t *testing.T) {
	org := "test"
	hash := "1234"
	path := fmt.Sprintf("/sourceimages/%s/%s/meta/user", org, hash)
	r := test.NewResponse(http.StatusNoContent, "")
	ts := test.NewMockAPI(t, test.Routes{"DELETE " + path: r})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	err := c.DeleteUserMetadata(org, hash)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateUserMetadataByName(t *testing.T) {
	org := "test"
	hash := "1234"
	metaName := "test"
	path := fmt.Sprintf("/sourceimages/%s/%s/meta/user/%s", org, hash, metaName)
	r := test.NewResponse(http.StatusNoContent, "")
	ts := test.NewMockAPI(t, test.Routes{"PUT " + path: r})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	err := c.UpdateUserMetadataByName(org, hash, metaName, bytes.NewBufferString(`"testing"`))
	if err != nil {
		t.Error(err)
	}
}
func TestDeleteUserMetadataByName(t *testing.T) {
	org := "test"
	hash := "1234"
	metaName := "test"
	path := fmt.Sprintf("/sourceimages/%s/%s/meta/user/%s", org, hash, metaName)
	r := test.NewResponse(http.StatusNoContent, "")
	ts := test.NewMockAPI(t, test.Routes{"DELETE " + path: r})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	err := c.DeleteUserMetadataByName(org, hash, metaName)
	if err != nil {
		t.Error(err)
	}
}
