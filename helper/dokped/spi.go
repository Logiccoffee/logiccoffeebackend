package dokped

import (
	"bytes"
	"errors"
	"os"

	"github.com/gocroot/model"
	"github.com/jung-kurt/gofpdf"
)

func GenerateSPI(project model.Project, strkey string) (filecontent []byte, err error) {
	if len(project.Members) == 0 {
		err = errors.New("penulis belum di set pada project ini")
		return
	}
	user := project.Members[0]
	// Define the AES key (must be 32 bytes for AES-256)
	key := []byte(strkey) // Replace with a secure key
	// Download the logo image from a URL
	imageURL := "http://naskah.bukupedia.co.id/template/picture1.enc" // Replace with the actual image URL
	imageData, err := downloadImage(imageURL)
	if err != nil {
		return
	}
	// Save the image data to a temporary file
	encryptedFile := "image.enc"
	err = os.WriteFile(encryptedFile, imageData, 0644)
	if err != nil {
		return
	}
	defer os.Remove(encryptedFile) // Clean up the temp file after use
	// Step 2: Decrypt the image back
	decryptedFile := "picture1.png" // Path to save the decrypted image
	err = DecryptImage(encryptedFile, decryptedFile, key)
	if err != nil {
		return
	}
	defer os.Remove(decryptedFile) // Clean up the temp file after use

	// Download the logo image from a URL
	imageURL = "http://naskah.bukupedia.co.id/template/image.png" // Replace with the actual image URL
	imageData, err = downloadImage(imageURL)
	if err != nil {
		return
	}
	// Save the image data to a temporary file
	imageFile := "image.png"
	err = os.WriteFile(imageFile, imageData, 0644)
	if err != nil {
		return
	}
	defer os.Remove(imageFile) // Clean up the temp file after use

	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Set up header to apply on every page
	pdf.SetHeaderFunc(func() {
		// Position the logo at the top-right corner
		pdf.ImageOptions("image.png", 150, 10, 50, 0, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
		pdf.Ln(20) // Add some vertical spacing to move the content down
	})

	// Set up footer to apply on every page
	pdf.SetFooterFunc(func() {
		pdf.SetY(-32) // Position 32mm from the bottom of the page
		pdf.SetFont("Arial", "B", 9)
		pdf.SetTextColor(0, 0, 255) // Set text color to blue (RGB: 0, 0, 255)

		// Add some spacing above the "PT. Penerbit Buku Pedia" text
		pdf.Ln(14) // Increase the gap between footer and content above

		// Footer title
		pdf.Cell(0, 7, "PT. Penerbit Buku Pedia")
		pdf.Ln(5) // Line space between title and address

		// Footer content with reduced line height
		pdf.SetTextColor(0, 0, 0) // Reset to black color
		pdf.SetFont("Arial", "", 8)
		pdf.MultiCell(0, 3, `Athena Residence No. E1 Ciwaruga, Bandung Barat 40559
e-mail: penerbit@bukupedia.co.id
Telp: 087752000300
www.bukupedia.co.id`, "", "L", false)
	})

	// Add a new page to the document
	pdf.AddPage()

	// Set font for body content
	pdf.SetFont("Arial", "", 12)

	// Add the document content (like the image provided)

	// Title and date alignment
	pdf.CellFormat(30, 6, "NO", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(100, 6, "SPI"+generateNomorSurat(), "", 0, "L", false, 0, "")
	pdf.CellFormat(50, 6, "Bandung Barat, "+getTodayFormattedDate(), "", 1, "R", false, 0, "")

	pdf.CellFormat(30, 6, "Lampiran", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(100, 6, "1 Berkas", "", 0, "L", false, 0, "")
	pdf.CellFormat(50, 6, "", "", 1, "R", false, 0, "")

	pdf.CellFormat(30, 6, "Perihal", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(100, 6, "Permohonan ISBN/Barcode", "", 0, "L", false, 0, "")
	// Align "untuk ebook" under "Permohonan" by indenting to the same level
	pdf.Ln(6)
	pdf.CellFormat(35, 6, "", "", 0, "L", false, 0, "") // Indent to align with "Permohonan"
	pdf.CellFormat(100, 6, "untuk ebook", "", 1, "L", false, 0, "")
	pdf.Ln(10) // Line break for spacing

	// Add "Kepada" section
	pdf.Cell(0, 6, "Kepada :")
	pdf.Ln(10) // Line break for spacing

	// Add recipient details
	pdf.MultiCell(0, 6, `Yth. Kepala Pusat Bibliografi dan Pengolahan Bahan Perpustakaan
	Perpustakaan Nasional RI`, "", "L", false)
	pdf.Ln(10)

	// Add body content
	pdf.MultiCell(0, 6, "Bersama ini kami atas nama,", "", "L", false)
	pdf.Cell(40, 6, "Penerbit")
	pdf.Cell(5, 6, ":")
	pdf.Cell(100, 6, "PT. Penerbit Buku Pedia")
	pdf.Ln(6)
	pdf.Cell(40, 6, "Penanggung jawab")
	pdf.Cell(5, 6, ":")
	pdf.Cell(100, 6, "Rolly Maulana Awangga")
	pdf.Ln(6)
	pdf.Cell(40, 6, "Admin")
	pdf.Cell(5, 6, ":")
	pdf.Cell(100, 6, "Mila Anisa")
	pdf.Ln(10)

	// Add the ISBN application request
	pdf.MultiCell(0, 6, "Mengajukan permohonan ISBN untuk,", "", "L", false)
	pdf.Cell(40, 6, "Judul")
	pdf.Cell(5, 6, ":")
	pdf.MultiCell(0, 6, project.Title, "", "L", false)
	pdf.Ln(2)
	pdf.Cell(40, 6, "Kepengarangan")
	pdf.Cell(5, 6, ":")
	var listpenulis string
	for _, penulis := range project.Members {
		listpenulis += penulis.Name + ","
	}
	pdf.MultiCell(0, 6, listpenulis+";"+project.Editor.Name, "", "L", false)
	pdf.Ln(2)
	pdf.Cell(40, 6, "Link/akses")
	pdf.Cell(5, 6, ":")
	pdf.MultiCell(0, 6, project.URLKatalog, "", "L", false)
	pdf.Ln(10)

	// Add closing remarks
	pdf.MultiCell(0, 6, `Bersama ini kami lampirkan dummy buku dan Surat Pernyataan Keaslian Karya dari Penulis.
	
	Demikian permohonan ini kami ajukan, atas perhatian dan kerja samanya diucapkan terima kasih.`, "", "L", false)

	// Signature Section
	pdf.Ln(20) // Add some space before the signature section
	// Set font for the "Hormat kami" text
	pdf.SetFont("Arial", "", 12)

	// Add "Hormat kami" and center it
	pdf.CellFormat(0, 6, "Hormat kami,", "", 1, "R", false, 0, "")

	//pdf.Ln(3) // Add some space between "Hormat kami" and the signature

	// Add the signature image below "Hormat kami"
	pdf.ImageOptions("picture1.png", 160, pdf.GetY(), 40, 0, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")

	pdf.Ln(25) // Add space after the signature

	// Add the name under the signature and center it
	pdf.CellFormat(0, 6, "Rolly Maulana Awangga", "", 1, "R", false, 0, "")

	// After completing the first page, add the following for the second page:

	// Add a new page to the document for the second page
	pdf.AddPage()

	// Title of the second page
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, "SURAT PERNYATAAN KEASLIAN KARYA", "", 1, "C", false, 0, "")
	pdf.Ln(10) // Add space below the title

	// Set font for body content and start the formal text
	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 6, "Yang bertandatangan di bawah ini :", "", "L", false)
	pdf.Ln(5) // Line break for spacing

	// Add the detailed sections (Name, Address, etc.)
	pdf.Cell(40, 6, "Nama")
	pdf.Cell(5, 6, ":")
	pdf.Cell(100, 6, user.Name)
	pdf.Ln(6)
	pdf.Cell(40, 6, "Alamat")
	pdf.Cell(5, 6, ":")
	pdf.MultiCell(0, 6, user.AlamatRumah, "", "L", false)
	pdf.Ln(2)
	pdf.Cell(40, 6, "NIK")
	pdf.Cell(5, 6, ":")
	pdf.Cell(100, 6, user.NIK)
	pdf.Ln(6)
	pdf.Cell(40, 6, "Telp. /HP")
	pdf.Cell(5, 6, ":")
	pdf.Cell(100, 6, user.PhoneNumber)
	pdf.Ln(10) // Line break for spacing

	// Statement section
	pdf.MultiCell(0, 6, "menyatakan dengan sesungguhnya, bahwa :", "", "L", false)
	pdf.Ln(5)

	// Add the content for "Judul" and "Penulis"
	pdf.Cell(40, 6, "Judul")
	pdf.Cell(5, 6, ":")
	pdf.MultiCell(0, 6, project.Title, "", "L", false)
	pdf.Ln(2)
	pdf.Cell(40, 6, "Penulis")
	pdf.Cell(5, 6, ":")
	pdf.MultiCell(0, 6, listpenulis, "", "L", false)
	pdf.Ln(10) // Line break for spacing

	// Publishing information
	pdf.MultiCell(0, 6, "adalah benar merupakan karya asli yang dibuat untuk diterbitkan dan disebarluaskan secara umum, melalui :", "", "L", false)
	pdf.Ln(5)
	pdf.Cell(40, 6, "Penerbit")
	pdf.Cell(5, 6, ":")
	pdf.Cell(100, 6, "PT. Penerbit Buku Pedia")
	pdf.Ln(6)
	pdf.Cell(40, 6, "Alamat")
	pdf.Cell(5, 6, ":")
	pdf.MultiCell(0, 6, "Komp. Athena, Jl. Athena Raya No.E1 RT.04/13, Ciwaruga, Parongpong, Kab. Bandung Barat, Jawa Barat 40559", "", "L", false)
	pdf.Ln(10)

	// Final statement
	pdf.MultiCell(0, 6, `Demikian surat ini dibuat dengan sebenar-benarnya serta akan menjadi pertanggungjawaban kami jika terdapat penyalahgunaan dan akibat yang ditimbulkannya.`, "", "L", false)
	pdf.Ln(7) // Line break for spacing

	// Location and date
	pdf.CellFormat(0, 6, "Bandung Barat, "+getTodayFormattedDate(), "", 1, "R", false, 0, "")
	pdf.Ln(7) // Line break for spacing

	// Signatures
	pdf.CellFormat(0, 6, "Penanggung jawab Penerbit,", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, "Penulis,", "", 1, "R", false, 0, "")
	//pdf.Ln(10) // Add space between signature labels and actual signatures

	// Add the signatures and seal image
	pdf.ImageOptions("picture1.png", 10, pdf.GetY(), 40, 0, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
	//pdf.ImageOptions("picture1.png", 150, pdf.GetY(), 40, 0, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
	//pdf.ImageOptions("picture1.png", 80, pdf.GetY(), 30, 0, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")

	pdf.Ln(25) // Add space after the signatures

	// Add the names below signatures
	pdf.CellFormat(0, 6, "Rolly Maulana Awangga", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, user.Name, "", 1, "R", false, 0, "")

	// Create an in-memory buffer to store the PDF
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return
	}
	filecontent = buf.Bytes()
	return
}
