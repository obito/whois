package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {
	domains := os.Args[1:]

	for _, domain := range domains {
		ext := extension(domain)

		log.Print(fmt.Sprintf("# whois.nic.%s", ext))
		result, err := query(fmt.Sprintf("whois.nic.%s", ext), domain)
		if err != nil {
			log.Fatal(err)
		}

		log.Print(result)

		whoisServer := server(result)

		log.Print("# " + whoisServer)

		result, err = query(whoisServer, domain)
		if err != nil {
			return
		}

		log.Print(result)
	}
}

func query(server, domain string) (string, error) {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(server, "43"), time.Second*30)
	if err != nil {
		return "", fmt.Errorf("connetion failed: %v", err)
	}

	defer conn.Close()

	conn.SetWriteDeadline(time.Now().Add(time.Second * 30))
	_, err = conn.Write([]byte(domain + "\r\n"))
	if err != nil {
		return "", fmt.Errorf("sending failed: %v", err)
	}

	conn.SetReadDeadline(time.Now().Add(time.Second * 30))
	buff, err := ioutil.ReadAll(conn)
	if err != nil {
		return "", fmt.Errorf("read failed: %v", err)
	}

	return string(buff), nil
}

func server(result string) string {
	var reg = regexp.MustCompile(`Registrar WHOIS Server: (.*)`)

	return strings.SplitN(reg.FindString(result), ": ", 2)[1]
}

func extension(domain string) string {
	var ext string

	if net.ParseIP(domain) == nil {
		extSplited := strings.Split(domain, ".")
		ext = extSplited[len(extSplited)-1]
	} else {
		ext = domain
	}

	return ext
}
