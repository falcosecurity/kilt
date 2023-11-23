package cfnpatcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/Jeffail/gabs/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	diff "github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
)

var optInTests = [...]string{
	"respect_ignores/opt_in_check",
	"respect_ignores/opt_in_ignored",
	"respect_ignores/opt_in_multiple_containers",
	"respect_ignores/opt_in_single_container",
}

var optPanicTests = [...]string{
	"respect_ignores/panic_opt_exclude",
	"respect_ignores/panic_opt_exclude_container",
	"respect_ignores/panic_opt_include",
	"respect_ignores/panic_opt_include_container",
}

var defaultTests = [...]string{
	"respect_ignores/opt_out_default",
	"respect_ignores/opt_out_ignored",
	"respect_ignores/opt_out_ignore_multiple_containers",
	"respect_ignores/opt_out_ignore_single_container",

	"patching/command",
	"patching/entrypoint",
	"patching/ref",
	"patching/ref_command",
	"patching/ref_env",
	"patching/ref_tags",
	"patching/tags",
	"patching/volumes_from",
}

var parameterizedEnvarsTests = [...]string{
	"patching/parameterize_env_add",
	"patching/parameterize_env_merge",
}

const defaultConfig = `
build {
	entry_point: ["/kilt/run", "--", ${?original.metadata.captured_tag}]
	command: [] ${?original.entry_point} ${?original.command}
	mount: [
		{
			name: "KiltImage"
			image: "KILT:latest"
			volumes: ["/kilt"]
			entry_point: ["/kilt/wait"]
		}
	]
}
`

const parameterizeEnvarsConfig = `
build {
	entry_point: ["/kilt/run", "--", ${?original.metadata.captured_tag}]
	command: [] ${?original.entry_point} ${?original.command}
	environment_variables: {
		"SO_LONG_AND_THANKS": "ForAllTheFish"
	}
	mount: [
		{
			name: "KiltImage"
			image: "KILT:latest"
			volumes: ["/kilt"]
			entry_point: ["/kilt/wait"]
		}
	]
}
`

func runTest(t *testing.T, name string, context context.Context, config Configuration) {
	fragment, err := ioutil.ReadFile("fixtures/" + name + ".json")
	if err != nil {
		t.Fatalf("cannot find fixtures/%s.json", name)
	}
	templateParameters := make([]byte, 0)
	result, err := Patch(context, &config, fragment, templateParameters)
	if err != nil {
		t.Fatalf("error patching: %s", err.Error())
	}
	expected, err := ioutil.ReadFile("fixtures/" + name + ".patched.json")
	if err != nil {
		// To regenerate test simply delete patched variant
		_ = ioutil.WriteFile("fixtures/"+name+".patched.json", result, 0644)
		return
	}

	fmt.Printf("result: %s\n", result)
	fmt.Printf("expected: %s\n", expected)

	differ := diff.New()
	d, err := differ.Compare(expected, result)

	if d.Modified() {
		var expectedJson map[string]interface{}
		t.Log("Found differences!")
		_ = json.Unmarshal(result, &expectedJson) // would error during diff
		formatter := formatter.NewAsciiFormatter(expectedJson, formatter.AsciiFormatterConfig{
			ShowArrayIndex: true,
			Coloring:       true,
		})
		diffString, _ := formatter.Format(d)
		fmt.Println(diffString)
		t.Fail()
	}
}

func TestPatchingOptIn(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	l := log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	for _, testName := range optInTests {
		t.Run(testName, func(t *testing.T) {
			runTest(t, testName, l.WithContext(context.Background()),
				Configuration{
					Kilt:         defaultConfig,
					OptIn:        true,
					RecipeConfig: "{}",
					UseRepositoryHints: false,
				})
		})
	}
}

func TestPatching(t *testing.T) {
	l := log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	for _, testName := range defaultTests {
		t.Run(testName, func(t *testing.T) {
			runTest(t, testName, l.WithContext(context.Background()),
				Configuration{
					Kilt:         defaultConfig,
					OptIn:        false,
					RecipeConfig: "{}",
					UseRepositoryHints: false,
				})
		})
	}
}

func TestPatchingForParameterizingEnvars(t *testing.T) {
	l := log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	for _, testName := range defaultTests {
		t.Run(testName, func(t *testing.T) {
			runTest(t, testName, l.WithContext(context.Background()),
				Configuration{
					Kilt:         defaultConfig,
					OptIn:        false,
					RecipeConfig: "{}",
					UseRepositoryHints: false,
					ParameterizeEnvars: true,
				})
		})
	}

	for _, testName := range parameterizedEnvarsTests {
		t.Run(testName, func(t *testing.T) {
			runTest(t, testName, l.WithContext(context.Background()),
				Configuration{
					Kilt:         parameterizeEnvarsConfig,
					OptIn:        false,
					RecipeConfig: "{}",
					UseRepositoryHints: false,
					ParameterizeEnvars: true,
				})
		})
	}
}

func TestPatchingForLogGroup(t *testing.T) {
	l := log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	tests := []struct {
		Name   string
		Config Configuration
	}{
		{
			"patching/log_group_empty",
			Configuration{
				Kilt:               defaultConfig,
				OptIn:              false,
				RecipeConfig:       "{}",
				UseRepositoryHints: false,
			},
		},
		{
			"patching/log_group",
			Configuration{
				Kilt:               defaultConfig,
				OptIn:              false,
				RecipeConfig:       "{}",
				UseRepositoryHints: false,
				LogGroup:           "test_logs",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			runTest(t, test.Name, l.WithContext(context.Background()), test.Config)
		})
	}
}

func TestOptTagPanic(t *testing.T) {
	l := log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	for _, testName := range optPanicTests {
		t.Run(testName, func(t *testing.T) {
			defer func() { recover() }()

			runTest(t, testName, l.WithContext(context.Background()),
				Configuration{
					Kilt:         defaultConfig,
					OptIn:        true,
					RecipeConfig: "{}",
					UseRepositoryHints: false,
				})
		})
	}
}

func TestIsFuncOptKey(t *testing.T) {
 	tests := []struct {
		key string
		out bool
	}{
		{
			"kilt-ignore",
			true,
		},
		{
			"kilt-include",
			true,
		},
		{
			"kilt-ignore-containers",
			true,
		},
		{
			"kilt-include-containers",
			true,
		},
		{
			"so-long-and-thanks-for-all-the-fish",
			false,
		},
	}
	for _, test := range tests {
		t.Run(test.key, func(t *testing.T) {
			assert.Equal(t, test.out, isOptTagKey(test.key), "OptIn/Out key not recognized")
		})
	}
}

func TestGetOptTags(t *testing.T) {
 	tests := []struct {
		name     string
		json     string
		expected map[string]string
	}{
		{
			name: `no-properties`,
			json: `{
"McGuffin": {
	"Tags":[
		{
			"Key": "SoLong",
			"Value": "AndThanksForAllTheFish"
		}
	]}
}`,
			expected: make(map[string]string),
		},
		{
			name: `no-tags`,
			json: `{
"Properties": {
	"Accio":[
		{
			"Key": "SoLong",
			"Value": "AndThanksForAllTheFish"
		}
	]}
}`,
			expected: make(map[string]string),
		},
		{
			name: `no-opt-tags`,
			json: `{
"Properties": {
	"Tags":[
		{
			"Key": "SoLong",
			"Value": "AndThanksForAllTheFish"
		},
		{
			"Key": "TimeIsAnIllusion",
			"Value": "LunchtimeDoublySo"
		}
	]}
}`,
			expected: make(map[string]string),
		},
		{
			name: `all-opt-tags`,
			json: `{
"Properties": {
	"Tags": [
		{
			"Key": "SoLong",
			"Value": "AndThanksForAllTheFish"
		},
		{
			"Key": "TimeIsAnIllusion",
			"Value": "LunchtimeDoublySo"
		},
		{
			"Key": "kilt-ignore",
			"Value": "nanananananaBatman"
		},
		{
			"Key": "kilt-include",
			"Value": "gimmeGimmeGimmeFriedChicken"
		},
		{
			"Key": "kilt-ignore-containers",
			"Value": "expelliarmus"
		},
		{
			"Key": "kilt-include-containers",
			"Value": "accioContainer"
		}
	]}
}`,
			expected: map[string]string {
				"kilt-ignore": "nanananananaBatman",
				"kilt-include": "gimmeGimmeGimmeFriedChicken",
				"kilt-ignore-containers": "expelliarmus",
				"kilt-include-containers": "accioContainer",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			jsonParsed, err := gabs.ParseJSON([]byte(tc.json))
			if err != nil {
				panic(err)
			}
			mm := getOptTags(jsonParsed)
			eq := reflect.DeepEqual(tc.expected, mm)
			if !eq {
				assert.Fail(t, "maps do not match")
			}
		})
	}
}

func TestGetParameterName(t *testing.T) {
 	tests := []struct {
		name     string
		expected string
	}{
		// No changes if there are no non-alphanumeric chars
		{
			name: `SOLONGANDTHANKSFORALLTHEFISH12345`,
			expected: `SOLONGANDTHANKSFORALLTHEFISH12345`,
		},
		{
			name: `solongandthanksforallthefish12345`,
			expected: `solongandthanksforallthefish12345`,
		},
		{
			name: `soLongAndThanksForAllTheFish12345`,
			expected: `soLongAndThanksForAllTheFish12345`,
		},
		// Tries to make the parameter name more readable if there are non-alphanumeric chars
		{
			name: `SOLONGANDTHANKSFORALLTHEFISH_`,
			expected: `solongandthanksforallthefish`,
		},
		{
			name: `SOLONG_ANDTHANKSFORALLTHEFISH`,
			expected: `solongAndthanksforallthefish`,
		},
		{
			name: `SO_LONG_AND_THANKS_FOR_ALL_THE_FISH`,
			expected: `soLongAndThanksForAllTheFish`,
		},
		{
			name: `_SO_LONG_AND_THANKS_FOR_ALL_THE_FISH_`,
			expected: `SoLongAndThanksForAllTheFish`,
		},
		{
			name: `__SO__LONG__AND__THANKS__FOR__ALL__THE__FISH__`,
			expected: `SoLongAndThanksForAllTheFish`,
		},
		{
			name: `solongandthanksforallthefish_`,
			expected: `solongandthanksforallthefish`,
		},
		{
			name: `solong_andthanksforallthefish`,
			expected: `solongAndthanksforallthefish`,
		},
		{
			name: `so_long_and_thanks_for_all_the_fish`,
			expected: `soLongAndThanksForAllTheFish`,
		},
		{
			name: `_so_long_and_thanks_for_all_the_fish_`,
			expected: `SoLongAndThanksForAllTheFish`,
		},
		{
			name: `__so__long__and__thanks__for__all__the__fish__`,
			expected: `SoLongAndThanksForAllTheFish`,
		},
		{
			name: `soLong_AndThanksForAllTheFish`,
			expected: `solongAndthanksforallthefish`,
		},
		{
			name: `so_Long_And_Thanks_For_All_The_Fish`,
			expected: `soLongAndThanksForAllTheFish`,
		},
		{
			name: `_so_Long_And_Thanks_For_All_The_Fish_`,
			expected: `SoLongAndThanksForAllTheFish`,
		},
		{
			name: `__so__Long__And__Thanks__For__All__The__Fish__`,
			expected: `SoLongAndThanksForAllTheFish`,
		},
		// Won't happen, actually
		{
			name: `soLong-ANDTHANKS_forAllTheFish___`,
			expected: `solongAndthanksForallthefish`,
		},
		{
			name: `soLong-ANDTHANKS-forAllTheFish!!!`,
			expected: `solongAndthanksForallthefish`,
		},
		{
			name: `soLong-ANDTHANKS-forAllTheFish!!!`,
			expected: `solongAndthanksForallthefish`,
		},
		{
			name: `soLongAndThanksForAllTheFish!!!`,
			expected: `solongandthanksforallthefish`,
		},
		{
			name: `***so___Long---And!!!Thanks???For+++All***The:::Fish|||`,
			expected: `SoLongAndThanksForAllTheFish`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, getParameterName(tc.name))
		})
	}
}
