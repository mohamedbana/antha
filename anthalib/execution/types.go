package execution

// holds types - concrete and interface

// antha config type
type AnthaConfig map[string]interface{}

// a map structure for defining requests to the stock manager
type StockRequest map[string]interface{}

// data structure for defining a request to communicate
// with the sceduler
type ScheduleRequest map[string]interface{}

// map data structure defining a request to find a piece of equipment
type EquipmentManagerRequest map[string]interface{}

// map data structure defining a request for an object
// to enter the waste stream
type GarbageCollectionRequest map[string]interface{}

// data structure for defining a request to the logger
type LogRequest map[string]interface{}

// data structure defining sample requests
type SampleRequest map[string]interface{}
