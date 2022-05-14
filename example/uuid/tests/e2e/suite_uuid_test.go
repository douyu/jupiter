package e2e

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	mocks "uuid/gen/mocks/grpc"
)

var _ = Describe("exampleService", func() {

	mockUuidGrpc := &mocks.UuidInterface{}

	uuidService := CreateUuidService()
	uuidService.UuidGrpc = mockUuidGrpc

	Context("List", func() {
		It("normal case", func() {
			Expect(uuidService).ShouldNot(BeNil())
		})
	})
})
