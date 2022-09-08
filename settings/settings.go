package settings

import "github.com/spf13/viper"

func init() {
	err := SetupSetting()
	if err != nil {
		panic(err)
	}
}

// general operation for Viper
type Setting struct {
	vp *viper.Viper
}

func NewSetting(configs ...string) (*Setting, error) {
	vp := viper.New()
	// find config.yaml file
	vp.SetConfigName("config")
	for _, config := range configs {
		if config != "" {
			vp.AddConfigPath(config)
		}
	}
	vp.SetConfigType("yaml")
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}

	s := &Setting{vp}
	return s, nil
}

var sections = make(map[string]interface{})

func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}

	if _, ok := sections[k]; !ok {
		sections[k] = v
	}
	return nil
}

// add config struct
type CrawlerSettingS struct {
	PageStart        int
	PageEnd          int
	ProxyHost        string
	EnableProxy      bool
	TagConsumerCount int
	HonConsumerCount int
	HonBuffer        int
	TagBuffer        int
	EnableFilter     bool
	WantedTags       []string
}

var CrawlerSetting *CrawlerSettingS
var WantedTagsSet = make(map[string]bool)

func SetupSetting() error {
	s, err := NewSetting("./")
	if err != nil {
		return err
	}
	err = s.ReadSection("CrawlerSetting", &CrawlerSetting)
	if err != nil {
		return err
	}
	for _, v := range CrawlerSetting.WantedTags {
		WantedTagsSet[v] = true
	}
	return nil
}
