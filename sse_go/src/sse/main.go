/*
Package sse implements the primary inner workings of the SSE Reporter.

The primary function is Run(), which starts a scheduler after initialization and registration of the
reporter with the mothership.
*/

package sse

import (
	"io/ioutil"
	"bytes"
	"collector"
	"encoding/json"
	"helper"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

var (
	mothership_url = "http://mothership.serverstatusmonitoring.com"
	register_uri   = "/register-service"
	collector_uri  = "/collector"
	status_uri     = "/status"

	collect_frequency_in_seconds = 1        // When to collect a snapshot and store in cache
	report_frequency_in_seconds  = 60       // When to report all snapshots in cache
	version                      = "1.0.0"  // The version of SSE this is

	hostname  = ""
	ipAddress = ""

	log_file           = "/var/log/sphire-sse.log"
	configuration_file = "/etc/sse/sse.conf"
	configuration      = new(Configuration)
)

/*
 Configuration struct is a direct map to the configuration located in the configuration JSON file.
*/
type Configuration struct {
	Identification struct {
		AccountID        string `json:"account_id"`
		OrganizationID   string `json:"organization_id"`
		OrganizationName string `json:"organization_name"`
		MachineNickname  string `json:"machine_nickname"`
	} `json:"identification"`
}

/*
 Snapshot struct is a collection of other structs which are relayed from the different segments
 of the collector package.
*/
type Snapshot struct {
	CPU     *collector.CPU     `json:"cpu"`
	Disks   *collector.Disks   `json:"disks"`
	Memory  *collector.Memory  `json:"memory"`
	Network *collector.Network `json:"network"`
	System  *collector.System  `json:"system"`
	Time    time.Time          `json:"system_time"`
}

/*
 Cache struct implements multiple Snapshot structs. This is cleared after it is reported to the mothership.
 Also includes the program Version and AccountId - the latter of which is gleaned from the configuration.
*/
type Cache struct {
	Node      []*Snapshot `json:"node"`
	AccountId string      `json:"account_id"`
	Version   string      `json:"version"`

	OrganizationID   string `json:"organization_id"`
	OrganizationName string `json:"organization_name"`
	MachineNickname  string `json:"machine_nickname"`
}

/*
Run Program entry point which initializes, registers and runs the main scheduler of the
program. Also handles initialization of the global logger.
*/
func Run() {
	// Define the global logger
	logger, err := os.OpenFile(log_file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(helper.Trace("Unable to secure log: "+log_file, "ERROR"))
		os.Exit(1)
	}
	defer logger.Close()
	log.SetOutput(logger)

	log.Println(helper.Trace("**** Starting program ****", "OK"))

	// Perform system initialization
	_, err = Initialize()
	if err != nil {
		log.Println(helper.Trace("Exiting.", "ERROR"))
		os.Exit(1)
	}

	// Perform the system registration
	log.Println(helper.Trace("Performing registration.", "OK"))
	body, err := Register()
	if err != nil {
		log.Println(helper.Trace("Unable to register this machine", "ERROR"))
		os.Exit(1)
	}

	ticker := time.NewTicker(time.Duration(collect_frequency_in_seconds) * time.Second)

	var counter int = 0
	var cache Cache = Cache{
		AccountId:        configuration.Identification.AccountID,
		OrganizationID:   configuration.Identification.OrganizationID,
		OrganizationName: configuration.Identification.OrganizationName,
		MachineNickname:  configuration.Identification.MachineNickname,
		Version:          version}

	for {
		<-ticker.C
		if counter > 0 && counter%report_frequency_in_seconds == 0 {
			cache.Sender()
			cache.Node = nil // Clear the Node Cache
			runtime.GC()
			counter = 0
		} else {
			var snapshot Snapshot = Snapshot{}
			cache.Node = append(cache.Node, snapshot.Collector()) // fill in the Snapshot struct and add to the cache
			counter++
			ticker = updateTicker()
		}
	}
}

/*
updateTicker updates the ticker in order to know when to run the codeblock next
*/
func updateTicker() *time.Ticker {
	var updatedSeconds int = time.Now().Second() + collect_frequency_in_seconds
	nextTick := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(),
		time.Now().Minute(), updatedSeconds, time.Now().Nanosecond(), time.Local)
	return time.NewTicker(nextTick.Sub(time.Now()))
}

/*
Initialize attempts to gather all the data for correct program initialization. Loads config, etc.
returns bool and error - if ever false, error will be set, otherwise if bool is true, error is nil.
*/
func Initialize() (bool, error) {
	var err error = nil

	// Attempt to get the server IP address
	ipAddress, err = helper.GetServerExternalIPAddress()
	if err != nil {
		log.Println(helper.Trace("Initialization failed, IP Address unattainable.", "ERROR"))
		return false, err
	}

	// Load and parse configuration file
	file, _ := os.Open(configuration_file)
	err = json.NewDecoder(file).Decode(configuration)

	if err != nil {
		log.Println(helper.Trace("Initialization failed - could not load configuration.", "ERROR"))
		return false, err
	}

	log.Println(helper.Trace("Initialization complete.", "OK"))
	return true, err
}

/*
Register performs a registration of this instance with the mothership
*/
func Register() (string, error) {
	var jsonStr = []byte(`{}`)

	// local struct
	registrationObject := map[string]interface{}{
		"configuration":     configuration,
		"mothership_url":    mothership_url,
		"register_uri":      register_uri,
		"version":           version,
		"collect_frequency": collect_frequency_in_seconds,
		"report_frequency":  report_frequency_in_seconds,
		"hostname":          hostname,
		"ip_address":        ipAddress,
		"log_file":          log_file,
		"config_file":       configuration_file,
	}

	jsonStr, _ = json.Marshal(registrationObject)
	req, err := http.NewRequest("POST", mothership_url+register_uri, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "REG")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	log.Println(helper.Trace("Registration complete.", "OK"))
	return string(body), nil
}

/*
Collector collects a snapshot of the system at the time of calling and stores it in
Snapshot struct.
*/
func (Snapshot *Snapshot) Collector() *Snapshot {
	Snapshot.Time = time.Now().Local()

	var CPU collector.CPU = collector.CPU{}
	Snapshot.CPU = CPU.Collect()

	var Disks collector.Disks = collector.Disks{}
	Snapshot.Disks = Disks.Collect()

	var Memory collector.Memory = collector.Memory{}
	Snapshot.Memory = Memory.Collect()

	var Network collector.Network = collector.Network{}
	Snapshot.Network = Network.Collect()

	var System collector.System = collector.System{}
	Snapshot.System = System.Collect()

	return Snapshot
}

/*
Sender sends the data in Cache to the mothership, then clears the Cache struct so that it can
accept new data.
*/
func (Cache *Cache) Sender() bool {
	var jsonStr = []byte(`{}`)

	jsonStr, _ = json.Marshal(Cache)
	req, err := http.NewRequest("POST", mothership_url+collector_uri, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "REG")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(helper.Trace("Unable to complete request", "ERROR"))
		return false
	}
	defer resp.Body.Close()

	read_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(helper.Trace("Unable to complete request" + string(read_body), "ERROR"))
		return false
	}

	return true
}
