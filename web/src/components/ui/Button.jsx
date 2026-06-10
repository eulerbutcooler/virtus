import { useState } from 'react'
import './Button.css'
import Spinner from './Spinner.jsx'

const Button = ({
  variant = 'primary',
  size = 'md',
  children,
  onClick,
  disabled,
  loading,
  type = 'button',
  className = '',
  style,
  ...rest
}) => {
  const [pressed, setPressed] = useState(false)

  const classes = [
    'btn',
    `btn--${size}`,
    `btn--${variant}`,
    loading ? 'loading' : '',
    className,
  ].filter(Boolean).join(' ')

  return (
    <button
      type={type}
      className={classes}
      onClick={onClick}
      disabled={disabled || loading}
      style={style}
      onMouseDown={() => setPressed(true)}
      onMouseUp={() => setPressed(false)}
      onMouseLeave={() => setPressed(false)}
      {...rest}
    >
      {loading && <Spinner size="sm" />}
      {children}
    </button>
  )
}

export default Button
