import { useState, useEffect } from 'react'
import { NavLink, useLocation } from 'react-router-dom'
import { HugeiconsIcon } from '@hugeicons/react'
import { Activity01Icon } from '@hugeicons/core-free-icons'
import { Analytics01Icon } from '@hugeicons/core-free-icons'
import { Add01Icon } from '@hugeicons/core-free-icons'
import { AddMoneyCircleIcon } from '@hugeicons/core-free-icons'
import { Queue01Icon } from '@hugeicons/core-free-icons'
import { ChartLineData01Icon } from '@hugeicons/core-free-icons'
import { Settings01Icon } from '@hugeicons/core-free-icons'
import { ViewIcon } from '@hugeicons/core-free-icons'
import { User02Icon } from '@hugeicons/core-free-icons'
import { SecurityIcon } from '@hugeicons/core-free-icons'
import { CheckmarkCircle01Icon } from '@hugeicons/core-free-icons'
import { ListChevronsDownUpIcon } from '@hugeicons/core-free-icons'
import { Cancel01Icon } from '@hugeicons/core-free-icons'
import { PanelRightCloseIcon } from '@hugeicons/core-free-icons'
import { PanelRightOpenIcon } from '@hugeicons/core-free-icons'
import './Sidebar.css'

const memberNav = [
  { path: '/dashboard', label: 'Dashboard', icon: Activity01Icon },
  { path: '/dashboard/requests', label: 'My Requests', icon: ListChevronsDownUpIcon },
  { path: '/dashboard/queue', label: 'Queue', icon: Queue01Icon },
  { path: '/dashboard/contributions', label: 'Contributions', icon: AddMoneyCircleIcon },
  { path: '/dashboard/impact', label: 'Impact', icon: ChartLineData01Icon },
]

const adminNav = [
  { path: '/admin/users', label: 'Users', icon: User02Icon },
  { path: '/admin/requests', label: 'Requests', icon: SecurityIcon },
  { path: '/admin/fulfillments', label: 'Fulfillments', icon: CheckmarkCircle01Icon },
]

const bottomNav = [
  { path: '/transparency', label: 'Transparency', icon: ViewIcon },
  { path: '/dashboard/settings', label: 'Settings', icon: Settings01Icon },
]

function SidebarItem({ item, collapsed }) {
  return (
    <NavLink
      to={item.path}
      className={({ isActive }) =>
        ['sidebar__item', isActive ? 'sidebar__item--active' : ''].join(' ')
      }
      end={item.path === '/dashboard'}
    >
      <span className="sidebar__item-icon">
        <HugeiconsIcon icon={item.icon} size={20} color="currentColor" strokeWidth={1.5} />
      </span>
      <span className="sidebar__item-label">{item.label}</span>
    </NavLink>
  )
}

export default function Sidebar({ collapsed, onToggle, mobileOpen, onMobileClose }) {
  const location = useLocation()
  const isAdmin = location.pathname.startsWith('/admin')
  const isInstitution = location.pathname.startsWith('/institutions')

  return (
    <>
      {/* Mobile overlay */}
      {mobileOpen && (
        <div className="sidebar__overlay" onClick={onMobileClose} />
      )}

      <aside className={[
        'sidebar',
        collapsed ? 'sidebar--collapsed' : '',
        mobileOpen ? 'sidebar--mobile-open' : 'sidebar--hidden',
      ].join(' ')}>
        <div className="sidebar__header">
          <span className="sidebar__wordmark">Virtus</span>
        </div>

        <nav className="sidebar__nav">
          {/* Member section */}
          {!isAdmin && !isInstitution && (
            <div className="sidebar__section">
              {memberNav.map(item => (
                <SidebarItem key={item.path} item={item} collapsed={collapsed} />
              ))}
            </div>
          )}

          {/* Admin section */}
          {isAdmin && (
            <div className="sidebar__section">
              <div className="sidebar__section-label">Admin</div>
              {adminNav.map(item => (
                <SidebarItem key={item.path} item={item} collapsed={collapsed} />
              ))}
            </div>
          )}

          {/* Institution section */}
          {isInstitution && (
            <div className="sidebar__section">
              <SidebarItem
                item={{ path: '/institutions/dashboard', label: 'Dashboard', icon: Analytics01Icon }}
                collapsed={collapsed}
              />
            </div>
          )}

          <div className="sidebar__divider" />

          {/* Bottom section */}
          <div className="sidebar__section">
            {bottomNav.map(item => (
              <SidebarItem key={item.path} item={item} collapsed={collapsed} />
            ))}
          </div>
        </nav>

        <div className="sidebar__footer">
          <button className="sidebar__toggle" onClick={onToggle} aria-label={collapsed ? 'Expand sidebar' : 'Collapse sidebar'}>
            <HugeiconsIcon
              icon={collapsed ? PanelRightOpenIcon : PanelRightCloseIcon}
              size={20}
              color="currentColor"
              strokeWidth={1.5}
            />
          </button>
        </div>
      </aside>
    </>
  )
}
