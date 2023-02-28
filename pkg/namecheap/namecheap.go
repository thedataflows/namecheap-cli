package namecheap

import "encoding/xml"

// ApiResponse was generated 2023-02-16 09:03:21 by https://xml-to-go.github.io/ in Ukraine.
type ApiResponse struct {
	XMLName xml.Name `xml:"ApiResponse"`
	// Text             string   `xml:",chardata"`
	Status string `xml:"Status,attr"`
	Xmlns  string `xml:"xmlns,attr"`
	Errors struct {
		// Text  string `xml:",chardata"`
		Error []struct {
			Text   string `xml:",chardata"`
			Number string `xml:"Number,attr"`
		} `xml:"Error"`
	} `xml:"Errors"`
	Warnings struct {
		// Text  string `xml:",chardata"`
		Warning []struct {
			Text   string `xml:",chardata"`
			Number string `xml:"Number,attr"`
		} `xml:"Warning"`
	} `xml:"Warnings"`
	RequestedCommand string `xml:"RequestedCommand"`
	CommandResponse  struct {
		// Text                    string `xml:",chardata"`
		Type                    string `xml:"Type,attr"`
		DomainDNSGetHostsResult struct {
			// Text          string `xml:",chardata"`
			Domain        string `xml:"Domain,attr"`
			EmailType     string `xml:"EmailType,attr"`
			IsUsingOurDNS string `xml:"IsUsingOurDNS,attr"`
			Host          []Host `xml:"host"`
		} `xml:"DomainDNSGetHostsResult"`
	} `xml:"CommandResponse"`
	Server            string `xml:"Server"`
	GMTTimeDifference string `xml:"GMTTimeDifference"`
	ExecutionTime     string `xml:"ExecutionTime"`
}

type Host struct {
	// Text               string `xml:",chardata"`
	HostId             string `xml:"HostId,attr"`
	Name               string `xml:"Name,attr"`
	Type               string `xml:"Type,attr"`
	Address            string `xml:"Address,attr"`
	MXPref             string `xml:"MXPref,attr"`
	TTL                string `xml:"TTL,attr"`
	AssociatedAppTitle string `xml:"AssociatedAppTitle,attr"`
	FriendlyName       string `xml:"FriendlyName,attr"`
	IsActive           string `xml:"IsActive,attr"`
	IsDDNSEnabled      string `xml:"IsDDNSEnabled,attr"`
}
