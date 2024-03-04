package model_test

import (
	"github.com/Nerzal/gocloak/v13"
	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/go-backend/pkg/model"
)

var _ = Describe("Claims model", func() {
	Describe("HasRole", func() {
		It("should work", func() {
			role := &gocloak.Role{
				Name: gocloak.StringP(faker.Username()),
			}
			role2 := &gocloak.Role{
				Name: gocloak.StringP(faker.Email()),
			}
			claims := &model.SROClaims{
				RealmRoles: model.ClaimRoles{
					Roles: []string{faker.Username(), *role.Name},
				},
			}

			Expect(claims.HasRole(role)).To(BeTrue())
			Expect(claims.HasRole(role2)).To(BeFalse())
		})
	})

	Describe("HasResourceRole", func() {
		It("should work", func() {
			role := &gocloak.Role{
				Name: gocloak.StringP(faker.Username()),
			}

			role2 := &gocloak.Role{
				Name: gocloak.StringP(faker.Email()),
			}

			clientId := faker.Username()
			claims := &model.SROClaims{
				ResourceAccess: model.ClaimResourceAccess{
					clientId: model.ClaimRoles{
						Roles: []string{faker.Username(), *role.Name},
					},
				},
			}

			Expect(claims.HasResourceRole(role, clientId)).To(BeTrue())
			Expect(claims.HasResourceRole(role2, clientId)).To(BeFalse())
			Expect(claims.HasResourceRole(role, clientId+"a")).To(BeFalse())
		})
	})
})
