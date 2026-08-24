import { Routes, Route } from 'react-router'
import { LayoutBase } from '../layouts/LayoutBase'
import { Inicio } from '../pages/Inicio'
import { Catalogo } from '../pages/Catalogo'
import { DetalleVehiculo } from '../pages/DetalleVehiculo'
import { InicioSesion } from '../pages/InicioSesion'
import { Registro } from '../pages/Registro'
import { PanelAdministracion } from '../pages/PanelAdministracion'
import { AdminVehiculos } from '../pages/AdminVehiculos'
import { FormularioVehiculo } from '../pages/FormularioVehiculo'
import { NoEncontrada } from '../pages/NoEncontrada'
import { RutaProtegida } from '../components/RutaProtegida'
import { MisConsultas } from '../pages/MisConsultas'
import { BandejaEntrada } from '../pages/BandejaEntrada'
import { ChatVendedor } from '../pages/ChatVendedor'
import { FormularioTestDrive } from '../pages/FormularioTestDrive'
import { MisTestDrives } from '../pages/MisTestDrives'
import { GestionTestDrives } from '../pages/GestionTestDrives'
import { FormularioReserva } from '../pages/FormularioReserva'
import { MisReservas } from '../pages/MisReservas'
import { GestionReservas } from '../pages/GestionReservas'
import { GestionUsuarios } from '../pages/GestionUsuarios'
import { MisCotizaciones } from '../pages/MisCotizaciones'
import { BandejaCotizaciones } from '../pages/BandejaCotizaciones'
import { ChatCotizacionVendedor } from '../pages/ChatCotizacionVendedor'

export function Rutas() {
  return (
    <Routes>
      <Route element={<LayoutBase />}>
        <Route path="/" element={<Inicio />} />
        <Route path="/catalogo" element={<Catalogo />} />
        <Route path="/catalogo/:id" element={<DetalleVehiculo />} />
        <Route path="/login" element={<InicioSesion />} />
        <Route path="/registro" element={<Registro />} />
        <Route
          path="/admin"
          element={
            <RutaProtegida rol="administrador">
              <PanelAdministracion />
            </RutaProtegida>
          }
        />
        <Route
          path="/admin/vehiculos"
          element={
            <RutaProtegida rol="administrador">
              <AdminVehiculos />
            </RutaProtegida>
          }
        />
        <Route
          path="/admin/vehiculos/nuevo"
          element={
            <RutaProtegida rol="administrador">
              <FormularioVehiculo />
            </RutaProtegida>
          }
        />
        <Route
          path="/admin/vehiculos/:id/editar"
          element={
            <RutaProtegida rol="administrador">
              <FormularioVehiculo />
            </RutaProtegida>
          }
        />
        <Route
          path="/admin/usuarios"
          element={
            <RutaProtegida rol="administrador">
              <GestionUsuarios />
            </RutaProtegida>
          }
        />
        <Route
          path="/mis-consultas"
          element={
            <RutaProtegida rol="cliente">
              <MisConsultas />
            </RutaProtegida>
          }
        />
        <Route
          path="/mis-consultas/:id"
          element={
            <RutaProtegida rol="cliente">
              <MisConsultas />
            </RutaProtegida>
          }
        />
        <Route
          path="/vendedor/bandeja"
          element={
            <RutaProtegida rol="vendedor">
              <BandejaEntrada />
            </RutaProtegida>
          }
        />
        <Route
          path="/vendedor/bandeja/:id"
          element={
            <RutaProtegida rol="vendedor">
              <ChatVendedor />
            </RutaProtegida>
          }
        />
        <Route
          path="/vendedor/cotizaciones"
          element={
            <RutaProtegida rol="vendedor">
              <BandejaCotizaciones />
            </RutaProtegida>
          }
        />
        <Route
          path="/vendedor/cotizaciones/:id"
          element={
            <RutaProtegida rol="vendedor">
              <ChatCotizacionVendedor />
            </RutaProtegida>
          }
        />
        <Route
          path="/vendedor/test-drives"
          element={
            <RutaProtegida rol="vendedor">
              <GestionTestDrives />
            </RutaProtegida>
          }
        />
        <Route
          path="/vendedor/reservas"
          element={
            <RutaProtegida rol="vendedor">
              <GestionReservas />
            </RutaProtegida>
          }
        />
        <Route
          path="/catalogo/:id/test-drive"
          element={
            <RutaProtegida rol="cliente">
              <FormularioTestDrive />
            </RutaProtegida>
          }
        />
        <Route
          path="/catalogo/:id/reservar"
          element={
            <RutaProtegida rol="cliente">
              <FormularioReserva />
            </RutaProtegida>
          }
        />
        <Route
          path="/mis-test-drives"
          element={
            <RutaProtegida rol="cliente">
              <MisTestDrives />
            </RutaProtegida>
          }
        />
        <Route
          path="/mis-cotizaciones"
          element={
            <RutaProtegida rol="cliente">
              <MisCotizaciones />
            </RutaProtegida>
          }
        />
        <Route
          path="/mis-cotizaciones/:id"
          element={
            <RutaProtegida rol="cliente">
              <MisCotizaciones />
            </RutaProtegida>
          }
        />
        <Route
          path="/mis-reservas"
          element={
            <RutaProtegida rol="cliente">
              <MisReservas />
            </RutaProtegida>
          }
        />
        <Route path="*" element={<NoEncontrada />} />
      </Route>
    </Routes>
  )
}
