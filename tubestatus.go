package main

import(
    "net/http"
    "fmt"
    "log"
    "encoding/json"
    "os"
    "errors"
    "io/ioutil"
)

type TubeLine struct {
    ID string
    FullName string
    Shortcuts []string
}

var Bakerloo = TubeLine{ID: "bakerloo", FullName: "Bakerloo", Shortcuts: []string{"b", "bl"}}
var Central = TubeLine{ID: "central", FullName: "Central", Shortcuts: []string{"ce"}}
var Circle = TubeLine{ID: "circle", FullName: "Circle", Shortcuts: []string{"ci"}}
var District = TubeLine{ID: "district", FullName: "District", Shortcuts: []string{"d"}}
var HammersmithCity = TubeLine{ID: "hammersmith-city", FullName: "Hammersmith and City", Shortcuts: []string{"h", "hc", "hs", "hsc"}}
var Jubilee = TubeLine{ID: "jubilee", FullName: "Jubilee", Shortcuts: []string{"j"}}
var Metropolitan = TubeLine{ID: "metropolitan", FullName: "Metropolitan", Shortcuts: []string{"m"}}
var Northern = TubeLine{ID: "northern", FullName: "Nothern", Shortcuts: []string{"n"}}
var Piccadilly = TubeLine{ID: "piccadilly", FullName: "Piccadilly", Shortcuts: []string{"p"}}
var Victoria = TubeLine{ID: "victoria", FullName: "Victoria", Shortcuts: []string{"v"}}
var WaterlooCity = TubeLine{ID: "waterloo-city", FullName: "Wanterloo and City", Shortcuts: []string{"w", "wc", "wl", "wlc"}}
var Overground = TubeLine{ID: "london-overground", FullName: "London Overground", Shortcuts: []string{"o", "lo", "og", "log"}}
var TflRail = TubeLine{ID: "tfl-rail", FullName: "TFL Rail", Shortcuts: []string{"r", "tflr"}}

var lines = []TubeLine{Bakerloo,
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
                       TflRail}

func GetLine(l string) (TubeLine, error) {
    numlines := 0
    var selectedLine TubeLine
    arglen := len(l)
    for _, line := range(lines) {
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

    if numlines > 1 {
        return selectedLine, errors.New("Too many lines start with " + l + ". Try adding an extra letter.")
    } else if numlines < 1 {
        return selectedLine, errors.New("There are no lines which start with " + l + ".")
    }
    return selectedLine, nil
}

func main() {
    if len(os.Args) < 2 {
        log.Fatal("No arguments provided, exiting")
    }

    app_id := os.Args[1]
    app_key := os.Args[2]

    args := os.Args[3:]

    var selectedLines []TubeLine

    for _, l := range(args) {
        line, err := GetLine(l)
        if err != nil {
            log.Fatal(err)
        }

        selectedLines = append(selectedLines, line)
    }

    for _, line := range(selectedLines) {
        fmt.Println(line.FullName)
        resp, err := http.Get("https://api.tfl.gov.uk/Line/" + line.ID + "/Disruption?app_id=" + app_id + "&app_key=" + app_key)
        if err != nil {
            log.Fatal(err)
        }

        var unmarshalledstatus interface{}

        linejson, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            log.Fatal(err)
        }

        if string(linejson) != "[]" {
            err = json.Unmarshal(linejson, &unmarshalledstatus)
            if err != nil {
                log.Fatal(err)
            }
            status := unmarshalledstatus.([]interface{})[0]
            lineStatuses := status.(map[string]interface{})
            fmt.Println(lineStatuses["description"])
        }
    }
}
