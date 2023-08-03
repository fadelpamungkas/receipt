package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/signintech/gopdf"
)

const (
	quantityColumnOffset = 360
	rateColumnOffset     = 405
	amountColumnOffset   = 480
)

const (
	RECEIPT_LABEL    = "RECEIPT"
	NOTES_LABEL      = "NOTES"
	RECEIPT_NO_LABEL = "RECEIPT NO"
	DATE_LABEL       = "DATE"
	BILL_TO_LABEL    = "BILL TO"
	SUBTOTAL_LABEL   = "SUBTOTAL"
	TAX_LABEL        = "TAX"
	DISCOUNT_LABEL   = "DISCOUNT"
	TOTAL_LABEL      = "TOTAL"
	ITEM_LABEL       = "ITEM"
	QUANTITY_LABEL   = "QTY"
	RATE_LABEL       = "RATE"
	PRICE_LABEL      = "PRICE"
)

func setCompany(pdf *gopdf.GoPdf, c Company) {
	if c.Logo != "" {
		width, height := getImageDimension(c.Logo)
		scaledWidth := 100.0
		scaledHeight := float64(height) * scaledWidth / float64(width)
		_ = pdf.Image(c.Logo, pdf.GetX(), pdf.GetY(), &gopdf.Rect{W: scaledWidth, H: scaledHeight})
		pdf.Br(scaledHeight + 24)
	}
	_ = pdf.SetFont("Inter", "", 12)
	pdf.SetTextColor(55, 55, 55)
	_ = pdf.Cell(nil, c.Name)
	pdf.Br(36)
	pdf.SetStrokeColor(225, 225, 225)
	pdf.Line(pdf.GetX(), pdf.GetY(), 100, pdf.GetY())
	pdf.Br(36)
}

func getImageDimension(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	defer file.Close()

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", imagePath, err)
	}
	return image.Width, image.Height
}

func setHeader(pdf *gopdf.GoPdf) {
	err := pdf.SetFont("Inter", "", 32)
	if err != nil {
		log.Print(err.Error())
		return
	}
	pdf.SetXY(rateColumnOffset, 40)
	pdf.SetTextColor(0, 0, 0)
	_ = pdf.Cell(nil, RECEIPT_LABEL)
	pdf.Br(36)

	pdf.Br(24)
	rightItem(pdf, RECEIPT_NO_LABEL, "083")
	rightItem(pdf, DATE_LABEL, "Aug 3, 2023")
}

func setNotes(pdf *gopdf.GoPdf, notes string) {
	pdf.SetY(650)

	_ = pdf.SetFont("Inter", "", 10)
	pdf.SetTextColor(55, 55, 55)
	_ = pdf.Cell(nil, NOTES_LABEL)
	pdf.Br(18)
	_ = pdf.SetFont("Inter", "", 9)
	pdf.SetTextColor(0, 0, 0)

	formattedNotes := strings.ReplaceAll(notes, `\n`, "\n")
	notesLines := strings.Split(formattedNotes, "\\n")

	for i := 0; i < len(notesLines); i++ {
		_ = pdf.Cell(nil, notesLines[i])
		pdf.Br(15)
	}

	pdf.Br(48)
}

func setBillTo(pdf *gopdf.GoPdf, c Company) {
	pdf.SetTextColor(75, 75, 75)
	_ = pdf.SetFont("Inter", "", 9)
	_ = pdf.Cell(nil, BILL_TO_LABEL)
	pdf.Br(18)
	pdf.SetTextColor(75, 75, 75)
	_ = pdf.SetFont("Inter", "", 15)
	_ = pdf.Cell(nil, c.Name)
	pdf.Br(24)
	pdf.SetTextColor(75, 75, 75)
	_ = pdf.SetFont("Inter", "", 12)
	_ = pdf.Cell(nil, c.Email)
	pdf.Br(64)
}

func setTotals(pdf *gopdf.GoPdf, subtotal float64, tax float64, discount float64) {
	pdf.SetY(650)

	rightItem(pdf, SUBTOTAL_LABEL, toDollar(subtotal))
	if tax > 0 {
		rightItem(pdf, TAX_LABEL, toDollar(tax))
	}
	if discount > 0 {
		rightItem(pdf, DISCOUNT_LABEL, toDollar(discount))
	}
	rightItem(pdf, TOTAL_LABEL, toDollar(subtotal+tax-discount))
}

func rightItem(pdf *gopdf.GoPdf, label string, value string) {
	if label == TOTAL_LABEL {
		_ = pdf.SetFont("Inter-Bold", "", 11.5)
	} else {
		_ = pdf.SetFont("Inter", "", 9)
	}
	pdf.SetTextColor(75, 75, 75)
	pdf.SetX(rateColumnOffset)
	_ = pdf.Cell(nil, label)
	pdf.SetTextColor(0, 0, 0)
	_ = pdf.SetFontSize(12)
	pdf.SetX(amountColumnOffset)
	_ = pdf.Cell(nil, value)
	pdf.Br(24)
}

func toDollar(total float64) string {
	return fmt.Sprintf("$%.2f", total)
}

func setItems(pdf *gopdf.GoPdf, items []Item) float64 {
	_ = pdf.SetFont("Inter", "", 9)
	pdf.SetTextColor(55, 55, 55)
	_ = pdf.Cell(nil, ITEM_LABEL)
	pdf.SetX(quantityColumnOffset)
	_ = pdf.Cell(nil, QUANTITY_LABEL)
	pdf.SetX(rateColumnOffset)
	_ = pdf.Cell(nil, PRICE_LABEL)
	pdf.SetX(amountColumnOffset)
	_ = pdf.Cell(nil, TOTAL_LABEL)
	pdf.Br(24)

	subtotal := 0.0

	for _, item := range items {
		_ = pdf.SetFont("Inter", "", 11)
		pdf.SetTextColor(0, 0, 0)

		_ = pdf.Cell(nil, item.name)
		pdf.SetX(quantityColumnOffset)
		_ = pdf.Cell(nil, strconv.Itoa(item.quantity))
		pdf.SetX(rateColumnOffset)
		_ = pdf.Cell(nil, toDollar(item.price))
		pdf.SetX(amountColumnOffset)
		total := float64(item.quantity) * item.price
		_ = pdf.Cell(nil, toDollar(total))
		pdf.Br(24)
		subtotal += total
	}

	return subtotal
}
