package dokped

import (
	"bytes"
	"errors"
	"os"

	"github.com/gocroot/model"
	"github.com/jung-kurt/gofpdf"
)

func GenerateSPKT(project model.Project, strkey string) (filecontent []byte, err error) {
	if len(project.Members) == 0 {
		err = errors.New("penulis belum di set pada project ini")
		return
	}
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
	pdf.CellFormat(100, 6, "SPKT"+generateNomorSurat(), "", 0, "L", false, 0, "")
	pdf.CellFormat(50, 6, "Bandung Barat, "+getTodayFormattedDate(), "", 1, "R", false, 0, "")

	pdf.CellFormat(30, 6, "Lampiran", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(100, 6, "1 Berkas", "", 0, "L", false, 0, "")
	pdf.CellFormat(50, 6, "", "", 1, "R", false, 0, "")

	pdf.CellFormat(30, 6, "Perihal", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(100, 6, "Penyerahan Karya Terbitan", "", 0, "L", false, 0, "")
	// Align "untuk ebook" under "Permohonan" by indenting to the same level
	pdf.Ln(6)
	pdf.CellFormat(35, 6, "", "", 0, "L", false, 0, "") // Indent to align with "Permohonan"
	pdf.CellFormat(100, 6, "Buku Elektronik", "", 1, "L", false, 0, "")
	pdf.Ln(10) // Line break for spacing

	// Add "Kepada" section
	pdf.Cell(0, 6, "Kepada :")
	pdf.Ln(10) // Line break for spacing

	// Add recipient details
	pdf.MultiCell(0, 6, `Yth. Kepala Dinas Perpustakaan dan Kearsipan Daerah
	Provinsi Jawa Barat`, "", "L", false)
	pdf.Ln(10)

	// Add body content
	pdf.MultiCell(0, 6, "Berdasarkan Undang Undang Republik Indonesia Nomor 13 tahun 2018 Tentang Serah Simpan Karya Cetak dan Karya Rekam Pasal 4 dan 5. Bahwa setiap penerbit buku menyerahkan buku yang telah diterbitkan kepada Perpustakaan provinsi tempat domisili. Bersama ini kami atas nama,", "", "L", false)
	pdf.Cell(40, 6, "Penerbit")
	pdf.Cell(5, 6, ":")
	pdf.Cell(100, 6, "PT. Penerbit Buku Pedia")
	pdf.Ln(6)
	pdf.Cell(40, 6, "Penanggung jawab")
	pdf.Cell(5, 6, ":")
	pdf.Cell(100, 6, "Rolly Maulana Awangga")
	pdf.Ln(6)
	pdf.Cell(40, 6, "Jabatan")
	pdf.Cell(5, 6, ":")
	pdf.Cell(100, 6, "Direktur Utama")
	pdf.Ln(10)

	// Add the ISBN application request
	pdf.MultiCell(0, 6, "Menyerahkan hasil terbitan dari penerbit kami yang diterbitkan secara elektronik. Adapun data buku elektronik yang kami serahkan kepada Perpustakaan Daerah Jawa Barat adalah,", "", "L", false)
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
	pdf.MultiCell(0, 6, `Bersama dengan surat ini, kami sertakan sumber buku elektronik yang sudah kami terbitkan dalam bentuk Cakram Padat/Compact Disk(CD).
Demikian surat pengantar ini kami buat, atas perhatian dan kerjasamanya diucapkan terima kasih.`, "", "L", false)

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

	// Create an in-memory buffer to store the PDF
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return
	}
	filecontent = buf.Bytes()
	return
}
