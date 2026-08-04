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
        <Route path="*" element={<NoEncontrada />} />
      </Route>
    </Routes>
  )
}
