package main

import (
    "fmt"
    "os"
    "time"
    "github.com/tarm/serial"
    "strconv"
)

var parity_map = map[byte] serial.Parity {
    'n': serial.ParityNone,
    'o': serial.ParityOdd,
    'e': serial.ParityEven,
    'm': serial.ParityMark,
    's': serial.ParitySpace,
}

var stopbits_map = map[int] serial.StopBits {
    1:  serial.Stop1,
    15: serial.Stop1Half,
    2:  serial.Stop2,
}

func append2log (log *os.File, t0 time.Time, t1 time.Time, n int, line string) {
    var fullline string = fmt.Sprintf("%d.%09d,%d.%09d,%d,%s", t0.Unix(), t0.Nanosecond(), t1.Unix(), t1.Nanosecond(), n, line)
    
    if _, err := log.Write([]byte(fullline)) ; err != nil {
        fmt.Println(err)
        os.Exit(3)
    }
}

func main () {
    // guard: command line args
    if len(os.Args) != 6 {
        fmt.Printf("Syntax: %s DEVICE BAUDRATE PARITY STOPBITS FILENAME\n", os.Args[0])
        fmt.Printf("        %s /dev/ttyACM0 9600 n 1 log.csv\n", os.Args[0])
        fmt.Printf("        %s COM1 9600 n 1 log.csv\n", os.Args[0])
        os.Exit(1)
    }
    var dev_path   string = os.Args[1]
    var dev_baud_s string = os.Args[2]
    var dev_par_s  string = os.Args[3]
    var dev_stop_s string = os.Args[4]
    var log_path   string = os.Args[5]
    
    // guard: stop sanity
    dev_baud_i, err := strconv.Atoi(dev_baud_s)
    if err != nil {
        fmt.Println("Error converting baudrate number")
        fmt.Println(err)
        os.Exit(1)
    }
    
    // guard: parity sanity
    if len(dev_par_s)!=1 {
        fmt.Println("Wrong parity length. Try one of [n,o,e,m,s] ...")
        os.Exit(1)
    }
    var dev_par_c byte = dev_par_s[0]
    if _, exists := parity_map[dev_par_c]; !exists {
        fmt.Println("Unknown parity. Try one of [n,o,e,m,s] ...")
        os.Exit(1)
    }
    
    // guard: stop sanity
    dev_stop_i, err := strconv.Atoi(dev_stop_s)
    if err != nil {
        fmt.Println("Error converting number of stopbits")
        fmt.Println(err)
        os.Exit(1)
    }
    if _, exists := stopbits_map[dev_stop_i]; !exists {
        fmt.Println("Unknown number of stop bits. Try one of [1,15,2] ...")
        os.Exit(1)
    }
    
    // print out configuration
    fmt.Printf("serial-siggen: %s[%d,%s,%d] -> %s\n", dev_path, dev_baud_i, dev_par_s, dev_stop_i, log_path)
    
    // parse serial options
    var dev_par serial.Parity   = parity_map[dev_par_c]
    var dev_stop serial.StopBits = stopbits_map[dev_stop_i]
    
    // open log file
    log, err := os.OpenFile(log_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println(err)
        os.Exit(2)
    }
    
    // open serial port
    c := &serial.Config{Name: dev_path, Baud: dev_baud_i, Parity: dev_par, StopBits: dev_stop}
    s, err := serial.OpenPort(c)
    if err != nil {
        fmt.Println(err)
        os.Exit(4)
    }
    
    // service loop
    var i int = 0
    for {
        // produce string to transmit
        var data string = fmt.Sprintf("%d\n", i)
        
        // transmit string
        var t0 time.Time = time.Now()
        n, err := s.Write([]byte(data)) // TODO: Add loop to ensure that all is transmitted
        if err != nil {
            fmt.Println(err)
            os.Exit(4)
        }
        var t1 time.Time = time.Now()
        
        // log transmitted string
        append2log(log, t0, t1, n, data)
        
        // sleep
        time.Sleep(1*time.Second)
        
        // increment counter
        i++
    }
}

