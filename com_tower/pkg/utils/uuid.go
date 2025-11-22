package utils

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type UUID uuid.UUID

func (u *UUID) UnmarshalJSON(b []byte) error {
    s := strings.Trim(string(b), `"`)
    parsedUUID, err := uuid.Parse(s)
    if err != nil {
        return err
    }
    *u = UUID(parsedUUID)
    return nil
}

func (u *UUID) MarshalJSON() ([]byte, error) {
    return fmt.Appendf(nil, "\"%s\"", uuid.UUID(*u).String()), nil
}

func (u *UUID) String() string {
	return uuid.UUID(*u).String()
}
