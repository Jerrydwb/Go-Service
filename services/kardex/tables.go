package kardex

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"

	"kardex-pdf-service/internal/shared"
)

// =============================================================================
// CONSTANTS
// =============================================================================

const (
	lineHt     = 2.4  // Altura de cada línea de texto en datos
	headerHt   = 2.8  // Altura de cada línea en headers
	cellGap    = 0.0  // Padding interno de cada celda
	safetyY    = 7.0  // Margen de seguridad para page break (espacio para headers)
)

// =============================================================================
// HELPERS
// =============================================================================

// maxContentY retorna la coordenada Y máxima usable antes de necesitar page break.
// Se calcula dinámicamente según el tamaño de página (portrait vs landscape).
func (g *Generator) maxContentY(pdf *gofpdf.Fpdf) float64 {
	_, pageHeight := pdf.GetPageSize()
	_, _, _, bottomMargin := pdf.GetMargins()
	return pageHeight - bottomMargin - safetyY
}

// wrapText divide un texto en múltiples líneas que caben dentro de maxWidth.
// Usa SplitLines del propio gofpdf (respeta la font actual).
func wrapText(pdf *gofpdf.Fpdf, text string, maxWidth float64) []string {
	if text == "" {
		return []string{""}
	}
	lines := pdf.SplitLines([]byte(text), maxWidth-cellGap-cellGap)
	if len(lines) == 0 {
		return []string{""}
	}
	result := make([]string, len(lines))
	for i, l := range lines {
		result[i] = string(l)
	}
	return result
}

// rowHeight calcula la altura de una fila dadas las líneas de texto
// que debe contener la celda más alta.
func rowHeight(lineCount int) float64 {
	return float64(lineCount)*lineHt + cellGap*2
}

// =============================================================================
// PRODUCT SUMMARY TABLE
// =============================================================================

// addProductSummaryPage agrega la tabla de resumen de productos al PDF actual.
func (g *Generator) addProductSummaryPage(pdf *gofpdf.Fpdf, insumos []KardexInsumo) {
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(0, 6, "RESUMEN DE PRODUCTOS", "", 1, "C", false, 0, "")
	pdf.Ln(2)

	_, _, colWidths := g.summaryColWidths(pdf)
	maxY := g.maxContentY(pdf)

	// Headers
	g.renderSummaryHeaders(pdf, colWidths)

	// Datos
	pdf.SetFont("Arial", "", 5.5)
	for i, insumo := range insumos {
		initial := insumo.Kardex.Initial
		finish := insumo.Kardex.Finish

		denomination := "Sin denominacion"
		barcode := "S/C"
		if insumo.InsumoDetail.Insumo != nil {
			if insumo.InsumoDetail.Insumo.Denomination != "" {
				denomination = shared.ToLatin(insumo.InsumoDetail.Insumo.Denomination)
			}
			if insumo.InsumoDetail.Insumo.Barcode != "" {
				barcode = shared.ToLatin(insumo.InsumoDetail.Insumo.Barcode)
			}
		}

		hasNegative := false
		var initialQty, initialPrice, initialTotal float64
		var finalQty, finalPrice, finalTotal float64

		if initial != nil {
			initialQty = initial.Cantidad
			initialPrice = initial.Precio
			initialTotal = initial.Parcial
			if initialQty < 0 || initialPrice < 0 {
				hasNegative = true
			}
		}
		if finish != nil {
			finalQty = finish.Cantidad
			finalPrice = finish.Precio
			finalTotal = finish.Parcial
			if finalQty < 0 || finalPrice < 0 {
				hasNegative = true
			}
		}

		// Pre-calcular líneas de denominación para saber la altura de la fila
		pdf.SetFont("Arial", "", 5.5)
		denomLines := wrapText(pdf, denomination, colWidths[2])
		rh := rowHeight(len(denomLines))

		// Page break dinámico
		if pdf.GetY()+rh > maxY {
			pdf.AddPage()
			g.renderSummaryHeaders(pdf, colWidths)
			pdf.SetFont("Arial", "", 5.5)
		}

		// Aplicar estilo de fila (zebra)
		g.style.ApplyData(pdf, i)

		// Color para valores negativos
		if hasNegative {
			colorRedSoft.text(pdf)
		}

		y := pdf.GetY()
		leftMargin, _, _, _ := pdf.GetMargins()
		x := leftMargin

		// Dibujar bordes de todas las celdas con la altura de la fila
		cx := x
		rectFlag := g.style.DataRectFlag()
		for j := 0; j < 9; j++ {
			pdf.Rect(cx, y, colWidths[j], rh, rectFlag)
			cx += colWidths[j]
		}

		// Contenido de celdas de una sola línea (centrado verticalmente)
		singleLineY := y + cellGap + (rh-lineHt)/2

		// Nro
		pdf.SetXY(x+cellGap, singleLineY)
		pdf.CellFormat(colWidths[0]-cellGap*2, lineHt, fmt.Sprintf("%d", i+1), "", 0, "C", false, 0, "")

		// Codigo
		pdf.SetXY(x+colWidths[0]+cellGap, singleLineY)
		pdf.CellFormat(colWidths[1]-cellGap*2, lineHt, barcode, "", 0, "C", false, 0, "")

		// Denominación (multi-línea)
		denomX := x + colWidths[0] + colWidths[1] + cellGap
		denomLineY := y + cellGap
		for _, line := range denomLines {
			pdf.SetXY(denomX, denomLineY)
			pdf.CellFormat(colWidths[2]-cellGap*2, lineHt, line, "", 0, "L", false, 0, "")
			denomLineY += lineHt
		}

		// Saldo Inicial (cols 3,4,5)
		siX := x + colWidths[0] + colWidths[1] + colWidths[2]
		pdf.SetXY(siX+cellGap, singleLineY)
		pdf.CellFormat(colWidths[3]-cellGap*2, lineHt, shared.FormatNumber(initialQty), "", 0, "R", false, 0, "")
		pdf.SetXY(siX+colWidths[3]+cellGap, singleLineY)
		pdf.CellFormat(colWidths[4]-cellGap*2, lineHt, shared.FormatNumber(initialPrice), "", 0, "R", false, 0, "")
		pdf.SetXY(siX+colWidths[3]+colWidths[4]+cellGap, singleLineY)
		pdf.CellFormat(colWidths[5]-cellGap*2, lineHt, shared.FormatNumber(initialTotal), "", 0, "R", false, 0, "")

		// Saldo Final (cols 6,7,8)
		sfX := siX + colWidths[3] + colWidths[4] + colWidths[5]
		pdf.SetXY(sfX+cellGap, singleLineY)
		pdf.CellFormat(colWidths[6]-cellGap*2, lineHt, shared.FormatNumber(finalQty), "", 0, "R", false, 0, "")
		pdf.SetXY(sfX+colWidths[6]+cellGap, singleLineY)
		pdf.CellFormat(colWidths[7]-cellGap*2, lineHt, shared.FormatNumber(finalPrice), "", 0, "R", false, 0, "")
		pdf.SetXY(sfX+colWidths[6]+colWidths[7]+cellGap, singleLineY)
		pdf.CellFormat(colWidths[8]-cellGap*2, lineHt, shared.FormatNumber(finalTotal), "", 0, "R", false, 0, "")

		pdf.SetXY(leftMargin, y+rh)
	}

	colorBlack.text(pdf)
}

func (g *Generator) renderSummaryHeaders(pdf *gofpdf.Fpdf, colWidths []float64) {
	pdf.SetFont("Arial", "B", 6)
	g.style.ApplyHeader(pdf)
	hbf := g.style.HeaderBorderFlag()

	// Fila 1
	pdf.CellFormat(colWidths[0], headerHt, "Nro", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[1], headerHt, "Codigo", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[2], headerHt, "Producto / Insumo", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[3]+colWidths[4]+colWidths[5], headerHt, "Saldo Inicial", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[6]+colWidths[7]+colWidths[8], headerHt, "Saldo Final", hbf, 0, "C", true, 0, "")
	pdf.Ln(-1)

	// Fila 2
	pdf.CellFormat(colWidths[0], headerHt, "", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[1], headerHt, "", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[2], headerHt, "", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[3], headerHt, "Cant.", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[4], headerHt, "Precio", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[5], headerHt, "Total", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[6], headerHt, "Cant.", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[7], headerHt, "Precio", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[8], headerHt, "Total", hbf, 0, "C", true, 0, "")
	pdf.Ln(-1)
}

// =============================================================================
// INSUMO DETAIL PAGE
// =============================================================================

// addInsumoPage agrega la tabla de detalle de un insumo al PDF actual.
func (g *Generator) addInsumoPage(pdf *gofpdf.Fpdf, insumo *KardexInsumo, index, total int) {
	denomination := "Sin denominacion"
	barcode := "S/C"
	if insumo.InsumoDetail.Insumo != nil {
		if insumo.InsumoDetail.Insumo.Denomination != "" {
			denomination = shared.ToLatin(insumo.InsumoDetail.Insumo.Denomination)
		}
		if insumo.InsumoDetail.Insumo.Barcode != "" {
			barcode = shared.ToLatin(insumo.InsumoDetail.Insumo.Barcode)
		}
	}

	title := fmt.Sprintf("%d/%d: %s (%s)", index, total, denomination, barcode)

	// Separador entre productos
	g.style.ApplySeparator(pdf)
	leftMargin, _, _, _ := pdf.GetMargins()
	pageWidth, _ := pdf.GetPageSize()
	pdf.Line(leftMargin, pdf.GetY(), pageWidth-leftMargin, pdf.GetY())
	pdf.Ln(1)

	pdf.SetFont("Arial", "B", 6.5)
	colorBlack.text(pdf)
	pdf.CellFormat(0, 5, title, "", 1, "L", false, 0, "")
	pdf.Ln(1)

	// Filtrar proformas
	validMovements := make([]KardexMovimiento, 0, len(insumo.Kardex.Movimientos))
	for _, m := range insumo.Kardex.Movimientos {
		if m.Asunto != "proforma" {
			validMovements = append(validMovements, m)
		}
	}

	g.buildKardexTable(pdf, validMovements, insumo.Kardex.Initial, insumo.Kardex.Finish)
}

// =============================================================================
// KARDEX TABLE
// =============================================================================

func (g *Generator) buildKardexTable(pdf *gofpdf.Fpdf, movements []KardexMovimiento, initial *KardexInitial, finish *KardexFinish) {
	colWidths := g.kardexColWidths(pdf)
	maxY := g.maxContentY(pdf)
	rh := rowHeight(1) // Todas las filas de movimientos son 1 línea (datos numéricos cortos)

	// Headers - Fila 1
	g.renderKardexHeaderRow1(pdf, colWidths)
	pdf.Ln(-1)

	// Headers - Fila 2
	g.renderKardexHeaderRow2(pdf, colWidths)
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 5.5)

	// Fila de saldo inicial
	if initial != nil {
		if pdf.GetY()+rh > maxY {
			pdf.AddPage()
			g.renderKardexHeaderRow1(pdf, colWidths)
			pdf.Ln(-1)
			g.renderKardexHeaderRow2(pdf, colWidths)
			pdf.Ln(-1)
			pdf.SetFont("Arial", "", 5.5)
		}
		g.style.ApplySaldoInit(pdf)
		g.renderKardexDataRowContent(pdf, colWidths, "", "", "", "", "", "S. INICIAL", true,
			"", "", "", "", "", "",
			shared.FormatNumber(initial.Cantidad), shared.FormatNumber(initial.Precio), shared.FormatNumber(initial.Parcial))
	}

	// Movimientos
	var totals struct {
		cantidadEntradas float64
		costoEntradas    float64
		cantidadSalidas  float64
		costoSalidas     float64
	}

	for idx, m := range movements {
		// Page break dinámico
		if pdf.GetY()+rh > maxY {
			pdf.AddPage()
			g.renderKardexHeaderRow1(pdf, colWidths)
			pdf.Ln(-1)
			g.renderKardexHeaderRow2(pdf, colWidths)
			pdf.Ln(-1)
			pdf.SetFont("Arial", "", 5.5)
		}

		isSalida := m.Asunto == "salida"
		tipoOp := shared.ToLatin(fmt.Sprintf("%s-VEN", m.TipoDocumento))

		var entCant, entPrecio, entTotal string
		var salCant, salPrecio, salTotal string

		if !isSalida {
			cantIngreso := shared.GetFloatValue(m.CantidadIngreso)
			precioIngreso := shared.GetFloatValue(m.PrecioIngreso)
			parcialIngreso := shared.GetFloatValue(m.ParcialIngreso)
			entCant = shared.FormatNumber(cantIngreso)
			entPrecio = shared.FormatNumber(precioIngreso)
			entTotal = shared.FormatNumber(parcialIngreso)
			totals.cantidadEntradas += cantIngreso
			totals.costoEntradas += parcialIngreso
		}

		if isSalida {
			cantSalida := shared.GetFloatValue(m.CantidadSalida)
			precioSalida := shared.GetFloatValue(m.PrecioSalida)
			parcialSalida := shared.GetFloatValue(m.ParcialSalida)
			salCant = shared.FormatNumber(cantSalida)
			salPrecio = shared.FormatNumber(precioSalida)
			salTotal = shared.FormatNumber(parcialSalida)
			totals.cantidadSalidas += cantSalida
			totals.costoSalidas += parcialSalida
		}

		// Aplicar estilo según tipo de fila (salida vs entrada, con zebra)
		if isSalida {
			g.style.ApplySalida(pdf, idx)
		} else {
			g.style.ApplyData(pdf, idx)
		}

		g.renderKardexDataRowContent(pdf, colWidths,
			fmt.Sprintf("%d", idx+1),
			shared.ToLatin(m.FechaKardex),
			shared.ToLatin(m.TipoDocumento),
			shared.ToLatin(m.SerieDocumento),
			shared.ToLatin(m.NumeroDocumento),
			tipoOp,
			false,
			entCant, entPrecio, entTotal,
			salCant, salPrecio, salTotal,
			shared.FormatNumber(shared.GetFloatValue(m.CantidadFinal)),
			shared.FormatNumber(shared.GetFloatValue(m.PrecioFinal)),
			shared.FormatNumber(shared.GetFloatValue(m.CostoFinal)),
		)
	}

	// Fila de saldo final
	if finish != nil {
		if pdf.GetY()+rh > maxY {
			pdf.AddPage()
			g.renderKardexHeaderRow1(pdf, colWidths)
			pdf.Ln(-1)
			g.renderKardexHeaderRow2(pdf, colWidths)
			pdf.Ln(-1)
			pdf.SetFont("Arial", "", 5.5)
		}
		g.style.ApplySaldoFinal(pdf)
		g.renderKardexDataRowContent(pdf, colWidths, "", "", "", "", "", "S. FINAL", true,
			"", "", "", "", "", "",
			shared.FormatNumber(finish.Cantidad), shared.FormatNumber(finish.Precio), shared.FormatNumber(finish.Parcial))
	}

	// Fila de totales — COSTO DE VENTAS por ecuación contable: Ini + Compras - Fin
	if pdf.GetY()+rh > maxY {
		pdf.AddPage()
		g.renderKardexHeaderRow1(pdf, colWidths)
		pdf.Ln(-1)
		g.renderKardexHeaderRow2(pdf, colWidths)
		pdf.Ln(-1)
	}
	var costoNeto float64
	if initial != nil && finish != nil {
		costoNeto = initial.Parcial + totals.costoEntradas - finish.Parcial
	} else {
		costoNeto = totals.costoEntradas - totals.costoSalidas
	}
	pdf.SetFont("Arial", "B", 5.5)
	g.style.ApplyTotales(pdf)
	dbf := g.style.DataBorderFlag()
	pdf.CellFormat(colWidths[0]+colWidths[1]+colWidths[2]+colWidths[3]+colWidths[4], lineHt, "", dbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[5], lineHt, "TOTALES:", dbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[6], lineHt, shared.FormatNumber(totals.cantidadEntradas), dbf, 0, "R", true, 0, "")
	pdf.CellFormat(colWidths[7], lineHt, "", dbf, 0, "R", true, 0, "")
	pdf.CellFormat(colWidths[8], lineHt, shared.FormatNumber(totals.costoEntradas), dbf, 0, "R", true, 0, "")
	pdf.CellFormat(colWidths[9], lineHt, shared.FormatNumber(totals.cantidadSalidas), dbf, 0, "R", true, 0, "")
	pdf.CellFormat(colWidths[10], lineHt, "", dbf, 0, "R", true, 0, "")
	pdf.CellFormat(colWidths[11], lineHt, shared.FormatNumber(totals.costoSalidas), dbf, 0, "R", true, 0, "")
	pdf.CellFormat(colWidths[12], lineHt, "", dbf, 0, "R", true, 0, "")
	pdf.CellFormat(colWidths[13], lineHt, "", dbf, 0, "R", true, 0, "")
	pdf.CellFormat(colWidths[14], lineHt, shared.FormatNumber(costoNeto), dbf, 0, "R", true, 0, "")
	pdf.Ln(-1)
}

// renderKardexDataRowContent dibuja una fila completa de la tabla kardex con bordes.
// Los strings vacíos generan celdas vacías. Si boldLabel != "" se usa font bold.
// NOTA: El color/fondo debe aplicarse ANTES de llamar este método.
func (g *Generator) renderKardexDataRowContent(pdf *gofpdf.Fpdf, cw []float64,
	nro, fecha, tipoDoc, serie, numero, tipoOp string, boldLabel bool,
	entCant, entPrecio, entTotal string,
	salCant, salPrecio, salTotal string,
	saldoCant, saldoPrecio, saldoTotal string,
) {
	if boldLabel {
		pdf.SetFont("Arial", "B", 5.5)
	} else {
		pdf.SetFont("Arial", "", 5.5)
	}

	leftMargin, _, _, _ := pdf.GetMargins()
	x := leftMargin
	y := pdf.GetY()
	rh := rowHeight(1)

	// Columnas: [0]Nro [1]Fecha [2]TipoDoc [3]Serie [4]Numero [5]TipoOp
	// [6]EntCant [7]EntPrecio [8]EntTotal [9]SalCant [10]SalPrecio [11]SalTotal
	// [12]SaldoCant [13]SaldoPrecio [14]SaldoTotal
	texts := []string{nro, fecha, tipoDoc, serie, numero, tipoOp,
		entCant, entPrecio, entTotal, salCant, salPrecio, salTotal,
		saldoCant, saldoPrecio, saldoTotal}
	aligns := []string{"C", "C", "C", "C", "C", "C",
		"R", "R", "R", "R", "R", "R",
		"R", "R", "R"}

	// Dibujar bordes y texto
	cx := x
	rectFlag := g.style.DataRectFlag()
	for i, txt := range texts {
		pdf.Rect(cx, y, cw[i], rh, rectFlag)
		if txt != "" {
			pdf.SetXY(cx+cellGap, y+cellGap)
			pdf.CellFormat(cw[i]-cellGap*2, lineHt, txt, "", 0, aligns[i], false, 0, "")
		}
		cx += cw[i]
	}

	pdf.SetXY(x, y+rh)
}

// =============================================================================
// HEADER HELPERS (DRY)
// =============================================================================

func (g *Generator) renderKardexHeaderRow1(pdf *gofpdf.Fpdf, colWidths []float64) {
	pdf.SetFont("Arial", "B", 5.5)
	g.style.ApplyHeader(pdf)
	hbf := g.style.HeaderBorderFlag()

	pdf.CellFormat(colWidths[0], headerHt, "Nro", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[1], headerHt, "Fecha", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[2]+colWidths[3]+colWidths[4], headerHt, "Comprobante", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[5], headerHt, "Tipo Op.", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[6]+colWidths[7]+colWidths[8], headerHt, "Entradas", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[9]+colWidths[10]+colWidths[11], headerHt, "Salidas", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[12]+colWidths[13]+colWidths[14], headerHt, "Saldo Final", hbf, 0, "C", true, 0, "")
}

func (g *Generator) renderKardexHeaderRow2(pdf *gofpdf.Fpdf, colWidths []float64) {
	pdf.SetFont("Arial", "B", 5.5)
	g.style.ApplyHeader(pdf)
	hbf := g.style.HeaderBorderFlag()

	pdf.CellFormat(colWidths[0], headerHt, "", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[1], headerHt, "", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[2], headerHt, "Tipo", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[3], headerHt, "Serie", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[4], headerHt, "Numero", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[5], headerHt, "", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[6], headerHt, "Cant.", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[7], headerHt, "P.Unit.", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[8], headerHt, "Total", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[9], headerHt, "Cant.", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[10], headerHt, "P.Unit.", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[11], headerHt, "Total", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[12], headerHt, "Cant.", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[13], headerHt, "P.Unit.", hbf, 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[14], headerHt, "Total", hbf, 0, "C", true, 0, "")
}

// =============================================================================
// DYNAMIC COLUMN WIDTHS
// =============================================================================

// kardexColWidths calcula los anchos de las 15 columnas adaptándose a la orientación.
func (g *Generator) kardexColWidths(pdf *gofpdf.Fpdf) []float64 {
	leftMargin, _, rightMargin, _ := pdf.GetMargins()
	pageWidth, _ := pdf.GetPageSize()
	available := pageWidth - leftMargin - rightMargin

	ratios := []float64{
		3.3, 6.2, 5.0, 5.0, 6.2, 9.1,
		6.2, 6.2, 7.4,
		6.2, 6.2, 7.4,
		6.2, 6.2, 7.4,
	}

	widths := make([]float64, 15)
	for i, r := range ratios {
		widths[i] = (r / 100.0) * available
	}
	return widths
}

// summaryColWidths calcula los anchos de las 9 columnas del resumen.
func (g *Generator) summaryColWidths(pdf *gofpdf.Fpdf) (fixedWidths []float64, denominationWidth float64, colWidths []float64) {
	leftMargin, _, rightMargin, _ := pdf.GetMargins()
	pageWidth, _ := pdf.GetPageSize()
	available := pageWidth - leftMargin - rightMargin

	fixedRatios := []float64{3.5, 7.1, 5.7, 5.7, 7.1, 5.7, 5.7, 7.1}
	totalFixedRatio := float64(0)
	for _, r := range fixedRatios {
		totalFixedRatio += r
	}

	fixedWidths = make([]float64, 8)
	for i, r := range fixedRatios {
		fixedWidths[i] = (r / 100.0) * available
	}

	denominationWidth = (100.0 - totalFixedRatio) / 100.0 * available

	colWidths = []float64{
		fixedWidths[0], fixedWidths[1], denominationWidth,
		fixedWidths[2], fixedWidths[3], fixedWidths[4],
		fixedWidths[5], fixedWidths[6], fixedWidths[7],
	}
	return
}
