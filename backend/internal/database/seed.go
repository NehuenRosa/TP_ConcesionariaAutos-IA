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

// vehiculosPorDefecto son unidades de prueba para operar el catálogo y los
// filtros de búsqueda sin cargar datos a mano.
var vehiculosPorDefecto = []vehiculoPorDefecto{
	{marca: "Toyota", modelo: "Corolla", anio: 2022, kilometraje: 18000, combustible: "Nafta", transmision: "Automática", tipo: "sedán", precio: 32000000, condicion: models.CondicionUsado, imagenes: []string{"https://picsum.photos/seed/corolla/800/600"}},
	{marca: "Ford", modelo: "Ranger", anio: 2023, kilometraje: 5000, combustible: "Diésel", transmision: "Manual", tipo: "pick-up", precio: 45000000, condicion: models.CondicionUsado, imagenes: []string{"https://picsum.photos/seed/ranger/800/600"}},
	{marca: "Volkswagen", modelo: "Taos", anio: 2024, kilometraje: 0, combustible: "Nafta", transmision: "Automática", tipo: "suv", precio: 38000000, condicion: models.CondicionNuevo, imagenes: []string{"https://picsum.photos/seed/taos/800/600"}},
	{marca: "Chevrolet", modelo: "Onix", anio: 2021, kilometraje: 35000, combustible: "Nafta", transmision: "Manual", tipo: "hatchback", precio: 15000000, condicion: models.CondicionUsado, imagenes: []string{"https://picsum.photos/seed/onix/800/600"}},
	{marca: "Fiat", modelo: "Cronos", anio: 2024, kilometraje: 0, combustible: "Nafta", transmision: "Automática", tipo: "sedán", precio: 22000000, condicion: models.CondicionNuevo, imagenes: []string{"https://picsum.photos/seed/cronos/800/600"}},
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
