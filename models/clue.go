package models

import (
	"database/sql"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/nathanhollows/AmazingTrace/helpers"
	"github.com/yeqown/go-qrcode"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"gorm.io/gorm"
)

// Clue stores a simple riddle based clue for a location
type Clue struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
	Code      string `gorm:"uniqueIndex:clue_code,sort:desc;not null;primarykey"`
	Location  string `gorm:"not null"`
	Clue      string `gorm:"not null"`
	Longitude string
	Latitude  string
	Points    int
}

func CleanUpPosters(db *gorm.DB) {
	// Find every clue
	clues := []Clue{}
	db.Find(&clues)
	// Delete any posters that don't have a clue
	files, _ := ioutil.ReadDir("assets/img/posters")
	for _, f := range files {
		if f.Name() == ".gitkeep" {
			continue
		}
		found := false
		for _, c := range clues {
			if f.Name() == c.Code+".png" {
				found = true
				break
			}
		}
		if !found {
			os.Remove("assets/img/posters/" + f.Name())
		}
	}
	// Generate posters for any clues that don't have one
	for _, c := range clues {
		if !helpers.FileExists("assets/img/posters/" + c.Code + ".png") {
			go c.GeneratePoster()
		}
	}

}

// BeforeCreate generates a random string for the clue to be identified by
func (c *Clue) BeforeCreate(tx *gorm.DB) (err error) {
	c.Code = helpers.NewCode(4)
	return
}

// BeforeDelete will make sure every ClueLog has been soft deleted too
func (c *Clue) BeforeDelete(tx *gorm.DB) (err error) {
	result := tx.Where("clue_code = ?", c.Code).Delete(&ClueLog{})
	os.Remove("web/static/img/posters/" + c.Code + ".png")
	return result.Error
}

// AfterCreate generate a new poster
func (c Clue) AfterCreate(tx *gorm.DB) (err error) {
	go func(c Clue) {
		err = c.GeneratePoster()
	}(Clue{Code: c.Code})
	return err
}

// GeneratePoster pre-emptively generates the poster for the new clue
func (c Clue) GeneratePoster() error {
	imgb, _ := os.Open("assets/poster.png")
	img, _ := png.Decode(imgb)
	defer imgb.Close()

	background := color.RGBA{255, 213, 79, 255}
	foreground := color.RGBA{35, 35, 35, 255}
	// TODO: Factor out the hard coded link
	qrc, err := qrcode.New("https://trace.co.nz/"+c.Code,
		qrcode.WithBgColor(background),
		qrcode.WithFgColor(foreground),
		qrcode.WithBuiltinImageEncoder(qrcode.PNG_FORMAT))
	if err != nil {
		fmt.Printf("could not generate QRCode: %v", err)
		return err
	}
	if err := qrc.Save("assets/img/temp/" + c.Code + ".png"); err != nil {
		fmt.Printf("could not save image: %v", err)
		return err
	}

	wmb, _ := os.Open("assets/img/temp/" + c.Code + ".png")
	watermark, _ := png.Decode(wmb)
	defer wmb.Close()

	offset := image.Pt(463, 1075)
	b := img.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, img, image.Point{}, draw.Src)
	draw.Draw(m, watermark.Bounds().Add(offset), watermark, image.Point{}, draw.Over)

	addLabel(m, 440, 2050, fmt.Sprint("trace.co.nz/", c.Code))

	imgw, _ := os.Create("assets/img/posters/" + c.Code + ".png")
	png.Encode(imgw, m)
	defer imgw.Close()

	os.Remove("assets/img/temp/" + c.Code + ".png")

	return nil
}

var (
	dpi      = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
	fontfile = flag.String("fontfile", "assets/fonts/RobotoMono-Bold.ttf", "RobotoMono-Bold")
	hinting  = flag.String("hinting", "none", "none | full")
	size     = flag.Float64("size", 72, "font size in points")
)

func addLabel(img *image.RGBA, x, y int, label string) {
	flag.Parse()
	col := color.RGBA{254, 214, 79, 255}
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}

	// Read the font data.
	fontBytes, err := ioutil.ReadFile(*fontfile)
	if err != nil {
		log.Println(err)
		return
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}
	// Draw the text.
	h := font.HintingNone
	switch *hinting {
	case "full":
		h = font.HintingFull
	}
	d := &font.Drawer{
		Dst: img,
		Src: image.NewUniform(col),
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    *size,
			DPI:     *dpi,
			Hinting: h,
		}),
		Dot: point,
	}
	d.DrawString(label)
}
