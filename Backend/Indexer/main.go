package main

import (
	"Project/zincShare"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	_ "net/http/pprof"
	"net/mail"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"time"
)

const (
	// pathDataBase = "C:/Users/wpjav/Documents/convertbd/particionada"
	pathDataBase        = "D:/Descargas/new enron/enron_mail_20110402"
	emailMaxSize        = 1000000
	fileEmailSize       = 700
	emailZincSize       = 700
	sizeProcessEmail    = 10
	sizeProcessDataZinc = 10
	batchEmail          = 50
)

var emails = make(chan string, fileEmailSize)
var concurrentProcessEmail = make(chan struct{}, sizeProcessEmail)
var concurrentProcessZinc = make(chan struct{}, sizeProcessDataZinc)
var dataZinc = make(chan string, emailZincSize)
var waitGroupForEmails, waitGroupForZinc sync.WaitGroup
var countBigEmails, countFormaterFalls, inserted, rejected int
var mutex sync.Mutex

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	_ = zincShare.DeleteIndex()
	err := zincShare.CreateIndex()
	if err != nil {
		panic(err)
	}
	currentTime := time.Now()
	waitGroupForZinc.Add(1)

	// Goroutine for send and index with Zincsearch
	go func() {
		defer waitGroupForZinc.Done()
		emailBuilder := strings.Builder{}
		i := 0
		for dataEmail := range dataZinc {
			if emailBuilder.Len()+len(dataEmail) > emailMaxSize || i == batchEmail {
				waitGroupForZinc.Add(1)
				go zincCreate(emailBuilder.String(), i)
				i = 0
				emailBuilder.Reset()
			}
			emailBuilder.WriteString(dataEmail)
			emailBuilder.WriteByte(10)
			i++
		}
		if emailBuilder.Len() != 0 {
			waitGroupForZinc.Add(1)
			go zincCreate(emailBuilder.String(), i)
		}
	}()

	// Goroutine for emails transformation
	go func() {
		for email := range emails {
			go processEmail(email)
		}
	}()

	//Read files of local databas
	findPathFile(pathDataBase)
	close(emails)
	waitGroupForEmails.Wait()
	close(dataZinc)
	waitGroupForZinc.Wait()
	//showing problems in the process
	fmt.Println("Duracion: ", time.Since(currentTime))
	fmt.Println("Email to big: ", countBigEmails)
	fmt.Println("Error in the format: ", countFormaterFalls)
	fmt.Println("Inserted: ", inserted)
	fmt.Println("Rechazados batch: ", rejected)

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}

}

// Find the file in the database
func findPathFile(startPath string) {
	path, err := os.Open(startPath)
	if err != nil {
		log.Println(err)
	}
	defer path.Close()

	files, err := path.ReadDir(-1)
	if err != nil {
		log.Println(err)
	}
	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			findPathFile(startPath + "/" + fileInfo.Name())
		} else {
			file, err := fileInfo.Info()
			if err != nil {
				log.Println(err)
			}
			if file.Size() > emailMaxSize {
				countBigEmails++
				continue
			}
			waitGroupForEmails.Add(1)
			emails <- startPath + "/" + fileInfo.Name()

		}
	}

}

func processEmail(email string) {
	defer waitGroupForEmails.Done()
	concurrentProcessEmail <- struct{}{}
	fileContent, err := os.ReadFile(email)
	if err != nil {
		log.Println(err)
	}
	emailToBytes := bytes.NewReader(fileContent)
	emailConvert, err := mail.ReadMessage(emailToBytes)
	if err != nil {
		handleError(err)
		return
	}
	messageBody, err := io.ReadAll(emailConvert.Body)
	if err != nil {
		handleError(err)
		return
	}
	dataZinc <- fmt.Sprintf(`{"_id": "%s", "from": %s, "to": %s, "subject": %s, "content": %s}`,
		email, fmt.Sprintf("%q", emailConvert.Header.Get("From")), fmt.Sprintf("%q", emailConvert.Header.Get("To")),
		fmt.Sprintf("%q", emailConvert.Header.Get("Subject")), fmt.Sprintf("%q", strings.ReplaceAll(string(messageBody), "\"", "'")))
	<-concurrentProcessEmail
}

func zincCreate(dataEmail string, elements int) {
	defer waitGroupForZinc.Done()
	concurrentProcessZinc <- struct{}{}

	count, err := zincShare.CreateData(dataEmail)
	if err != nil {
		panic(err)
	}

	mutex.Lock()
	inserted += count
	rejected += elements - count
	mutex.Unlock()
	<-concurrentProcessZinc
}

func handleError(err error) {
	mutex.Lock()
	defer mutex.Unlock()
	countFormaterFalls++
	<-concurrentProcessEmail
}
