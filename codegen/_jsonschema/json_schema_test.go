package genschema_test

import (
	"goa.design/goa.v2/design"
	. "goa.design/goa.v2/design/apidsl"
	"goa.design/goa.v2/eval"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TypeSchema", func() {
	var typ design.DataType

	var s *genschema.JSONSchema

	BeforeEach(func() {
		typ = nil
		s = nil
		eval.Reset()
		design.ProjectedMediaTypes = make(design.MediaTypeRoot)
	})

	JustBeforeEach(func() {
		s = genschema.TypeSchema(design.Design, typ)
	})

	Context("with a media type", func() {
		BeforeEach(func() {
			MediaType("application/foo.bar", func() {
				Attributes(func() {
					Attribute("bar")
				})
				View("default", func() {
					Attribute("bar")
				})
			})

			Ω(eval.Run()).ShouldNot(HaveOccurred())
			typ = design.Design.MediaTypes["application/foo.bar"]
		})

		It("returns a proper JSON schema type", func() {
			Ω(s).ShouldNot(BeNil())
			Ω(s.Ref).Should(Equal("#/definitions/FooBar"))
		})
	})

	Context("with a media type with self-referencing attributes", func() {
		BeforeEach(func() {
			MediaType("application/vnd.menu+json", func() {
				Attributes(func() {
					Attribute("name", design.String, "The name of an application")
					Attribute("children", CollectionOf("application/vnd.menu+json"), func() {
						View("nameonly")
					})

				})
				View("default", func() {
					Attribute("name")
					Attribute("children", func() {
						View("nameonly")
					})
				})
				View("nameonly", func() {
					Attribute("name")
				})
			})

			Ω(func() { eval.Run() }).ShouldNot(Panic())
			Ω(eval.Errors).ShouldNot(HaveOccurred())
			typ = design.Design.MediaTypes["application/vnd.menu"]
		})

		It("returns a proper JSON schema type", func() {
			Ω(s).ShouldNot(BeNil())
			Ω(s.Ref).Should(Equal("#/definitions/Menu"))
		})

	})
})