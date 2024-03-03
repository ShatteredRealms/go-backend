package helpers_test

import (
	"context"
	"fmt"

	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"

	"github.com/ShatteredRealms/go-backend/pkg/helpers"
)

var _ = Describe("Error helpers", func() {
	var (
		ctx       context.Context
		err       error
		errString string
		wasFatal  bool
	)

	log.StandardLogger().ExitFunc = func(int) { wasFatal = true }

	BeforeEach(func() {
		ctx = context.Background()
		err = fmt.Errorf("test error")
		errString = faker.Name()
		wasFatal = false
	})

	Describe("check", func() {
		It("should do nothing if err is nil", func() {
			helpers.Check(ctx, nil, errString)
			Expect(wasFatal).To(BeFalse())
			helpers.Check(nil, nil, errString)
			Expect(wasFatal).To(BeFalse())
			helpers.Check(ctx, nil, "")
			Expect(wasFatal).To(BeFalse())
			helpers.Check(nil, nil, "")
			Expect(wasFatal).To(BeFalse())
		})
		It("should fatal if err", func() {
			helpers.Check(ctx, err, errString)
			Expect(wasFatal).To(BeTrue())

			wasFatal = false
			helpers.Check(nil, err, errString)
			Expect(wasFatal).To(BeTrue())

			wasFatal = false
			helpers.Check(ctx, err, "")
			Expect(wasFatal).To(BeTrue())

			wasFatal = false
			helpers.Check(nil, err, "")
			Expect(wasFatal).To(BeTrue())
		})
	})
})
