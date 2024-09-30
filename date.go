package gohelpers

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/araddon/dateparse"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type APIDate struct {
	time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (ad *APIDate) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	t, err := dateparse.ParseAny(s)
	if err != nil {
		return err
	}
	ad.Time = t
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (ad APIDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(ad.Time.Format(time.RFC3339))
}

// Scan implements the sql.Scanner interface
func (ad *APIDate) Scan(value interface{}) error {
	if value == nil {
		ad.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		ad.Time = v
	case string:
		t, err := dateparse.ParseAny(v)
		if err != nil {
			return err
		}
		ad.Time = t
	case []byte:
		t, err := dateparse.ParseAny(string(v))
		if err != nil {
			return err
		}
		ad.Time = t
	default:
		return fmt.Errorf("cannot scan type %T into CustomDate", v)
	}

	return nil
}

// Value implements the driver.Valuer interface
func (ad APIDate) Value() (driver.Value, error) {
	if ad.Time.IsZero() {
		return nil, nil
	}
	return ad.Time, nil
}

// GormDataType implements the gorm.GormDataTypeInterface interface
func (APIDate) GormDataType() string {
	return "date"
}

// GormDBDataType implements the schema.GormDBDataTypeInterface interface
func (APIDate) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql", "sqlite":
		return "DATE"
	case "postgres":
		return "DATE"
	}
	return ""
}
