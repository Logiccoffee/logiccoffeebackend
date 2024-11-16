package fpdf

import (
	"bytes"
	"io"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pkg/errors"
)

// MergePDFBytes merges two PDF files provided as []byte and returns the merged result as []byte.
func MergePDFBytes(pdf1, pdf2 []byte) ([]byte, error) {
	// Create in-memory buffers for input PDFs
	input1 := bytes.NewReader(pdf1)
	input2 := bytes.NewReader(pdf2)

	// Create temporary files to save in-memory PDFs (pdfcpu works with file paths)
	tmpFile1, err := os.CreateTemp("", "pdf1_*.pdf")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create temp file for pdf1")
	}
	defer os.Remove(tmpFile1.Name()) // Clean up the temporary file

	tmpFile2, err := os.CreateTemp("", "pdf2_*.pdf")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create temp file for pdf2")
	}
	defer os.Remove(tmpFile2.Name()) // Clean up the temporary file

	// Write the in-memory bytes to temporary files
	if _, err := io.Copy(tmpFile1, input1); err != nil {
		return nil, errors.Wrap(err, "failed to copy pdf1 data to temp file")
	}
	if _, err := io.Copy(tmpFile2, input2); err != nil {
		return nil, errors.Wrap(err, "failed to copy pdf2 data to temp file")
	}

	// Close the files so they can be read later by pdfcpu
	if err := tmpFile1.Close(); err != nil {
		return nil, errors.Wrap(err, "failed to close temp file for pdf1")
	}
	if err := tmpFile2.Close(); err != nil {
		return nil, errors.Wrap(err, "failed to close temp file for pdf2")
	}

	// Create another temporary file to store the merged output
	mergedFile, err := os.CreateTemp("", "merged_*.pdf")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create temp file for merged PDF")
	}
	defer os.Remove(mergedFile.Name()) // Clean up the temporary file after reading it

	// Prepare the input files for merging
	inputFiles := []string{tmpFile1.Name(), tmpFile2.Name()}

	// Use the os.Create to ensure w is a valid io.Writer (output file writer)
	outFile, err := os.Create(mergedFile.Name()) // Create the output file for the merged PDF
	if err != nil {
		return nil, errors.Wrap(err, "failed to create output file for merged PDF")
	}
	defer outFile.Close() // Ensure the file is closed after the merge

	// Call the pdfcpu.Merge function to merge the PDFs
	err = api.Merge(mergedFile.Name(), inputFiles, outFile, nil, false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to merge PDFs")
	}

	// Read the merged PDF into memory and return it as []byte
	mergedPDF, err := os.ReadFile(mergedFile.Name())
	if err != nil {
		return nil, errors.Wrap(err, "failed to read merged PDF file")
	}

	return mergedPDF, nil
}
