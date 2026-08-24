package services

import (
	"context"
	"errors"
	"testing"

	"concesionaria/backend/internal/googleid"
	"concesionaria/backend/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func nuevoServicioAutenticacion(repo *fakeUsuarioRepository) AutenticacionService {
	return NuevoAutenticacionService(repo, "secreto-de-prueba", DuracionToken, nil)
}

func TestRegistrarExitoso(t *testing.T) {
	servicio := nuevoServicioAutenticacion(nuevoFakeUsuarioRepository())

	usuario, err := servicio.Registrar(context.Background(), EntradaRegistro{
		Nombre:   "  Ana Pérez  ",
		Email:    "  ANA@Ejemplo.com ",
		Password: "secreto123",
	})
	if err != nil {
		t.Fatalf("Registrar devolvió error: %v", err)
	}
	if usuario.Email != "ana@ejemplo.com" {
		t.Errorf("email sin normalizar: %q", usuario.Email)
	}
	if usuario.Nombre != "Ana Pérez" {
		t.Errorf("nombre sin recortar: %q", usuario.Nombre)
	}
	if usuario.Rol != models.RolCliente {
		t.Errorf("rol esperado %q, obtenido %q", models.RolCliente, usuario.Rol)
	}
	if usuario.Password == "secreto123" {
		t.Error("la contraseña se guardó en texto plano")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(usuario.Password), []byte("secreto123")); err != nil {
		t.Error("el hash de la contraseña no coincide")
	}
}

func TestRegistrarEmailRepetido(t *testing.T) {
	servicio := nuevoServicioAutenticacion(nuevoFakeUsuarioRepository())

	primero, err := servicio.Registrar(context.Background(), EntradaRegistro{
		Nombre:   "Ana",
		Email:    "ana@ejemplo.com",
		Password: "secreto123",
	})
	if err != nil {
		t.Fatalf("primer registro falló: %v", err)
	}

	_ = primero
	_, err = servicio.Registrar(context.Background(), EntradaRegistro{
		Nombre:   "Otra Ana",
		Email:    "ANA@ejemplo.com",
		Password: "secreto456",
	})
	if !errors.Is(err, ErrEmailEnUso) {
		t.Errorf("se esperaba ErrEmailEnUso, se obtuvo %v", err)
	}
}

func TestRegistrarDatosInvalidos(t *testing.T) {
	casos := []EntradaRegistro{
		{Nombre: "", Email: "ana@ejemplo.com", Password: "secreto123"},
		{Nombre: "Ana", Email: "no-es-un-email", Password: "secreto123"},
		{Nombre: "Ana", Email: "ana@ejemplo.com", Password: "corta"},
	}
	for i, caso := range casos {
		servicio := nuevoServicioAutenticacion(nuevoFakeUsuarioRepository())
		if _, err := servicio.Registrar(context.Background(), caso); !errors.Is(err, ErrDatosRegistroInvalidos) {
			t.Errorf("caso %d: se esperaba ErrDatosRegistroInvalidos, se obtuvo %v", i, err)
		}
	}
}

func TestIniciarSesionExitoso(t *testing.T) {
	repo := nuevoFakeUsuarioRepository()
	hash, _ := bcrypt.GenerateFromPassword([]byte("secreto123"), bcrypt.DefaultCost)
	repo.porEmail["ana@ejemplo.com"] = &models.Usuario{
		ID:       1,
		Nombre:   "Ana",
		Email:    "ana@ejemplo.com",
		Password: string(hash),
		Rol:      models.RolCliente,
	}
	repo.porID[1] = repo.porEmail["ana@ejemplo.com"]

	servicio := nuevoServicioAutenticacion(repo)
	usuario, tokenString, err := servicio.IniciarSesion(context.Background(), EntradaLogin{
		Email:    "ANA@Ejemplo.com",
		Password: "secreto123",
	})
	if err != nil {
		t.Fatalf("IniciarSesion devolvió error: %v", err)
	}
	if usuario.ID != 1 {
		t.Errorf("usuario esperado con ID 1, obtenido %d", usuario.ID)
	}
	if tokenString == "" {
		t.Error("se esperaba un token JWT no vacío")
	}
}

func TestIniciarSesionPasswordIncorrecta(t *testing.T) {
	repo := nuevoFakeUsuarioRepository()
	hash, _ := bcrypt.GenerateFromPassword([]byte("secreto123"), bcrypt.DefaultCost)
	repo.porEmail["ana@ejemplo.com"] = &models.Usuario{
		ID:       1,
		Email:    "ana@ejemplo.com",
		Password: string(hash),
		Rol:      models.RolCliente,
	}
	repo.porID[1] = repo.porEmail["ana@ejemplo.com"]

	servicio := nuevoServicioAutenticacion(repo)
	_, _, err := servicio.IniciarSesion(context.Background(), EntradaLogin{
		Email:    "ana@ejemplo.com",
		Password: "incorrecta",
	})
	if !errors.Is(err, ErrCredencialesInvalidas) {
		t.Errorf("se esperaba ErrCredencialesInvalidas, se obtuvo %v", err)
	}
}

func TestIniciarSesionEmailInexistente(t *testing.T) {
	servicio := nuevoServicioAutenticacion(nuevoFakeUsuarioRepository())
	_, _, err := servicio.IniciarSesion(context.Background(), EntradaLogin{
		Email:    "nadie@ejemplo.com",
		Password: "secreto123",
	})
	if !errors.Is(err, ErrCredencialesInvalidas) {
		t.Errorf("se esperaba ErrCredencialesInvalidas, se obtuvo %v", err)
	}
}

func TestObtenerPorID(t *testing.T) {
	repo := nuevoFakeUsuarioRepository()
	repo.porID[5] = &models.Usuario{ID: 5, Nombre: "Ana", Email: "ana@ejemplo.com", Rol: models.RolCliente}

	servicio := nuevoServicioAutenticacion(repo)

	usuario, err := servicio.ObtenerPorID(context.Background(), 5)
	if err != nil || usuario.ID != 5 {
		t.Errorf("ObtenerPorID(5) = %v, %v; se esperaba usuario 5", usuario, err)
	}

	_, err = servicio.ObtenerPorID(context.Background(), 99)
	if !errors.Is(err, ErrUsuarioNoEncontrado) {
		t.Errorf("se esperaba ErrUsuarioNoEncontrado, se obtuvo %v", err)
	}
}

// fakeVerificadorGoogle simula la verificación de credenciales de Google.
type fakeVerificadorGoogle struct {
	identidad *googleid.Identidad
	err       error
}

func (f *fakeVerificadorGoogle) Verificar(ctx context.Context, idToken string) (*googleid.Identidad, error) {
	return f.identidad, f.err
}

func TestIniciarSesionConGoogleCreaCliente(t *testing.T) {
	repo := nuevoFakeUsuarioRepository()
	servicio := NuevoAutenticacionService(repo, "secreto-de-prueba", DuracionToken, &fakeVerificadorGoogle{
		identidad: &googleid.Identidad{Sub: "sub-google-1", Email: "Nuevo@Ejemplo.com", Nombre: "Nuevo Cliente"},
	})

	usuario, tokenString, err := servicio.IniciarSesionConGoogle(context.Background(), "credencial")
	if err != nil {
		t.Fatalf("IniciarSesionConGoogle devolvió error: %v", err)
	}
	if tokenString == "" {
		t.Error("se esperaba un token JWT no vacío")
	}
	if usuario.Rol != models.RolCliente || usuario.Proveedor != models.ProveedorGoogle {
		t.Errorf("usuario creado incorrecto: rol=%q proveedor=%q", usuario.Rol, usuario.Proveedor)
	}
	if usuario.GoogleSub == nil || *usuario.GoogleSub != "sub-google-1" {
		t.Errorf("google_sub no guardado: %v", usuario.GoogleSub)
	}
	if usuario.Password != "" {
		t.Error("una cuenta creada con Google no debe tener contraseña local")
	}
	if _, ok := repo.porEmail["nuevo@ejemplo.com"]; !ok {
		t.Error("el usuario no quedó persistido con el email normalizado")
	}
}

func TestIniciarSesionConGoogleVinculaPorEmail(t *testing.T) {
	repo := nuevoFakeUsuarioRepository()
	hash, _ := bcrypt.GenerateFromPassword([]byte("secreto123"), bcrypt.DefaultCost)
	repo.porEmail["ana@ejemplo.com"] = &models.Usuario{
		ID: 1, Nombre: "Ana", Email: "ana@ejemplo.com",
		Password: string(hash), Rol: models.RolVendedor, Proveedor: models.ProveedorLocal,
	}
	repo.porID[1] = repo.porEmail["ana@ejemplo.com"]

	servicio := NuevoAutenticacionService(repo, "secreto-de-prueba", DuracionToken, &fakeVerificadorGoogle{
		identidad: &googleid.Identidad{Sub: "sub-ana", Email: "ana@ejemplo.com", Nombre: "Otro Nombre"},
	})

	usuario, _, err := servicio.IniciarSesionConGoogle(context.Background(), "credencial")
	if err != nil {
		t.Fatalf("IniciarSesionConGoogle devolvió error: %v", err)
	}
	if len(repo.porEmail) != 1 {
		t.Errorf("se esperaba un único usuario, hay %d", len(repo.porEmail))
	}
	if usuario.ID != 1 {
		t.Errorf("debió vincularse al usuario existente ID 1, se obtuvo %d", usuario.ID)
	}
	if usuario.Rol != models.RolVendedor {
		t.Errorf("el rol original debió conservarse, se obtuvo %q", usuario.Rol)
	}
	if usuario.Nombre != "Ana" {
		t.Errorf("el nombre original debió conservarse, se obtuvo %q", usuario.Nombre)
	}
	if usuario.Proveedor != models.ProveedorGoogle || usuario.GoogleSub == nil || *usuario.GoogleSub != "sub-ana" {
		t.Errorf("vinculación no aplicada: proveedor=%q google_sub=%v", usuario.Proveedor, usuario.GoogleSub)
	}
}

func TestIniciarSesionConGoogleRecurrenteNoDuplica(t *testing.T) {
	repo := nuevoFakeUsuarioRepository()
	servicio := NuevoAutenticacionService(repo, "secreto-de-prueba", DuracionToken, &fakeVerificadorGoogle{
		identidad: &googleid.Identidad{Sub: "sub-fijo", Email: "recurrente@ejemplo.com", Nombre: "Recurrente"},
	})

	for i := 0; i < 3; i++ {
		if _, _, err := servicio.IniciarSesionConGoogle(context.Background(), "credencial"); err != nil {
			t.Fatalf("ingreso %d falló: %v", i+1, err)
		}
	}
	if len(repo.porEmail) != 1 {
		t.Errorf("se esperaba un único usuario tras ingresos repetidos, hay %d", len(repo.porEmail))
	}
}

func TestIniciarSesionConGoogleSinConfigurar(t *testing.T) {
	servicio := nuevoServicioAutenticacion(nuevoFakeUsuarioRepository())
	if servicio.GoogleHabilitado() {
		t.Error("sin verificador, Google debe estar deshabilitado")
	}
	_, _, err := servicio.IniciarSesionConGoogle(context.Background(), "credencial")
	if !errors.Is(err, ErrGoogleNoDisponible) {
		t.Errorf("se esperaba ErrGoogleNoDisponible, se obtuvo %v", err)
	}
}

func TestIniciarSesionConGoogleCredencialInvalida(t *testing.T) {
	repo := nuevoFakeUsuarioRepository()
	servicio := NuevoAutenticacionService(repo, "secreto-de-prueba", DuracionToken, &fakeVerificadorGoogle{
		err: googleid.ErrTokenInvalido,
	})

	_, _, err := servicio.IniciarSesionConGoogle(context.Background(), "credencial-falsa")
	if !errors.Is(err, ErrCredencialGoogleInvalida) {
		t.Errorf("se esperaba ErrCredencialGoogleInvalida, se obtuvo %v", err)
	}
	if len(repo.porEmail) != 0 {
		t.Error("no debe crearse ningún usuario con credencial inválida")
	}
}
