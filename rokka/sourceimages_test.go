package rokka

import (
	"net/http"
	"testing"

	"github.com/rokka-io/rokka-go/test"
)

func TestListSourceImages(t *testing.T) {
	ts := test.NewMockAPI("./fixtures/ListSourceImages.json", http.StatusOK)
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.ListSourceImages("test", nil)
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestListSourceImagesWithLimitAndOffset(t *testing.T) {
	ts := test.NewMockAPI("./fixtures/ListSourceImagesWithLimitAndOffset.json", http.StatusOK)
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.ListSourceImages("test", map[string]string{"limit": "10", "offset": "20"})
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
