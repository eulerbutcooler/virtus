import './Avatar.css'

const Avatar = ({ name, src, size = 'md', className = '', style, ...rest }) => {
  const classes = ['avatar', `avatar--${size}`, className].filter(Boolean).join(' ')

  const initials = name
    ? name
        .split(' ')
        .map((n) => n[0])
        .join('')
        .slice(0, 2)
        .toUpperCase()
    : '?'

  return (
    <div className={classes} style={style} {...rest} aria-label={name || 'Avatar'}>
      {src ? (
        <img src={src} alt={name || 'Avatar'} loading="lazy" />
      ) : (
        <span>{initials}</span>
      )}
    </div>
  )
}

export default Avatar
