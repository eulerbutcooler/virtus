import { useState } from 'react'
import './Tooltip.css'

const Tooltip = ({ children, content, position = 'bottom', className = '', style, ...rest }) => {
  const [visible, setVisible] = useState(false)

  const classes = ['tooltip-wrapper', className].filter(Boolean).join(' ')
  const contentClasses = ['tooltip__content', `tooltip__content--${position}`].filter(Boolean).join(' ')

  return (
    <div
      className={classes}
      style={style}
      onMouseEnter={() => setVisible(true)}
      onMouseLeave={() => setVisible(false)}
      onFocus={() => setVisible(true)}
      onBlur={() => setVisible(false)}
      {...rest}
    >
      {children}
      {visible && (
        <div className={contentClasses} role="tooltip">
          {content}
        </div>
      )}
    </div>
  )
}

export default Tooltip
