package group

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgtype"
	"github.com/pkg/errors"
)

// GrupoInts es el tipo que uso para guardar slices de ints en las structs.
// En la base de datos se guardan como []int.
// Cuando se codifican como JSON, se envían como slices de strings, para
// evitar el redondeo que hace javascripts los int64. Obviamente cuando
// lee los JSON los vuelve a convertir a números.
type Int []int

// InClause returns the string needed to be used in where InClause
// fmt.Sprintf("id IN %v", InClause(ii)) => "id IN (34,543,23,13)"
func (g Int) InClause() (out string, err error) {
	if len(g) == 0 {
		return out, errors.Errorf("empty slice")
	}
	str := []string{}
	for _, v := range g {
		str = append(str, fmt.Sprint(v))
	}
	out = "(" + strings.Join(str, ",") + ")"
	return
}

// EncodeBinary implement pgx binary protocol interface.
func (dst *Int) DecodeBinary(ci *pgtype.ConnInfo, src []byte) error {
	deInts := pgtype.Int8Array{}
	err := deInts.DecodeBinary(ci, src)
	if err != nil {
		return errors.Wrap(err, "decoding GrupoInts")
	}

	*dst = Int{}
	for _, v := range deInts.Elements {
		*dst = append(*dst, int(v.Int))
	}
	return nil
}

// EncodeBinary implement pgx binary protocol interface.
func (src Int) EncodeBinary(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {

	d := pgtype.Int8Array{}

	d.Dimensions = []pgtype.ArrayDimension{
		{
			Length:     int32(len(src)),
			LowerBound: -2147483648,
		},
	}
	d.Status = pgtype.Present
	for _, v := range src {
		d.Elements = append(d.Elements, pgtype.Int8{Int: int64(v), Status: pgtype.Present})
	}
	return d.EncodeBinary(ci, buf)
}

// Value implements SQL interface.
func (g Int) Value() (driver.Value, error) {

	str := []string{}
	for _, v := range g {
		str = append(str, strconv.Itoa(v))
	}

	out := []byte("{" + strings.Join(str, ",") + "}")

	return out, nil
}

// Scan implements SQL interface.
func (g *Int) Scan(src interface{}) error {
	// p := pq.Int64Array(g)

	// Chequeo que el campo sea []byte
	source, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion error .([]byte)")
	}
	if source == nil {
		return nil
	}

	str := strings.Replace(string(source), " ", "", -1)
	if str == "{}" {
		return nil
	}

	*g = Int{}
	vals := strings.Split(str[1:len(str)-1], ",")
	for _, v := range vals {
		if v == "" {
			continue
		}
		integer, err := strconv.Atoi(v)
		if err != nil {
			return errors.Wrapf(err, "converting '%v' into int", v)
		}
		*g = append(*g, integer)
	}

	return nil
}

func (g *Int) UnmarshalJSON(b []byte) error {
	ss := []string{}
	if err := json.Unmarshal(b, &ss); err != nil {
		// Si no puedo como strings, pruebo si está viniendo como []int
		ii := []int{}
		err = json.Unmarshal(b, &ii)
		if err != nil {
			return errors.Wrap(err, "couldn't be read as []int nor []string")
		}
		*g = ii
		return nil
	}

	arr := []int{}
	*g = arr
	for _, v := range ss {
		enInt, err := strconv.Atoi(v)
		if err != nil {
			return errors.Wrapf(err, "convirtiendo '%v' into int", v)
		}
		*g = append(*g, int(enInt))
	}

	return nil
}

func (g Int) MarshalJSON() ([]byte, error) {

	out := []string{}
	for _, v := range g {
		out = append(out, fmt.Sprint(v))
	}

	return json.Marshal(out)
}
