package main

import (
	"context"
	"encoding/json"
	"go.elastic.co/apm/module/apmhttp/v2"
	"go.elastic.co/apm/v2"
	"net/http"
	"os"
	"time"
)

const (
	apmServer = "https://apm-opty-kalbe.apps.lab.i-3.my.id"
	apmName = "tester-infra"
	apmSecret = "OfX0K5s0i8M37XQ49cnYh77W"
	verifyServerCert = "false"
)


func main() {
	os.Setenv("ELASTIC_APM_SERVICE_NAME", apmName)
	os.Setenv("ELASTIC_APM_SERVER_URL", apmServer)
	os.Setenv("ELASTIC_APM_SECRET_TOKEN", apmSecret)
	mux := http.NewServeMux()
	mux.HandleFunc("/", baseHandler)
	http.ListenAndServe(":8081", apmhttp.Wrap(mux))
}

func baseHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span, ctx := apm.StartSpan(ctx, "baseHandler", "custom")
	defer span.End()
	processingRequest(ctx)
	todo, err := getTodoFromAPI(ctx)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	data, _ := json.Marshal(todo)
	w.Write(data)
}

func processingRequest(ctx context.Context) {
	span, ctx := apm.StartSpan(ctx, "processingRequest", "custom")
	defer span.End()
	doSomething(ctx)
	// time sleep simulate some processing time
	time.Sleep(15 * time.Millisecond)
	return
}

func doSomething(ctx context.Context) {
	span, ctx := apm.StartSpan(ctx, "doSomething", "custom")
	defer span.End()
	// time sleep simulate some processing time
	time.Sleep(20 * time.Millisecond)
	return
}

func getTodoFromAPI(ctx context.Context) (map[string]interface{}, error) {
	span, ctx := apm.StartSpan(ctx, "getTodoFromAPI", "custom")
	defer span.End()
	var result map[string]interface{}
	resp, err := http.Get("https://jsonplaceholder.typicode.com/todos/1")
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, err
}
