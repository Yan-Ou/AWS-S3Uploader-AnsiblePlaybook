package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func webHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("view.html")
	t.Execute(w, "Hello World")
}

func upload(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		// GET
		t, _ := template.ParseFiles("view.html")

		t.Execute(w, nil)

	} else if r.Method == "POST" {
		// Post
		infile, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer infile.Close()

		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()

		io.Copy(f, infile)

		accessKey := r.FormValue("accessKey")
		secretKey := r.FormValue("secretKey")
		token := ""
		creds := credentials.NewStaticCredentials(accessKey, secretKey, token)
		cfg := aws.NewConfig().WithRegion("ap-southeast-2").WithCredentials(creds)

		svc := s3.New(session.New(), cfg)
		file, _ := os.Open(handler.Filename)
		defer file.Close()

		fileInfo, _ := file.Stat()
		size := fileInfo.Size()
		buffer := make([]byte, size)
		file.Read(buffer)
		fileBytes := bytes.NewReader(buffer)
		fileType := http.DetectContentType(buffer)
		path := "/media/" + file.Name()

		params := &s3.PutObjectInput{
			Bucket:        aws.String("wukong-bucket2018"),
			Key:           aws.String(path),
			Body:          fileBytes,
			ContentLength: aws.Int64(size),
			ContentType:   aws.String(fileType),
		}

		resp, _ := svc.PutObject(params)
		fmt.Printf("response %s", awsutil.StringValue(resp))

	} else {
		fmt.Println("Unknown HTTP " + r.Method + "  Method")
	}
}

func main() {
	// server := http.Server{
	// 	Addr: "127.0.0.1:8001",
	// }

	// http.HandleFunc("/", webHandler)
	// http.HandleFunc("/upload", upload)
	// server.ListenAndServe()
	r1 := http.NewServeMux()
	r1.HandleFunc("/", webHandler)

	r2 := http.NewServeMux()
	r2.HandleFunc("/upload", upload)
	// log.Fatal(http.ListenAndServe(":8001", r1))
	// log.Fatal(http.ListenAndServe(":8002", r2))
	go func() { log.Fatal(http.ListenAndServe(":8001", r1)) }()
	go func() { log.Fatal(http.ListenAndServe(":8002", r2)) }()
	select {}
}
