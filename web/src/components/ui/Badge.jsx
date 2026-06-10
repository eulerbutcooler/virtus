import './Badge.css'

const Badge = ({ variant = 'neutral', children, className = '', style, ...rest }) => {
  const classes = ['badge', `badge--${variant}`, className].filter(Boolean).join(' ')

  return (
    <span className={classes} style={style} {...rest}>
      {children}
    </span>
  )
}

export default Badge
