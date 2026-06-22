package kardex

import "github.com/jung-kurt/gofpdf"

// =============================================================================
// TABLE STYLES — Sistema de estilos visuales para las tablas PDF
// =============================================================================

// ColorRGB representa un color en formato RGB (0-255).
type ColorRGB struct {
	R, G, B int
}

// ApplySetFillColor aplica el color como fill del PDF.
func (c ColorRGB) fill(pdf *gofpdf.Fpdf) {
	pdf.SetFillColor(c.R, c.G, c.B)
}

// ApplySetDrawColor aplica el color como draw (bordes) del PDF.
func (c ColorRGB) draw(pdf *gofpdf.Fpdf) {
	pdf.SetDrawColor(c.R, c.G, c.B)
}

// ApplySetTextColor aplica el color como texto del PDF.
func (c ColorRGB) text(pdf *gofpdf.Fpdf) {
	pdf.SetTextColor(c.R, c.G, c.B)
}

// Colores predefinidos
var (
	// Blancos y negros
	colorBlack       = ColorRGB{0, 0, 0}
	colorWhite       = ColorRGB{255, 255, 255}
	colorGrayDark    = ColorRGB{80, 80, 80}
	colorGrayMedium  = ColorRGB{160, 160, 160}
	colorGrayLight   = ColorRGB{204, 204, 204}
	colorGrayBG      = ColorRGB{238, 238, 238}
	colorGrayAlt     = ColorRGB{248, 248, 248}

	// Azul corporativo
	colorBlueDark    = ColorRGB{26, 54, 93}     // #1a365d
	colorBlueMedium  = ColorRGB{44, 82, 130}    // #2c5282
	colorBlueLight   = ColorRGB{235, 245, 255}  // #ebf5ff
	colorBlueHeader  = ColorRGB{232, 240, 254}  // #e8f0fe
	colorBlueBorder  = ColorRGB{160, 190, 220}  // #a0bedc

	// Colores semánticos
	colorGreenLight  = ColorRGB{240, 255, 240}  // #f0fff0 — saldo final
	colorRedSoft     = ColorRGB{197, 48, 48}    // #c53030 — salidas
)

// TableStyle define el estilo visual completo de una tabla.
type TableStyle struct {
	Name string

	// Header
	HeaderBg       ColorRGB // Fondo del header
	HeaderText     ColorRGB // Color del texto del header
	HeaderBorder   ColorRGB // Color del borde del header
	HeaderBorderW  float64  // Grosor del borde del header

	// Datos
	DataBorder     ColorRGB // Color del borde de las celdas de datos
	DataBorderW    float64  // Grosor del borde de datos
	DataAltBg      ColorRGB // Fondo alternado (filas pares) — mismo que DataBg si no hay zebra
	DataText       ColorRGB // Color del texto de datos
	DataTextSalida ColorRGB // Color del texto para salidas

	// Filas especiales
	SaldoInitBg    ColorRGB // Fondo saldo inicial
	SaldoFinalBg   ColorRGB // Fondo saldo final
	TotalesBg      ColorRGB // Fondo fila de totales
	TotalesText    ColorRGB // Color texto totales

	// Separadores entre productos
	SeparatorColor ColorRGB // Color de la línea separadora
	SeparatorWidth float64  // Grosor de la línea separadora
}

// ApplyHeader aplica los colores del header al PDF.
func (s *TableStyle) ApplyHeader(pdf *gofpdf.Fpdf) {
	s.HeaderBg.fill(pdf)
	s.HeaderText.text(pdf)
	s.HeaderBorder.draw(pdf)
	pdf.SetLineWidth(s.HeaderBorderW)
}

// ApplyData aplica los colores de datos normales al PDF.
func (s *TableStyle) ApplyData(pdf *gofpdf.Fpdf, rowIndex int) {
	s.DataText.text(pdf)
	s.DataBorder.draw(pdf)
	pdf.SetLineWidth(s.DataBorderW)

	// Zebra striping
	if rowIndex%2 == 1 {
		s.DataAltBg.fill(pdf)
	} else {
		colorWhite.fill(pdf)
	}
}

// ApplySalida aplica estilo para filas de salida.
func (s *TableStyle) ApplySalida(pdf *gofpdf.Fpdf, rowIndex int) {
	s.DataTextSalida.text(pdf)
	s.DataBorder.draw(pdf)
	pdf.SetLineWidth(s.DataBorderW)

	if rowIndex%2 == 1 {
		s.DataAltBg.fill(pdf)
	} else {
		colorWhite.fill(pdf)
	}
}

// ApplySaldoInit aplica estilo para la fila de saldo inicial.
func (s *TableStyle) ApplySaldoInit(pdf *gofpdf.Fpdf) {
	s.SaldoInitBg.fill(pdf)
	s.DataBorder.draw(pdf)
	pdf.SetLineWidth(s.DataBorderW)
}

// ApplySaldoFinal aplica estilo para la fila de saldo final.
func (s *TableStyle) ApplySaldoFinal(pdf *gofpdf.Fpdf) {
	s.SaldoFinalBg.fill(pdf)
	s.DataBorder.draw(pdf)
	pdf.SetLineWidth(s.DataBorderW)
}

// ApplyTotales aplica estilo para la fila de totales.
func (s *TableStyle) ApplyTotales(pdf *gofpdf.Fpdf) {
	s.TotalesBg.fill(pdf)
	s.TotalesText.text(pdf)
	s.DataBorder.draw(pdf)
	pdf.SetLineWidth(s.DataBorderW)
}

// ApplySeparator aplica estilo para la línea separadora entre productos.
func (s *TableStyle) ApplySeparator(pdf *gofpdf.Fpdf) {
	if s.SeparatorWidth <= 0 {
		return
	}
	s.SeparatorColor.draw(pdf)
	pdf.SetLineWidth(s.SeparatorWidth)
}

// HeaderBorderFlag retorna "1" si los bordes del header deben dibujarse, "0" si no.
// Se usa como parámetro border en pdf.CellFormat — gofpdf interpreta SetLineWidth(0)
// como hairline de 1px, por eso hay que omitir el flag para no dibujar nada.
func (s *TableStyle) HeaderBorderFlag() string {
	if s.HeaderBorderW > 0 {
		return "1"
	}
	return "0"
}

// DataBorderFlag retorna "1" si los bordes de datos deben dibujarse, "0" si no.
func (s *TableStyle) DataBorderFlag() string {
	if s.DataBorderW > 0 {
		return "1"
	}
	return "0"
}

// DataRectFlag retorna "FD" (fill+draw) si los bordes deben dibujarse, "F" (solo fill) si no.
// Se usa en pdf.Rect — gofpdf dibuja el borde siempre con "FD" aunque SetLineWidth sea 0.
func (s *TableStyle) DataRectFlag() string {
	if s.DataBorderW > 0 {
		return "FD"
	}
	return "F"
}

// =============================================================================
// ESTILOS PREDEFINIDOS
// =============================================================================

// StyleSimple es el estilo minimalista — sin bordes en celdas, sin colores de fondo,
// solo una línea separadora entre productos.
var StyleSimple = TableStyle{
	Name: "simple",

	HeaderBg:       colorWhite,
	HeaderText:     colorBlack,
	HeaderBorder:   colorBlack,
	HeaderBorderW:  0,

	DataBorder:     colorBlack,
	DataBorderW:    0,
	DataAltBg:      colorWhite,
	DataText:       colorBlack,
	DataTextSalida: colorBlack,

	SaldoInitBg:    colorWhite,
	SaldoFinalBg:   colorWhite,
	TotalesBg:      colorWhite,
	TotalesText:    colorBlack,

	SeparatorColor: colorBlack,
	SeparatorWidth: 0.15,
}

// StyleFormal es el estilo corporativo — azul oscuro, zebra, colores semánticos.
var StyleFormal = TableStyle{
	Name: "formal",

	HeaderBg:       colorBlueDark,
	HeaderText:     colorWhite,
	HeaderBorder:   colorBlueMedium,
	HeaderBorderW:  0.2,

	DataBorder:     colorBlueBorder,
	DataBorderW:    0.1,
	DataAltBg:      colorGrayAlt,
	DataText:       colorBlack,
	DataTextSalida: colorRedSoft,

	SaldoInitBg:    colorBlueLight,
	SaldoFinalBg:   colorGreenLight,
	TotalesBg:      colorBlueHeader,
	TotalesText:    colorBlueDark,

	SeparatorColor: colorBlueMedium,
	SeparatorWidth: 0.4,
}

// StyleContador es el estilo para impresión B/N — alto contraste, mínima tinta.
// Optimizado para impresoras láser monocromáticas de contadores.
var StyleContador = TableStyle{
	Name: "contador",

	HeaderBg:       colorBlack,
	HeaderText:     colorWhite,
	HeaderBorder:   colorBlack,
	HeaderBorderW:  0.3,

	DataBorder:     colorGrayLight,
	DataBorderW:    0.1,
	DataAltBg:      ColorRGB{245, 245, 245}, // zebra muy sutil
	DataText:       colorBlack,
	DataTextSalida: colorGrayDark, // gris oscuro en vez de rojo — B/N friendly

	SaldoInitBg:    ColorRGB{230, 230, 230}, // gris claro
	SaldoFinalBg:   ColorRGB{220, 220, 220}, // gris un poco más oscuro
	TotalesBg:      ColorRGB{210, 210, 210}, // gris medio
	TotalesText:    colorBlack,

	SeparatorColor: colorBlack,
	SeparatorWidth: 0.3,
}

// GetStyleByName retorna el estilo por nombre. Default: simple.
func GetStyleByName(name string) TableStyle {
	switch name {
	case "formal":
		return StyleFormal
	case "contador":
		return StyleContador
	default:
		return StyleSimple
	}
}
