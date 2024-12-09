package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"
)

type Ticket struct {
	Number            string `json:"number"`
	ShortDescription  string `json:"short_description"`
	OpenedAt          string `json:"opened_at"`
	ClosedAt          string `json:"closed_at"`
	BusinessDuration  string `json:"business_duration"`
	Description       string `json:"description"`
	CloseNotes        string `json:"close_notes"`
	State             string `json:"state"`
	ReassignmentCount string `json:"reassignment_count"`
	AssignmentGroup   string `json:"assignment_group"`
	UpdatedBy         string `json:"sys_updated_by"`
	ReportedBy        string `json:"u_reported_by"`
	WorkNotes         string `json:"work_notes"`
	ClosedBy          string `json:"closed_by"`
	ResolvedBy        string `json:"resolved_by"`
	KBNumber          string `json:"kb_number"`
}

// will use json as it's easier to parse
func ParseData(data [][]string) []Ticket {

	// for matching kb numbers
	re := regexp.MustCompile(`KB\d{7}|KB0`)

	var tickets []Ticket
	for i, line := range data {
		if i > 0 {
			var ticket Ticket
			for i, field := range line {
				switch i {
				case 0:
					ticket.Number = field
				case 1:
					ticket.OpenedAt = field
				case 2:
					ticket.ShortDescription = field
				case 3:
					ticket.Description = field
					// search worknotes for  kb number
					find := re.Find([]byte(field))
					if string(find) != "" {
						ticket.KBNumber = string(find)
					}
				case 7:
					// if type Closed -- IGNORE, CLOSED BY SYSTEM
					ticket.State = field
				case 10:
					ticket.AssignmentGroup = field
				case 13:
					ticket.UpdatedBy = field
				case 14:
					ticket.ReportedBy = field
				case 18:
					ticket.WorkNotes = field
					find := re.Find([]byte(field))
					if string(find) != "" {
						ticket.KBNumber = string(find)
					}
				case 21:
					ticket.CloseNotes = field
					find := re.Find([]byte(field))
					if string(find) != "" {
						ticket.KBNumber = string(find)
					}
				case 35:
					ticket.BusinessDuration = field
				case 44:
					ticket.ClosedAt = field
				case 45:
					ticket.ClosedBy = field
				case 111:
					ticket.ReassignmentCount = field
				case 116:
					ticket.ResolvedBy = field
				}
			}
			tickets = append(tickets, ticket)
		}
	}
	return tickets
}

func OutputTickets(output string, tickets []Ticket) {
	filename := fmt.Sprintf(output+"Solution Center Metrics-%s.csv", time.Now().Format("20060102"))

	f, err := os.Create(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	header := []string{
		"Number", "ShortDescription", "OpenedAt", "ClosedAt", "BusinessDuration",
		"Description", "CloseNotes", "State", "ReassignmentCount", "AssignmentGroup",
		"UpdatedBy", "ReportedBy", "WorkNotes", "ClosedBy", "ResolvedBy", "KBNumber",
	}
	writer.Write(header)

	for _, ticket := range tickets {
		var row []string
		row = append(row, ticket.Number, ticket.ShortDescription, ticket.OpenedAt, ticket.ClosedAt, ticket.BusinessDuration,
			ticket.Description, ticket.CloseNotes, ticket.State, ticket.ReassignmentCount, ticket.AssignmentGroup,
			ticket.UpdatedBy, ticket.ReportedBy, ticket.WorkNotes, ticket.ClosedBy, ticket.ResolvedBy, ticket.KBNumber)
		writer.Write(row)
	}

}

func main() {
	input := `C:\Users\sc_counter\Iowa State University\Solution Center - Documents\Metrics\Data\Solution Center Metrics.csv`
	output := "C:\\Users\\sc_counter\\Iowa State University\\Solution Center - Documents\\Metrics\\Sorted Data\\"

	// header
	fmt.Println("Solution Center Metrics Daemon")
	fmt.Println("Cameron Wilson @ sc_cameronw\n")

	fmt.Println("Input File:", input)
	fmt.Println("Output Path:", output)

	for {
		f, err := os.Open(input)
		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()

		csvReader := csv.NewReader(f)
		data, err := csvReader.ReadAll()
		if err != nil {
			log.Fatal(err)
		}

		tickets := ParseData(data)

		var validTickets []Ticket

		// sort through tickets for SOLCTR's tickets
		for _, ticket := range tickets {
			if (ticket.AssignmentGroup == "Solution Center") && (ticket.State != "Closed") {
				validTickets = append(validTickets, ticket)
			}
		}

		OutputTickets(output, validTickets)

		lastRun := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("\rLast run: %s || Tickets sorted: %d || Tickets returned: %d\x1b[?25l", lastRun, len(tickets), len(validTickets))

		time.Sleep(time.Hour * 24)
	}
}
