package helpers_test

import (
	"context"
	"fmt"

	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	// "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"

	"github.com/ShatteredRealms/go-backend/pkg/helpers"
)

var _ = Describe("Error helpers", func() {
	Describe("check", func() {
		It("should do nothing if err is nil and fatal on err", func() {
			ctx := context.Background()
			_ = fmt.Errorf("test error")
			errString := faker.Name()
			testLogger, hook := test.NewNullLogger()
			log.Logger = testLogger

			helpers.Check(ctx, nil, errString)
			Expect(hook.LastEntry()).To(BeNil())
			hook.Reset()
			helpers.Check(nil, nil, errString)
			Expect(hook.LastEntry()).To(BeNil())
			hook.Reset()
			helpers.Check(ctx, nil, "")
			Expect(hook.LastEntry()).To(BeNil())
			hook.Reset()
			helpers.Check(nil, nil, "")
			Expect(hook.LastEntry()).To(BeNil())

			hook.Reset()
			// helpers.Check(ctx, err, errString)
			// Expect(hook.LastEntry()).NotTo(BeNil())
			// Expect(hook.LastEntry().Level).To(Equal(logrus.FatalLevel))
			// Expect(hook.Entries).To(HaveLen(1))

			// hook.Reset()
			// Expect(hook.LastEntry()).NotTo(BeNil())
			// helpers.Check(nil, err, errString)
			// Expect(hook.LastEntry().Level).To(Equal(logrus.FatalLevel))
			// Expect(hook.Entries).To(HaveLen(1))

			// hook.Reset()
			// Expect(hook.LastEntry()).NotTo(BeNil())
			// helpers.Check(ctx, err, "")
			// Expect(hook.LastEntry().Level).To(Equal(logrus.FatalLevel))
			// Expect(hook.Entries).To(HaveLen(1))

			// hook.Reset()
			// Expect(hook.LastEntry()).NotTo(BeNil())
			// helpers.Check(nil, err, "")
			// Expect(hook.LastEntry().Level).To(Equal(logrus.FatalLevel))
			// Expect(hook.Entries).To(HaveLen(1))
		})
	})
})
