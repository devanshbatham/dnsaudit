package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"
	"log"
)

type ZonemasterResponse struct {
	Messages       []interface{} `json:"messages"`
	SeeAlso        []interface{} `json:"see_also"`
	Version        string        `json:"version"`
	DataCallName   string        `json:"data_call_name"`
	DataCallStatus string        `json:"data_call_status"`
	Cached         bool          `json:"cached"`
	Data           struct {
		Result []struct {
			ID string `json:"id"`
		} `json:"result"`
	} `json:"data"`
}

type ZonemasterDetailsResponse struct {
	Data struct {
		Result struct {
			Results []struct {
				Module  string `json:"module"`
				Level   string `json:"level"`
				Message string `json:"message"`
			} `json:"results"`
		} `json:"result"`
	} `json:"data"`
}

func getZonemasterResults(resourceID string) (*ZonemasterResponse, error) {
	baseURL := "https://stat.ripe.net/data/zonemaster/data.json"

	resp, err := http.Get(fmt.Sprintf("%s?resource=%s", baseURL, resourceID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var zonemasterResponse ZonemasterResponse
	if err := json.NewDecoder(resp.Body).Decode(&zonemasterResponse); err != nil {
		return nil, err
	}

	return &zonemasterResponse, nil
}

func getZonemasterDetailsResults(resourceID string) (*ZonemasterDetailsResponse, error) {
	baseURL := "https://stat.ripe.net/data/zonemaster/data.json"

	resp, err := http.Get(fmt.Sprintf("%s?resource=%s&method=details", baseURL, resourceID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var zonemasterDetailsResponse ZonemasterDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&zonemasterDetailsResponse); err != nil {
		return nil, err
	}

	return &zonemasterDetailsResponse, nil
}

func printColoredText(colorCode int, text string) {
	fmt.Printf("\x1b[%dm%s\x1b[0m", colorCode, text)
}

func main() {
	domainPtr := flag.String("domain", "", "Domain name to query")
	updatePtr := flag.Bool("update", false, "Update cached results")

	flag.Parse()
	    // Print tool banner
	log.SetFlags(0)
    log.Print(`


       __                           ___ __ 
  ____/ /___  __________ __  ______/ (_) /_
 / __  / __ \/ ___/ __  / / / / __  / / __/
/ /_/ / / / (__  ) /_/ / /_/ / /_/ / / /_  
\____/_/ /_/____/\____/\____/\____/_/\__/  
                                           

 
`)

	if *domainPtr == "" {
		fmt.Println("Please provide a domain using the -domain flag.")
		return
	}

	if *updatePtr {
		// Make a new request to start scan
		_, err := http.Get(fmt.Sprintf("https://stat.ripe.net/data/zonemaster/data.json?resource=%s&method=test", *domainPtr))
		if err != nil {
			fmt.Println("Error starting new scan:", err)
			return
		}

		fmt.Printf("\x1b[32m[INFO]\x1b[0m Making new request to start scan for %s\n", *domainPtr)
		time.Sleep(30 * time.Millisecond)
		fmt.Printf("\x1b[32m[INFO]\x1b[0m New Scan started\n")
		time.Sleep(30 * time.Millisecond)
		fmt.Printf("\x1b[32m[INFO]\x1b[0m Scan results should be available in ~30 secs\n")
		time.Sleep(30 * time.Millisecond)
		fmt.Printf("\x1b[32m[INFO]\x1b[0m Run 'dnsaudit -domain %s' again after 30 secs\n", *domainPtr)
		time.Sleep(30 * time.Millisecond)
		fmt.Printf("\x1b[32m[INFO]\x1b[0m Quitting...\n")
		return
	}

	results, err := getZonemasterResults(*domainPtr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(results.Data.Result) == 0 {
		fmt.Printf("\x1b[32m[INFO]\x1b[0m No cached results found. Making new request to start scan for %s\n", *domainPtr)
		time.Sleep(20 * time.Millisecond)
		// Make a new request to start scan
		_, err = http.Get(fmt.Sprintf("https://stat.ripe.net/data/zonemaster/data.json?resource=%s&method=test", *domainPtr))
		if err != nil {
			fmt.Println("[ERROR] Error starting new scan:", err)
			return
		}

		fmt.Printf("\x1b[32m[INFO]\x1b[0m New Scan started\n")
		time.Sleep(30 * time.Millisecond)
		fmt.Printf("\x1b[32m[INFO]\x1b[0m Scan results should be available in ~30 secs\n")
		time.Sleep(30 * time.Millisecond)
		fmt.Printf("\x1b[32m[INFO]\x1b[0m Run 'dnsaudit -d %s' again after 30 secs\n", *domainPtr)
		time.Sleep(30 * time.Millisecond)
		fmt.Printf("\x1b[32m[INFO]\x1b[0m Quitting...\n")
		time.Sleep(30 * time.Millisecond)
		return
	}

	resourceID := results.Data.Result[0].ID
	detailsResults, err := getZonemasterDetailsResults(resourceID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(detailsResults.Data.Result.Results) > 0 {
		for i, item := range detailsResults.Data.Result.Results {
			if strings.ToUpper(item.Module) != "SYSTEM" {
				var levelColor int
				switch strings.ToUpper(item.Level) {
				case "INFO":
					levelColor = 32 // Green
				case "WARNING":
					levelColor = 33 // Yellow
				case "ERROR":
					levelColor = 31 // Red
				case "NOTICE":
					levelColor = 35 // Magenta
				default:
					levelColor = 37 // White
				}

				moduleColor := 36 // Cyan for all modules
				fmt.Print("[")
				printColoredText(levelColor, item.Level)
				fmt.Print("] [")
				printColoredText(moduleColor, item.Module)
				fmt.Printf("] \x1b[33m%s\x1b[0m", item.Message)

				if i < len(detailsResults.Data.Result.Results)-1 {
					fmt.Print(" ") // Space between messages
				}
				time.Sleep(20 * time.Millisecond)
			}
		}
		fmt.Println() // Print a newline after all output
	} else {
		fmt.Println("No 'results' field found in the JSON.")
	}
}
