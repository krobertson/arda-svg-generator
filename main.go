package main

import (
	"fmt"
	"path/filepath"
)

type Mapping struct {
	Name     string
	File     string
	Document string
	Count    int
}

var cardMapping = map[string]*Mapping{
	`3” x 5”`: {
		Name:     "index-cards",
		File:     "templates/card-3x5.svg",
		Document: "templates/indexcard.svg",
		Count:    1,
	},
	`3.5" x 2"`: {
		Name:     "business-cards",
		File:     "templates/card-2x35.svg",
		Document: "templates/doc-avery-business-cards.svg",
		Count:    10,
	},
}

var labelMapping = map[string]*Mapping{
	`1” x 3”`: {
		Name:     "uline",
		File:     "templates/label-1x3.svg",
		Document: "templates/doc-uline-1x3.svg",
		Count:    20,
	},
	`3.5" x 2"`: {
		Name:     "business-cards",
		File:     "templates/label-2x35.svg",
		Document: "templates/doc-avery-business-cards.svg",
		Count:    10,
	},
	`3.5" x 2" Notes`: {
		Name:     "business-cards",
		File:     "templates/label-2x35.svg",
		Document: "templates/doc-avery-business-cards.svg",
		Count:    10,
	},
}

type document struct {
	mapping *Mapping
	entries []string
}

var documents = make(map[string]*document)

func getDocument(m *Mapping) *document {
	d, exists := documents[m.Document]
	if exists {
		return d
	}

	d = &document{
		mapping: m,
		entries: make([]string, 0),
	}
	documents[m.Document] = d
	return d
}

func main() {
	items, err := loadItems()
	if err != nil {
		panic(err)
	}

	for _, item := range items {
		cm, cexists := cardMapping[item.CardSize]
		lm, lexists := labelMapping[item.LabelSize]

		if cexists {
			entry := filepath.Join("cards", item.ID()+".svg")
			if err := generateEntry(item, cm.File, entry); err != nil {
				panic(err)
			}

			doc := getDocument(cm)
			doc.entries = append(doc.entries, filepath.Join("..", entry))
		}

		if lexists {
			entry := filepath.Join("labels", item.ID()+".svg")
			if err := generateEntry(item, lm.File, entry); err != nil {
				panic(err)
			}

			doc := getDocument(lm)
			doc.entries = append(doc.entries, filepath.Join("..", entry))
		}
	}

	fmt.Printf("Finished generating entries\n")

	for _, doc := range documents {
		fmt.Printf("Processing %q, %d entries\n", doc.mapping.Document, len(doc.entries))
		num := 1
		finished := false

		for {
			f := filepath.Join("out", fmt.Sprintf("%s-%d.svg", doc.mapping.Name, num))
			items := doc.entries
			if len(items) > doc.mapping.Count {
				items = items[0:doc.mapping.Count]
				doc.entries = doc.entries[doc.mapping.Count:]
			} else {
				finished = true
			}
			fmt.Printf("  Generate %q with %d\n", doc.mapping.Document, len(items))
			err := generateDocument(items, doc.mapping.Document, f)
			if err != nil {
				panic(err)
			}
			num += 1
			if finished {
				break
			}
		}
	}

}
