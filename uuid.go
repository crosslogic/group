package group

import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

type UUID []uuid.UUID

func (u UUID) InClause() (out string, err error) {
	if len(u) == 0 {
		return out, errors.Errorf("empty slice")
	}

	arr := []string{}
	for _, v := range u {
		arr = append(arr, v.String())
	}
	out = fmt.Sprintf("(%v)", strings.Join(arr, ","))

	return

}
