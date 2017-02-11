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

var lines = []string{"bakerloo",
              "central",
              "circle",
              "district",
              "hammersmith-city",
              "jubilee",
              "metropolitan",
              "northern",
              "piccadilly",
              "victoria",
              "waterloo-city",
              "london-overground",
              "tfl-rail"}

func GetLine(l string, lines []string) (string, error) {
    numlines := 0
    var selectedLine string
    arglen := len(l)
    for _, line := range(lines) {
        if line[:arglen] == l {
            selectedLine = line
            numlines += 1
        }
    }

    if numlines > 1 {
        return "", errors.New("Too many lines start with " + l + ". Try adding an extra letter.")
    } else if numlines < 1 {
        return "", errors.New("There are no lines which start with " + l + ".")
    } else {
        return selectedLine, nil
    }

    return "", errors.New("Could not find associated line")
}

func main() {
    if len(os.Args) < 2 {
        log.Fatal("No arguments provided, exiting")
    }

    app_id := os.Args[1]
    app_key := os.Args[2]

    args := os.Args[3:]

    var selectedLines []string

    for _, l := range(args) {
        line, err := GetLine(l, lines)
        if err != nil {
            log.Fatal(err)
        }

        selectedLines = append(selectedLines, line)
    }

    for _, line := range(selectedLines) {
        fmt.Println(line)
        resp, err := http.Get("https://api.tfl.gov.uk/Line/" + line + "/Disruption?app_id=" + app_id + "&app_key=" + app_key)
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
