package repository_test

import (
	"context"

	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/go-backend/pkg/model"
)

var _ = Describe("Inventory repository", func() {
	createInventory := func() *model.CharacterInventory {
		inv := &model.CharacterInventory{}
		Expect(faker.FakeData(&inv)).To(Succeed())
		Expect(invRepo.UpdateInventory(nil, inv)).To(Succeed())
		Expect(inv).NotTo(BeNil())

		return inv
	}

	Describe("GetInventory", func() {
		When("given valid input", func() {
			It("should work", func() {
				inv := createInventory()
				out, err := invRepo.GetInventory(context.Background(), inv.CharacterId)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out).To(BeEquivalentTo(inv))

				out, err = invRepo.GetInventory(nil, inv.CharacterId)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out).To(BeEquivalentTo(inv))
			})
		})

		When("given invalid input", func() {
			It("", func() {
				out, err := invRepo.GetInventory(nil, 1e19)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("UpdateInventory", func() {
		When("given valid input", func() {
			It("should work with new values", func() {
				inv := createInventory()
				inv.CharacterId = inv.CharacterId + 1
				Expect(invRepo.UpdateInventory(context.Background(), inv)).To(Succeed())

				out, err := invRepo.GetInventory(nil, inv.CharacterId)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out).To(BeEquivalentTo(inv))
			})

			It("should replace existing", func() {
				inv := createInventory()

				newInv := &model.CharacterInventory{}
				Expect(faker.FakeData(newInv)).To(Succeed())
				newInv.CharacterId = inv.CharacterId

				Expect(invRepo.UpdateInventory(context.Background(), newInv)).To(Succeed())

				out, err := invRepo.GetInventory(nil, inv.CharacterId)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out).NotTo(BeEquivalentTo(inv))
				Expect(out).To(BeEquivalentTo(newInv))
			})
		})

		When("given invalid input", func() {
			It("should throw and error", func() {
				out, err := invRepo.GetInventory(nil, 1e19)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})
})
