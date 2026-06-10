import { useEffect, useRef } from 'react'
import './Dropdown.css'

const Dropdown = ({ trigger, children, open, onToggle, className = '', style, ...rest }) => {
  const ref = useRef(null)

  useEffect(() => {
    const handleClick = (e) => {
      if (ref.current && !ref.current.contains(e.target)) {
        if (open) onToggle?.()
      }
    }
    document.addEventListener('mousedown', handleClick)
    return () => document.removeEventListener('mousedown', handleClick)
  }, [open, onToggle])

  const classes = ['dropdown', className].filter(Boolean).join(' ')

  return (
    <div className={classes} style={style} ref={ref} {...rest}>
      <div onClick={onToggle}>{trigger}</div>
      {open && (
        <div className="dropdown__menu" role="menu">
          {children}
        </div>
      )}
    </div>
  )
}

const DropdownItem = ({ children, onClick, className = '', ...rest }) => {
  const classes = ['dropdown__item', className].filter(Boolean).join(' ')
  return (
    <button className={classes} onClick={onClick} role="menuitem" {...rest}>
      {children}
    </button>
  )
}

const DropdownDivider = ({ className = '', ...rest }) => {
  return <div className="dropdown__divider" {...rest} />
}

Dropdown.Item = DropdownItem
Dropdown.Divider = DropdownDivider

export default Dropdown
