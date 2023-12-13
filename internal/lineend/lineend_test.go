package lineend

import (
	"testing"
)

type testDataType struct {
	Input  string
	Expect string
}

func TestGetLineEnd(t *testing.T) {
	testData := []testDataType{
		{
			Input:  "",
			Expect: "",
		},
		{
			Input:  "my line",
			Expect: "",
		},
		{
			Input:  "my\r\nline\r\n",
			Expect: "\r\n",
		},
		{
			Input:  "my\rline\r",
			Expect: "\r",
		},
		{
			Input:  "my\nline\n",
			Expect: "\n",
		},
	}

	for i, test := range testData {
		if GetLineEnd(test.Input) != test.Expect {
			t.Fatal("Number of failed test:", i)
		}
	}
}

func TestUnknownToDos(t *testing.T) {
	testData := []testDataType{
		{
			Input:  "",
			Expect: "",
		},
		{
			Input:  "my line",
			Expect: "my line",
		},
		{
			Input:  "my\r\nline\r\n",
			Expect: "my\r\nline\r\n",
		},
		{
			Input:  "my\rline\r",
			Expect: "my\r\nline\r\n",
		},
		{
			Input:  "my\nline\n",
			Expect: "my\r\nline\r\n",
		},
	}

	for i, test := range testData {
		if UnknownToDos(test.Input) != test.Expect {
			t.Fatal("Number of failed test:", i)
		}
	}
}

func TestUnknownToOldMac(t *testing.T) {
	testData := []testDataType{
		{
			Input:  "",
			Expect: "",
		},
		{
			Input:  "my line",
			Expect: "my line",
		},
		{
			Input:  "my\r\nline\r\n",
			Expect: "my\rline\r",
		},
		{
			Input:  "my\rline\r",
			Expect: "my\rline\r",
		},
		{
			Input:  "my\nline\n",
			Expect: "my\rline\r",
		},
	}

	for i, test := range testData {
		if UnknownToOldMac(test.Input) != test.Expect {
			t.Fatal("Number of failed test:", i)
		}
	}
}

func TestUnknownToUnix(t *testing.T) {
	testData := []testDataType{
		{
			Input:  "",
			Expect: "",
		},
		{
			Input:  "my line",
			Expect: "my line",
		},
		{
			Input:  "my\r\nline\r\n",
			Expect: "my\nline\n",
		},
		{
			Input:  "my\rline\r",
			Expect: "my\nline\n",
		},
		{
			Input:  "my\nline\n",
			Expect: "my\nline\n",
		},
	}

	for i, test := range testData {
		if UnknownToUnix(test.Input) != test.Expect {
			t.Fatal("Number of failed test:", i)
		}
	}
}
