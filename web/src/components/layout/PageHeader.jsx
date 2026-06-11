import './PageHeader.css'

export default function PageHeader({ title, subtitle, actions, action }) {
  const content = actions || action;
  return (
    <div className="page-header">
      <div className="page-header__top">
        <div className="page-header__titles">
          <h1 className="page-header__title">{title}</h1>
          {subtitle && <p className="page-header__subtitle">{subtitle}</p>}
        </div>
        {content && <div className="page-header__actions">{content}</div>}
      </div>
    </div>
  )
}
