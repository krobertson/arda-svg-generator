# Arda SVG Generator

This is meant to be used with [Arda Cards](https://arda.cards/) to generate a
pure SVG based version of the cards/labels.

Arda is an awesome platform for implementing Kanban for inventory management.
I've begin using it for all of the supplies in my manufacturing business. They
bootstrapped using Coda for the system and printing functionality with Coda is
sort of limited. It is largely HTML based, and there are definite limitations in
getting HTML output to tightly confirm to rigid parameters, such as when
printing to an Avery Business Card sheet. Scaling, margins, and strict physical
dimensioning isn't really HTML's thing.

Hence SVGs, and a pure SVG approach, it is powerful alternative. Still markup
based, vector based so scaling is seamless, and easily adapts to printed media,
where positioning and boundaries matter.

The philosophy here is to generate 2 levels of SVGs:

1. Generate an individual card/label SVG for each item based on a template for
   the specified size.
2. Generate a document from a template, which contains multiple individual
   items, positioned correctly for the print media.

You can think of #1 as file for an individual 3x2.5" business card, and #2 as
the collection of cards to be printed on the specific Avery template, or
whatever the print format is.

With this structure, the print document uses an `<image>` tag to load an
individual tile. That way, all styling and positioning are localized to each
tile. One tile's CSS won't clash with another tile. Text oberlapping won't throw
off all remaining tiles formatting.

One gotcha though, Chrome only recurses one level for remote `<image>`
references. So each individual card needs to have child images, including its QR
code, embedded in data URI format.

### Usage

To use this program, you will need to have Go installed. It was simply the
language I was most comfortable with and had the majority of what I needed with
minimal external dependencies. The only external dependency is a library for SVG
QR code generation.

Follow these steps:

1. You will need to Unlock your Items table to customize columns and export a
   CSV. For me, this is in the top right next to Share and Insert buttons within
   Arda.
2. Hover the Items table and go to Columns. Ensure the following non-standard
   columns are shown:
    * Color
    * Image (URL)
    * QR Code URL
    * Department
3. Hover the Items table and go to Options. Click on the 3 dots in the top right
   and go to "Download as CSV".
4. Remember to relock the screen.

Put the `Items.csv` file you downloaded into the current folder, then run the program.

```
$ go run *.go
```

It will output all the individual printable documents into `out/`. The separate
cards will be in `cards/` and labels will be in `labels/`. It will download any
images and store them in data URI format within the `cache/` folder.

To only process a subset of your items at a time, simply edit the `Items.csv`
file. Leave the header row, as it is how it maps column positions.