package service

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"rest-app/internal/app/ocr/model"
	"rest-app/internal/app/ocr/port"

	"github.com/otiai10/gosseract/v2"
	"gocv.io/x/gocv"
)

type ocr struct {
	OCRClient *gosseract.Client
	//HuggingFaceRepo port.IHuggingFaceHTTP
	GoogleAIRepo port.IGoogleAIHTTP
}

func NewOCRService(OCRClient *gosseract.Client, GoogleAIRepo port.IGoogleAIHTTP) port.IOCRService {
	return &ocr{
		OCRClient: OCRClient,
		//HuggingFaceRepo: HuggingFaceRepo,
		GoogleAIRepo: GoogleAIRepo,
	}
}

func (o *ocr) ReceiptDataGenerator(ctx context.Context, imgBytes []byte) (*model.ReceiptTransaction, error) {

	defer o.OCRClient.Close()

	var (
		receiptData         model.ReceiptTransaction
		optimizedImageBytes []byte
		err                 error
	)

	optimizedImageBytes, err = o.optimizeImageFromBytes(imgBytes)

	if err != nil {
		return nil, fmt.Errorf("failed to optimize image: %w", err)
	}

	// Set the optimized image for OCR
	err = o.OCRClient.SetImageFromBytes(optimizedImageBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to set optimized image: %w", err)
	}

	text, err := o.OCRClient.Text()
	if err != nil {
		return nil, fmt.Errorf("OCR processing failed: %w", err)
	}

	fmt.Println("OCR Result:", text)

	// Parse generated text from OCR using AI for JSON Result
	resByte, err := o.GoogleAIRepo.ProceedTxtToJSONGeneratorPrompt(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("AI Text processing failed: %w", err)
	}

	err = json.Unmarshal(resByte, &receiptData)
	if err != nil {
		return nil, fmt.Errorf("parsing AI Result Text failed: %w", err)
	}

	return &receiptData, nil
}

// optimizeImageFromBytes loads image from byte array and applies preprocessing
func (o *ocr) optimizeImageFromBytes(imageBytes []byte) ([]byte, error) {
	// Decode image from bytes
	img, err := gocv.IMDecode(imageBytes, gocv.IMReadColor)
	if err != nil {
		return nil, fmt.Errorf("unable to decode image from bytes: %w", err)
	}
	if img.Empty() {
		return nil, fmt.Errorf("decoded image is empty")
	}
	defer img.Close()

	return o.preprocessImage(img)
}

// preprocessImage applies various image optimization techniques for better OCR
func (o *ocr) preprocessImage(src gocv.Mat) ([]byte, error) {
	// Create working matrices
	gray := gocv.NewMat()
	blurred := gocv.NewMat()
	thresh := gocv.NewMat()
	morphed := gocv.NewMat()
	denoised := gocv.NewMat()
	resized := gocv.NewMat()

	defer gray.Close()
	defer blurred.Close()
	defer thresh.Close()
	defer morphed.Close()
	defer denoised.Close()
	defer resized.Close()

	// Step 1: Resize image if too large (improve processing speed and sometimes accuracy)
	if src.Cols() > 2000 || src.Rows() > 2000 {
		newWidth := src.Cols()
		newHeight := src.Rows()

		// Scale down while maintaining aspect ratio
		if src.Cols() > src.Rows() {
			newWidth = 2000
			newHeight = int(float64(src.Rows()) * (2000.0 / float64(src.Cols())))
		} else {
			newHeight = 2000
			newWidth = int(float64(src.Cols()) * (2000.0 / float64(src.Rows())))
		}

		gocv.Resize(src, &resized, image.Pt(newWidth, newHeight), 0, 0, gocv.InterpolationLinear)
		src = resized.Clone() // Use resized image for further processing
		defer src.Close()
	}

	// Step 2: Convert to grayscale
	gocv.CvtColor(src, &gray, gocv.ColorBGRToGray)

	// Step 3: Apply Gaussian blur to reduce noise
	gocv.GaussianBlur(gray, &blurred, image.Pt(3, 3), 0, 0, gocv.BorderDefault)

	// Step 4: Apply adaptive threshold for better text extraction
	// This works better than simple threshold for receipts with varying lighting
	gocv.AdaptiveThreshold(blurred, &thresh, 255, gocv.AdaptiveThresholdMean, gocv.ThresholdBinary, 11, 2)

	// Step 5: Morphological operations to clean up the image
	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(2, 2))
	defer kernel.Close()

	// Opening operation (erosion followed by dilation) to remove noise
	gocv.MorphologyEx(thresh, &morphed, gocv.MorphOpen, kernel)

	// Step 6: Optional - Apply median blur for additional noise reduction
	gocv.MedianBlur(morphed, &denoised, 3)

	// Step 7: Encode the processed image back to bytes
	buf, err := gocv.IMEncode(".png", denoised) // Use PNG for lossless compression
	if err != nil {
		return nil, fmt.Errorf("failed to encode processed image: %w", err)
	}

	return buf.GetBytes(), nil
}
