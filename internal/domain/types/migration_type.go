package dtypes

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type MigrationType uint8

const (
	MigrationTypeNone MigrationType = iota
	MigrationTypeNoneRequired
	MigrationTypeNoneDeferred
)

var (
	migrationTypes = map[MigrationType]string{
		MigrationTypeNone:         "",
		MigrationTypeNoneRequired: "required",
		MigrationTypeNoneDeferred: "deferred",
	}

	migrationTypesByStringValue = map[string]MigrationType{
		"":         MigrationTypeNone,
		"required": MigrationTypeNoneRequired,
		"deferred": MigrationTypeNoneDeferred,
	}
)

// String returns the value as a string.
func (v MigrationType) String() string {
	return migrationTypes[v]
}

// Ptr returns the pointer.
func (v MigrationType) Ptr() *MigrationType {
	return &v
}

// Exists returns true if the MigrationType is valid.
func (v MigrationType) Exists() bool {
	_, exists := migrationTypes[v]
	return exists
}

// ParseMigrationType converting string to MigrationType
func ParseMigrationType(s string) (MigrationType, error) {
	migrationType, exists := migrationTypesByStringValue[strings.ToLower(s)]
	if !exists {
		return MigrationTypeNone, errors.New("invalid value for MigrationType")
	}
	return migrationType, nil
}

// ValidateMigrationType checks if the field value is a valid MigrationType enum.
func ValidateMigrationType(fl validator.FieldLevel) bool {
	_, err := ParseMigrationType(fl.Field().String())
	return err == nil
}
