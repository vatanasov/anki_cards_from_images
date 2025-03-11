package anki_cards_from_images

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

var dname string

func TestMain(m *testing.M) {
	var err error
	dname, err = os.MkdirTemp("", "examples_dir")
	if err != nil {
		panic(err)
	}

	result := m.Run()

	os.RemoveAll(dname)

	os.Exit(result)
}

func TestProcess(t *testing.T) {
	w := &bytes.Buffer{}
	t.Run("Error when invalid file", func(t *testing.T) {
		fname := filepath.Join(dname, "invalid")
		err := os.WriteFile(fname, []byte{}, 0600)
		if err != nil {
			t.Fatal(err)
		}
		err = Process(dname, w)
		if err == nil {
			t.Fatal("Did not error out when processing invalid file")
		}
		os.Remove(fname)
	})

	t.Run("Error when missing answer file", func(t *testing.T) {
		fname := filepath.Join(dname, "1_q.jpg")
		err := os.WriteFile(fname, []byte{}, 0600)
		if err != nil {
			t.Fatal(err)
		}
		err = Process(dname, w)
		if err == nil {
			t.Fatal("Did not error out when processing invalid file")
		}
		os.Remove(fname)
	})

	t.Run("Success", func(t *testing.T) {
		fname := filepath.Join(dname, "1_q.jpg")
		err := os.WriteFile(fname, []byte{}, 0600)
		if err != nil {
			t.Fatal(err)
		}

		aFname := filepath.Join(dname, "1_a.jpg")
		err = os.WriteFile(aFname, []byte{}, 0600)
		if err != nil {
			t.Fatal(err)
		}
		err = Process(dname, w)
		if err != nil {
			t.Fatal(err)
		}
		resultBuffer := &bytes.Buffer{}
		resultWriter := csv.NewWriter(resultBuffer)
		resultWriter.WriteAll([][]string{
			{fmt.Sprintf("<img src=\"%s\">", fname),
				fmt.Sprintf("<img src=\"%s\">", aFname)},
		})
		expected := resultBuffer.String()
		if w.String() != expected {
			t.Fatalf("Wrong result, expected \n%s, got \n%s", expected, w.String())
		}
		os.Remove(fname)
	})
}
