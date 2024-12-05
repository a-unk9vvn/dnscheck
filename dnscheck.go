package main

import (
 "fmt"
 "net"
 "os"
 "strings"
)

// Function to lookup and return detailed DNS records
func lookupDNS(domain string) string {
 var builder strings.Builder
 builder.WriteString(fmt.Sprintf("=== DNS Lookup Results for Domain: %s ===\n\n", domain))

 // A Records
 builder.WriteString("[+] A Records:\n")
 ips, err := net.LookupIP(domain)
 if err != nil {
  builder.WriteString(fmt.Sprintf("  [!] Error fetching A records: %v\n", err))
 } else {
  for _, ip := range ips {
   if ip.To4() != nil {
    builder.WriteString(fmt.Sprintf("  - IP Address: %s\n", ip))
   }
  }
 }

 // AAAA Records
 builder.WriteString("\n[+] AAAA Records:\n")
 if err == nil {
  for _, ip := range ips {
   if ip.To16() != nil && ip.To4() == nil {
    builder.WriteString(fmt.Sprintf("  - IPv6 Address: %s\n", ip))
   }
  }
 } else {
  builder.WriteString(fmt.Sprintf("  [!] Error fetching AAAA records: %v\n", err))
 }

 // CNAME Records
 builder.WriteString("\n[+] CNAME Records:\n")
 cname, err := net.LookupCNAME(domain)
 if err != nil {
  builder.WriteString(fmt.Sprintf("  [!] Error fetching CNAME records: %v\n", err))
 } else {
  builder.WriteString(fmt.Sprintf("  - Canonical Name: %s\n", cname))
 }

 // MX Records
 builder.WriteString("\n[+] MX Records:\n")
 mxRecords, err := net.LookupMX(domain)
 if err != nil {
  builder.WriteString(fmt.Sprintf("  [!] Error fetching MX records: %v\n", err))
 } else {
  for _, mx := range mxRecords {
   builder.WriteString(fmt.Sprintf("  - Mail Server: %s, Priority: %d\n", mx.Host, mx.Pref))
  }
 }

 // NS Records
 builder.WriteString("\n[+] NS Records:\n")
 nsRecords, err := net.LookupNS(domain)
 if err != nil {
  builder.WriteString(fmt.Sprintf("  [!] Error fetching NS records: %v\n", err))
 } else {
  for _, ns := range nsRecords {
   builder.WriteString(fmt.Sprintf("  - Name Server: %s\n", ns.Host))
  }
 }

 // TXT Records
 builder.WriteString("\n[+] TXT Records:\n")
 txtRecords, err := net.LookupTXT(domain)
 if err != nil {
  builder.WriteString(fmt.Sprintf("  [!] Error fetching TXT records: %v\n", err))
 } else {
  for _, txt := range txtRecords {
   builder.WriteString(fmt.Sprintf("  - %s\n", txt))
  }
 }

 builder.WriteString("\n=== End of Results ===\n")
 return builder.String()
}

func main() {
 // Check for sufficient arguments
 if len(os.Args) < 3 {
  fmt.Println("Usage: go run main.go -d <domain> [-o <output file>]")
  os.Exit(1)
 }

 var domain, outputFile string
 for i := 1; i < len(os.Args); i++ {
  switch os.Args[i] {
  case "-d":
   if i+1 < len(os.Args) {
    domain = os.Args[i+1]
    i++
   } else {
    fmt.Println("Error: No domain provided after -d")
    os.Exit(1)
   }
  case "-o":
   if i+1 < len(os.Args) {
    outputFile = os.Args[i+1]
    i++
   } else {
    fmt.Println("Error: No output file provided after -o")
    os.Exit(1)
   }
  default:
   fmt.Printf("Unknown option: %s\n", os.Args[i])
   os.Exit(1)
  }
 }

 // Perform DNS lookup
 if domain == "" {
  fmt.Println("Error: Domain not provided")
  os.Exit(1)
 }

 results := lookupDNS(domain)

 // Output results
 if outputFile != "" {
  err := os.WriteFile(outputFile, []byte(results), 0644)
  if err != nil {
   fmt.Printf("Error writing to file: %v\n", err)
   os.Exit(1)
  }
  fmt.Printf("Results saved to %s\n", outputFile)
 } else {
  fmt.Println(results)
 }
}
