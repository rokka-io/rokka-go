package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/rokka-io/rokka-go/rokka"
	"github.com/rokka-io/rokka-go/test"
)

func TestRunCodeGenerator(t *testing.T) {
	ts := test.NewMockAPI(test.Routes{"GET /operations": test.Response{http.StatusOK, "./testdata/GetOperations.json", nil}})
	defer ts.Close()

	snapshotFileName := "./testdata/operations_object.go.snapshot"
	_, updateSnapshot := os.LookupEnv("UPDATE_SNAPSHOT")

	if _, err := os.Stat(snapshotFileName); os.IsNotExist(err) {
		if !updateSnapshot {
			t.Errorf("Snapshot file %s does not exist. Run `UPDATE_SNAPSHOT=1 go test ./... -run '^TestRunCodeGenerator$'` to create it.", snapshotFileName)
			return
		}
		_, err := os.Create(snapshotFileName)
		if err != nil {
			panic(err)
		}
	}

	snapshot, err := os.Open(snapshotFileName)
	if err != nil {
		panic(err)
	}
	defer snapshot.Close()

	f, err := ioutil.TempFile(os.TempDir(), "operations_object.go")
	if err != nil {
		panic(err)
	}
	defer os.Remove(f.Name())

	generate(&rokka.Config{APIAddress: ts.URL}, f.Name())

	if updateSnapshot {
		// TODO: continue
		os.Rename(f.Name(), snapshotFileName)
		fmt.Println("Updated snapshot")
		return
	}

	snapshotScanner := bufio.NewScanner(snapshot)
	fScanner := bufio.NewScanner(f)

	success := true
	line := 1
	for snapshotScanner.Scan() {
		line++
		snapshotLine := snapshotScanner.Text()
		if !fScanner.Scan() {
			t.Logf("(%d) - %s", line, snapshotLine)
			success = false
			continue
		}
		fLine := fScanner.Text()
		if snapshotLine != fLine && !strings.HasPrefix(snapshotLine, "// This file was generated at") {
			t.Logf("(%d) - %s", line, snapshotLine)
			t.Logf("(%d) + %s", line, fLine)
			success = false
		}
	}
	// in case newly generated file is longer than existing one.
	for fScanner.Scan() {
		line++
		success = false
		t.Logf("(%d) + %s", line, fScanner.Text())
	}

	if !success {
		t.Errorf("Snapshot isn't equal to generated code. Please check the diff written. To update the snapshot, execute `UPDATE_SNAPSHOT=1 go test ./... -run '^TestRunCodeGenerator$'`.")
	}
}
