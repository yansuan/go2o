// Autogenerated by Thrift Compiler (0.9.3)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"go2o/core/service/thrift/idl/gen-go/define"
	"os"
)

func Usage() {
	fmt.Fprintln(os.Stderr, "Usage of ", os.Args[0], " [-h host:port] [-u url] [-f[ramed]] function [arg1 [arg2...]]:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "\nFunctions:")
	fmt.Fprintln(os.Stderr, "   Login(string user, string pwd, bool update)")
	fmt.Fprintln(os.Stderr)
	os.Exit(0)
}

func main() {

}