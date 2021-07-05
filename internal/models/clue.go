package models

import (
	"database/sql"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"time"

	"github.com/nathanhollows/AmazingTrace/internal/helpers"
	"github.com/yeqown/go-qrcode"
	"gorm.io/gorm"
)

// Clue stores a simple riddle based clue for a location
type Clue struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
	Code      string `gorm:"uniqueIndex:clue_code,sort:desc;not null;primarykey"`
	Location  string `gorm:"index:unique;not null"`
	Clue      string `gorm:"not null"`
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

// GeneratePoster pre-emptively generates the poster for the new clue
func (c *Clue) GeneratePoster() error {
	// TODO: Add URL to bottom of poster
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
	if err := qrc.Save("assets/" + c.Code + ".png"); err != nil {
		fmt.Printf("could not save image: %v", err)
		return err
	}

	wmb, _ := os.Open("assets/" + c.Code + ".png")
	watermark, _ := png.Decode(wmb)
	defer wmb.Close()

	offset := image.Pt(463, 1075)
	b := img.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, img, image.ZP, draw.Src)
	draw.Draw(m, watermark.Bounds().Add(offset), watermark, image.ZP, draw.Over)

	imgw, _ := os.Create("web/static/img/posters/" + c.Code + ".png")
	png.Encode(imgw, m)
	defer imgw.Close()

	os.Remove("assets/" + c.Code + ".png")

	return nil
}
