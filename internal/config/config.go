package config

import (
	"errors"
	"github.com/caarlos0/env/v6"
	"github.com/iancoleman/strcase"
	"github.com/spf13/pflag"
	"log"
	"reflect"
	"sort"
	"strings"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost:8080" flagName:"server-address" flagShortName:"a" flagUsage:"server address"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080" flagName:"base-url" flagShortName:"b" flagUsage:"base url"`
	StorageFilePath string `env:"FILE_STORAGE_PATH" flagName:"storage-file-path" flagShortName:"f" flagUsage:"storage file path"`
}

func LoadFromEnv() *Config {
	cfg := new(Config)
	err := env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

func DefineFlags(cfg *Config) *pflag.FlagSet {
	return DefineFlagSet(pflag.CommandLine, cfg)
}
func DefineFlagSet(fs *pflag.FlagSet, cfg *Config) *pflag.FlagSet {
	if fs == nil {
		panic(errors.New("non-nil pflag.FlagSet expected"))
	}

	// because I'll take care of the sorting myself
	fs.SortFlags = false

	// RV = ReflectValue
	cfgRV := reflect.Indirect(reflect.ValueOf(cfg))
	cfgRType := cfgRV.Type()

	flagNames := make([]string, 0, cfgRType.NumField())
	flagShortNames := make([]string, 0, cfgRType.NumField())
	type AnyFlagKit struct {
		fieldInd       int
		flagName       string
		flagShortName  string
		flagUsageValue string
	}
	anyFlagName2Kit := make(map[string]AnyFlagKit, cfgRType.NumField())
	for fieldInd := 0; fieldInd < cfgRType.NumField(); fieldInd++ {
		cfgStructField := cfgRType.Field(fieldInd)
		cfgStructFieldTag := cfgStructField.Tag

		flagName := cfgStructFieldTag.Get("flagName")
		flagName = strings.TrimSpace(flagName)
		if flagName == "-" {
			continue
		}
		if flagName == "" {
			flagName = strcase.ToDelimited(cfgStructField.Name, '-')
		}

		flagShortName := cfgStructFieldTag.Get("flagShortName")
		flagShortName = strings.TrimSpace(flagShortName)

		flagUsageValue := cfgStructFieldTag.Get("flagUsage")
		flagUsageValue = strings.TrimSpace(flagUsageValue)

		afk := AnyFlagKit{
			fieldInd:       fieldInd,
			flagName:       flagName,
			flagShortName:  flagShortName,
			flagUsageValue: flagUsageValue,
		}

		if flagShortName != "" {
			flagShortNames = append(flagShortNames, flagShortName)
			anyFlagName2Kit[flagShortName] = afk
		} else {
			flagNames = append(flagNames, flagName)
			anyFlagName2Kit[flagName] = afk
		}
	}

	sort.Strings(flagShortNames)
	sort.Strings(flagNames)
	anyFlagNames := append(flagShortNames, flagNames...)

	for _, anyFlagName := range anyFlagNames {
		afk := anyFlagName2Kit[anyFlagName]

		cfgFieldRV := cfgRV.Field(afk.fieldInd)
		cfgFieldValueInterface := cfgFieldRV.Interface()

		cfgFieldValueAddrRV := cfgFieldRV.Addr()
		cfgFieldValueAddrInterface := cfgFieldValueAddrRV.Interface()

		switch cfgFieldValue := cfgFieldValueInterface.(type) {
		case bool:
			fs.BoolVarP(cfgFieldValueAddrInterface.(*bool), afk.flagName, afk.flagShortName, cfgFieldValue, afk.flagUsageValue)
		case int:
			fs.IntVarP(cfgFieldValueAddrInterface.(*int), afk.flagName, afk.flagShortName, cfgFieldValue, afk.flagUsageValue)
		case string:
			fs.StringVarP(cfgFieldValueAddrInterface.(*string), afk.flagName, afk.flagShortName, cfgFieldValue, afk.flagUsageValue)
		default:
			panic(errors.New("unexpected field type"))
		}

	}

	return fs
}
