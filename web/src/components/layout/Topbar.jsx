import { HugeiconsIcon } from '@hugeicons/react'
import { Menu01Icon, Logout01Icon } from '@hugeicons/core-free-icons'
import Avatar from '../ui/Avatar.jsx'
import ThemeToggle from '../ui/ThemeToggle.jsx'
import './Topbar.css'

export default function Topbar({ pageTitle, onMenuToggle, user, onLogout, collapsed }) {
  return (
    <header className={['topbar', collapsed ? 'topbar--collapsed' : ''].filter(Boolean).join(' ')}>
      <div className="topbar__left">
        <button
          className="topbar__hamburger"
          onClick={onMenuToggle}
          aria-label="Toggle navigation menu"
        >
          <HugeiconsIcon icon={Menu01Icon} size={20} color="currentColor" strokeWidth={1.5} />
        </button>
        {pageTitle && <span className="topbar__page-title">{pageTitle}</span>}
      </div>

      <div className="topbar__right">

        <ThemeToggle />

        {user && (
          <div className="topbar__user">
            <Avatar name={user.name} size="sm" />
            <span className="topbar__user-name">{user.name}</span>
          </div>
        )}
        <button className="topbar__logout" aria-label="Logout" onClick={onLogout}>
          <HugeiconsIcon icon={Logout01Icon} size={18} color="currentColor" strokeWidth={1.5} />
        </button>
      </div>
    </header>
  )
}
