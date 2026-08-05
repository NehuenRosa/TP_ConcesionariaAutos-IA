package database

import (
	"errors"

	"concesionaria/backend/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// usuarioPorDefecto describe un usuario sembrado al arrancar el sistema.
type usuarioPorDefecto struct {
	nombre   string
	email    string
	password string
	rol      string
}

// usuariosPorDefecto son las cuentas de desarrollo creadas en el primer
// arranque para poder operar el sistema sin una pantalla de gestión de
// usuarios.
var usuariosPorDefecto = []usuarioPorDefecto{
	{nombre: "Administrador", email: "administrador@concesionaria.local", password: "Admin123!", rol: models.RolAdministrador},
	{nombre: "Vendedor", email: "vendedor@concesionaria.local", password: "Vendedor123!", rol: models.RolVendedor},
	{nombre: "Cliente Genérico", email: "cliente@concesionaria.local", password: "Cliente123!", rol: models.RolCliente},
}

// vehiculoPorDefecto describe un vehículo sembrado al arrancar el sistema.
type vehiculoPorDefecto struct {
	marca       string
	modelo      string
	anio        int
	kilometraje int
	combustible string
	transmision string
	tipo        string
	precio      float64
	condicion   string
	imagenes    []string
}

// fotoWikimedia construye la URL de una imagen real alojada en Wikimedia
// Commons para el auto específico de cada modelo del seed.
func fotoWikimedia(nombre string) string {
	return "https://commons.wikimedia.org/wiki/Special:FilePath/" + nombre + "?width=1600"
}

// vehiculosPorDefecto son unidades de prueba para operar el catálogo y los
// filtros de búsqueda sin cargar datos a mano.
var vehiculosPorDefecto = []vehiculoPorDefecto{
	{marca: "Toyota", modelo: "Corolla", anio: 2022, kilometraje: 18000, combustible: "Nafta", transmision: "Automática", tipo: "sedán", precio: 32000000, condicion: models.CondicionUsado, imagenes: []string{
		fotoWikimedia("2018%20Toyota%20Corolla%20%28ZRE172R%29%20Ascent%20sedan%20%282018-11-02%29%2002.jpg"),
		fotoWikimedia("CorollaAltisZRE172V14PearlWhiteFR.jpg"),
	}},
	{marca: "Ford", modelo: "Ranger", anio: 2023, kilometraje: 5000, combustible: "Diésel", transmision: "Manual", tipo: "pick-up", precio: 45000000, condicion: models.CondicionUsado, imagenes: []string{
		fotoWikimedia("Ford%20Ranger%20%28T6%2C%20P703%29%20Wildtrak%20IMG%207320.jpg"),
		fotoWikimedia("Ford%20Ranger%20Raptor%20-%2020231003-P1003802.jpg"),
	}},
	{marca: "Volkswagen", modelo: "Taos", anio: 2024, kilometraje: 0, combustible: "Nafta", transmision: "Automática", tipo: "suv", precio: 38000000, condicion: models.CondicionNuevo, imagenes: []string{
		fotoWikimedia("Volkswagen%20Taos%201X7A0365.jpg"),
		fotoWikimedia("Volkswagen%20Taos%201X7A6748.jpg"),
	}},
	{marca: "Chevrolet", modelo: "Onix", anio: 2021, kilometraje: 35000, combustible: "Nafta", transmision: "Manual", tipo: "hatchback", precio: 15000000, condicion: models.CondicionUsado, imagenes: []string{
		fotoWikimedia("Chevrolet%20Onix%20RS%2C%20EVA-01.jpg"),
		fotoWikimedia("Chevrolet%20Onix%20RS%2C%20EVA-01-2.jpg"),
	}},
	{marca: "Fiat", modelo: "Cronos", anio: 2024, kilometraje: 0, combustible: "Nafta", transmision: "Automática", tipo: "sedán", precio: 22000000, condicion: models.CondicionNuevo, imagenes: []string{
		fotoWikimedia("Fiat%20Cronos%201.8%2016V%20E.Torq%20Precision.jpg"),
		fotoWikimedia("Fiat%20Cronos%202026%20front%20end%20with%20headlights%20on.jpg"),
	}},
	{marca: "Ford", modelo: "Mustang", anio: 2019, kilometraje: 28000, combustible: "Nafta", transmision: "Manual", tipo: "coupe", precio: 65000000, condicion: models.CondicionUsado, imagenes: []string{
		fotoWikimedia("Photography%20by%20David%20Adam%20Kess%20Ford%20Mustang%20%28sixth%20generation%29%20fastback.jpg"),
		fotoWikimedia("Ford%20Mustang%20Sixth%20generation%2C%20convertible%20with%20Racing%20Stripes.jpg"),
	}},
	{marca: "Porsche", modelo: "911 Carrera", anio: 2024, kilometraje: 0, combustible: "Nafta", transmision: "Automática", tipo: "coupe", precio: 185000000, condicion: models.CondicionNuevo, imagenes: []string{
		fotoWikimedia("Porsche%20992%20Carrera%20S%20coupe%20IMG%205832.jpg"),
		fotoWikimedia("Porsche%20992%20Carrera%20S%20coupe%20IMG%205843.jpg"),
	}},
	{marca: "Jeep", modelo: "Wrangler", anio: 2020, kilometraje: 42000, combustible: "Nafta", transmision: "Manual", tipo: "suv", precio: 48000000, condicion: models.CondicionUsado, imagenes: []string{
		fotoWikimedia("Jeep%20Wrangler%20Rubicon%20%28JL%29%204xe%201X7A0285.jpg"),
		fotoWikimedia("Jeep%20Wrangler%20Rubicon%20392%20IMG%207673.jpg"),
	}},
	{marca: "Mercedes-Benz", modelo: "Clase C", anio: 2023, kilometraje: 0, combustible: "Nafta", transmision: "Automática", tipo: "sedán", precio: 92000000, condicion: models.CondicionNuevo, imagenes: []string{
		fotoWikimedia("MERCEDES-BENZ%20C-CLASS%20LWB%20%28W206%29%20China%20%2817%29.jpg"),
		fotoWikimedia("MERCEDES-BENZ%20C-CLASS%20LWB%20%28W206%29%20China%20%287%29.jpg"),
	}},
	{marca: "BMW", modelo: "M4", anio: 2023, kilometraje: 0, combustible: "Nafta", transmision: "Automática", tipo: "coupe", precio: 120000000, condicion: models.CondicionNuevo, imagenes: []string{
		fotoWikimedia("BMW%20M4%2C%20EMS%2023%2C%20Essen%20%28P1170092%29.jpg"),
		fotoWikimedia("BMW%20M4%20in%20der%20BMW%20Welt.jpg"),
	}},
	{marca: "Nissan", modelo: "GT-R", anio: 2018, kilometraje: 30000, combustible: "Nafta", transmision: "Automática", tipo: "coupe", precio: 95000000, condicion: models.CondicionUsado, imagenes: []string{
		fotoWikimedia("Nissan%20GT-R%20R35%20at%20the%202026%20Adelaide%20Motorsport%20Festival%20%28DSCF2567%29.jpg"),
		fotoWikimedia("NISSAN%20GT-R%20%28R35%2C%202011%20FACELIFT%29%20China.jpg"),
	}},
	{marca: "Chevrolet", modelo: "Camaro", anio: 2020, kilometraje: 22000, combustible: "Nafta", transmision: "Manual", tipo: "coupe", precio: 58000000, condicion: models.CondicionUsado, imagenes: []string{
		fotoWikimedia("Chevrolet%20Camaro%20Hirschaid%202022-20220709-RM-111908.jpg"),
		fotoWikimedia("Chevrolet%20Camaro%20Hirschaid%202022-20220709-RM-112013.jpg"),
	}},
	{marca: "Land Rover", modelo: "Defender", anio: 2024, kilometraje: 0, combustible: "Diésel", transmision: "Automática", tipo: "suv", precio: 130000000, condicion: models.CondicionNuevo, imagenes: []string{
		fotoWikimedia("Land%20Rover%20Defender%20%28L663%29%20V8%20IMG%206604.jpg"),
		fotoWikimedia("Land%20Rover%20Defender%20110%20First%20Edition%202020%20-%20rear.jpg"),
	}},
	{marca: "Ferrari", modelo: "Roma", anio: 2023, kilometraje: 0, combustible: "Nafta", transmision: "Automática", tipo: "coupe", precio: 380000000, condicion: models.CondicionNuevo, imagenes: []string{
		fotoWikimedia("Ferrari%20Roma%201X7A0309.jpg"),
		fotoWikimedia("Ferrari%20Roma%20IMG%209620.jpg"),
	}},
}

// SembrarVehiculos crea los vehículos por defecto si no existen (idempotente)
// y completa el campo tipo de los vehículos existentes que aún no lo tienen.
func SembrarVehiculos(base *gorm.DB) error {
	var total int64
	if err := base.Model(&models.Vehiculo{}).Count(&total).Error; err != nil {
		return err
	}
	if total == 0 {
		for _, v := range vehiculosPorDefecto {
			imagenes := make([]models.Imagen, 0, len(v.imagenes))
			for _, url := range v.imagenes {
				imagenes = append(imagenes, models.Imagen{URL: url})
			}
			nuevo := models.Vehiculo{
				Marca:       v.marca,
				Modelo:      v.modelo,
				Anio:        v.anio,
				Kilometraje: v.kilometraje,
				Combustible: v.combustible,
				Transmision: v.transmision,
				Tipo:        v.tipo,
				Precio:      v.precio,
				Condicion:   v.condicion,
				Estado:      models.EstadoDisponible,
				Imagenes:    imagenes,
			}
			if err := base.Create(&nuevo).Error; err != nil {
				return err
			}
		}
		return nil
	}

	for _, v := range vehiculosPorDefecto {
		if err := base.Model(&models.Vehiculo{}).
			Where("marca = ? AND modelo = ? AND (tipo = '' OR tipo IS NULL)", v.marca, v.modelo).
			Update("tipo", v.tipo).Error; err != nil {
			return err
		}
	}
	return nil
}

// SembrarUsuarios crea los usuarios por defecto si no existen (idempotente).
func SembrarUsuarios(base *gorm.DB) error {
	for _, u := range usuariosPorDefecto {
		var usuario models.Usuario
		err := base.Where("email = ?", u.email).First(&usuario).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		nuevo := models.Usuario{
			Nombre:   u.nombre,
			Email:    u.email,
			Password: string(hash),
			Rol:      u.rol,
		}
		if err := base.Create(&nuevo).Error; err != nil {
			return err
		}
	}
	return nil
}
