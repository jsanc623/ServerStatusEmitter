package main

var (

	// URIRegister is the uri to be used to register the system this program will run on
	//URI = {
	//	Register = "/register"
	//
	//	// URICollector is the uri where collected data will be sent to
	//	URICollector = "/collector"
	//
	//	// URIStatus is the uri to check the upstatus of URL
	//	URIStatus = "/status"
	//}



	// Hostname is the hostname of the system this program will run on
	//Hostname = ""

	// IPAddress is the IP address of the system this program will run on
	//IPAddress = ""

	// LogFile is the file where we want to log event data and errors
	//LogFile = "/var/log/sphire-sse.log"

	// ConfigurationFile is the configuration file we want to use
	//ConfigurationFile = "/etc/sse/config.json"

	// Configuration is the configuration instance (loads the above LogFile)
	//Configuration = new(Config)

	// CollectFrequencySeconds is value which tells us
	// to collect a snapshot and store in cache
	// every X seconds where X is a non negative integer
	//CollectFrequencySeconds = 1

	// ReportFrequencySeconds tells us the frequency
	// in seconds to report all snapshots in cache
	//ReportFrequencySeconds = 1

	//// CPU is an instance of collector.CPU
	//CPU = collector.CPU{}
	//
	//// Disks is an instance of collector.Disks
	//Disks = collector.Disks{}
	//
	//// Memory is an instance of collector.Memory
	//Memory = collector.Memory{}
	//
	//// Network is an instance of collector.Network
	//Network = collector.Network{}
	//
	//// System is an instance of collector.System
	//System = collector.System{}

	// Version denotes the version of this program
	//Version = "1.0.1"
)

func main() {
	// Define the global logger
	//logger, err := os.OpenFile(LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//HandleError(err)
	//defer logger.Close()
	//log.SetOutput(logger)
	//
	//// Load and parse configuration file
	//file, err := os.Open(ConfigurationFile)
	//HandleError(err)
	//err = json.NewDecoder(file).Decode(Configuration)
	//HandleError(err)
	//
	//// Set any parameters that need to be set
	//CollectFrequencySeconds = Configuration.Settings.Reporting.CollectFrequencySeconds
	//ReportFrequencySeconds = Configuration.Settings.Reporting.ReportFrequencySeconds
	//
	//var status = helper.Status{}
	//var statusResult = status.CheckStatus(URL + URIStatus)
	//if statusResult == false {
	//	HandleError(errors.New("mothership unreachable - check your internet connection"))
	//}
	//
	//// Perform system initialization
	//var serverObj = sse.Server{}
	//server, ipAddress, hostname, version, error := serverObj.Initialize()
	//HandleError(error)
	//
	//IPAddress = ipAddress
	//Hostname = hostname
	//Version = version
	//
	//// Perform registration
	//body, err := sse.Register(map[string]interface{}{
	//	"configuration":     Configuration,
	//	"mothership_url":    URL,
	//	"register_uri":      URIRegister,
	//	"version":           Version,
	//	"collect_frequency": CollectFrequencySeconds,
	//	"report_frequency":  ReportFrequencySeconds,
	//	"hostname":          Hostname,
	//	"ip_address":        IPAddress,
	//	"log_file":          LogFile,
	//	"config_file":       ConfigurationFile,
	//}, URL+URIRegister+"/"+Version)
	//if err != nil {
	//	HandleError(errors.New("Unable to register this machine" + string(body)))
	//}
	//
	//// Set up our collector
	//var counter int
	//var snapshot = sse.Snapshot{}
	//var cache = sse.Cache{
	//	AccountID:        Configuration.Identification.AccountID,
	//	OrganizationID:   Configuration.Identification.OrganizationID,
	//	OrganizationName: Configuration.Identification.OrganizationName,
	//	MachineNickname:  Configuration.Identification.MachineNickname,
	//	Version:          Version,
	//	Server:           server}
	//
	//ticker := time.NewTicker(time.Duration(CollectFrequencySeconds) * time.Second)
	//death := make(chan os.Signal, 1)
	//signal.Notify(death, os.Interrupt, os.Kill)
	//
	//go func() {
	//	for {
	//		select {
	//		case <-ticker.C: // send the updated time back via to the channel
	//			// reset the snapshot to an empty struct
	//			snapshot = sse.Snapshot{}
	//
	//			// fill in the Snapshot struct and add to the cache
	//			cache.Node = append(cache.Node, snapshot.Collector(Configuration.Settings.Disk.IncludePartitionData,
	//				Configuration.Settings.System.IncludeUsers))
	//			counter++
	//
	//			if counter > 0 && counter%ReportFrequencySeconds == 0 {
	//				cache.Sender(URL + URICollector)
	//				cache.Node = nil // Clear the Node Cache
	//				counter = 0
	//			}
	//		case <-death:
	//			fmt.Println("died")
	//			return
	//		}
	//	}
	//}()

	return

}