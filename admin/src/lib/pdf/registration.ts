import QRCode from "qrcode";
import { PDFDocument, rgb, StandardFonts } from "pdf-lib";

type Options = {
	code: string;
	qrUrl: string; // what the QR encodes (often same as a deep link)
	title?: string;
	subtitle?: string;
	instructions?: string[];
	footer?: string;
};

export async function buildRegistrationPdf({
											   code,
											   qrUrl,
											   title = "Registration Code",
											   subtitle = "Scan the QR code to register",
											   instructions = [
												   "Open your camera or QR scanner.",
												   "Scan the code above.",
												   "Follow the link to complete registration."
											   ],
											   footer = "If you need help, contact support."
										   }: Options) {
	const pdf = await PDFDocument.create();
	const page = pdf.addPage([595.28, 841.89]); // A4 in points (72dpi): 210×297mm
	const { width, height } = page.getSize();

	const font = await pdf.embedFont(StandardFonts.Helvetica);
	const fontBold = await pdf.embedFont(StandardFonts.HelveticaBold);

	// Generate a high-res QR PNG (simple + reliable)
	// scale controls pixel size; errorCorrectionLevel improves scan reliability
	const qrPngDataUrl = await QRCode.toDataURL(qrUrl, {
		errorCorrectionLevel: "M",
		margin: 1,
		scale: 12
	});

	const qrBytes = dataUrlToBytes(qrPngDataUrl);
	const qrImage = await pdf.embedPng(qrBytes);

	// Layout constants
	const margin = 56; // ~20mm
	const cardW = width - margin * 2;
	const cardH = height - margin * 2;

	// Card background (subtle)
	page.drawRectangle({
		x: margin,
		y: margin,
		width: cardW,
		height: cardH,
		color: rgb(1, 0.98, 0.995), // very soft pink-white
		borderColor: rgb(0.92, 0.7, 0.82),
		borderWidth: 1.2
	});

	// Header
	page.drawText(title, {
		x: margin + 28,
		y: height - margin - 54,
		size: 28,
		font: fontBold,
		color: rgb(0.2, 0.12, 0.18)
	});

	page.drawText(subtitle, {
		x: margin + 28,
		y: height - margin - 84,
		size: 13,
		font,
		color: rgb(0.35, 0.25, 0.32)
	});

	// QR placement (big + centered-ish)
	const qrSize = 320; // points
	const qrX = margin + (cardW - qrSize) / 2;
	const qrY = margin + cardH - 84 - qrSize - 48;

	page.drawImage(qrImage, { x: qrX, y: qrY, width: qrSize, height: qrSize });

	// Code pill
	const pillW = 420;
	const pillH = 44;
	const pillX = margin + (cardW - pillW) / 2;
	const pillY = qrY - 66;

	page.drawRectangle({
		x: pillX,
		y: pillY,
		width: pillW,
		height: pillH,
		color: rgb(1, 0.92, 0.96),
		borderColor: rgb(0.92, 0.6, 0.78),
		borderWidth: 1
	});

	const codeLabel = `Code: ${code}`;
	const codeTextWidth = fontBold.widthOfTextAtSize(codeLabel, 16);
	page.drawText(codeLabel, {
		x: pillX + (pillW - codeTextWidth) / 2,
		y: pillY + 14,
		size: 16,
		font: fontBold,
		color: rgb(0.22, 0.1, 0.18)
	});

	// Instructions
	let y = pillY - 36;
	page.drawText("Instructions:", {
		x: margin + 28,
		y,
		size: 14,
		font: fontBold,
		color: rgb(0.2, 0.12, 0.18)
	});

	y -= 18;
	for (const [i, line] of instructions.entries()) {
		page.drawText(`${i + 1}. ${line}`, {
			x: margin + 36,
			y,
			size: 12,
			font,
			color: rgb(0.35, 0.25, 0.32)
		});
		y -= 16;
	}

	// Footer
	page.drawText(footer, {
		x: margin + 28,
		y: margin + 26,
		size: 10.5,
		font,
		color: rgb(0.45, 0.35, 0.42)
	});

	const bytes = await pdf.save();
	return new Blob([bytes], { type: "application/pdf" });
}

function dataUrlToBytes(dataUrl: string) {
	const base64 = dataUrl.split(",")[1]!;
	const bin = atob(base64);
	const bytes = new Uint8Array(bin.length);
	for (let i = 0; i < bin.length; i++) bytes[i] = bin.charCodeAt(i);
	return bytes;
}
