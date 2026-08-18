package dbentity

import (
	"reflect"
	"testing"
)

type testEntity struct {
	ID        int    `db:"id" select:"true" insert:"true" issearch:"true" pfield:"id"`
	Name      string `db:"name" select:"true" insert:"true" issearch:"true" pfield:"name"`
	Email     string `db:"email" select:"true" insert:"false"`
	Password  string `db:"password" select:"false" insert:"true"`
	CreatedAt string `db:"created_at" select:"true" insert:"false"`
	Internal  string `select:"true" insert:"false"`
}

type testEntityBase = BaseEntity[testEntity]

func TestBaseEntityFieldName(t *testing.T) {
	ent := testEntity{}
	base := testEntityBase{}

	tests := []struct {
		name     string
		fieldPtr any
		want     string
	}{
		{
			name:     "ID",
			fieldPtr: &ent.ID,
			want:     "id",
		},
		{
			name:     "Name",
			fieldPtr: &ent.Name,
			want:     "name",
		},
		{
			name:     "Email",
			fieldPtr: &ent.Email,
			want:     "email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := base.FieldName(&ent, tt.fieldPtr)

			if got != tt.want {
				t.Errorf("FieldName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBaseEntityFieldNameWithAlias(t *testing.T) {
	ent := testEntity{}
	base := testEntityBase{}

	tests := []struct {
		name  string
		field any
		alias string
		want  string
	}{
		{
			name:  "without alias",
			field: &ent.ID,
			alias: "",
			want:  "id",
		},
		{
			name:  "with alias",
			field: &ent.ID,
			alias: "e",
			want:  "e.id",
		},
		{
			name:  "with another alias",
			field: &ent.CreatedAt,
			alias: "entity",
			want:  "entity.created_at",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := base.FieldNameWithAlias(&ent, tt.field, tt.alias)

			if got != tt.want {
				t.Errorf("FieldNameWithAlias() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBaseEntityPaginateFieldName(t *testing.T) {
	ent := testEntity{}
	base := testEntityBase{}

	tests := []struct {
		name  string
		field any
		want  string
	}{
		{
			name:  "ID",
			field: &ent.ID,
			want:  "id",
		},
		{
			name:  "Name",
			field: &ent.Name,
			want:  "name",
		},
		{
			name:  "Email without pagination tag",
			field: &ent.Email,
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := base.PaginateFieldName(&ent, tt.field)

			if got != tt.want {
				t.Errorf("PaginateFieldName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBaseEntityFieldsAll(t *testing.T) {
	base := testEntityBase{}

	got := base.FieldsAll()

	want := []string{
		"id",
		"name",
		"email",
		"created_at",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("FieldsAll() = %v, want %v", got, want)
	}
}

func TestBaseEntitySearchableFields(t *testing.T) {
	base := testEntityBase{}

	got := base.SearchableFields()

	want := []string{
		"id",
		"name",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("SearchableFields() = %v, want %v", got, want)
	}
}

func TestBaseEntityFieldsAllWithAlias(t *testing.T) {
	base := testEntityBase{}

	got := base.FieldsAllWithAlias("e")

	want := []string{
		"e.id",
		"e.name",
		"e.email",
		"e.created_at",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("FieldsAllWithAlias() = %v, want %v", got, want)
	}
}

func TestBaseEntityFields(t *testing.T) {
	base := testEntityBase{}

	got := base.Fields()

	want := []string{
		"id",
		"name",
		"password",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Fields() = %v, want %v", got, want)
	}
}

func TestBaseEntityValues(t *testing.T) {
	base := testEntityBase{}

	ent := testEntity{
		ID:        42,
		Name:      "John",
		Email:     "john@example.com",
		Password:  "secret",
		CreatedAt: "2026-08-18",
		Internal:  "internal",
	}

	got := base.Values(&ent)

	want := []any{
		42,
		"John",
		"secret",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Values() = %v, want %v", got, want)
	}
}

func TestBaseEntityFieldPointers(t *testing.T) {
	base := testEntityBase{}

	ent := testEntity{}

	got, err := base.FieldPointers(&ent)
	if err != nil {
		t.Fatalf("FieldPointers() error = %v", err)
	}

	if len(got) != 6 {
		t.Fatalf("FieldPointers() returned %d pointers, want 6", len(got))
	}

	if reflect.ValueOf(got[0]).Pointer() != reflect.ValueOf(&ent.ID).Pointer() {
		t.Errorf("FieldPointers()[0] does not point to ID")
	}

	if reflect.ValueOf(got[1]).Pointer() != reflect.ValueOf(&ent.Name).Pointer() {
		t.Errorf("FieldPointers()[1] does not point to Name")
	}
}

func TestBaseEntityFieldPointer(t *testing.T) {
	base := testEntityBase{}

	ent := testEntity{}

	tests := []struct {
		name      string
		fieldName string
		want      any
	}{
		{
			name:      "ID",
			fieldName: "id",
			want:      &ent.ID,
		},
		{
			name:      "Name",
			fieldName: "name",
			want:      &ent.Name,
		},
		{
			name:      "CreatedAt",
			fieldName: "created_at",
			want:      &ent.CreatedAt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := base.FieldPointer(&ent, tt.fieldName)
			if err != nil {
				t.Fatalf("FieldPointer() error = %v", err)
			}

			if reflect.ValueOf(got).Pointer() != reflect.ValueOf(tt.want).Pointer() {
				t.Errorf("FieldPointer() points to %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseEntityFieldsMap(t *testing.T) {
	base := testEntityBase{}

	ent := testEntity{
		ID:        42,
		Name:      "John",
		Email:     "john@example.com",
		Password:  "secret",
		CreatedAt: "2026-08-18",
		Internal:  "internal",
	}

	got := base.FieldsMap(&ent)

	want := map[string]any{
		"id":       42,
		"name":     "John",
		"password": "secret",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("FieldsMap() = %v, want %v", got, want)
	}
}
