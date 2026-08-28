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
	{marca: "Toyota", modelo: "Corolla Cross", anio: 2023, kilometraje: 0, combustible: "Nafta", transmision: "Automática", tipo: "suv", precio: 45000000, condicion: models.CondicionNuevo, imagenes: []string{
		fotoWikimedia("Toyota%20Corolla%20Cross%20Hybrid%201X7A1861.jpg"),
		fotoWikimedia("Toyota%20Corolla%20Cross%20Hybrid%20Z%20rear.jpg"),
	}},
	{marca: "Toyota", modelo: "Hilux SRX", anio: 2022, kilometraje: 25000, combustible: "Diésel", transmision: "Automática", tipo: "pick-up", precio: 62000000, condicion: models.CondicionUsado, imagenes: []string{
		fotoWikimedia("Toyota%20Hilux%20%2834602333401%29.jpg"),
		fotoWikimedia("Toyota%20Hilux%20twincab%20%2830733462156%29.jpg"),
	}},
	{marca: "Toyota", modelo: "Yaris XLS", anio: 2023, kilometraje: 15000, combustible: "Nafta", transmision: "Automática", tipo: "hatchback", precio: 17000000, condicion: models.CondicionUsado, imagenes: []string{
		fotoWikimedia("2025%20Toyota%20Yaris%20XLS%2B%20%282nd%20facelift%29%20in%20Argentina.jpg"),
		fotoWikimedia("2025%20Toyota%20Yaris%20S%201.5.jpg"),
	}},
	{marca: "Ford", modelo: "Fiesta Kinetic", anio: 2016, kilometraje: 42000, combustible: "Nafta", transmision: "Manual", tipo: "sedán", precio: 9200000, condicion: models.CondicionUsado, imagenes: []string{
		fotoWikimedia("2015%20Ford%20Fiesta%20sedan%201.6%20SE%20Plus.jpg"),
	}},
	{marca: "Ford", modelo: "Focus Titanium", anio: 2017, kilometraje: 38000, combustible: "Nafta", transmision: "Automática", tipo: "hatchback", precio: 14000000, condicion: models.CondicionUsado, imagenes: []string{
		fotoWikimedia("2018%20Ford%20Focus%20Titanium%201.0%20%28Front%29.jpg"),
		fotoWikimedia("2018%20Ford%20Focus%20Titanium%201.0%20%28Rear%29.jpg"),
	}},
	{marca: "Volkswagen", modelo: "Gol Trend", anio: 2014, kilometraje: 85000, combustible: "Nafta", transmision: "Manual", tipo: "hatchback", precio: 7500000, condicion: models.CondicionUsado, imagenes: []string{
		fotoWikimedia("2014%20Volkswagen%20Gol%20Trend%201.6%20Cup.jpg"),
	}},
	{marca: "Volkswagen", modelo: "Amarok V6", anio: 2022, kilometraje: 30000, combustible: "Diésel", transmision: "Automática", tipo: "pick-up", precio: 52000000, condicion: models.CondicionUsado, imagenes: []string{
		fotoWikimedia("Volkswagen%20Amarok%20V6%20TDi%20Extreme%202022.jpg"),
		fotoWikimedia("2023%20Volkswagen%20Amarok%20Extreme%20V6%204X4.jpg"),
	}},
	{marca: "Fiat", modelo: "Pulse Drive", anio: 2022, kilometraje: 0, combustible: "Nafta", transmision: "Automática", tipo: "suv", precio: 22500000, condicion: models.CondicionNuevo, imagenes: []string{
		fotoWikimedia("2022%20Fiat%20Pulse%201.3%20Drive.jpg"),
	}},
	{marca: "Fiat", modelo: "Palio Fire", anio: 2018, kilometraje: 60000, combustible: "Nafta", transmision: "Manual", tipo: "hatchback", precio: 5800000, condicion: models.CondicionUsado, imagenes: []string{
		fotoWikimedia("Fiat%20Palio%20Fire%202017.jpg"),
	}},
	{marca: "Chevrolet", modelo: "Cruze LTZ", anio: 2019, kilometraje: 32000, combustible: "Nafta", transmision: "Automática", tipo: "sedán", precio: 18500000, condicion: models.CondicionUsado, imagenes: []string{
		fotoWikimedia("CHEVROLET%20CRUZE%20%28J400%29%20China.jpg"),
		fotoWikimedia("Chevrolet%20Cruze%20%28third%20generation%29%201X7A0417.jpg"),
	}},
	{marca: "Chevrolet", modelo: "S10 LTZ", anio: 2021, kilometraje: 45000, combustible: "Diésel", transmision: "Automática", tipo: "pick-up", precio: 39000000, condicion: models.CondicionUsado, imagenes: []string{
		fotoWikimedia("ChevroletS10-Carilo-06280.jpg"),
		fotoWikimedia("Chevrolet%20S10%202.4%20MPFI.jpg"),
	}},
}

// SembrarVehiculos crea los vehículos por defecto si no existen (idempotente)
// y completa el campo tipo de los vehículos existentes que aún no lo tienen.
// Con la base vacía crea todo el catálogo; con datos ya cargados crea los
// modelos por defecto que falten para que los nuevos vehículos del seed
// queden disponibles aunque la base ya esté poblada.
func SembrarVehiculos(base *gorm.DB) error {
	var total int64
	if err := base.Model(&models.Vehiculo{}).Count(&total).Error; err != nil {
		return err
	}
	if total == 0 {
		return crearVehiculosPorDefecto(base)
	}

	for _, v := range vehiculosPorDefecto {
		if err := base.Model(&models.Vehiculo{}).
			Where("marca = ? AND modelo = ? AND (tipo = '' OR tipo IS NULL)", v.marca, v.modelo).
			Update("tipo", v.tipo).Error; err != nil {
			return err
		}

		var totalModelo int64
		if err := base.Model(&models.Vehiculo{}).
			Where("marca = ? AND modelo = ?", v.marca, v.modelo).Count(&totalModelo).Error; err != nil {
			return err
		}
		if totalModelo > 0 {
			continue
		}
		nuevo := vehiculoDesdeDefecto(v)
		if err := base.Create(&nuevo).Error; err != nil {
			return err
		}
	}
	return nil
}

// crearVehiculosPorDefecto inserta todo el catálogo por defecto.
func crearVehiculosPorDefecto(base *gorm.DB) error {
	for _, v := range vehiculosPorDefecto {
		nuevo := vehiculoDesdeDefecto(v)
		if err := base.Create(&nuevo).Error; err != nil {
			return err
		}
	}
	return nil
}

// vehiculoDesdeDefecto arma la entidad Vehiculo a partir del registro por
// defecto del seed.
func vehiculoDesdeDefecto(v vehiculoPorDefecto) models.Vehiculo {
	imagenes := make([]models.Imagen, 0, len(v.imagenes))
	for _, url := range v.imagenes {
		imagenes = append(imagenes, models.Imagen{URL: url})
	}
	return models.Vehiculo{
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
