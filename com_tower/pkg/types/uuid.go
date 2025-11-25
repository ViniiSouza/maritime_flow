package types

import (
	"fmt"

	"github.com/google/uuid"
)

type UUID uuid.UUID

func (u *UUID) UnmarshalJSON(b []byte) error {
	id, err := uuid.Parse(string(b[:]))
	if err != nil {
		return err
	}
	*u = UUID(id)
	return nil
}

func (u UUID) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, "\"%s\"", uuid.UUID(u).String()), nil
}

func (u *UUID) String() string {
	return uuid.UUID(*u).String()
}
