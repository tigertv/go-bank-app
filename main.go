package main

import (
	"net/http"
	"strconv"
	"encoding/json"
	"log"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"html/template"
	"syscall"
)

type Page struct {
	Title string
	Msg string
}

type Result struct {
	Result string
}

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

var bank = &Bank{}

func printOutJson(w http.ResponseWriter, json []byte) {
	w.Header().Set("Content-type", "application/json;charset=UTF-8")
	w.Write(json)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	if r.URL.Path == "/"  {
		t, _ := template.ParseFiles("templates/index.html")
		t.Execute(w, &Page{Title: "Index page", Msg: "Message"})
	}
}

func transferHandler(w http.ResponseWriter, r *http.Request) {
	fromIdStr:= r.FormValue("fromid")
	toIdStr:= r.FormValue("toid")
	sumStr := r.FormValue("sum")

	const msgText = "transferHandler: %s\n"
	fromId, err := strconv.Atoi(fromIdStr)
	if err != nil {
		Error.Printf(msgText, err)
		return
	}

	toId, err := strconv.Atoi(toIdStr)
	if err != nil {
		Error.Printf(msgText, err)
		return
	}

	sum, err := strconv.ParseFloat(sumStr, 64)
	if err != nil {
		Error.Printf(msgText, err)
		return
	}

	isTransferEnded := bank.transfer(fromId, toId, sum)

	var result Result

	msg:= "Transfer from "+fromIdStr+" to "+toIdStr+" sum: "+sumStr
	if (isTransferEnded) {
		result = Result{Result: "OK"}
		Info.Println(msg)
	} else {
		result = Result{Result: "Error"}
		Error.Println(msg)
	}

	jsonData, err := json.Marshal(result)

	if err != nil {
		Error.Printf(msgText, err)
		return
	}

	printOutJson(w, jsonData) 
}

func addAccountHandler(w http.ResponseWriter, r *http.Request) {
	sum := r.FormValue("sum")
	accountId := bank.addAccount(sum)

	const header = "addAccountHandler"

	if accountId < 0 {
		Error.Printf(header + ": can't create an account \n")
		return
	}

	account := bank.getAccount(accountId)
	if account == nil{
		Error.Printf(header + ": can't get an account\n")
		return
	}

	jsonData, err := json.Marshal(account)
	if err != nil {
		Error.Printf(header + ": %s\n", err)
		return
	}

	Info.Println("Created an account with id: "+strconv.Itoa(account.Id)+" balance: "+sum)

	printOutJson(w, jsonData) 
}

func getAccountHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Error.Printf("getAccountHandler: %s\n", err)
		return
	}

	account := bank.getAccount(id)
	msg := "Getting the account with id: "+idStr
	if account == nil {
		Error.Println(msg)
		return
	}

	Info.Println(msg)

	jsonData, err := json.Marshal(account)

	if err != nil {
		Error.Println(err)
		return
	}

	printOutJson(w, jsonData) 
}

func cssHandler(w http.ResponseWriter, r *http.Request) {
	path := "templates" + r.URL.Path

	file, err := ioutil.ReadFile(path)

	if err != nil {
		Error.Println(err)
		return
	}

	w.Header().Set("Content-type", "text/css")
	w.Write(file)
}

func Init(
    traceHandle io.Writer,
    infoHandle io.Writer,
    warningHandle io.Writer,
    errorHandle io.Writer) {

    Trace = log.New(traceHandle,
        "TRACE: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Info = log.New(infoHandle,
        "INFO: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Warning = log.New(warningHandle,
        "WARNING: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Error = log.New(errorHandle,
        "ERROR: ",
        log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	osSignal := make(chan os.Signal, 1)

	signal.Notify(osSignal, 
			os.Interrupt,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)

	go func() {
		<-osSignal
		Info.Println("Server is stopping")
		os.Exit(1)
	}()

	f, err := os.OpenFile("server.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
	    log.Fatalf("main() error opening file: %v", err)
	}
	defer f.Close()

	Init(f, f, f, f)
	//Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	Info.Println("Server is starting")

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/css/", cssHandler)
	http.HandleFunc("/api/accounts", getAccountHandler)
	http.HandleFunc("/api/accounts/add", addAccountHandler)
	http.HandleFunc("/api/accounts/transfer", transferHandler)
	http.ListenAndServe(":8090", nil)
}
