import './Card.css'

const Card = ({ children, interactive, className = '', style, onClick, ...rest }) => {
  const classes = [
    'card',
    interactive ? 'card--interactive' : '',
    className,
  ].filter(Boolean).join(' ')

  return (
    <div className={classes} style={style} onClick={onClick} {...rest}>
      {children}
    </div>
  )
}

export default Card
