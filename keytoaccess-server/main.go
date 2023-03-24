package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
)

var passwordHashes = []string{
	"1311f8fc80a7ea28d78dd7723f09c44c1754cd35160ca8e7133ae3d7f636a19a",
}

func main() {
	http.HandleFunc("/", handleAccess)

	err := http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil)
	if err != nil {
		panic(err)
	}
}

func handleAccess(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		//r.ParseForm()
		// Get the password from the form
		password := r.FormValue("key")
		// log.Println("unhashed password '", password, "'")

		hash := computeHash(password)

		// Compare the password hash against the known hash
		if isValidHash(hash) {
			// Get the IP address of the client
			remoteAddr, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Printf("Error getting IP address: %v", err)
				return
			}

			// Check if the IP address is valid
			if net.ParseIP(remoteAddr) == nil {
				http.Error(w, "Invalid IP address", http.StatusBadRequest)
				log.Printf("Invalid IP address: %s", remoteAddr)
				return
			}

			fmt.Println(remoteAddr)

			reader := bufio.NewReader(os.Stdin)
			text, err := reader.ReadString('\n')
			if err != nil {
				log.Println("Error getting backend response")
				return
			}

			if text != "OK\n" {
				http.Error(w, "Internal Server Error", http.StatusBadRequest)
				log.Printf("Error adding rule, backed said: '%s'", text)
				return
			}

			// // Add a firewall rule to allow the client's IP address for 24 hours on ports 80 and 443
			// cmd := exec.Command("iptables", "-I", "INPUT", "-p", "tcp", "-s", remoteAddr, "--dport", "80", "-m", "comment", "--comment", "keytoaccess-"+hash, "-j", "ACCEPT")
			// err = cmd.Run()
			// if err != nil {
			// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			// 	log.Printf("Error adding firewall rule: %v", err)
			// 	return
			// }

			// cmd = exec.Command("iptables", "-I", "INPUT", "-p", "tcp", "-s", remoteAddr, "--dport", "443", "-m", "comment", "--comment", "keytoaccess", "-j", "ACCEPT")
			// err = cmd.Run()
			// if err != nil {
			// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			// 	log.Printf("Error adding firewall rule: %v", err)
			// 	return
			// }

			// Send a response to the client
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))

			// Log the successful access
			log.Printf("Successful access from IP address: %s and hash: %s", remoteAddr, hash)

			// Log the successful access to file
			logToFile("Successful access from IP address: " + remoteAddr + " and hash: " + hash + "\n")
			return
		}

		// If the password is incorrect, send an error response
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		// log.Println("Expected hash ", passwordHash, " received hash ", passwordHashed)
		return
	}

	// Render the login form template
	t, err := template.ParseFiles("login.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error parsing template: %v", err)
		return
	}
	t.Execute(w, nil)
}

func isValidHash(userHash string) bool {
	for _, hash := range passwordHashes {
		if userHash == hash {
			return true
		}
	}
	return false
}

func computeHash(key string) string {
	hasher := sha256.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func logToFile(msg string) {
	file, err := os.OpenFile("keytoaccess.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening log file: %v", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(msg); err != nil {
		log.Printf("Error writing to log file: %v", err)
		return
	}
}

func isValidIP(ip string) bool {
	// Check if the IP address is valid using regular expression
	ipRegexp := regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	return ipRegexp.MatchString(ip)
}
