package timetable_parser_config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config is the top-level structure
type Config struct {
	Version   string    `yaml:"version"`
	Timetable Timetable `yaml:"timetable"`
}

// Timetable contains seasons and routes
type Timetable struct {
	Seasons []Season `yaml:"seasons"`
	Routes  []Route  `yaml:"routes"`
}

// Route represents a single transport route
type Route struct {
	Start  string `yaml:"start"`
	Finish string `yaml:"finish"`
}

// Season represents a single season in the timetable
type Season struct {
	Name      string    `yaml:"name"`
	Start     time.Time `yaml:"start"`
	End       time.Time `yaml:"end"`
	FetchDate time.Time `yaml:"fetch_date"`
	Warning   Warning   `yaml:"warning"`
}

func (s *Season) UnmarshalYAML(value *yaml.Node) error {
	// Create an intermediate struct to unmarshal into
	aux := &struct {
		Name      string  `yaml:"name"`
		Start     string  `yaml:"start"`
		End       string  `yaml:"end"`
		FetchDate string  `yaml:"fetch_date"`
		Warning   Warning `yaml:"warning"`
	}{}

	// Unmarshal the YAML node into the intermediate struct
	if err := value.Decode(aux); err != nil {
		return fmt.Errorf("failed to unmarshal Season: %w", err)
	}

	// Parse the date strings
	start, err := parseDate(aux.Start)
	if err != nil {
		return fmt.Errorf("parseDate [Start='%s']: %w", aux.Start, err)
	}

	end, err := parseDate(aux.End)
	if err != nil {
		return fmt.Errorf("parseDate [End='%s']: %w", aux.End, err)
	}

	fetchDate, err := parseDate(aux.FetchDate)
	if err != nil {
		return fmt.Errorf("parseDate [FetchDate='%s']: %w", aux.FetchDate, err)
	}

	// Assign parsed values
	s.Name = aux.Name
	s.Start = start
	s.End = end
	s.FetchDate = fetchDate
	s.Warning = aux.Warning

	return nil
}

const dateLayout = "2006-01-02"

func parseDate(input string) (time.Time, error) {
	if len(input) == 0 {
		return time.Time{}, nil
	}
	date, err := time.Parse(dateLayout, input)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid Start date format (expected YYYY-MM-DD): %w", err)
	}
	return date, nil
}

type Warning struct {
	Be string `yaml:"be"`
	De string `yaml:"de"`
	En string `yaml:"en"`
	Hr string `yaml:"hr"`
	Ru string `yaml:"ru"`
	Sk string `yaml:"sk"`
	Sr string `yaml:"sr"`
	Tr string `yaml:"tr"`
	Uk string `yaml:"uk"`
}

func Load(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("os.Open: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	var config Config
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, fmt.Errorf("decoder.Decode: %w", err)
	}
	return config, nil
}
