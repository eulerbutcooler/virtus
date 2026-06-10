import { HugeiconsIcon } from '@hugeicons/react'
import { Menu01Icon } from '@hugeicons/core-free-icons'
import { Logout01Icon } from '@hugeicons/core-free-icons'
import Avatar from '../ui/Avatar.jsx'
import './Topbar.css'

export default function Topbar({ pageTitle, onMenuToggle, user }) {
  return (
    <header className="topbar">
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
        {user && (
          <div className="topbar__user">
            <Avatar name={user.name} size="sm" />
            <span className="topbar__user-name">{user.name}</span>
          </div>
        )}
        <button className="topbar__logout" aria-label="Logout">
          <HugeiconsIcon icon={Logout01Icon} size={18} color="currentColor" strokeWidth={1.5} />
        </button>
      </div>
    </header>
  )
}
