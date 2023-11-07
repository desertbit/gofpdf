package tiff_test

import (
	"github.com/desertbit/gofpdf"
	"github.com/desertbit/gofpdf/contrib/tiff"
	"github.com/desertbit/gofpdf/internal/example"
)

// ExampleRegisterFile demonstrates the loading and display of a TIFF image.
func ExampleRegisterFile() {
	pdf, _ := gofpdf.New("L", "mm", "A4", "")
	pdf.SetFont("Helvetica", "", 12)
	pdf.SetFillColor(200, 200, 220)
	pdf.AddPageFormat("L", gofpdf.SizeType{Wd: 200, Ht: 200})
	opt := gofpdf.ImageOptions{ImageType: "tiff", ReadDpi: false}
	_ = tiff.RegisterFile(pdf, "sample", opt, "../../image/golang-gopher.tiff")
	pdf.Image("sample", 0, 0, 200, 200, false, "", 0, "")
	fileStr := example.Filename("Fpdf_Contrib_Tiff")
	err := pdf.OutputFileAndClose(fileStr)
	example.Summary(err, fileStr)
	// Output:
	// Successfully generated ../../pdf/Fpdf_Contrib_Tiff.pdf
}
