package parser

//
//import "time"
//
//type (
//	Config struct {
//		LogType        string       `yaml:"log_type"`
//		TimeMultiplier float64      `yaml:"time_multiplier"`
//		CSV            CSVConfig    `yaml:"csv_config"`
//		LTSV           LTSVConfig   `yaml:"ltsv_config"`
//		RegExp         RegExpConfig `yaml:"regexp_config"`
//	}
//
//	CSVConfig struct {
//		Delimiter rune   `yaml:"delimiter"`
//		Format    string `yaml:"format"`
//	}
//
//	LTSVConfig struct {
//		FieldDelimiter byte              `yaml:"field_delimiter"`
//		ValueDelimiter byte              `yaml:"value_delimiter"`
//		Mapping        map[string]string `yaml:"mapping"`
//	}
//
//	RegExpConfig struct {
//		Pattern string `yaml:"pattern"`
//	}
//)
//
//const (
//	TypeAuto   = "auto"
//	TypeCSV    = "csv"
//	TypeLTSV   = "ltsv"
//	TypeRegExp = "regexp"
//)
//
//var (
//	DefaultConfig = Config{
//		LogType:        TypeAuto,
//		TimeMultiplier: time.Second.Seconds(),
//		CSV: CSVConfig{
//			Delimiter: ' ',
//		},
//		LTSV: LTSVConfig{
//			FieldDelimiter: '\t',
//			ValueDelimiter: ':',
//		},
//		RegExp: RegExpConfig{},
//	}
//)
