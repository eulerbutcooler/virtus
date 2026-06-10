import './Input.css'

const Input = ({
  type = 'text',
  placeholder,
  value,
  onChange,
  disabled,
  textarea,
  rows = 3,
  className = '',
  style,
  ...rest
}) => {
  const props = {
    className: `input ${className}`.trim(),
    placeholder,
    value,
    onChange,
    disabled,
    style,
    ...rest,
  }

  if (textarea) {
    return <textarea {...props} rows={rows} />
  }

  return <input type={type} {...props} />
}

export default Input
