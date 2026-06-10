import './EmptyState.css'

const EmptyState = ({ title, description, action, className = '', style, ...rest }) => {
  const classes = ['empty-state', className].filter(Boolean).join(' ')

  return (
    <div className={classes} style={style} {...rest}>
      {title && <h3 className="empty-state__title">{title}</h3>}
      {description && <p className="empty-state__description">{description}</p>}
      {action && <div>{action}</div>}
    </div>
  )
}

export default EmptyState
