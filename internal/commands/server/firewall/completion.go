package serverfirewall

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// completeCommonPorts returns completion for well-known port numbers from /etc/services
func completeCommonPorts(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	ports := parseServicesFile()

	var suggestions []string
	for _, port := range ports {
		if strings.HasPrefix(port, toComplete) {
			suggestions = append(suggestions, port)
		}
		if len(suggestions) >= 50 { // Limit suggestions
			break
		}
	}

	return suggestions, cobra.ShellCompDirectiveNoFileComp
}

// parseServicesFile parses /etc/services to extract common port numbers
func parseServicesFile() []string {
	var ports []string
	seen := make(map[string]bool)

	// Add commonly used ports first
	commonPorts := []struct {
		port string
		desc string
	}{
		{"22", "SSH"},
		{"80", "HTTP"},
		{"443", "HTTPS"},
		{"21", "FTP"},
		{"25", "SMTP"},
		{"53", "DNS"},
		{"110", "POP3"},
		{"143", "IMAP"},
		{"3306", "MySQL"},
		{"5432", "PostgreSQL"},
		{"6379", "Redis"},
		{"8080", "HTTP-Alt"},
		{"8443", "HTTPS-Alt"},
		{"3000", "Dev-Server"},
		{"5000", "Dev-Server"},
		{"8000", "Dev-Server"},
	}

	for _, cp := range commonPorts {
		entry := cp.port + "\t" + cp.desc
		ports = append(ports, entry)
		seen[cp.port] = true
	}

	// Try to read /etc/services for additional ports
	file, err := os.Open("/etc/services")
	if err != nil {
		return ports // Return common ports if file can't be read
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse line: service-name port/protocol [aliases] [# comment]
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		serviceName := fields[0]
		portProto := fields[1]

		// Extract port number (before the /)
		parts := strings.Split(portProto, "/")
		if len(parts) < 2 {
			continue
		}

		port := parts[0]
		protocol := parts[1]

		// Only include TCP and UDP ports
		if protocol != "tcp" && protocol != "udp" {
			continue
		}

		// Skip if we've already added this port
		if seen[port] {
			continue
		}

		// Validate port number
		portNum, err := strconv.Atoi(port)
		if err != nil || portNum < 1 || portNum > 65535 {
			continue
		}

		seen[port] = true
		entry := port + "\t" + serviceName + "/" + protocol
		ports = append(ports, entry)

		if len(ports) >= 200 { // Limit total entries
			break
		}
	}

	return ports
}

// completeIPAddress returns completion suggestions for IP addresses
func completeIPAddress(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var suggestions []string

	// Common private network ranges
	commonRanges := []struct {
		addr string
		desc string
	}{
		{"10.0.0.0/8", "Private-ClassA"},
		{"172.16.0.0/12", "Private-ClassB"},
		{"192.168.0.0/16", "Private-ClassC"},
		{"192.168.1.0/24", "Private-Common"},
		{"0.0.0.0/0", "Any-IPv4"},
		{"::/0", "Any-IPv6"},
		{"127.0.0.1", "Localhost-IPv4"},
		{"::1", "Localhost-IPv6"},
	}

	for _, range_ := range commonRanges {
		if strings.HasPrefix(range_.addr, toComplete) {
			entry := range_.addr + "\t" + range_.desc
			suggestions = append(suggestions, entry)
		}
	}

	// If user is typing a partial IP, suggest completing octets
	if len(toComplete) > 0 && (strings.Contains(toComplete, ".") || strings.Contains(toComplete, ":")) {
		// Don't auto-complete partial IPs - let user type them
		// Just provide the common ranges
	}

	return suggestions, cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveNoSpace
}

// completeSkipConfirmation returns completion suggestions for skip-confirmation flag
func completeSkipConfirmation(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	suggestions := []string{
		"0\tAlways require confirmation",
		"10\tSkip confirmation for up to 10 rules",
	}

	var filtered []string
	for _, s := range suggestions {
		if strings.HasPrefix(s, toComplete) {
			filtered = append(filtered, s)
		}
	}

	return filtered, cobra.ShellCompDirectiveNoFileComp
}
