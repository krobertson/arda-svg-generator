package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Item represents an individual item in your Arda table.
type Item struct {
	Name             string
	CardSize         string
	LabelSize        string
	OrderMechanism   string
	InternalSupplier string
	Supplier         string
	Location         string
	MinQty           int
	OrderQty         int
	Notes            string
	QrUrl            string
	ImageUrl         string
	Color            string
	Department       string

	id         string
	imageCache string
	qrData     string
}

// ID is the md5sum of the link URL, which we use as a unique identifer for the
// item when storing cards/labels.
func (item *Item) ID() string {
	if item.id != "" {
		return item.id
	}

	h := md5.Sum([]byte(item.QrUrl))
	item.id = hex.EncodeToString(h[:])
	return item.id
}

// QrCode returns an optimized SVG of the QR code URL. Typically, this is then
// embedded into the cards.
func (item *Item) QrCode() string {
	if item.qrData != "" {
		return item.qrData
	}

	var err error
	item.qrData, err = generateQrCode(item.QrUrl)
	if err != nil {
		fmt.Printf("Error generating QR code: %q\n", err)
	}
	return item.qrData
}

// QrCodeSvg returns the embedded data for the QR code SVG generated by QrCode.
func (item *Item) QrCodeSvg() string {
	qrcode := item.QrCode()
	return "data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString([]byte(qrcode))
}

// HasImage returns a boolean to indicate whether the ImageUrl is specified.
func (item *Item) HasImage() bool {
	return item.ImageUrl != ""
}

// ImageSvg returns the item's image to be embedded into the cards/labels, or a
// blank string if no image is set.
func (item *Item) ImageSvg() string {
	if !item.HasImage() {
		return ""
	}

	if item.imageCache == "" {
		var err error
		item.imageCache, err = cacheImage(item.ImageUrl)
		if err != nil {
			fmt.Printf("Err in image: %q\n", err)
			return ""
		}
	}

	// if we still didn't get an image, just return
	if item.imageCache == "" {
		return ""
	}

	b, err := os.ReadFile(item.imageCache)
	if err != nil {
		fmt.Printf("Err in image: %q\n", err)
	}
	return string(b)
}

// NoteList splits the Notes fields on semicolons (';') so some things could be
// multilined.
func (item *Item) NoteList() []string {
	return strings.Split(item.Notes, ";")
}

// CssClass returns the Color field as it could be used for css color settings.
func (item *Item) CssClass() string {
	switch item.Color {
	case "Light Blue":
		return "blue"
	default:
		return strings.ToLower(item.Color)
	}
}

// Production returns a boolean if the Order Mechanism is set to Production.
func (item *Item) Production() bool {
	return item.OrderMechanism == "Production"
}

// loadItems reads the "Items.csv" file from the current directory to load all
// the items table.
func loadItems() ([]*Item, error) {
	f, err := os.Open("Items.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)

	// read header row and get header mapping
	record, err := r.Read()
	if err != nil {
		return nil, err
	}
	mapping := make(map[int]string)
	for i, r := range record {
		mapping[i] = r
	}

	items := make([]*Item, 0)

	// read items
	for {
		record, err := r.Read()
		if err != nil && err != io.EOF {
			return nil, err
		}
		if record == nil {
			break
		}

		rr := make(map[string]string)
		for i, r := range record {
			rr[mapping[i]] = r
		}

		item, err := readItem(rr)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// readItem converts the CSV row into an Item object.
func readItem(record map[string]string) (*Item, error) {
	item := &Item{
		Name:             record["Item"],
		CardSize:         record["Card Size"],
		LabelSize:        record["Label Size"],
		OrderMechanism:   record["Order Mechanism"],
		InternalSupplier: record["Internal Supplier"],
		Supplier:         record["Primary Supplier"],
		Location:         record["Location"],
		Notes:            record["Notes"],
		QrUrl:            record["QR Code URL"],
		ImageUrl:         record["Image (URL)"],
		Color:            record["Color"],
		Department:       record["Department"],
	}

	var err error
	if record["Min Qty"] != "" {
		item.MinQty, err = strconv.Atoi(record["Min Qty"])
		if err != nil {
			return nil, err
		}
	}

	if record["Order Qty"] != "" {
		item.OrderQty, err = strconv.Atoi(record["Order Qty"])
		if err != nil {
			return nil, err
		}
	}

	return item, nil
}
