package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

var GlobalConfig Config

const (
	BF_EXPECTED_EL         = 1000  // broj ocekivanih elemenata u bloom filteru
	BF_FALSE_POSITIVE_RATE = 0.001 // bloom filter false positive
	CMS_EPSILON            = 0.001
	CMS_DELTA              = 0.001
	STRUCTURE_TYPE         = "skiplist"
	SKIP_LIST_HEIGHT       = 10
	CRC_SIZE               = 4
	TIMESTAMP_SIZE         = 16
	TOMBSTONE_SIZE         = 1
	KEY_SIZE_SIZE          = 8
	VALUE_SIZE_SIZE        = 8
	CRC_START              = 0
	TIMESTAMP_START        = CRC_START + CRC_SIZE
	TOMBSTONE_START        = TIMESTAMP_START + TIMESTAMP_SIZE
	KEY_SIZE_START         = TOMBSTONE_START + TOMBSTONE_SIZE
	VALUE_SIZE_START       = KEY_SIZE_START + KEY_SIZE_SIZE
	KEY_START              = VALUE_SIZE_START + VALUE_SIZE_SIZE
)

type Config struct {
	BFExpectedElements  int     `yaml:"BFExpectedElements"`
	BFFalsePositiveRate float64 `yaml:"bloomFalsePositive"`
	CmsEpsilon          float64 `yaml:"cmsEpsilon"`
	CmsDelta            float64 `yaml:"cmsDelta"`
	StructureType       string  `yaml:"StructureType"`
	SkipListHeight      int     `yaml:"skipListHeight"`
	CrcSize             int     `yaml:"crcSize"`
	TimestampSize       int     `yaml:"timestampSize"`
	TombstoneSize       int     `yaml:"tombstoneSize"`
	KeySizeSize         int     `yaml:"keySizeSize"`
	ValueSizeSize       int     `yaml:"valueSizeSize"`
	CrcStart            int     `yaml:"crcStart"`
	TimestampStart      int     `yaml:"timestampStart"`
	TombstoneStart      int     `yaml:"tombstoneStart"`
	KeySizeStart        int     `yaml:"keySizeStart"`
	ValueSizeStart      int     `yaml:"ValueSizeStart"`
	KeyStart            int     `yaml:"keyStart"`
}

func NewConfig(filename string) *Config {
	var config Config
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		config.BFExpectedElements = BF_EXPECTED_EL
		config.BFFalsePositiveRate = BF_FALSE_POSITIVE_RATE
		config.CmsDelta = CMS_DELTA
		config.CmsEpsilon = CMS_EPSILON
		config.StructureType = STRUCTURE_TYPE
		config.SkipListHeight = SKIP_LIST_HEIGHT
		config.CrcSize = CRC_SIZE
		config.TimestampSize = TIMESTAMP_SIZE
		config.TombstoneSize = TOMBSTONE_SIZE
		config.KeySizeSize = KEY_SIZE_SIZE
		config.ValueSizeSize = VALUE_SIZE_SIZE
		config.CrcStart = CRC_START
		config.TimestampStart = TIMESTAMP_START
		config.TombstoneStart = TOMBSTONE_START
		config.KeySizeStart = KEY_SIZE_START
		config.ValueSizeStart = VALUE_SIZE_START
		config.KeyStart = KEY_START
	} else {
		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			fmt.Printf("Unmarshal: %v", err)
		}
	}

	return &config
}

func Init() {
	PATH := "config\\config.yaml"

	GlobalConfig = *NewConfig(PATH)

	if _, err := os.Stat(PATH); errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(PATH)
		defer f.Close()
		if err != nil {
			panic(err)
		}

		out, err := yaml.Marshal(GlobalConfig)
		if err != nil {
			panic(err)
		}

		f.Write(out)
	}
}
