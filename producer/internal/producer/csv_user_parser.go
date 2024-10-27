package producer

import (
	"strconv"
	"time"
)

const (
	IgnoredValue = "-1"
)

type emailCryptor interface {
	Encrypt(val string) (string, error)
}

type csvUserParser struct {
	cryptor emailCryptor
}

func NewCsvUserParser(cryptor emailCryptor) *csvUserParser {
	return &csvUserParser{
		cryptor: cryptor,
	}
}

func (p *csvUserParser) Parse(r Record) (*User, error) {
	id, err := toInt64(r[0])
	if err != nil {
		return nil, err
	}

	email, err := p.cryptor.Encrypt(r[3])
	if err != nil {
		return nil, err
	}

	createdAt, err := toTimeUTC(r[4])
	if err != nil {
		return nil, err
	}

	deleteAt, err := toTimeUTC(r[5])
	if err != nil {
		return nil, err
	}

	mergedAt, err := toTimeUTC(r[6])
	if err != nil {
		return nil, err
	}

	parent, err := toInt64(r[7])
	if err != nil {
		return nil, err
	}

	return &User{
		ID:           *id,
		FirstName:    r[1],
		LastName:     r[2],
		Email:        email,
		CreatedAt:    *createdAt,
		DeletedAt:    deleteAt,
		MergedAt:     mergedAt,
		ParentUserID: parent,
	}, nil
}

type User struct {
	ID           int64      `json:"id"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Email        string     `json:"email"`
	CreatedAt    time.Time  `json:"created_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
	MergedAt     *time.Time `json:"merged_at"`
	ParentUserID *int64     `json:"parent_user_id"`
}

func toInt64(v string) (*int64, error) {
	if v == IgnoredValue {
		return nil, nil
	}

	id, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func toTimeUTC(v string) (*time.Time, error) {
	if v == IgnoredValue {
		return nil, nil
	}

	millis, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return nil, err
	}

	// Convert milliseconds to seconds and nanoseconds
	seconds := millis / 1000
	nanoseconds := (millis % 1000) * 1000000

	t := time.Unix(seconds, nanoseconds).UTC()

	return &t, nil
}
