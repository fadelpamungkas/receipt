package main

import (
	_ "embed"
	"fmt"
	"log"
	"strings"

	"github.com/signintech/gopdf"
)

//go:embed "Inter/Inter Variable/Inter.ttf"
var interFont []byte

//go:embed "Inter/Inter Hinted for Windows/Desktop/Inter-Bold.ttf"
var interBoldFont []byte

type Receipt struct {
	Number   string
	Date     string
	BillFrom Company
	BillTo   Company
	Items    []Item
	Notes    string
	Subtotal float64
	Tax      float64
	Discount float64
	Total    float64
}

type Company struct {
	Name  string
	Logo  string
	Email string
}

type Item struct {
	name     string
	quantity int
	price    float64
}

func main() {
	generateReceipt(Receipt{
		Number: "083",
		Date:   "Aug 3, 2023",
		BillFrom: Company{
			Name:  "Gopher Inc.",
			Logo:  "gopher.png",
			Email: "gopher@gmail.com",
		},
		BillTo: Company{
			Name:  "John Doe",
			Email: "john@gmail.com",
		},
		Items: []Item{
			{name: "Mystic Staff", quantity: 8, price: 320},
			{name: "Aghanim's Scepter", quantity: 4, price: 420},
			{name: "Perseverance", quantity: 1, price: 100},
			{name: "Vitality Booster", quantity: 19, price: 110},
			{name: "Sacred Relic", quantity: 2, price: 380},
		},
		Notes:    "Thank you for your business",
		Tax:      130,
		Discount: 80,
	})
}

func generateReceipt(receipt Receipt) {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{
		PageSize: *gopdf.PageSizeA4,
	})
	pdf.SetMargins(40, 40, 40, 40)
	pdf.AddPage()

	err := pdf.AddTTFFontData("Inter", interFont)
	if err != nil {
		log.Print(err.Error())
		return
	}

	err = pdf.AddTTFFontData("Inter-Bold", interBoldFont)
	if err != nil {
		log.Print(err.Error())
		return
	}

	setCompany(&pdf, receipt.BillFrom)
	setBillTo(&pdf, receipt.BillTo)
	subtotal := setItems(&pdf, receipt.Items)
	setNotes(&pdf, receipt.Notes)
	setTotals(&pdf, subtotal, receipt.Tax, receipt.Discount)
	setHeader(&pdf)

	fmt.Println("PDF generated!")

	name := strings.ReplaceAll(receipt.BillTo.Name, " ", "-")
	filename := fmt.Sprintf("%s-%s.pdf", receipt.Number, name)
	err = pdf.WritePdf("receipts/" + filename)
	if err != nil {
		log.Print(err.Error())
		return
	}
}
