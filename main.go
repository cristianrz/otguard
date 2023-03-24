package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"time"
)

const (
	listenAddress = "0.0.0.0:8080"
	correctKey    = "mySecretKey"
	ruleComment   = "keytoaccess"
)

func main() {
	http.HandleFunc("/", handleRequest)
	err := http.ListenAndServe(listenAddress, nil)
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Request was not POST")

		showInputPage(w)
		return
	}

	fmt.Println("Received POST")

	key := r.FormValue("key")
	if key != correctKey {
		fmt.Println("Received key but was incorrect: ", key)

		showInputPage(w)
		return
	}

	fmt.Println("Received key and it was correct!")

	ip := r.RemoteAddr[:len(r.RemoteAddr)-6] // remove port number
	err := allowIP(ip)
	if err != nil {
		fmt.Fprintln(w, "Failed to create firewall rule")
		return
	}

	fmt.Fprintln(w, "OK")
}

func showInputPage(w http.ResponseWriter) {
	fmt.Fprintln(w, `
		<!DOCTYPE html>
		<html>
			<head>
				<title>Enter Key</title>
			</head>
			<body>
				<form method="post">
					<label for="key">Enter key:</label>
					<input type="password" name="key" id="key" required>
					<br>
					<button type="submit">Submit</button>
				</form>
			</body>
		</html>
	`)
}

func allowIP(ip string) error {
	cmd := exec.Command("iptables", "-A", "INPUT", "-s", ip, "-p", "tcp", "--dport", "80", "-j", "ACCEPT", "-m", "comment", "--comment", ruleComment)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run iptables command: %v", err)
	}

	// Remove the rule when the program exits
	// defer func() {
	// 	cmd := exec.Command("iptables", "-D", "INPUT", "-s", ip, "-j", "ACCEPT", "-m", "comment", "--comment", ruleComment)
	// 	cmd.Run()
	// }()

	// Remove the rule after 24 hours
	time.AfterFunc(24*time.Hour, func() {
		cmd := exec.Command("iptables", "-D", "INPUT", "-s", ip, "-p", "tcp", "--dport", "80", "-j", "ACCEPT", "-m", "comment", "--comment", ruleComment)
		cmd.Run()
	})

	return nil
}
