package main

import(
    "net/http"
    "fmt"
    "encoding/json"
    "os"
    "io/ioutil"
    "flag"
    "os/user"
)

type TubeLine struct {
    ID string
    FullName string
    Shortcuts []string
}

var Bakerloo = TubeLine{ID: "bakerloo", FullName: "Bakerloo", Shortcuts: []string{"b", "bl"}}
var Central = TubeLine{ID: "central", FullName: "Central", Shortcuts: []string{"ce"}}
var Circle = TubeLine{ID: "circle", FullName: "Circle", Shortcuts: []string{"ci"}}
var District = TubeLine{ID: "district", FullName: "District", Shortcuts: []string{"di"}}
var HammersmithCity = TubeLine{ID: "hammersmith-city", FullName: "Hammersmith and City", Shortcuts: []string{"h", "hc", "hs", "hsc"}}
var Jubilee = TubeLine{ID: "jubilee", FullName: "Jubilee", Shortcuts: []string{"j"}}
var Metropolitan = TubeLine{ID: "metropolitan", FullName: "Metropolitan", Shortcuts: []string{"m"}}
var Northern = TubeLine{ID: "northern", FullName: "Nothern", Shortcuts: []string{"n"}}
var Piccadilly = TubeLine{ID: "piccadilly", FullName: "Piccadilly", Shortcuts: []string{"p"}}
var Victoria = TubeLine{ID: "victoria", FullName: "Victoria", Shortcuts: []string{"v"}}
var WaterlooCity = TubeLine{ID: "waterloo-city", FullName: "Wanterloo and City", Shortcuts: []string{"w", "wc", "wl", "wlc"}}
var Overground = TubeLine{ID: "london-overground", FullName: "London Overground", Shortcuts: []string{"o", "lo", "og", "log"}}
var TflRail = TubeLine{ID: "tfl-rail", FullName: "TFL Rail", Shortcuts: []string{"r", "tflr"}}
var DLR = TubeLine{ID: "dlr", FullName: "DLR", Shortcuts: []string{"dl", "dlr"}}
var Tram = TubeLine{ID: "tram", FullName: "Tram", Shortcuts: []string{"t"}}

var lines = map[string]TubeLine{ "bakerloo": Bakerloo,
                        "central": Central,
                        "circle": Circle,
                        "district": District,
                        "hammersmith-city": HammersmithCity,
                        "jubilee": Jubilee,
                        "metropolitan": Metropolitan,
                        "northern": Northern,
                        "piccadilly": Piccadilly,
                        "victoria": Victoria,
                        "waterloo-city": WaterlooCity,
                        "london-overground": Overground,
                        "tfl-rail": TflRail,
                        "dlr": DLR,
                        "tram": Tram}

func main() {

    curruser, _ := user.Current()
    homedir := curruser.HomeDir

    var confjson interface{}
    var conf map[string]interface{}
    var app_id string
    var app_key string
    var conf_lines []interface{}

    conffile, err := ioutil.ReadFile(homedir + "/.config/tubestatus/config.json")
    if err == nil {
        _ = json.Unmarshal(conffile, &confjson)
        conf = confjson.(map[string]interface{})
        app_id = conf["app_id"].(string)
        app_key = conf["app_key"].(string)
        conf_lines = conf["lines"].([]interface{})
    }

    var flag_app_id = flag.String("app_id", "nil", "TFL Api app_id")
    var flag_app_key = flag.String("app_key", "nil", "TFL Api app_key")
    flag.Parse()

    if *flag_app_id != "nil" {
        app_id = *flag_app_id
    }

    if *flag_app_key != "nil" {
        app_id = *flag_app_key
    }

    var argsLines []string

    flagLines := flag.Args()

    if len(flagLines) == 0 {
        for _, v := range(conf_lines) {
            argsLines = append(argsLines, v.(string))
        }
    } else {
        argsLines = flagLines
    }

    if len(argsLines) == 0 {
        appname := os.Args[0]
        fmt.Println("Usage of " + appname + ":")
        flag.PrintDefaults()
        os.Exit(1)
    }

    resp, err := http.Get("https://api.tfl.gov.uk/Line/Mode/tube,tflrail,dlr,overground,tram/Status?app_id=" + app_id + "&app_key=" + app_key)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    var bodyJson []interface{}
    json.Unmarshal(bodyBytes, &bodyJson)
    for _,v := range(bodyJson) {
        apiline := v.(map[string]interface{})
        linename := apiline["id"].(string)
        Line := lines[linename]
        shortcuts := Line.Shortcuts
        var foundLine = 0
        for _, sl := range(argsLines) {
            if sl == linename {
                foundLine += 1
            } else {
                for _, sc := range(shortcuts) {
                    if sl == sc {
                        foundLine += 1
                    }
                }
            }
        }
        if foundLine == 1 {
            status := apiline["lineStatuses"].([]interface{})
            for _, v2 := range(status) {
                status2 := v2.(map[string]interface{})
                statusSeverity := status2["statusSeverity"].(float64)
                if statusSeverity != 10 {
                    disruption := status2["disruption"].(map[string]interface{})
                    description := disruption["description"].(string)
                    fmt.Println(description)
                }
            }
        }
    }
}
