package xsdgen

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

type testLogger testing.T

func grep(pattern, data string) bool {
	matched, err := regexp.MatchString(pattern, data)
	if err != nil {
		panic(err)
	}
	return matched
}

func (t *testLogger) Printf(format string, v ...interface{}) {
	t.Logf(format, v...)
}

func TestLibrarySchema(t *testing.T) {
	testGen(t, "-ns http://dyomedea.com/ns/library testdata/library.xsd")
}
func TestPurchasOrderSchema(t *testing.T) {
	testGen(t, "-ns http://www.example.com/PO1 testdata/po1.xsd")
}
func TestUSTreasureSDN(t *testing.T) {
	testGen(t, "-ns http://tempuri.org/sdnList.xsd testdata/sdn.xsd")
}
func TestSoap(t *testing.T) {
	testGen(t, "-ns http://schemas.xmlsoap.org/soap/encoding/ testdata/soap11.xsd")
}
func TestSimpleStruct(t *testing.T) {
	testGen(t, "-ns http://example.org/ns testdata/simple-struct.xsd")
}
func testGen(t *testing.T, command string) string {
	file, err := os.CreateTemp("", "xsdgen")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	var cfg Config
	cfg.Option(DefaultOptions...)
	cfg.Option(LogOutput((*testLogger)(t)))

	err = cfg.GenCLI(append([]string{"-o", file.Name()}, strings.Split(command, " ")...)...)
	if err != nil {
		t.Error(err)
	}
	data, err := os.ReadFile(file.Name())
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestMixedType(t *testing.T) {
	data := testGen(t, "-ns http://example.org testdata/mixed-complex.xsd")
	if !grep(`PositiveNumber[^}]*,chardata`, data) {
		t.Errorf("type decl for PositiveNumber did not contain chardata, got \n%s", data)
	} else {
		t.Logf("got \n%s", data)
	}
}

func TestBase64Binary(t *testing.T) {
	t.Logf("%s\n", testGen(t, "-ns http://example.org/ testdata/base64.xsd"))
}

func TestSimpleUnion(t *testing.T) {
	t.Logf("%s\n", testGen(t, "-ns http://example.org/ testdata/simple-union.xsd"))
}

func TestImports(t *testing.T) {
	t.Logf("%s\n", testGen(t, "-f -ns ns1 testdata/ns1.xsd"))
}
