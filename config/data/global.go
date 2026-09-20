package data

type GlobalConfig struct {
	// runtime
	Port     uint16 `default:"8080"`
	Timezone string `default:"UTC"`
	Debug    bool   `default:"false"`
}
