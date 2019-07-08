package event

import (
	"github.com/optimizely/go-sdk/optimizely/entities"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

type TestConfig struct {
}

func (TestConfig) GetEventByKey(string) (entities.Event, error) {
	return entities.Event{ExperimentIds: []string{"15402980349"}, ID: "15368860886", Key: "sample_conversion"}, nil
}

func (TestConfig)GetFeatureByKey(string) (entities.Feature, error) {
	return entities.Feature{}, nil
}

func (TestConfig)GetProjectID() string {
	return "15389410617"
}
func (TestConfig)GetRevision()  string {
	return "7"
}
func (TestConfig)GetAccountID() string {
	return "8362480420"
}
func (TestConfig)GetAnonymizeIP() bool {
	return true
}
func (TestConfig)GetAttributeID(key string) string { // returns "" if there is no id
	return ""
}
func (TestConfig)GetBotFiltering() bool {
	return false
}

func RandomString(len int) string {
	bytes := make([]byte, len)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + rand.Intn(25))  //A=65 and Z = 65+25
	}
	return string(bytes)
}

var userId = RandomString(10)

func BuildTestImpressionEvent() UserEvent {
	config := TestConfig{}

	experiment := entities.Experiment{}
	experiment.Key = "background_experiment"
	experiment.LayerID = "15399420423"
	experiment.ID = "15402980349"

	variation := entities.Variation{}
	variation.Key = "variation_a"
	variation.ID = "15410990633"

	logEvent := CreateImpressionUserEvent(config, experiment, variation, userId, make(map[string]interface{}))

	return logEvent
}

func BuildTestConversionEvent() UserEvent {
	config := TestConfig{}

	logEvent := CreateConversionUserEvent(config, "sample_conversion", userId, make(map[string]interface{}),make(map[string]interface{}))

	return logEvent
}

func TestCreateImpressionEvent(t *testing.T) {

	logEvent := BuildTestImpressionEvent()

	processor := NewEventProcessor(100, 100)

	processor.ProcessImpression(logEvent)

	result, ok := processor.(*DefaultEventProcessor)

	if ok {
		assert.Equal(t, 1, result.EventsCount())

		time.Sleep(2000 * time.Millisecond)

		assert.Equal(t, 0, result.EventsCount())
	}
}

func TestCreateConversionEvent(t *testing.T) {

	logEvent := BuildTestConversionEvent()

	processor := NewEventProcessor(100, 100)

	processor.ProcessImpression(logEvent)

	result, ok := processor.(*DefaultEventProcessor)

	if ok {
		assert.Equal(t, 1, result.EventsCount())

		time.Sleep(2000 * time.Millisecond)

		assert.Equal(t, 0, result.EventsCount())
	}
}