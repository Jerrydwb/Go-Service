package shared

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// FlexibleString puede deserializar tanto números como strings.
// Útil para MariaDB que a veces devuelve strings y a veces números.
type FlexibleString string

func (fs *FlexibleString) UnmarshalJSON(data []byte) error {
	// Intentar primero como string
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*fs = FlexibleString(s)
		return nil
	}

	// Intentar como número
	var n json.Number
	if err := json.Unmarshal(data, &n); err == nil {
		*fs = FlexibleString(n.String())
		return nil
	}

	// Intentar como int
	var i int64
	if err := json.Unmarshal(data, &i); err == nil {
		*fs = FlexibleString(strconv.FormatInt(i, 10))
		return nil
	}

	// Intentar como float
	var f float64
	if err := json.Unmarshal(data, &f); err == nil {
		*fs = FlexibleString(fmt.Sprintf("%.0f", f))
		return nil
	}

	return fmt.Errorf("cannot unmarshal %s into FlexibleString", string(data))
}

func (fs FlexibleString) String() string {
	return string(fs)
}
