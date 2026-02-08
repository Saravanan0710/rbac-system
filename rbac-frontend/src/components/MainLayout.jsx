import './MainLayout.css'

function MainLayout({ children }) {
  return (
    <main className="main-layout">
      <div className="main-content">
        {children}
      </div>
    </main>
  )
}

export default MainLayout
