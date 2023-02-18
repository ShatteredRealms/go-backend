package model_test

import (
	"fmt"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/kend/pkg/model"
	"github.com/kend/test/factory"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/bcrypt"
)

var (
	f = factory.NewFactory()
)

var _ = Describe("User", func() {
	var user *model.User

	BeforeEach(func() {
		user = f.NewBaseUser()
	})

	Context("Login", func() {
		var expectedError error
		var password string
		var passwordBytes []byte

		BeforeEach(func() {
			expectedError = nil
			password = user.Password
			passwordBytes, _ = bcrypt.GenerateFromPassword([]byte(password), 0)
			user.Password = string(passwordBytes)
		})

		It("should work if the password is correct", func() {
		})

		It("should fail if the db password isn't encrypted", func() {
			user.Password = password
			expectedError = fmt.Errorf("crypto/bcrypt: hashedSecret too short to be a bcrypted password")
		})

		It("should fail if the password does not match", func() {
			password = password + "a"
			expectedError = fmt.Errorf("invalid password")
		})

		AfterEach(func() {
			if expectedError == nil {
				Expect(user.Login(password)).To(BeNil())
			} else {
				Expect(user.Login(password)).To(Equal(expectedError))
			}
		})
	})

	Context("Validation", func() {
		var expectedError error

		BeforeEach(func() {
			expectedError = nil
		})

		It("should require an email", func() {
			user.Email = ""
			expectedError = fmt.Errorf("cannot create a user without an email")
		})

		It("should require a valid email", func() {
			user.Email = helpers.RandString(10)
			expectedError = fmt.Errorf("email is not valid")
		})

		It("should require a first name", func() {
			user.FirstName = ""
			expectedError = fmt.Errorf("first name cannot be empty")
		})

		It(fmt.Sprintf("should require a first name with max length %d", model.MaxFirstNameLength), func() {
			user.FirstName = helpers.RandString(model.MaxFirstNameLength + 1)
			expectedError = fmt.Errorf("first name cannot be longer than 50 characters")
		})

		It("should require a last name", func() {
			user.LastName = ""
			expectedError = fmt.Errorf("last name cannot be empty")
		})

		It(fmt.Sprintf("should require a last name with max length %d", model.MaxLastNameLength), func() {
			user.LastName = helpers.RandString(model.MaxLastNameLength + 1)
			expectedError = fmt.Errorf("last name cannot be longer than 50 characters")
		})

		It("should require a username", func() {
			user.Username = ""
			expectedError = fmt.Errorf("cannot create a user without a username")
		})

		It(fmt.Sprintf("should require a username with min length %d", model.MinUsernameLength), func() {
			user.Username = helpers.RandString(model.MinUsernameLength - 1)
			expectedError = fmt.Errorf("username less than minimum length of 3")
		})

		It(fmt.Sprintf("should require a username with max length %d", model.MaxUsernameLength), func() {
			user.Username = helpers.RandString(model.MaxUsernameLength + 1)
			expectedError = fmt.Errorf("username exeeded maximum length of 25")
		})

		It("should require a password", func() {
			user.Password = ""
			expectedError = fmt.Errorf("cannot create a user without a password")
		})

		It(fmt.Sprintf("should require a password with minimum length of %d", model.MinPasswordLength), func() {
			user.Password = helpers.RandString(model.MinPasswordLength - 1)
			expectedError = fmt.Errorf("password less than minimum length of %d", model.MinPasswordLength)
		})

		It(fmt.Sprintf("should require a password with maximum length of %d", model.MaxPasswordLength), func() {
			user.Password = helpers.RandString(model.MaxPasswordLength + 1)
			expectedError = fmt.Errorf("password exeeded maximum length of %d", model.MaxPasswordLength)
		})

		It(fmt.Sprintf("should allow a password of length of %d", model.MaxPasswordLength), func() {
			user.Password = helpers.RandString(model.MaxPasswordLength)
		})

		It(fmt.Sprintf("should allow a password of length of %d", model.MinPasswordLength), func() {
			user.Password = helpers.RandString(model.MinPasswordLength)
		})

		AfterEach(func() {
			if expectedError == nil {
				Expect(user.Validate()).To(BeNil())
			} else {
				Expect(user.Validate()).To(Equal(expectedError))
			}
		})
	})

	Context("UpdateInfo", func() {
		var existingUser *model.User

		BeforeEach(func() {
			existingUser = f.NewBaseUser()
		})

		It("should require a valid first name", func() {
			user.FirstName = helpers.RandString(model.MaxFirstNameLength + 1)
			Expect(existingUser.UpdateInfo(*user)).NotTo(BeNil())
		})

		It("should require a valid last name", func() {
			user.LastName = helpers.RandString(model.MaxLastNameLength + 1)
			Expect(existingUser.UpdateInfo(*user)).NotTo(BeNil())
		})

		It("should require a valid username", func() {
			user.Username = helpers.RandString(model.MaxUsernameLength + 1)
			Expect(existingUser.UpdateInfo(*user)).NotTo(BeNil())
		})

		It("should require a valid email", func() {
			user.Email = "a"
			Expect(existingUser.UpdateInfo(*user)).NotTo(BeNil())
		})

		It("should update first name only if only that is given", func() {
			user = &model.User{
				FirstName: helpers.RandString(10),
			}
			Expect(existingUser.UpdateInfo(*user)).To(BeNil())
			Expect(existingUser.FirstName).To(Equal(user.FirstName))
		})

		It("should update last name only if only that is given", func() {
			user = &model.User{
				LastName: helpers.RandString(10),
			}
			Expect(existingUser.UpdateInfo(*user)).To(BeNil())
			Expect(existingUser.LastName).To(Equal(user.LastName))
		})

		It("should update username only if only that is given", func() {
			user = &model.User{
				Username: helpers.RandString(10),
			}
			Expect(existingUser.UpdateInfo(*user)).To(BeNil())
			Expect(existingUser.Username).To(Equal(user.Username))
		})

		It("should update email only if only that is given", func() {
			user = &model.User{
				Email: helpers.RandString(10) + "@example.com",
			}
			Expect(existingUser.UpdateInfo(*user)).To(BeNil())
			Expect(existingUser.Email).To(Equal(user.Email))
		})

		It("should update first name, last name, username and email if given valid info", func() {
			Expect(existingUser.UpdateInfo(*user)).To(BeNil())
			Expect(existingUser.FirstName).To(Equal(user.FirstName))
			Expect(existingUser.LastName).To(Equal(user.LastName))
			Expect(existingUser.Username).To(Equal(user.Username))
			Expect(existingUser.Email).To(Equal(user.Email))
		})
	})

	It("should know if it exists", func() {
		Expect(user.Exists()).To(BeFalse())
		user.ID = 1
		Expect(user.Exists()).To(BeTrue())
	})
})
