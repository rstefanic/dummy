package postgresql

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

func fakeData(datatype, udt string) (string, error) {
	switch datatype {
	case "ARRAY":
		underlyingDt, err := udtToPsqlDatatype(udt)
		if err != nil {
			return "", err
		}

		value, err := fakeData(underlyingDt, "")
		if err != nil {
			return "", err
		}

		var array strings.Builder
		array.WriteString("ARRAY[")
		array.WriteString(value)
		array.WriteString("]")
		return array.String(), nil
	case "bigint":
		bigIntVal := gofakeit.Int64()
		return strconv.FormatInt(bigIntVal, 10), nil
	case "bit":
		charset := "01"
		val := make([]byte, 8)
		for i := range val {
			val[i] = charset[gofakeit.IntRange(0, len(charset)-1)]
		}
		var bitString strings.Builder
		bitString.WriteString("B'")
		bitString.WriteString(string(val))
		bitString.WriteString("'")
		return bitString.String(), nil
	case "boolean":
		boolVal := gofakeit.Bool()
		if boolVal {
			return "true", nil
		} else {
			return "false", nil
		}
	case "numeric", "decimal":
		charset := "0123456789"
		maxPreDeicmalLen := 131072
		maxPostDecimalLen := 16383
		preDecimalLen := gofakeit.IntRange(1, maxPreDeicmalLen)
		postDecimalLen := gofakeit.IntRange(1, maxPostDecimalLen)

		val := make([]byte, preDecimalLen+1+postDecimalLen) // +1 for the decimal character
		i := 0
		j := 0

		for j < preDecimalLen {
			val[i] = charset[gofakeit.IntRange(0, len(charset)-1)]
			i += 1
			j += 1
		}

		val[i] = '.'
		i += 1

		j = 0
		for j < postDecimalLen {
			val[i] = charset[gofakeit.IntRange(0, len(charset)-1)]
			i += 1
			j += 1
		}

		return string(val), nil
	case "double precision":
		// PSQL double precision has 15 digits of precision
		val := strconv.FormatFloat(gofakeit.Float64(), 'f', 15, 64)
		return val, nil
	case "integer":
		intVal := gofakeit.Int16()
		return strconv.FormatInt(int64(intVal), 10), nil
	case "json", "jsonb":
		var jo gofakeit.JSONOptions

		// Use gofakeit to create random JSON fields
		err := gofakeit.Struct(&jo)
		if err != nil {
			return "", err
		}

		// Overwrite the fields to force this to be an object
		jo.Indent = false
		jo.RowCount = 1
		jo.Type = "object"

		jsonRaw, err := gofakeit.JSON(&jo)
		if err != nil {
			return "", err
		}

		var json strings.Builder
		json.WriteRune('\'')
		json.WriteString(strings.ReplaceAll(string(jsonRaw), "'", "''")) // escape single quotes
		json.WriteRune('\'')
		return json.String(), err
	case "real":
		// PSQL real has 6 digits of precision
		val := strconv.FormatFloat(gofakeit.Float64(), 'f', 6, 32)
		return val, nil
	case "serial":
		serialVal := gofakeit.IntRange(1, math.MaxInt32)
		return strconv.FormatInt(int64(serialVal), 10), nil
	case "smallserial":
		smallSerialVal := gofakeit.IntRange(1, math.MaxInt16)
		return strconv.FormatInt(int64(smallSerialVal), 10), nil
	case "text":
		var sentence strings.Builder
		sentence.WriteRune('\'')
		sentence.WriteString(strings.ReplaceAll(gofakeit.Sentence(1), "'", "''")) // escape single quotes
		sentence.WriteRune('\'')
		return sentence.String(), nil
	case "timestamp", "timestamp with time zone", "timestamp without time zone":
		var timestamp strings.Builder
		timestamp.WriteRune('\'')
		timestamp.WriteString(gofakeit.Date().Format(time.DateOnly))
		timestamp.WriteRune('\'')
		return timestamp.String(), nil
	case "uuid":
		var uuid strings.Builder
		uuid.WriteRune('\'')
		uuid.WriteString(gofakeit.UUID())
		uuid.WriteRune('\'')
		return uuid.String(), nil
	default:
		return "", errors.New("Datatype currently unsupported: " + datatype + "(" + udt + ")")
	}
}

func udtToPsqlDatatype(udt string) (string, error) {
	switch udt {
	case "_text", "text":
		return "text", nil
	default:
		return "", errors.New("Unknown UDT to datatype mapping: " + udt)
	}
}
