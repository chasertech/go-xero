package main

import (
	"crypto/x509"
	"encoding/pem"
	"encoding/xml"
	"flag"
	"github.com/garyburd/go-oauth/oauth"
	"io/ioutil"
	"net/url"
	"os"
)

func main() {
	Type := flag.String(
		"type",
		"ACCREC",
		"Invoice Type. See https://developer.xero.com/documentation/api/types/#InvoiceTypes.")
	ContactID := flag.String(
		"contact",
		"dd1dcb58-a767-4d49-b77c-eda94fbc2cf5",
		"Contact ID. See https://developer.xero.com/documentation/api/contacts/#title1.")
	Description := flag.String(
		"description",
		"Monthly rental for property at 56a Wilkins Avenue",
		"Line Item Description. See https://developer.xero.com/documentation/api/invoices/#title2.")
	LineAmount := flag.Float64(
		"amount",
		395.00,
		"Line Item Amount. See https://developer.xero.com/documentation/api/invoices/#title2.")
	PrivateKeyPath := flag.String(
		"private-key-path",
		"",
		"Path to X509 private key (REQUIRED)")
	ConsumerKey := flag.String(
		"consumer-key",
		"",
		"Consumer key (REQUIRED)")
	flag.Parse()
	if *PrivateKeyPath == "" || *ConsumerKey == "" {
		flag.PrintDefaults()
		os.Exit(2)
	}

	type LineItem struct {
		Description string
		LineAmount  float64
	}
	type Invoice struct {
		Type      string
		ContactID string     `xml:"Contact>ContactID"`
		LineItems []LineItem `xml:"LineItems>LineItem"`
	}
	type Invoices struct {
		Invoices Invoice `xml:"Invoice"`
	}
	v := &Invoices{Invoice{*Type, *ContactID, []LineItem{{*Description, *LineAmount}}}}
	output, _ := xml.Marshal(v)

	pemData, _ := ioutil.ReadFile(*PrivateKeyPath)
	block, _ := pem.Decode(pemData)
	privateKey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	client := oauth.Client{
		Credentials:     oauth.Credentials{Token: *ConsumerKey},
		SignatureMethod: oauth.RSASHA1,
		PrivateKey:      privateKey,
	}
	res, _ := client.Post(nil, &client.Credentials, "https://api.xero.com/api.xro/2.0/invoices", url.Values{"xml": {string(output)}})
	body, _ := ioutil.ReadAll(res.Body)
	os.Stdout.Write(body)
}
