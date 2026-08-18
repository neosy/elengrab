package etags

import (
	"reflect"
	"testing"

	"github.com/Masterminds/squirrel"
)

type testEntity struct {
	ID        int    `db:"id" select:"true"`
	Name      string `db:"name" select:"true"`
	Email     string `db:"email" select:"false"`
	Password  string `db:"password"`
	CreatedAt string `db:"created_at" select:"true" sqlexpr:"NOW()"`
	Internal  string `select:"true"`
}

func TestFields(t *testing.T) {
	ent := testEntity{}

	got := Fields(ent, TagNameSelect)
	want := []string{"id", "name", "password", "created_at"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Fields() = %v, want %v", got, want)
	}
}

func TestFieldsWithTrueTag(t *testing.T) {
	tests := []struct {
		name string
		ent  any
		want []string
	}{
		{
			name: "struct",
			ent:  testEntity{},
			want: []string{"id", "name", "created_at"},
		},
		{
			name: "pointer to struct",
			ent:  &testEntity{},
			want: []string{"id", "name", "created_at"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FieldsWithTrueTag(tt.ent, TagNameSelect)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FieldsWithTrueTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldsExceptFalseTag(t *testing.T) {
	tests := []struct {
		name string
		ent  any
		want []string
	}{
		{
			name: "struct",
			ent:  testEntity{},
			want: []string{"id", "name", "password", "created_at"},
		},
		{
			name: "pointer to struct",
			ent:  &testEntity{},
			want: []string{"id", "name", "password", "created_at"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FieldsExceptFalseTag(tt.ent, TagNameSelect)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FieldsExceptFalseTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValues(t *testing.T) {
	ent := testEntity{
		ID:        42,
		Name:      "John",
		Email:     "john@example.com",
		Password:  "secret",
		CreatedAt: "2026-08-18",
		Internal:  "internal",
	}

	got := Values(ent, TagNameSelect)

	if len(got) != 4 {
		t.Fatalf("Values() returned %d values, want 4", len(got))
	}

	if got[0] != 42 {
		t.Errorf("values[0] = %v, want 42", got[0])
	}

	if got[1] != "John" {
		t.Errorf("values[1] = %v, want John", got[1])
	}

	if got[2] != "secret" {
		t.Errorf("values[2] = %v, want secret", got[2])
	}

	assertSQLExpression(t, got[3], "NOW()")
}

func TestValuesWithTrueTag(t *testing.T) {
	ent := testEntity{
		ID:        42,
		Name:      "John",
		Email:     "john@example.com",
		Password:  "secret",
		CreatedAt: "2026-08-18",
		Internal:  "internal",
	}

	got := ValuesWithTrueTag(ent, TagNameSelect)

	if len(got) != 3 {
		t.Fatalf("ValuesWithTrueTag() returned %d values, want 3", len(got))
	}

	if got[0] != 42 {
		t.Errorf("values[0] = %v, want 42", got[0])
	}

	if got[1] != "John" {
		t.Errorf("values[1] = %v, want John", got[1])
	}

	assertSQLExpression(t, got[2], "NOW()")
}

func TestValuesExceptFalseTag(t *testing.T) {
	ent := testEntity{
		ID:        42,
		Name:      "John",
		Email:     "john@example.com",
		Password:  "secret",
		CreatedAt: "2026-08-18",
		Internal:  "internal",
	}

	got := ValuesExceptFalseTag(ent, TagNameSelect)

	if len(got) != 4 {
		t.Fatalf("ValuesExceptFalseTag() returned %d values, want 4", len(got))
	}

	if got[0] != 42 {
		t.Errorf("values[0] = %v, want 42", got[0])
	}

	if got[1] != "John" {
		t.Errorf("values[1] = %v, want John", got[1])
	}

	if got[2] != "secret" {
		t.Errorf("values[2] = %v, want secret", got[2])
	}

	assertSQLExpression(t, got[3], "NOW()")
}

func TestGetFieldTagByOffset(t *testing.T) {
	type offsetEntity struct {
		ID   int    `db:"id"`
		Name string `db:"name"`
	}

	ent := &offsetEntity{
		ID:   42,
		Name: "John",
	}

	other := 100

	tests := []struct {
		name     string
		fieldPtr any
		tagKey   string
		want     string
		wantOK   bool
	}{
		{
			name:     "ID field",
			fieldPtr: &ent.ID,
			tagKey:   "db",
			want:     "id",
			wantOK:   true,
		},
		{
			name:     "Name field",
			fieldPtr: &ent.Name,
			tagKey:   "db",
			want:     "name",
			wantOK:   true,
		},
		{
			name:     "missing tag",
			fieldPtr: &ent.ID,
			tagKey:   "missing",
			want:     "",
			wantOK:   true,
		},
		{
			name:     "field not found",
			fieldPtr: &other,
			tagKey:   "db",
			want:     "",
			wantOK:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := GetFieldTagByOffset(ent, tt.fieldPtr, tt.tagKey)

			if got != tt.want || ok != tt.wantOK {
				t.Errorf(
					"GetFieldTagByOffset() = (%q, %v), want (%q, %v)",
					got,
					ok,
					tt.want,
					tt.wantOK,
				)
			}
		})
	}
}

func assertSQLExpression(t *testing.T, value any, wantSQL string) {
	t.Helper()

	expr, ok := value.(squirrel.Sqlizer)
	if !ok {
		t.Fatalf("value has type %T, want squirrel.Sqlizer", value)
	}

	sql, args, err := expr.ToSql()
	if err != nil {
		t.Fatalf("ToSql() error = %v", err)
	}

	if sql != wantSQL {
		t.Errorf("expression SQL = %q, want %q", sql, wantSQL)
	}

	if len(args) != 0 {
		t.Errorf("expression args = %v, want none", args)
	}
}
