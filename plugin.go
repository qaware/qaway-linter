package qawaylinter

import "github.com/golangci/plugin-module-register/register"

func init() {
	register.Plugin("qawaylinter", New)
}

func New(conf any) (register.LinterPlugin, error) {
	// The configuration type will be map[string]any or []interface, it depends on your configuration.
	// You can use https://github.com/go-viper/mapstructure to convert map to struct.
	settings, err := register.DecodeSettings[Settings](conf)
	if err != nil {
		return nil, err
	}

	return &AnalyzerPlugin{Settings: settings}, nil
}
