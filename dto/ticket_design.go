package dto

type TicketDesign struct {
	HostName   string `json:"hostName" `
	Background string `json:"background" `
	Border     string `json:"border"  `
	Text       string `json:"text"`
	RoomId     string `json:"roomId" `
}
