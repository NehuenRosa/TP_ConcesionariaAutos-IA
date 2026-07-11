import { useEffect } from 'react'
import { Outlet, useLocation } from 'react-router-dom'
import { Navbar } from './Navbar'
import { Chatbot } from './Chatbot'

export function Layout() {
  const location = useLocation()

  useEffect(() => {
    const scrollToTop = () => {
      document.documentElement.scrollTop = 0
      document.body.scrollTop = 0
      window.scrollTo(0, 0)
    }
    scrollToTop()
    requestAnimationFrame(scrollToTop)
    setTimeout(scrollToTop, 50)
  }, [location.pathname, location.search])

  return (
    <div className="min-h-screen flex flex-col bg-surface">
      <div className="fixed top-0 right-0 w-[500px] h-[500px] bg-gradient-to-bl from-brand-200/15 via-transparent to-transparent rounded-full blur-[120px] pointer-events-none -z-10" />
      <div className="fixed bottom-0 left-0 w-[400px] h-[400px] bg-gradient-to-tr from-purple-200/15 via-transparent to-transparent rounded-full blur-[120px] pointer-events-none -z-10" />
      <Navbar />
      <main className="flex-1">
        <div key={location.pathname + location.search} className="animate-fade-in">
          <Outlet />
        </div>
      </main>
      <Chatbot />
    </div>
  )
}
