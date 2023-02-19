package model

import (
	"fmt"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"
	"net/mail"
	"time"

	"gopkg.in/nullbio/null.v4"
)

const (
	MinPasswordLength  = 6
	MaxPasswordLength  = 64
	MaxFirstNameLength = 50
	MaxLastNameLength  = MaxFirstNameLength
	MinUsernameLength  = 3
	MaxUsernameLength  = 25
)

// User Database model for a User
type User struct {
	Username  string    `gorm:"primarykey" json:"username"`
	FirstName string    `gorm:"not null" json:"first_name"`
	LastName  string    `gorm:"not null" json:"last_name"`
	Email     string    `gorm:"not null" json:"email"`
	Password  string    `gorm:"not null" json:"password"`
	Roles     Roles     `gorm:"many2many:user_roles" json:"roles"`
	BannedAt  null.Time `json:"banned_at"`

	// CurrentCharacterId The character that is currently being played. If nil then the account is not playing.
	// Otherwise, the account is connected to a server.
	// @TODO: Refactor to another model and service
	CurrentCharacter null.String `gorm:"unique" json:"currentCharacter"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Validate Checks if all user data fields are valid.
func (u *User) Validate() error {
	if u.Email == "" {
		return fmt.Errorf("cannot create a user without an email")
	}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return fmt.Errorf("email is not valid")
	}

	if err := u.validateFirstName(); err != nil {
		return err
	}

	if err := u.validateLastName(); err != nil {
		return err
	}

	if err := u.validatePassword(); err != nil {
		return err
	}

	if err := u.validateUsername(); err != nil {
		return err
	}

	return nil
}

// Login Checks if the given password belongs to the user
func (u *User) Login(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		if err.Error() == "crypto/bcrypt: hashedPassword is not the hash of the given password" {
			err = fmt.Errorf("invalid password")
		}
		return err
	}

	return nil
}

// UpdateInfo Updates the info of the user if the fields are present and valid. If any field is present but not valid
// then an error is returned. If there are no errors, then the non-nil fields for the FirstName, LastName, Email, and
// Username will be updated.
func (u *User) UpdateInfo(userDetails *pb.EditUserDetailsRequest) error {

	if userDetails.FirstName != nil {
		if err := u.updateFirstName(userDetails.FirstName.Value); err != nil {
			return err
		}
	}

	if userDetails.LastName != nil {
		if err := u.updateLastsName(userDetails.LastName.Value); err != nil {
			return err
		}
	}

	if userDetails.NewUsername != nil {
		if err := u.updateUsername(userDetails.NewUsername.Value); err != nil {
			return err
		}
	}

	if userDetails.Email != nil {
		if err := u.updateEmail(userDetails.Email.Value); err != nil {
			return err
		}
	}

	return nil
}

func (u *User) validateFirstName() error {
	if u.FirstName == "" {
		return fmt.Errorf("first name cannot be empty")
	}

	if len(u.FirstName) > MaxFirstNameLength {
		return fmt.Errorf("first name cannot be longer than 50 characters")
	}

	return nil
}

func (u *User) updateFirstName(val string) error {

	if err := u.validateFirstName(); err != nil {
		return err
	}

	u.FirstName = val
	return nil
}

func (u *User) validateLastName() error {
	if u.LastName == "" {
		return fmt.Errorf("last name cannot be empty")
	}

	if len(u.LastName) > MaxLastNameLength {
		return fmt.Errorf("last name cannot be longer than 50 characters")
	}

	return nil
}

func (u *User) updateLastsName(val string) error {
	if err := u.validateLastName(); err != nil {
		return err
	}

	u.LastName = val
	return nil
}
func (u *User) validatePassword() error {
	if u.Password == "" {
		return fmt.Errorf("cannot create a user without a password")
	}

	if len(u.Password) < MinPasswordLength {
		return fmt.Errorf("password less than minimum length of %d", MinPasswordLength)
	}

	if len(u.Password) > MaxPasswordLength {
		return fmt.Errorf("password exeeded maximum length of %d", MaxPasswordLength)
	}

	return nil
}

func (u *User) UpdatePassword(val string) error {
	if err := u.validatePassword(); err != nil {
		return err
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(u.Password), 0)
	if err != nil {
		return fmt.Errorf("password: %w", err)
	}

	u.Password = string(hashedPass)
	return nil
}
func (u *User) validateUsername() error {
	if u.Username == "" {
		return fmt.Errorf("cannot create a user without a username")
	}

	if len(u.Username) < MinUsernameLength {
		return fmt.Errorf("username less than minimum length of %d", MinUsernameLength)
	}

	if len(u.Username) > MaxUsernameLength {
		return fmt.Errorf("username exeeded maximum length of %d", MaxUsernameLength)
	}

	return nil
}

func (u *User) updateUsername(val string) error {
	if err := u.validateUsername(); err != nil {
		return err
	}

	u.Username = val
	return nil
}

func (u *User) validateEmail() error {
	_, err := mail.ParseAddress(u.Email)
	return err
}

func (u *User) updateEmail(val string) error {
	if err := u.validateEmail(); err != nil {
		return err
	}

	u.Email = val
	return nil
}

func (u *User) ToPb() *pb.UserMessage {
	return &pb.UserMessage{
		Username:         u.Username,
		Email:            u.Email,
		Roles:            u.Roles.ToPB().Roles,
		CreatedAt:        u.CreatedAt.Format("2006-01-02T15:04:05-0700"),
		BannedAt:         u.BannedAtWrapper(),
		CurrentCharacter: u.CurrentCharacterWrapper(),
	}
}

func (u *User) BannedAtWrapper() *wrapperspb.StringValue {
	var bannedAt *wrapperspb.StringValue
	if u.BannedAt.Valid {
		bannedAt = wrapperspb.String(u.BannedAt.Time.String())
	}

	return bannedAt
}

func (u *User) CurrentCharacterWrapper() *wrapperspb.StringValue {
	var currentCharacter *wrapperspb.StringValue
	if u.CurrentCharacter.Valid {
		currentCharacter = wrapperspb.String(u.CurrentCharacter.String)
	}

	return currentCharacter
}

func (u *User) ToVerbosePb(permissions UserPermissions) *pb.GetUserResponse {
	return &pb.GetUserResponse{
		Username:         u.Username,
		Email:            u.Email,
		FirstName:        u.FirstName,
		LastName:         u.LastName,
		Roles:            u.Roles.ToPB().Roles,
		Permissions:      permissions.ToPB().Permissions,
		CreatedAt:        u.CreatedAt.Format("2006-01-02T15:04:05-0700"),
		BannedAt:         u.BannedAtWrapper(),
		CurrentCharacter: u.CurrentCharacterWrapper(),
	}
}
