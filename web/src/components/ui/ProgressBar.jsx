import './ProgressBar.css'

const ProgressBar = ({ value, max = 100, label, variant = 'default', className = '', style, ...rest }) => {
  const pct = Math.min(Math.max((value / max) * 100, 0), 100)
  const classes = ['progressbar', className].filter(Boolean).join(' ')
  const fillClasses = ['progressbar__fill', `progressbar__fill--${variant}`].filter(Boolean).join(' ')

  return (
    <div className={classes} style={style} {...rest}>
      {label && <span className="progressbar__label">{label}</span>}
      <div className="progressbar__track">
        <div
          className={fillClasses}
          style={{ width: `${pct}%` }}
          role="progressbar"
          aria-valuenow={value}
          aria-valuemin={0}
          aria-valuemax={max}
        />
      </div>
    </div>
  )
}

export default ProgressBar
