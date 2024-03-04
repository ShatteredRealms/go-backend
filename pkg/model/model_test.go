package model_test

import (
	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/go-backend/pkg/model"
)

var _ = Describe("Model model", func() {
	Describe("IsCreated", func() {
		It("should work", func() {
			model := &model.Model{}
			Expect(model.IsCreated()).To(BeFalse())

			Expect(faker.FakeData(model)).To(Succeed())
			Expect(model.IsCreated()).To(BeTrue())
		})
	})
})
