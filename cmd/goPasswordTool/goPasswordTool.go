package main

import (
	"flag"
	"fmt"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/crypto"
	"io"
	"os"
)

const (
	APP     = "goPasswordTool"
	VERSION = "0.0.1"
)

func main() {
	os.Exit(realMain(os.Stdout))
}
func realMain(out io.Writer) int {
	fmt.Fprintf(out, "INFO: 'Starting %s v:%s'", APP, VERSION)

	if len(os.Args) < 2 {
		fmt.Fprintf(out, "💥💥 ERROR: 'missing parameter %s -password your_pass'\n", os.Args[0])
		return 1
	}

	password := flag.String("password", "", "Password in plain text to get the hash and salt version")
	flag.Parse()

	resSha256 := crypto.Sha256Hash(*password)
	resSha256Hashed, err := crypto.HashAndSalt(resSha256)
	if err != nil {
		fmt.Fprintf(out, "💥💥 ERROR: 'calling crypto.HashAndSalt(resSha256) got error: %v'\n", err)
		return 1
	}
	resHashed, err := crypto.HashAndSalt(*password)
	if err != nil {
		fmt.Fprintf(out, "💥💥 ERROR: 'calling crypto.HashAndSalt(*password) got error: %v'\n", err)
		return 1
	}

	fmt.Fprintf(out, "RESULT: 'password:  %10q\t sha256: %s'\n", *password, resSha256)
	fmt.Fprintf(out, "RESULT: 'password:  %10q\t hashSaltFromSha256: \t%s'\n", *password, resSha256Hashed)
	fmt.Fprintf(out, "RESULT: 'password:  %10q\t hashSaltFromPassword:\t%s'\n", *password, resHashed)

	resComparisonSha256 := crypto.ComparePasswords(resSha256Hashed, resSha256)
	fmt.Fprintf(out, "RESULT: 'crypto.ComparePasswords(resSha256Hashed, resSha256)  is %v'\n", resComparisonSha256)

	resComparison := crypto.ComparePasswords(resHashed, *password)
	fmt.Fprintf(out, "RESULT: 'crypto.ComparePasswords(resHashed, *password)  is %v'\n", resComparison)
	return 0
}
