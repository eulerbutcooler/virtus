import './Spinner.css'

const Spinner = ({ size = 'md', className = '', style, ...rest }) => {
  const classes = ['spinner', `spinner--${size}`, className].filter(Boolean).join(' ')

  return (
    <span className={classes} style={style} {...rest} role="status" aria-label="Loading">
      <span style={{ position: 'absolute', width: 1, height: 1, overflow: 'hidden', clip: 'rect(0,0,0,0)' }}>
        Loading
      </span>
    </span>
  )
}

export default Spinner
