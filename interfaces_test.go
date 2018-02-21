package main

import (
	"fmt"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type InterfacesSuite struct{}

func (s *InterfacesSuite) TestSimple(t sweet.T) {
	specs := getSpecs("simple")
	Expect(specs).To(HaveLen(2))
	Expect(specs).To(HaveKey("NoParams"))
	Expect(specs).To(HaveKey("Params"))

	//
	// NoparamsInterface

	noparams := specs["NoParams"]
	Expect(noparams.methods).To(HaveLen(3))

	Expect(noparams.methods).To(HaveKey("NoParamsNoResults"))
	Expect(noparams.methods["NoParamsNoResults"].results).To(BeEmpty())

	Expect(noparams.methods).To(HaveKey("NoParamsOneResult"))
	Expect(noparams.methods["NoParamsOneResult"].results).To(HaveLen(1))
	Expect(noparams.methods["NoParamsOneResult"].results[0].Text()).To(Equal("error"))

	Expect(noparams.methods).To(HaveKey("NoParamsMultipleResults"))
	Expect(noparams.methods["NoParamsMultipleResults"].results).To(HaveLen(2))
	Expect(noparams.methods["NoParamsMultipleResults"].results[0].Text()).To(Equal("[]string"))
	Expect(noparams.methods["NoParamsMultipleResults"].results[1].Text()).To(Equal("error"))

	//
	// paramsInterface

	params := specs["Params"]
	Expect(params.methods).To(HaveLen(4))

	Expect(params.methods).To(HaveKey("OneParam"))
	Expect(params.methods["OneParam"].params).To(HaveLen(1))
	Expect(params.methods["OneParam"].params[0].name).To(Equal("foo"))
	Expect(params.methods["OneParam"].params[0].typ.Text()).To(Equal("string"))

	Expect(params.methods).To(HaveKey("MultipleParams"))
	Expect(params.methods["MultipleParams"].params).To(HaveLen(3))
	Expect(params.methods["MultipleParams"].params[0].name).To(Equal("foo"))
	Expect(params.methods["MultipleParams"].params[0].typ.Text()).To(Equal("string"))
	Expect(params.methods["MultipleParams"].params[1].name).To(Equal("bar"))
	Expect(params.methods["MultipleParams"].params[1].typ.Text()).To(Equal("string"))
	Expect(params.methods["MultipleParams"].params[2].name).To(Equal("baz"))
	Expect(params.methods["MultipleParams"].params[2].typ.Text()).To(Equal("bool"))

	Expect(params.methods).To(HaveKey("Unnamed"))
	Expect(params.methods["Unnamed"].params).To(HaveLen(3))
	Expect(params.methods["Unnamed"].params[0].typ.Text()).To(Equal("string"))
	Expect(params.methods["Unnamed"].params[1].typ.Text()).To(Equal("string"))
	Expect(params.methods["Unnamed"].params[2].typ.Text()).To(Equal("bool"))

	Expect(params.methods).To(HaveKey("Variadic"))
	Expect(params.methods["Variadic"].params).To(HaveLen(2))
	Expect(params.methods["Variadic"].params[0].name).To(Equal("format"))
	Expect(params.methods["Variadic"].params[0].typ.Text()).To(Equal("string"))
	Expect(params.methods["Variadic"].params[1].name).To(Equal("params"))
	Expect(params.methods["Variadic"].params[1].typ.Text()).To(Equal("[]interface{}"))
	Expect(params.methods["Variadic"].variadic).To(BeTrue())
}

func (s *InterfacesSuite) TestComplex(t sweet.T) {
	specs := getSpecs("complex")
	Expect(specs).To(HaveLen(4))
	Expect(specs).To(HaveKey("SupInterface1"))
	Expect(specs).To(HaveKey("SupInterface2"))
	Expect(specs).To(HaveKey("SubInterface"))
	Expect(specs).To(HaveKey("EmbeddedTypes"))

	//
	// SubInterface

	sub := specs["SubInterface"]
	Expect(sub.methods).To(HaveLen(4))
	Expect(sub.methods).To(HaveKey("Foo"))
	Expect(sub.methods).To(HaveKey("Bar"))
	Expect(sub.methods).To(HaveKey("Baz"))
	Expect(sub.methods).To(HaveKey("String"))

	//
	// EmbeddedTypes

	embedded := specs["EmbeddedTypes"]
	Expect(embedded.methods).To(HaveLen(2))

	Expect(embedded.methods).To(HaveKey("Param"))
	Expect(embedded.methods["Param"].params).To(HaveLen(1))
	Expect(embedded.methods["Param"].params[0].typ.Text()).To(Equal("struct{X string; Y bool}"))

	Expect(embedded.methods).To(HaveKey("Result"))
	Expect(embedded.methods["Result"].results).To(HaveLen(2))
	Expect(embedded.methods["Result"].results[0].Text()).To(Equal("struct{Z int}"))
	Expect(embedded.methods["Result"].results[1].Text()).To(Equal("error"))
}

func (s *InterfacesSuite) TestLocaltypes(t sweet.T) {
	specs := getSpecs("localtypes")
	Expect(specs).To(HaveLen(3))
	Expect(specs).To(HaveKey("LocalRewrite"))
	Expect(specs).To(HaveKey("EmbeddedStruct"))
	Expect(specs).To(HaveKey("InterfaceStruct"))

	//
	// LocalRewrite

	basic := specs["LocalRewrite"]
	Expect(basic.methods).To(HaveLen(1))
	Expect(basic.methods).To(HaveKey("Test"))

	Expect(basic.methods["Test"].params).To(HaveLen(3))
	Expect(basic.methods["Test"].params[0].name).To(Equal("x"))
	Expect(basic.methods["Test"].params[0].typ.Text()).To(Equal(testPath + "/interfaces/localtypes.X"))
	Expect(basic.methods["Test"].params[1].name).To(Equal("y"))
	Expect(basic.methods["Test"].params[1].typ.Text()).To(Equal(testPath + "/interfaces/localtypes.Y"))
	Expect(basic.methods["Test"].params[2].name).To(Equal("z"))
	Expect(basic.methods["Test"].params[2].typ.Text()).To(Equal(testPath + "/interfaces/localtypes.Z"))

	Expect(basic.methods["Test"].results).To(HaveLen(3))
	Expect(basic.methods["Test"].results[0].Text()).To(Equal("*" + testPath + "/interfaces/localtypes.X"))
	Expect(basic.methods["Test"].results[1].Text()).To(Equal("*" + testPath + "/interfaces/localtypes.Y"))
	Expect(basic.methods["Test"].results[2].Text()).To(Equal("*" + testPath + "/interfaces/localtypes.Z"))

	//
	// EmbeddedStruct

	embeddedStruct := specs["EmbeddedStruct"]
	Expect(embeddedStruct.methods).To(HaveLen(1))
	Expect(embeddedStruct.methods).To(HaveKey("Foo"))
	Expect(embeddedStruct.methods["Foo"].params).To(HaveLen(1))
	Expect(embeddedStruct.methods["Foo"].params[0].typ.Text()).To(Equal("struct{z " + testPath + "/interfaces/localtypes.X}"))

	//
	// InterfaceStruct

	embeddedInterface := specs["InterfaceStruct"]
	Expect(embeddedInterface.methods).To(HaveLen(1))
	Expect(embeddedInterface.methods).To(HaveKey("Bar"))
	Expect(embeddedInterface.methods["Bar"].results).To(HaveLen(1))
	Expect(embeddedInterface.methods["Bar"].results[0].Text()).To(Equal("interface{Baz(y " + testPath + "/interfaces/localtypes.Y) struct{z " + testPath + "/interfaces/localtypes.Z}}"))
}

func (s *InterfacesSuite) TestImports(t sweet.T) {
	interfaces := []string{"Imports", "LocalImports", "RenamedImports"}

	specs := getSpecs("imports")
	Expect(specs).To(HaveLen(3))
	Expect(specs).To(HaveKey(interfaces[0]))
	Expect(specs).To(HaveKey(interfaces[1]))
	Expect(specs).To(HaveKey(interfaces[2]))

	for _, name := range interfaces {
		use := specs[name]
		Expect(use.methods).To(HaveLen(3))

		Expect(use.methods).To(HaveKey("GetX"))
		Expect(use.methods["GetX"].results).To(HaveLen(1))
		Expect(use.methods["GetX"].results[0].Text()).To(Equal(testPath + "/interfaces/localtypes.X"))

		Expect(use.methods).To(HaveKey("GetY"))
		Expect(use.methods["GetY"].results).To(HaveLen(1))
		Expect(use.methods["GetY"].results[0].Text()).To(Equal(testPath + "/interfaces/localtypes.Y"))

		Expect(use.methods).To(HaveKey("GetZ"))
		Expect(use.methods["GetZ"].results).To(HaveLen(1))
		Expect(use.methods["GetZ"].results[0].Text()).To(Equal(testPath + "/interfaces/localtypes.Z"))
	}
}

//
// Helpers

func getSpecs(name string) map[string]*wrappedInterface {
	packageName := fmt.Sprintf("%s/interfaces/%s", testPath, name)

	pkg, pkgType, err := parseDir(packageName)
	Expect(err).To(BeNil())

	v1 := newNameExtractor()
	walk(pkg, v1)

	v2 := newInterfaceExtractor(pkgType, packageName, v1.names)
	walk(pkg, v2)
	return v2.specs
}
