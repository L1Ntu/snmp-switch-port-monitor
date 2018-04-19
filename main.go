package main

import (
    "flag"
    "fmt"
    "os"
    "log"
    "bufio"
    g "github.com/soniah/gosnmp"
    "strconv"
    "strings"
    "time"
)

var ipFileName = flag.String("ip-file", "", "file containing ip addresses")
var outFileName = flag.String("out-file", "", "file for output write")

func main() {
    flag.Parse()

    if len(*ipFileName) == 0 || len(*outFileName) == 0 {
        printUsage()
    }

    data, err := readInputFile(*ipFileName)
    if err != nil {
        log.Fatalf("error: %s", err.Error())
    }

    if len(data) == 0 {
        log.Fatal("error: outFile is empty")
    }

    for ip := range data {
        getSnmpData(data[ip])
    }
}

// print usage
func printUsage() {
    fmt.Println("snmp-swith-port-monitor -- monitor port up/down state via SNMP")
    fmt.Println("")
    fmt.Println("Usage:")
    flag.PrintDefaults()
    fmt.Println("")
    fmt.Println("Examples:")
    fmt.Println("  snmp-switch-port-monitor --ip-file /file/ip-file.txt --out-file /file/result.dat")
    fmt.Println("  snmp-switch-port-monitor --help")
    os.Exit(0)
}

// read ip-file to get list of IP's
func readInputFile(path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }

    return lines, scanner.Err()
}

// write to file
func writeOutFile(data string) {
    if _, err := os.Stat(*outFileName); os.IsNotExist(err) {
        outFile, err := os.OpenFile(*outFileName, os.O_CREATE, 0644)
        if err != nil {
            log.Fatalf("error: cant create outFile, %s", err)
        }

        outFile.Close()
    }

    outFile, err := os.OpenFile(*outFileName, os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        log.Fatalf("error: can't open outFile, %s", err)
    }
    defer outFile.Close()

    length, err := outFile.WriteString(data)
    if err != nil {
        log.Fatalf("error: cant' write to file, %s%d", err, length)
    }
}

// getting SNMP data from host and write them into file
func getSnmpData(host string) {
    fmt.Printf("getting data from %s\n", host)

    g.Default.Target = host

    err := g.Default.Connect()
    if err != nil {
        fmt.Printf("error: connect() to host %s\n", host)
        return
    }
    defer g.Default.Conn.Close()

    oids := []string{}
    for idx := 1; idx <= 52; idx++ {
        oids = append(oids, strings.Join([]string{"1.3.6.1.2.1.2.2.1.8.", strconv.Itoa(idx)}, ""))
    }

    result, err2 := g.Default.Get(oids)
    if err2 != nil {
        log.Printf("error: %s %v", host, err2)
    } else {
        for i, variable := range result.Variables {
            if variable.Value == 1 || variable.Value == 2 {
                portNumber := i + 1
                t := time.Now().Format("2006-01-02 15:04:05")
                s := fmt.Sprintf("%s\t%d\t%d\t%s\n", host, portNumber, variable.Value, t)

                writeOutFile(s)
            }
        }
    }
}
