package services

import (
	"fmt"
	"strings"
)

// Condiciones comerciales de la concesionaria.
//
// Este archivo es la fuente ÚNICA de las condiciones de pago y financiación que
// el asistente puede ofrecer dentro de una cotización. La IA las recibe como
// texto y tiene prohibido inventar tasas, plazos, cuotas o descuentos que no
// figuren acá. Para cambiar las condiciones, editar estos datos (no prompts).

// metodoPago describe una forma de pago disponible para cualquier unidad.
type metodoPago struct {
	nombre  string
	detalle string
}

// planFinanciacion describe una alternativa de financiación de un vehículo.
type planFinanciacion struct {
	nombre      string
	adelantoMin float64 // porcentaje mínimo de adelanto sobre el precio (0..1)
	cuotas      int     // cantidad de cuotas del saldo financiado; 0 = no aplica
	tasaMensual float64 // interés mensual simple sobre el saldo (0..1); 0 = sin interés
	nota        string  // aclaración opcional
}

// porcentajeDescuentoContado es el descuento por pago de contado/transferencia.
const porcentajeDescuentoContado = 0.03

// metodosDePago son las formas de pago ofrecidas por la concesionaria.
var metodosDePago = []metodoPago{
	{
		nombre:  "Contado / transferencia bancaria",
		detalle: porcentaje(porcentajeDescuentoContado) + " de descuento sobre el precio publicado",
	},
	{nombre: "Tarjeta de débito", detalle: "precio publicado, sin recargo"},
	{nombre: "Tarjeta de crédito", detalle: "hasta 12 cuotas sin interés"},
}

// planesFinanciacion son los planes de financiación que ofrece la concesionaria
// sobre el precio publicado.
var planesFinanciacion = []planFinanciacion{
	{
		nombre:      "Plan 30/70",
		adelantoMin: 0.30,
		cuotas:      6,
		tasaMensual: 0,
		nota:        "sin interés sobre el saldo",
	},
	{
		nombre:      "12 cuotas sin interés",
		adelantoMin: 0,
		cuotas:      12,
		tasaMensual: 0,
	},
	{
		nombre:      "24 cuotas",
		adelantoMin: 0.40,
		cuotas:      24,
		tasaMensual: 0.012,
		nota:        "interés mensual del 1,2% sobre el saldo",
	},
	{
		nombre:      "36 cuotas",
		adelantoMin: 0.40,
		cuotas:      36,
		tasaMensual: 0.015,
		nota:        "interés mensual del 1,5% sobre el saldo",
	},
	{
		nombre:      "Plan de ahorro (prenda)",
		adelantoMin: 1,
		cuotas:      0,
		nota:        "todas las unidades; cuota mensual según el valor móvil del vehículo, a confirmar con un vendedor",
	},
}

// cuotaEstimada calcula el valor de cada cuota de un plan sobre un precio dado
// con interés simple mensual.
func cuotaEstimada(precio float64, plan planFinanciacion) float64 {
	monto := precio * (1 - plan.adelantoMin)
	interes := monto * plan.tasaMensual * float64(plan.cuotas)
	return (monto + interes) / float64(plan.cuotas)
}

// porcentaje convierte una fracción (0..1) en su representación porcentual
// entera (ej. 0.03 -> "3%").
func porcentaje(fraccion float64) string {
	return fmt.Sprintf("%.0f%%", fraccion*100)
}

// textoCondicionesComerciales arma el bloque de condiciones comerciales que el
// asistente usa para responder dentro de una cotización. Los valores se
// calculan en código sobre el precio del vehículo: la IA nunca genera montos.
func textoCondicionesComerciales(precio float64) string {
	var partes []string

	partes = append(partes, "Métodos de pago:")
	for _, metodo := range metodosDePago {
		partes = append(partes, "- "+metodo.nombre+": "+metodo.detalle)
	}

	partes = append(partes, "Planes de financiación sobre el precio publicado:")
	for _, plan := range planesFinanciacion {
		linea := "- " + plan.nombre + ":"
		if plan.adelantoMin > 0 {
			linea += " " + porcentaje(plan.adelantoMin) + " de adelanto"
		}
		if plan.cuotas > 0 {
			linea += " y el saldo en " + fmt.Sprintf("%d cuotas", plan.cuotas)
			if plan.tasaMensual > 0 {
				linea += " a " + porcentaje(plan.tasaMensual) + " mensual"
			} else {
				linea += " sin interés"
			}
			if precio > 0 {
				linea += " => cuota estimada " + formatearARS(cuotaEstimada(precio, plan))
			}
		} else {
			linea += ", la unidad entera"
		}
		if plan.nota != "" {
			linea += " (" + plan.nota + ")"
		}
		partes = append(partes, linea)
	}

	if precio > 0 {
		descuento := precio * porcentajeDescuentoContado
		partes = append(partes, fmt.Sprintf(
			"Nota: los montos se calculan sobre el precio publicado (%s); pagando de contado/transferencia queda en %s con un descuento de %s. Confirmar condiciones finales con un vendedor.",
			formatearARS(precio), formatearARS(precio-descuento), formatearARS(descuento),
		))
	} else {
		partes = append(partes, "Nota: los montos se calculan en el sistema sobre el precio del vehículo; confirmar condiciones finales con un vendedor.")
	}

	partes = append(partes, "La IA nunca inventa tasas, plazos ni montos: solo usa estos datos.")

	return strings.Join(partes, "\n")
}