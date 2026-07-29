import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider } from './context/AuthContext'
import { ProtectedRoute } from './components/ProtectedRoute'
import { Layout } from './components/Layout'
import { Login } from './pages/Login'
import { Register } from './pages/Register'
import { Catalog } from './pages/Catalog'
import { VehicleDetail } from './pages/VehicleDetail'
import { ContactSeller } from './pages/ContactSeller'
import { TestDriveRequest } from './pages/TestDriveRequest'
import { ReserveVehicle } from './pages/ReserveVehicle'
import { AdminDashboard } from './pages/admin/Dashboard'
import { VehicleManagement } from './pages/admin/VehicleManagement'
import { VehicleForm } from './pages/admin/VehicleForm'
import { ConsultationInbox } from './pages/seller/ConsultationInbox'
import { MyConsultations } from './pages/MyConsultations'
import { TestDriveManagement } from './pages/seller/TestDriveManagement'
import { ReservationManagement } from './pages/seller/ReservationManagement'

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route element={<Layout />}>
            <Route path="/" element={<Navigate to="/catalogo" replace />} />
            <Route path="/catalogo" element={<Catalog />} />
            <Route path="/vehiculos/:id" element={<VehicleDetail />} />
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />

            <Route element={<ProtectedRoute />}>
              <Route path="/consultar/:id" element={<ContactSeller />} />
              <Route path="/test-drive/:id" element={<TestDriveRequest />} />
              <Route path="/reservar/:id" element={<ReserveVehicle />} />
              <Route path="/mis-consultas" element={<MyConsultations />} />
            </Route>

            <Route element={<ProtectedRoute allowedRoles={['administrador']} />}>
              <Route path="/admin/dashboard" element={<AdminDashboard />} />
              <Route path="/admin/vehiculos" element={<VehicleManagement />} />
              <Route path="/admin/vehiculos/nuevo" element={<VehicleForm />} />
              <Route path="/admin/vehiculos/:id" element={<VehicleForm />} />
            </Route>

            <Route element={<ProtectedRoute allowedRoles={['vendedor', 'administrador']} />}>
              <Route path="/seller/consultas" element={<ConsultationInbox />} />
              <Route path="/seller/test-drives" element={<TestDriveManagement />} />
              <Route path="/seller/reservas" element={<ReservationManagement />} />
            </Route>
          </Route>
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  )
}
