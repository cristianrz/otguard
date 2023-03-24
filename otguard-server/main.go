package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/pquerna/otp/totp"
)

var secrets = map[string]string{}
var successMsg *string
var failMsg *string

func main() {
	portPtr := flag.Int("p", 8443, "listen port")
	successMsg = flag.String("s", "You are now logged in", "message on auth success")
	failMsg = flag.String("f", "Incorrect username or OTP", "message on auth failure")

	flag.Parse()

	this, err := os.Executable()
	if err != nil {
		log.Fatalln(err)
	}

	here := filepath.Dir(this)

	secrets, err = readSecretsFile(here + "/../etc/otguard/secrets")
	if err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc("/", handleAccess)

	err = http.ListenAndServeTLS(":"+strconv.Itoa(*portPtr), here+"/../etc/otguard/cert.pem", here+"/../etc/otguard/key.pem", nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func handleAccess(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("login.html")
	if err != nil {
		log.Println("failed parsing template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data := struct {
		Body string
	}{
		Body: `<label for="username">Username:</label>
		<input type="text" id="username" name="username" required>
		<label for="password">OTP:</label>
		<input type="password" id="key" name="key" required>
		<input type="submit" value="Login">`,
	}

	switch r.Method {
	case "GET":
		err := tmpl.Execute(w, data)
		if err != nil {
			log.Println("failed to execute template")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	case "POST":
		// Get the otp from the form
		otp := r.FormValue("key")
		username := r.FormValue("username")

		valid := false

		userSecret, ok := secrets[username]
		if ok {
			valid = totp.Validate(otp, userSecret)
		}

		if !valid {
			log.Printf("Access denied for user '%s' and OTP '%s'\n", username, otp)
			// If the password is incorrect, send an error response
			data.Body = "<p class=\"error\">" + *failMsg + "</p>" + data.Body
			err := tmpl.Execute(w, data)
			if err != nil {
				log.Println("failed to execute template")
			}
			return
		}

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

		if text == "EXISTS\n" {
			data.Body = "<p>You were already logged in </p>"
			err := tmpl.Execute(w, data)
			if err != nil {
				log.Println("failed to execute template")
			}
			return
		}

		if text != "OK\n" {

			http.Error(w, "Internal Server Error", http.StatusBadRequest)
			log.Printf("Error adding rule, backed said: '%s'", text)
			return
		}

		// Send a response to the client
		data.Body = "<p>" + *successMsg + "</p>"
		err = tmpl.Execute(w, data)
		if err != nil {
			log.Println("failed to execute template")
		}

		// Log the successful access
		log.Printf("Successful access from IP address: %s and hash: %s", remoteAddr, otp)

	default:
		w.WriteHeader(404)
	}

}

func readSecretsFile(filename string) (map[string]string, error) {
	secrets := make(map[string]string)

	file, err := os.Open(filename)
	if err != nil {
		return secrets, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			return secrets, fmt.Errorf("invalid line in secrets file: %s", line)
		}
		username := strings.TrimSpace(parts[0])
		secret := strings.TrimSpace(parts[1])
		secrets[username] = secret
	}

	if err := scanner.Err(); err != nil {
		return secrets, err
	}

	return secrets, nil
}

// func acceptedHash(userHash string) bool {
// 	for _, hash := range passwordHashes {
// 		if userHash == hash {
// 			return true
// 		}
// 	}
// 	return false
// }

// func acceptedOTP(userOTP string) bool {
// 	for _, secret := range secrets {
// 		refOTP, err := generateOTP(secret)
// 		if err != nil {
// 			return false
// 		}

// 		if userOTP != refOTP {
// 			return false
// 		}
// 	}
// 	return true
// }

// func computeHash(key string) string {
// 	hasher := sha256.New()
// 	hasher.Write([]byte(key))
// 	return hex.EncodeToString(hasher.Sum(nil))
// }

// func generateOTP(key string) (string, error) {
// 	// Generate a TOTP code using the key and current time.
// 	now := time.Now()
// 	totpCode, err := totp.GenerateCode(key, now)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to generate TOTP code: %v", err)
// 	}

// 	return totpCode, nil
// }