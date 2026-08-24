package services

import "testing"

func TestNormalizarRespuestaConversacionalQuitaMarkdown(t *testing.T) {
	respuesta := "```markdown\n## Toyota Corolla\n\n**Precio:** $20.000\n- Año 2022\n- [Ver ficha](https://ejemplo.com/ficha)\n```"

	esperada := "Toyota Corolla\n\nPrecio: $20.000\nAño 2022\nVer ficha"
	if obtenida := normalizarRespuestaConversacional(respuesta); obtenida != esperada {
		t.Fatalf("respuesta normalizada inesperada: %q", obtenida)
	}
}

func TestLimpiarMarcadorCotizacionNormalizaRespuesta(t *testing.T) {
	respuesta, vehiculoID := limpiarMarcadorCotizacion("**Te preparo la cotización.**\n[COTIZACION:7]")
	respuesta = normalizarRespuestaConversacional(respuesta)

	if respuesta != "Te preparo la cotización." {
		t.Fatalf("respuesta inesperada: %q", respuesta)
	}
	if vehiculoID != 7 {
		t.Fatalf("id de vehículo inesperado: %d", vehiculoID)
	}
}

func TestExtraerMarcadoresVehiculoDeduplicaYLimpia(t *testing.T) {
	ids, texto := extraerMarcadoresVehiculo("El Corolla te puede interesar [VEHICULO:3] y también el Civic [VEHICULO:7]. Otro Corolla [VEHICULO:3]")

	if len(ids) != 2 || ids[0] != 3 || ids[1] != 7 {
		t.Fatalf("ids inesperados: %v", ids)
	}
	if texto != "El Corolla te puede interesar y también el Civic. Otro Corolla" {
		t.Fatalf("texto limpio inesperado: %q", texto)
	}
}

func TestExtraerMarcadoresVehiculoSinMarcador(t *testing.T) {
	if ids, texto := extraerMarcadoresVehiculo("Hoy tenemos buen stock"); ids != nil || texto != "Hoy tenemos buen stock" {
		t.Fatalf("resultado inesperado sin marcadores: ids=%v texto=%q", ids, texto)
	}
}

func TestFiltrarIdsServidosDescartaFueraDeContexto(t *testing.T) {
	servidos := map[uint]struct{}{3: {}, 7: {}}

	filtrados := filtrarIdsServidos([]uint{3, 9, 7}, servidos)
	if len(filtrados) != 2 || filtrados[0] != 3 || filtrados[1] != 7 {
		t.Fatalf("ids filtrados inesperados: %v", filtrados)
	}

	if obtenidos := filtrarIdsServidos(nil, servidos); obtenidos != nil {
		t.Fatalf("se esperaba nil sin ids: %v", obtenidos)
	}
	if obtenidos := filtrarIdsServidos([]uint{3}, nil); obtenidos != nil {
		t.Fatalf("se esperaba nil sin contexto: %v", obtenidos)
	}
}
