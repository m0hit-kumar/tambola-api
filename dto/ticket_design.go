package dto

type TicketDesign struct {
	HostName   string `json:"hostName" `
	Background string `json:"background" `
	Border     string `json:"border"  `
	Text       string `json:"text"`
}

func TicketDesignRes() TicketDesign {
	return TicketDesign{
		HostName:   "Tambola",
		Background: "#ffffff",
		Border:     "#000000",
		Text:       "#000000",
	}
}
