package main

import (
    "flag"
    "fmt"
    "os"
    "log"
    "bufio"
)

func main() {
    mysqlHost := flag.String("host", "", "mysql hostname")
    mysqlPort := flag.Int("port", 3306, "mysql port")
    mysqlUser := flag.String("user", "", "mysql user")
    mysqlPass := flag.String("pass", "", "mysql password")
    fileName := flag.String("file", "", "file containing ip addresses")
    showResults := flag.Bool("results", false, "show results")

    flag.Parse()

    if len(*mysqlHost) == 0 || len(*mysqlUser) == 0 || len(*mysqlPass) == 0 || len(*fileName) == 0 || *mysqlPort == 0 {
        printUsage()
    }

    if *showResults == true {
        results()
    }

    data, err := readFile(*fileName)
    if err != nil {
        log.Fatalf("error: %s", err.Error())
    }

    fmt.Println(data)
}

func printUsage()  {
    fmt.Println("snmp-swith-port-monitor -- monitor port up/down state via SNMP and writing results into MySQL")
    fmt.Println("")
    fmt.Println("Usage:")
    flag.PrintDefaults()
    fmt.Println("")
    fmt.Println("Examples:")
    fmt.Println("  snmp-switch-port-monitor --host host --user user --pass --file /path/to/file.txt")
    fmt.Println("  snmp-switch-port-monitor --results")
    fmt.Println("  snmp-switch-port-monitor --help")
    os.Exit(0)
}

func readFile(path string) ([]string, error) {
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

func results() {
    fmt.Println("Here will be results")
    os.Exit(0)
}
