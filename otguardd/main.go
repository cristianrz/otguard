package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/user"
	"strconv"

	"path/filepath"
	"syscall"
)

var ports = []int{80, 443}

func main() {
	err := prepareLogFile()
	if err != nil {
		panic(err)
	}

	uidPtr := flag.Int("u", 0, "uid for otguard-web process")
	gidPtr := flag.Int("g", 0, "gid for otguard-web process")
	portPtr := flag.Int("p", 8443, "listen port")
	successMsg := flag.String("s", "You are now logged in", "message on auth success")
	failMsg := flag.String("f", "Incorrect username or OTP", "message on auth failure")

	flag.Parse()

	this, err := os.Executable()
	if err != nil {
		panic(err)
	}

	here := filepath.Dir(this)

	if *uidPtr == 0 || *gidPtr == 0 {
		log.Fatalln("Refusing to run with uid or guid 0")
	}

	// Set up the command to run otguard-web as the "otguard" user.
	cmd := exec.Command(here+"/otguard-web", "-p", strconv.Itoa(*portPtr), "-s", *successMsg, "-f", *failMsg)
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(*uidPtr), Gid: uint32(*gidPtr)}

	// Start the otguard-web process.
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating stdout pipe:", err)
		os.Exit(1)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating stdin pipe:", err)
		os.Exit(1)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalln("Error starting otguard-web: ", err)
	}

	defer syscall.Kill(cmd.Process.Pid, syscall.SIGHUP)

	log.Println("Started otguard-web")

	// Read the server's stdout and add firewall rules.
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		ip := scanner.Text()
		if !isValidIP(ip) {
			log.Println("Refusing to use invalid IP: ", ip)

			_, err = stdin.Write([]byte("FAIL\n"))
			if err != nil {
				log.Fatalln("Error writing to stdin:", err)
				return
			}
		}

		log.Println("Received valid ip: ", ip)

		exists, err := addFirewallRules(ip)
		if err != nil {
			log.Printf("Adding firewall rules failed for %s: %v", ip, err)

			if exists {
				_, err = stdin.Write([]byte("EXISTS\n"))
				if err != nil {
					log.Fatalln("Error writing to stdin:", err)
					return
				}

				continue
			}

			_, err = stdin.Write([]byte("FAIL\n"))
			if err != nil {
				log.Fatalln("Error writing to stdin:", err)
				return
			}
		}

		_, err = stdin.Write([]byte("OK\n"))
		if err != nil {
			log.Fatalln("Error writing to stdin:", err)
		}

		log.Println("Added rules for ip: ", ip)
	}

	// Check for errors during scanning.
	if err := scanner.Err(); err != nil {
		log.Fatalln("Error scanning server stdout:", err)
	}

	// Wait for the otguard-web process to exit.
	if err := cmd.Wait(); err != nil {
		log.Fatalln("Error waiting for otguard-web:", err)
	}
}

// addFirewallRules adds an iptables rule to callow ports 80 and 443 for the given IP.
func addFirewallRules(ip string) (bool, error) {

	cmd := exec.Command("iptables", "-C", "INPUT", "-s", ip, "-p", "tcp", "--dport", "80", "-j", "ACCEPT", "-m", "comment", "--comment", "otguard")

	err := cmd.Run()
	if err == nil {
		return true, fmt.Errorf("firewall rule already exists for ip: %s", ip)
	}

	for _, port := range ports {
		// Run the iptables command to add the rules.
		cmd = exec.Command("iptables", "-I", "INPUT", "-s", ip, "-p", "tcp", "--dport", strconv.Itoa(port), "-j", "ACCEPT", "-m", "comment", "--comment", "otguard")
		if err := cmd.Run(); err != nil {
			return false, fmt.Errorf("failed to add port %v rule: %v", port, err)
		}
	}

	return false, nil
}

// isValidIP checks if the given string is a valid IP address.
func isValidIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && !parsedIP.IsUnspecified()
}

func prepareLogFile() error {
	// Get the otguard user
	otguardUser, err := user.Lookup("otguard")
	if err != nil {
		return err
	}

	// Create the log file
	logFile, err := os.Create("/var/log/otguard-web.log")
	if err != nil {
		return err
	}
	defer logFile.Close()

	// Set the file permissions and owner
	err = os.Chmod("/var/log/otguard-web.log", 0644)
	if err != nil {
		return err
	}

	uid, err := strconv.Atoi(otguardUser.Uid)
	if err != nil {
		return err
	}
	gid, err := strconv.Atoi(otguardUser.Gid)
	if err != nil {
		return err
	}

	err = os.Chown("/var/log/otguard-web.log", uid, gid)
	if err != nil {
		return err
	}

	return nil
}
