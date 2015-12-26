package main
import(
  "gopkg.in/alecthomas/kingpin.v2"
  "fmt"
  "log"
  "net"
  "./ernyi"
  "github.com/hashicorp/memberlist"
  "strconv"
)

var (
	command = kingpin.Arg("command", "Command").Required().String()
	name = kingpin.Flag("name", "Name of the node").String()
	addr = kingpin.Flag("addr", "Address of Ernyi node in format host:port").String()
)

var (
	create = "create"
	join = "join"
	info = "info"
)

func CreateErnyi(){
	mconfig := memberlist.DefaultLANConfig()
	if *name == "" {
		log.Fatal("Name must be non-empty")
	}

	if *addr == "" {
		log.Fatal("Address must be non-empty")
	}
	shost, sport, err := net.SplitHostPort(*addr)
	if err != nil {
		log.Fatal(err)
	}

	mconfig.Name = *name
	mconfig.BindAddr = shost
	res, erratoi := strconv.Atoi(sport)
	if erratoi != nil {
		log.Fatal(erratoi)
	}
	mconfig.BindPort = res

	cfg := &ernyi.Config {
		MemberlistConfig: mconfig,
	}

	value := ernyi.CreateErnyi(cfg)
	log.Printf("Ernyi is started")
	value.Start()

}

func Join() {
	fmt.Println(*addr)
}

func ProcessCommands() {
	switch *command {
		case create:
			CreateErnyi()
		case join:
			Join()
		default:
			fmt.Println("Unknown command")
	}
}

func main() {
	kingpin.Version("0.0.1")
  	kingpin.Parse()
  	ProcessCommands()
}