package formater

import (
	"fmt"

	gabs "github.com/Jeffail/gabs/v2"
	"github.com/gookit/color"
)

// Available formaters
type Formater int

const (
	DefaultFormater Formater = iota
	ColoredFormater
	RawFormater
)

// LogFieldsMapping is used to map a json log to Data, Log, Service and extra fields.
type LogFieldsMapping struct {
	DateField    string
	LogField     string
	ServiceField string
	ExtraField   map[string]string
}

type LogFormater interface {
	PrintLog(log *gabs.Container)
}

// Create a new formater
func NewFormater(f Formater, mapping LogFieldsMapping) (LogFormater, error) {
	switch f {
	case DefaultFormater:
		return LogFormater(&defaultFormater{
			mapping: mapping,
		}), nil
	case ColoredFormater:
		return LogFormater(&coloredFormater{
			mapping: mapping,
		}), nil
	case RawFormater:
		return LogFormater(&rawFormater{}), nil
	}
	return nil, fmt.Errorf("unknown formater %v", f)
}

type defaultFormater struct {
	mapping LogFieldsMapping
}

func (f defaultFormater) PrintLog(log *gabs.Container) {
	fmt.Printf("[%v : %v] %v",
		log.Path(f.mapping.DateField).Data(), log.Path(f.mapping.ServiceField).Data(),
		log.Path(f.mapping.LogField).Data())
}

type coloredFormater struct {
	mapping LogFieldsMapping
}

func (f coloredFormater) PrintLog(log *gabs.Container) {
	fmt.Printf("[%v : %v] %v",
		color.Yellow.Sprint(log.Path(f.mapping.DateField).Data()),
		color.Blue.Sprint(log.Path(f.mapping.ServiceField).Data()),

		// Try to print colors to log itself if contains some keywords
		log.Path(f.mapping.LogField).Data())
}

type rawFormater struct{}

func (rawFormater) PrintLog(log *gabs.Container) {
	fmt.Printf("%v", log.String())
}
