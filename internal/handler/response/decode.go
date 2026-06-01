package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/eulerbutcooler/virtus/internal/domain"
	"github.com/eulerbutcooler/virtus/pkg/validator"
)

// Reads the JSON body into dst, then runs struct validation.
func DecodeAndValidate(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		if errors.Is(err, io.EOF) {
			return fmt.Errorf("%w: empty request body", domain.ErrInvalidInput)
		}
		return fmt.Errorf("%w: malformed JSON: %v", domain.ErrInvalidInput, err)
	}

	if err := validator.Validate(dst); err != nil {
		var vErrs validator.ValidationErrors
		if errors.As(err, &vErrs) {
			return vErrs
		}
		return err
	}
	return nil
}
