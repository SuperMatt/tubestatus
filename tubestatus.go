package main

import(
    "net/http"
    "fmt"
    "log"
    "encoding/json"
    "os"
    "errors"
    "io/ioutil"
    "flag"
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

var lines = []TubeLine{ Bakerloo,
                        Central,
                        Circle,
                        District,
                        HammersmithCity,
                        Jubilee,
                        Metropolitan,
                        Northern,
                        Piccadilly,
                        Victoria,
                        WaterlooCity,
                        Overground,
                        TflRail,
                        DLR,
                        Tram}

func GetLine(l string) (TubeLine, error) {
    numlines := 0
    var selectedLine TubeLine
    for _, line := range(lines) {
        if line.ID == l {
            selectedLine = line
            numlines += 1
        } else if len(l) <= len(line.ID) {
            arglen := len(l)
            if line.ID[:arglen] == l {
                selectedLine = line
                numlines += 1
            } else {
                for _, sc := range(line.Shortcuts) {
                    if sc == l {
                        selectedLine = line
                        numlines += 1
                    }
                }
            }
        }
    }

    if numlines > 1 {
        return selectedLine, errors.New("Too many lines start with " + l + ". Try adding an extra letter.")
    } else if numlines < 1 {
        return selectedLine, errors.New("There are no lines which start with " + l + ".")
    }
    return selectedLine, nil
}

func main() {
    var app_id = flag.String("app_id", "nil", "TFL Api app_id")
    var app_key = flag.String("app_key", "nil", "TFL Api app_key")
    flag.Parse()

    lines := flag.Args()

    if len(lines) == 0 {
        appname := os.Args[0]
        fmt.Println("Usage of " + appname + ":")
        flag.PrintDefaults()
        os.Exit(1)
    }

    var selectedLines []TubeLine

    for _, l := range(lines) {
        line, err := GetLine(l)
        if err != nil {
            log.Fatal(err)
        }

        selectedLines = append(selectedLines, line)
    }

    resp, err := http.Get("https://api.tfl.gov.uk/Line/Mode/tube,tflrail,dlr,overground,tram/Status?app_id=" + *app_id + "&app_key=" + *app_key)
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
        for _, selectedLine :=range(selectedLines) {
            if selectedLine.ID == apiline["id"] {
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
}
