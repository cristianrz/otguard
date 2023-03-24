package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	// Set up the command to run keytoaccess-server as the "keytoaccess" user.
	cmd := exec.Command("./keytoaccess-server")
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: 1001, Gid: 1001}

	// Start the keytoaccess-server process.
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating stdout pipe:", err)
		os.Exit(1)
	}

	// defer stdout.Close()

	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating stdin pipe:", err)
		os.Exit(1)
	}

	// defer stdin.Close()

	if err := cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "Error starting keytoaccess-server:", err)
		os.Exit(1)
	}

	// Read the server's stdout and add firewall rules.
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		ip := scanner.Text()
		if isValidIP(ip) {
			err := addFirewallRules(ip)
			if err == nil {
				// ip is valid
				_, err = stdin.Write([]byte("OK\n"))
				if err != nil {
					fmt.Println("Error writing to stdin:", err)
					return
				}
			} else {
				// adding firewall rules failed
				_, err = stdin.Write([]byte("OK\n"))
				if err != nil {
					fmt.Println("Error writing to stdin:", err)
					return
				}
			}
		} else {
			// ip is not valid
			_, err = stdin.Write([]byte("FAIL\n"))
			if err != nil {
				fmt.Println("Error writing to stdin:", err)
				return
			}

		}

	}

	// Check for errors during scanning.
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error scanning server stdout:", err)
		os.Exit(1)
	}

	// Wait for the keytoaccess-server process to exit.
	if err := cmd.Wait(); err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for keytoaccess-server:", err)
		os.Exit(1)
	}
}

// addFirewallRules adds an iptables rule to allow ports 80 and 443 for the given IP.
func addFirewallRules(ip string) error {

	cmd := exec.Command("iptables", "-C", "INPUT", "-s", ip, "-p", "tcp", "--dport", "80", "-j", "ACCEPT", "-m", "comment", "--comment", "keytoaccess")

	if err := cmd.Run(); err == nil {
		fmt.Printf("Firewall rule already exists for ip: %s\n", ip)
		return fmt.Errorf("Firewall rule already exists for ip: %s", ip)
	}

	// Run the iptables command to add the rules.
	cmd = exec.Command("iptables", "-A", "INPUT", "-s", ip, "-p", "tcp", "--dport", "80", "-j", "ACCEPT", "-m", "comment", "--comment", "keytoaccess")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add port 80 rule: %v", err)
	}
	cmd = exec.Command("iptables", "-A", "INPUT", "-s", ip, "-p", "tcp", "--dport", "443", "-j", "ACCEPT", "-m", "comment", "--comment", "keytoaccess")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add port 443 rule: %v", err)
	}

	fmt.Println("Added firewall rule for ip: ", ip)
	return nil
}

// isValidIP checks if the given string is a valid IP address.
func isValidIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && !parsedIP.IsUnspecified()
}
