import { Routes, Route } from 'react-router'
import { LayoutBase } from '../layouts/LayoutBase'
import { Inicio } from '../pages/Inicio'
import { Catalogo } from '../pages/Catalogo'
import { DetalleVehiculo } from '../pages/DetalleVehiculo'
import { InicioSesion } from '../pages/InicioSesion'
import { PanelAdministracion } from '../pages/PanelAdministracion'
import { AdminVehiculos } from '../pages/AdminVehiculos'
import { FormularioVehiculo } from '../pages/FormularioVehiculo'
import { NoEncontrada } from '../pages/NoEncontrada'

export function Rutas() {
  return (
    <Routes>
      <Route element={<LayoutBase />}>
        <Route path="/" element={<Inicio />} />
        <Route path="/catalogo" element={<Catalogo />} />
        <Route path="/catalogo/:id" element={<DetalleVehiculo />} />
        <Route path="/login" element={<InicioSesion />} />
        <Route path="/admin" element={<PanelAdministracion />} />
        <Route path="/admin/vehiculos" element={<AdminVehiculos />} />
        <Route path="/admin/vehiculos/nuevo" element={<FormularioVehiculo />} />
        <Route path="/admin/vehiculos/:id/editar" element={<FormularioVehiculo />} />
        <Route path="*" element={<NoEncontrada />} />
      </Route>
    </Routes>
  )
}
