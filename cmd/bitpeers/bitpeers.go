package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"github.com/RaghavSood/bitpeers"
	flag "github.com/spf13/pflag"
	"os"
)

type BitPeersDB bitpeers.PeersDB

var peersFilePath string
var formatOption string
var addressOnly bool

func init() {
	flag.StringVar(&peersFilePath, "filepath", "", "the path to peers.dat")
	flag.StringVar(&formatOption, "format", "json", "the output format {json|text}")
	flag.BoolVar(&addressOnly, "addressonly", false, "outputs only addresses if specified")
	flag.Parse()
}

func main() {
	if peersFilePath == "" {
		fmt.Fprintf(os.Stderr, "Invalid peers file %s\n", peersFilePath)
		os.Exit(1)
	}

	if formatOption != "json" && formatOption != "text" {
		fmt.Fprintf(os.Stderr, "Invalid output format %s\n", formatOption)
		os.Exit(1)
	}

	rawPeersDB, err := bitpeers.NewPeersDB(peersFilePath)
	if err != nil {
		fmt.Println(err)
	}

	peersDb := BitPeersDB(rawPeersDB)

	if addressOnly {
		addressArray := make([]string, peersDb.NTried+peersDb.NNew)
		var i uint32
		for i = 0; i < peersDb.NNew; i++ {
			addressArray[i] = peersDb.NewAddrInfo[i].Address.PeerAddress.String()
		}
		for i = 0; i < peersDb.NTried; i++ {
			addressArray[peersDb.NNew+i] = peersDb.TriedAddrInfo[i].Address.PeerAddress.String()
		}

		if formatOption == "text" {
			for _, i := range addressArray {
				//fmt.Printf("%s\n", i)
				name := "IP.txt"
				content := fmt.Sprintf("%s\n", i)
				WriteWithFileWrite(name, content)
			}
			return
		} else {
			encodedPeers, err := json.Marshal(addressArray)
			if err != nil {
				fmt.Sprintf("Error converting to JSON: %s", err)
				os.Exit(1)
			}
			fmt.Println(string(encodedPeers))
			return
		}
	}

	if formatOption == "text" {
		//peersDb.dump()
		name := "out.txt"
		peersDb.write_to_txt(name)
	} else {
		encodedPeers, err := json.Marshal(peersDb)
		if err != nil {
			fmt.Sprintf("Error converting to JSON: %s", err)
			os.Exit(1)
		}
		fmt.Println(string(encodedPeers))
	}

}

func (peersDB BitPeersDB) write_to_txt(name string) {
	content := "bitpeers\n"
	content += "--------\n"
	content += fmt.Sprintf("Path: %s\n", peersDB.Path)
	content += fmt.Sprintf("MessageBytes: 0x%s\n", hexstring(peersDB.MessageBytes))
	content += fmt.Sprintf("Version: %d\n", peersDB.Version)
	content += fmt.Sprintf("KeySize: %d\n", peersDB.KeySize)
	content += fmt.Sprintf("NKey: %s\n", peersDB.NKey)
	content += fmt.Sprintf("NNew: %d\n", peersDB.NNew)
	content += fmt.Sprintf("NTried: %d\n", peersDB.NTried)
	content += fmt.Sprintf("NewBuckets: %d\n\n", peersDB.NewBuckets)
	WriteWithFileWrite(name, content)

	var i uint32
	for i = 0; i < peersDB.NNew; i++ {
		content = fmt.Sprintf("%s\n", peersDB.NewAddrInfo[i])
		WriteWithFileWrite(name, content)
	}

	fmt.Println("Tried Addresses:")
	for i = 0; i < peersDB.NTried; i++ {
		content = fmt.Sprintf("%s\n", peersDB.TriedAddrInfo[i])
		WriteWithFileWrite(name, content)
	}
}

//使用ioutil.WriteFile方式写入文件,是将[]byte内容写入文件,如果content字符串中没有换行符的话，默认就不会有换行符
func WriteWithIoutil(name,content string) {
    data :=  []byte(content)
    if ioutil.WriteFile(name,data,/*0644*/os.ModeAppend) == nil {
        fmt.Println("\n",content)
    }
}

//使用os.OpenFile()相关函数打开文件对象，并使用文件对象的相关方法进行文件写入操作
func WriteWithFileWrite(name,content string){
    fileObj,err := os.OpenFile(name,os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
    if err != nil {
        fmt.Println("Failed to open the file",err.Error())
        os.Exit(2)
    }
    defer fileObj.Close()
    //if _,err := fileObj.WriteString(content);err == nil {
    //    fmt.Println("\n",content)
    //}
    contents := []byte(content)
    if _,err := fileObj.Write(contents);err == nil {
        fmt.Println("\n",content)
    }
}

func (peersDB BitPeersDB) dump() {
	fmt.Println("bitpeers")
	fmt.Println("--------")
	fmt.Printf("Path: %s\n", peersDB.Path)
	fmt.Printf("MessageBytes: 0x%s\n", hexstring(peersDB.MessageBytes))
	fmt.Printf("Version: %d\n", peersDB.Version)
	fmt.Printf("KeySize: %d\n", peersDB.KeySize)
	fmt.Printf("NKey: %s\n", hexstring(peersDB.NKey))
	fmt.Printf("NNew: %d\n", peersDB.NNew)
	fmt.Printf("NTried: %d\n", peersDB.NTried)
	fmt.Printf("NewBuckets: %d\n", peersDB.NewBuckets)
	fmt.Println("")
	var i uint32
	for i = 0; i < peersDB.NNew; i++ {
		fmt.Print(peersDB.NewAddrInfo[i])
	}
	fmt.Println("Tried Addresses:")
	for i = 0; i < peersDB.NTried; i++ {
		fmt.Print(peersDB.TriedAddrInfo[i])
	}
}

func hexstring(input []byte) string {
	return hex.EncodeToString(input)
}
