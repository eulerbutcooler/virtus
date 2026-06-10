import { useEffect, useRef, useState } from 'react'
import './Modal.css'

const Modal = ({ open, onClose, title, children, footer, className = '' }) => {
  const [exiting, setExiting] = useState(false)
  const [visible, setVisible] = useState(false)
  const overlayRef = useRef(null)

  useEffect(() => {
    if (open) {
      setExiting(false)
      setVisible(true)
    } else if (visible) {
      setExiting(true)
      const timer = setTimeout(() => setVisible(false), 200)
      return () => clearTimeout(timer)
    }
  }, [open])

  useEffect(() => {
    const handleKey = (e) => {
      if (e.key === 'Escape' && open && onClose) {
        onClose()
      }
    }
    window.addEventListener('keydown', handleKey)
    return () => window.removeEventListener('keydown', handleKey)
  }, [open, onClose])

  useEffect(() => {
    if (visible) {
      document.body.style.overflow = 'hidden'
    } else {
      document.body.style.overflow = ''
    }
    return () => {
      document.body.style.overflow = ''
    }
  }, [visible])

  if (!visible) return null

  const overlayClass = ['modal-overlay', exiting ? 'modal-overlay--exiting' : ''].filter(Boolean).join(' ')
  const containerClass = ['modal-container', exiting ? 'modal-container--exiting' : '', className].filter(Boolean).join(' ')

  return (
    <div
      className={overlayClass}
      ref={overlayRef}
      onClick={(e) => {
        if (e.target === overlayRef.current && onClose) {
          onClose()
        }
      }}
      role="dialog"
      aria-modal="true"
      aria-labelledby={title ? 'modal-title' : undefined}
    >
      <div className={containerClass}>
        {title && (
          <div className="modal-header">
            <h2 id="modal-title" className="modal-title">{title}</h2>
            {onClose && (
              <button className="modal-close" onClick={onClose} aria-label="Close dialog">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                  <line x1="18" y1="6" x2="6" y2="18" />
                  <line x1="6" y1="6" x2="18" y2="18" />
                </svg>
              </button>
            )}
          </div>
        )}
        <div className="modal-body">{children}</div>
        {footer && <div className="modal-footer">{footer}</div>}
      </div>
    </div>
  )
}

export default Modal
