package dokped

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gocroot/model"
	"github.com/jung-kurt/gofpdf"
)

func GenerateSPK(project model.Project, strkey string) (filecontent []byte, err error) {
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
		fmt.Println("Error downloading image:", err)
		return
	}
	// Save the image data to a temporary file
	encryptedFile := "image.enc"
	err = os.WriteFile(encryptedFile, imageData, 0644)
	if err != nil {
		fmt.Println("Error saving image:", err)
		return
	}
	defer os.Remove(encryptedFile) // Clean up the temp file after use
	// Step 2: Decrypt the image back
	decryptedFile := "picture1.png" // Path to save the decrypted image
	err = DecryptImage(encryptedFile, decryptedFile, key)
	if err != nil {
		fmt.Println("Error decrypting image:", err)
		return
	}
	fmt.Println("Image decrypted and saved as", decryptedFile)
	defer os.Remove(decryptedFile) // Clean up the temp file after use

	// Download the logo image from a URL
	imageURL = "http://naskah.bukupedia.co.id/template/image.png" // Replace with the actual image URL
	imageData, err = downloadImage(imageURL)
	if err != nil {
		fmt.Println("Error downloading image:", err)
		return
	}
	// Save the image data to a temporary file
	imageFile := "image.png"
	err = os.WriteFile(imageFile, imageData, 0644)
	if err != nil {
		fmt.Println("Error saving image:", err)
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

	// Add the image as a header (logo)
	//pdf.ImageOptions("image.png", 150, 10, 50, 0, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")

	// Add some vertical spacing before the title
	//pdf.Ln(20) // Create space between the logo and the text

	// Set font for the title and center the text
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "PERJANJIAN", "", 1, "C", false, 0, "")

	// Set a smaller font size for the SPK number and center it
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 0, "NO. SPK"+generateNomorSurat(), "", 1, "C", false, 0, "")
	pdf.Ln(15) // Add space before the content

	// Set font for body content to size 10
	pdf.SetFont("Arial", "", 10)

	// Check available space to prevent content from reaching the footer
	if pdf.GetY() > 240 {
		pdf.AddPage() // Add a new page if content reaches too close to the footer
	}

	// Add the content with tighter line spacing
	pdf.Cell(0, 5, "Yang bertanda tangan di bawah ini :")
	pdf.Ln(5) // Line break with 5 units of space

	// First entry (Rolly Maulana Awangga)
	pdf.Cell(10, 5, "1.")
	pdf.CellFormat(40, 5, "Nama", "0", 0, "L", false, 0, "")
	pdf.CellFormat(5, 5, ":", "0", 0, "L", false, 0, "")
	pdf.CellFormat(100, 5, "Rolly Maulana Awangga", "0", 1, "L", false, 0, "")

	pdf.Cell(10, 5, "")
	pdf.CellFormat(40, 5, "Jabatan", "0", 0, "L", false, 0, "")
	pdf.CellFormat(5, 5, ":", "0", 0, "L", false, 0, "")
	pdf.CellFormat(100, 5, "Direktur PT. Penerbit Buku Pedia", "0", 1, "L", false, 0, "")

	pdf.Cell(10, 5, "")
	pdf.CellFormat(40, 5, "Alamat", "0", 0, "L", false, 0, "")
	pdf.CellFormat(5, 5, ":", "0", 0, "L", false, 0, "")
	pdf.CellFormat(100, 5, "Athena Residence E1 Ciwaruga Kab. Bandung Barat 40559", "0", 1, "L", false, 0, "")

	// Add description text after the first entry
	pdf.MultiCell(0, 5, "Adalah bertindak atas nama Penerbit Buku Pedia yang selanjutnya disebut Pihak Kesatu.", "", "L", false)

	// Second entry (Muhammad Rizal Satria)
	pdf.Cell(10, 5, "2.")
	pdf.CellFormat(40, 5, "Nama Lengkap", "0", 0, "L", false, 0, "")
	pdf.CellFormat(5, 5, ":", "0", 0, "L", false, 0, "")
	pdf.CellFormat(100, 5, user.Name, "0", 1, "L", false, 0, "")

	pdf.Cell(10, 5, "")
	pdf.CellFormat(40, 5, "Alamat Rumah", "0", 0, "L", false, 0, "")
	pdf.CellFormat(5, 5, ":", "0", 0, "L", false, 0, "")
	pdf.CellFormat(100, 5, user.AlamatRumah, "0", 1, "L", false, 0, "")

	pdf.Cell(10, 5, "")
	pdf.CellFormat(40, 5, "Alamat Kantor", "0", 0, "L", false, 0, "")
	pdf.CellFormat(5, 5, ":", "0", 0, "L", false, 0, "")
	pdf.CellFormat(100, 5, user.AlamatKantor, "0", 1, "L", false, 0, "")

	pdf.Cell(10, 5, "")
	pdf.CellFormat(40, 5, "Telp. / HP", "0", 0, "L", false, 0, "")
	pdf.CellFormat(5, 5, ":", "0", 0, "L", false, 0, "")
	pdf.CellFormat(100, 5, user.PhoneNumber, "0", 1, "L", false, 0, "")

	pdf.Cell(10, 5, "")
	pdf.CellFormat(40, 5, "Alamat Email", "0", 0, "L", false, 0, "")
	pdf.CellFormat(5, 5, ":", "0", 0, "L", false, 0, "")
	pdf.CellFormat(100, 5, user.Email, "0", 1, "L", false, 0, "")

	pdf.Cell(10, 5, "")
	pdf.CellFormat(40, 5, "Pekerjaan", "0", 0, "L", false, 0, "")
	pdf.CellFormat(5, 5, ":", "0", 0, "L", false, 0, "")
	pdf.CellFormat(100, 5, user.Pekerjaan, "0", 1, "L", false, 0, "")

	pdf.MultiCell(0, 5, "Adalah Penulis buku yang berjudul : "+project.Title+" yang selanjutnya disebut sebagai Pihak Kedua.", "", "L", false)
	pdf.Ln(4)
	pdf.MultiCell(0, 5, "Kedua belah pihak dengan dasar saling percaya mengadakan perjanjian sebagai berikut :", "", "L", false)

	// Pasal 1 (Naskah Buku) Section - Adjusted
	pdf.Ln(4)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(0, 5, "Pasal 1 (Naskah Buku)", "", 1, "C", false, 0, "") // Center aligned
	pdf.SetFont("Arial", "", 10)

	// Add the numbered list with adjusted margins for numbers
	addNumberedSection(pdf, 1, "Pihak Kedua memberikan jaminan kepada Pihak Kesatu bahwa naskah buku tidak menyinggung atau merugikan hak penulis lain atau penerbit lain, tidak memuat plagiarism, hal-hal yang dapat sebagai fitnah, atau merugikan pihak ketiga.")
	addNumberedSection(pdf, 2, "Pihak Kedua bersedia mengganti setiap kerugian atau membayar segala ongkos jika timbul gugatan pihak ketiga karena jaminan di atas ternyata tidak atau kurang dipenuhi.")
	addNumberedSection(pdf, 3, "Pihak Kedua bila perlu menyiapkan koreksi-koreksi (revisi) atau tambahan yang diperlukan untuk persiapan cetakan.")

	// Pasal 2 (Teknik Penerbitan Buku) Section
	pdf.Ln(4)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(0, 5, "Pasal 2 (Teknik Penerbitan Buku)", "", 1, "C", false, 0, "") // Center aligned
	pdf.SetFont("Arial", "", 10)

	// Add the numbered list with adjusted margins for numbers
	addNumberedSection(pdf, 1, "Jumlah Buku yang didistribusikan ditetapkan oleh Pihak Kesatu sesuai dengan sudut pandang usaha dengan melihat peluang pasar kemudian diinformasikan pada Pihak Kedua. Segala biaya yang telah dan akan dikeluarkan untuk keperluan tersebut sepenuhnya menjadi tanggung jawab Pihak Kesatu.")
	addNumberedSection(pdf, 2, "Distribusi ulang buku dilakukan secara otomatis oleh Pihak Kesatu kecuali ada pemberitahuan dari Pihak Kedua bahwa buku yang telah diterbitkan sebelumnya ada perubahan, penambahan ataupun revisi.")
	addNumberedSection(pdf, 3, "Pemberitahuan perubahan, penambahan ataupun revisi dari Pihak Kedua dilakukan sejak buku tersebut diterbitkan sampai tiba saatnya untuk didistribusikan ulang (kurang lebih 1 tahun).")
	addNumberedSection(pdf, 4, "Untuk kebutuhan promosi dan pemasaran, diperlukan 50 (lima puluh) eksemplar buku contoh pada distribusi perdana dengan rincian 40 (empat puluh) eksemplar untuk kebutuhan promosi dan pemasaran Pihak Kesatu, dan 10 (sepuluh) eksemplar buku contoh untuk penulis.")
	addNumberedSection(pdf, 5, "Diperhitungkan sejumlah 2% (dua persen) buku untuk toleransi terhadap risiko kegagalan yang terjadi pada proses distribusi atau penjualan.")

	// Pasal 3 (Hak-hak Pihak Kedua) Section
	pdf.Ln(4)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(0, 5, "Pasal 3 (Hak-hak Pihak Kedua)", "", 1, "C", false, 0, "") // Center aligned
	pdf.SetFont("Arial", "", 10)

	// Add the numbered list with adjusted margins for numbers
	addNumberedSection(pdf, 1, "Pihak Kedua memiliki hak penuh atas naskah buku.")
	addNumberedSection(pdf, 2, "Pihak Kedua menerima minimal sebanyak 100 (seratus) eksemplar buku untuk dipasarkan. Hasil pemasarannya dapat disetorkan langsung ke Pihak Kesatu atau diperhitungkan dengan royalty.")
	addNumberedSection(pdf, 3, "Pihak Kedua berhak mendapat buku berikutnya dengan cara membeli dari Pihak Kesatu dan mendapatkan potongan harga sebesar 35% (tiga puluh lima persen) dari harga satuan yang ditetapkan.")

	// Pasal 4 (Pembayaran Hak Intelektual Pihak Kedua / Royalti) Section
	pdf.Ln(4)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(0, 5, "Pasal 4 (Pembayaran Hak Intelektual Pihak Kedua / Royalti)", "", 1, "C", false, 0, "") // Center aligned
	pdf.SetFont("Arial", "", 10)

	// Add the numbered list with adjusted margins for numbers
	addNumberedSection(pdf, 1, "Pihak Kesatu Membayar royalty kepada Pihak Kedua sebesar 10% (sepuluh persen) dari harga jual buku tanpa perantara, jika melalui perantara maka perhitungan 10% (sepuluh persen) diambil dari nilai yang diterima Pihak Kesatu setelah dipotong biaya distribusi perantara.")
	addNumberedSection(pdf, 2, "Pihak Kesatu berhak membayar royalty seluruhnya dalam bentuk buku, apabila buku yang diterbitkan kurang laku dalam jangka waktu 2 (dua) tahun sejak buku pertama kali dipasarkan.")

	// Pasal 5 (Hal-hal Lain) Section
	pdf.Ln(4)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(0, 5, "Pasal 5 (Hal-hal Lain)", "", 1, "C", false, 0, "") // Center aligned
	pdf.SetFont("Arial", "", 10)

	// Add the numbered list with adjusted margins for numbers
	addNumberedSection(pdf, 1, "Selama Pihak Kesatu masih menerbitkan buku yang dimaksud dalam perjanjian ini, Pihak Kedua tidak dapat menarik diri, membatalkan perjanjian, atau menerbitkannya di penerbit lain tanpa ada persetujuan tertulis dari Pihak Kesatu.")
	addNumberedSection(pdf, 2, "Surat perjanjian ini ditandatangani oleh penulis atau salah seorang penulis (jika penulis buku lebih dari satu orang).")
	addNumberedSection(pdf, 3, "Apabila buku yang ditulis merupakan hasil karya Bersama (tim), maka penulis yang menandatangani surat perjanjian bertanggung jawab terhadap rekan-rekan penulis lainnya dalam hal keontetikan naskah, materi revisi, dan royalty.")
	addNumberedSection(pdf, 4, "Bila terjadi perselisihan antara tim penulis (apabila buku yang ditulis merupakan karya bersama), Pihak Kesatu tidak akan turut campur di dalamnya. Pihak Kesatu hanya akan berhubungan dengan wakil dari Pihak Kedua yang menandatangani surat perjanjian.")
	addNumberedSection(pdf, 5, "Hal-hal yang tidak diatur atau belum sempurna diatur dalam perjanjian ini akan diputuskan oleh kedua belah pihak dengan persetujuan bersama.")

	// Pasal 6 (Akibat-akibat Hukum yang Timbul) Section
	pdf.Ln(4)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(0, 5, "Pasal 6 (Akibat-akibat Hukum yang Timbul)", "", 1, "C", false, 0, "") // Center aligned
	pdf.SetFont("Arial", "", 10)

	// Add the content without numbering
	pdf.MultiCell(0, 5, "Segala perbedaan yang mungkin timbul antara Kedua Pihak akan diselesaikan secara musyawarah untuk mufakat. Jika hal tersebut tidak tercapai, maka Pihak Kesatu dan Pihak Kedua memilih domisili di Bandung dan penyelesaian di Kantor Panitera Pengadilan Negeri Bandung.", "", "L", false)
	pdf.MultiCell(0, 5, "Surat Perjanjian ini dibuat secara elektronik bermaterai elektronik cukup serta memiliki kekuatan hukum.", "", "L", false)

	// Add space before the closing statement
	pdf.Ln(5)

	// Add the closing statement
	pdf.MultiCell(0, 5, "Demikian Surat Perjanjian ini dibuat di Bandung pada hari Senin tanggal Dua Puluh Dua bulan Desember tahun Dua Ribu Dua Puluh Empat.", "", "L", false)

	// Signature Section
	pdf.Ln(10) // Add some space before the signature section

	// Set font for the signature labels
	pdf.SetFont("Arial", "", 12)

	// Pihak Kedua (left)
	pdf.CellFormat(90, 5, "Pihak Kedua,", "", 0, "C", false, 0, "") // Centered label for Pihak Kedua

	// Pihak Kesatu (right)
	pdf.CellFormat(90, 5, "Pihak Kesatu,", "", 1, "C", false, 0, "") // Centered label for Pihak Kesatu

	// Add signature images and the stamp in the middle
	pdf.Ln(10) // Add space before the signatures

	// Add the signature for Pihak Kedua on the left
	//pdf.ImageOptions("image.png", 30, pdf.GetY(), 40, 0, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")

	// Add the stamp in the middle
	//pdf.ImageOptions("image.png", 90, pdf.GetY()-10, 40, 0, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")

	// Add the signature for Pihak Kesatu on the right
	pdf.ImageOptions("picture1.png", 125, pdf.GetY()-5, 40, 0, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")

	pdf.Ln(20) // Add some space after the signatures

	// Add the names under the signatures
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(90, 5, user.Name, "", 0, "C", false, 0, "")               // Name under Pihak Kedua
	pdf.CellFormat(90, 5, "Rolly Maulana Awangga", "", 1, "C", false, 0, "") // Name under Pihak Kesatu

	// Output the PDF to a file
	//err = pdf.OutputFileAndClose("perjanjian_document_with_pasal1_aligned_numbers.pdf")
	// Create an in-memory buffer to store the PDF
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return
	}
	filecontent = buf.Bytes()
	return
}

// Function to add numbered section with aligned text
func addNumberedSection(pdf *gofpdf.Fpdf, number int, text string) {
	pdf.CellFormat(10, 5, fmt.Sprintf("%d.", number), "0", 0, "L", false, 0, "")
	pdf.MultiCell(180, 5, text, "", "L", false)
}

// Function to download the image
func downloadImage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the image data into a byte slice
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func getTodayFormattedDate() string {
	// Daftar nama bulan dalam bahasa Indonesia
	months := []string{
		"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}

	// Load lokasi untuk Asia/Jakarta
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return ""
	}

	// Mengambil tanggal hari ini sesuai dengan zona waktu Asia/Jakarta
	today := time.Now().In(location)

	// Membentuk tanggal dengan format 22 Juli 2024
	formattedDate := fmt.Sprintf("%d %s %d", today.Day(), months[int(today.Month())-1], today.Year())

	return formattedDate
}

// Function to encrypt an image file
func EncryptImage(inputFile, outputFile string, key []byte) error {
	// Read the input image file
	imageData, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	// Encrypt the image data
	encryptedData, err := encryptAES(imageData, key)
	if err != nil {
		return err
	}

	// Write the encrypted data to a file
	err = os.WriteFile(outputFile, encryptedData, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Function to decrypt an image file
func DecryptImage(inputFile, outputFile string, key []byte) error {
	// Read the encrypted file
	encryptedData, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	// Decrypt the data
	decryptedData, err := decryptAES(encryptedData, key)
	if err != nil {
		return err
	}

	// Write the decrypted image back to a file
	err = os.WriteFile(outputFile, decryptedData, 0644)
	if err != nil {
		return err
	}

	return nil
}

// AES encryption function
func encryptAES(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// GCM is an authenticated encryption mode
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Seal encrypts the data and appends the nonce to the beginning of the ciphertext
	ciphertext := aesGCM.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// AES decryption function
func decryptAES(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// Function to generate the code in the format YYYYMMDDHHMMSS
func generateNomorSurat() string {
	// Get the current date and time
	now := time.Now()

	// Format the code as YYYYMMDDHHMMSS
	code := now.Format("20060102150405") // Format: YearMonthDayHourMinuteSecond
	return code
}
