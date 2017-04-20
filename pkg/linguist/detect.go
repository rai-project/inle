package linguist

import (
	"errors"

	"github.com/rai-project/linguist"
)

func Detect(prog string) (string, error) {
	lang := linguist.Detect(prog)
	if lang == "" {
		return "Unknown", errors.New("unable to determine language")
	}
	return lang, nil
}
