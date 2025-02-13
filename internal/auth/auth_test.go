package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestParseTimestamp(t *testing.T) {

	id, _ := uuid.NewRandom()

	tests := map[string]struct {
		input       uuid.UUID
		secret, err string
	}{
		`empty`: {
			input:  id,
			secret: "JWT Testing",
			err:    "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {

			token, err := MakeJWT(tt.input, tt.secret, time.Hour)
			if err != nil {
				t.Errorf("unable to make JWT: %s", err.Error())
			}

			decode, err := ValidateJWT(token, tt.secret)
			if err != nil {
				t.Errorf("unable to make JWT: %s", err.Error())
			}

			if tt.input != decode {
				t.Errorf("expected: %s got: %s", tt.input, decode)
			}

		})
	}

}
